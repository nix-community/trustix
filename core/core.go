// MIT License
//
// Copyright (c) 2020 Tweag IO
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/lazyledger/smt"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/config"
	vlog "github.com/tweag/trustix/log"
	"github.com/tweag/trustix/schema"
	"github.com/tweag/trustix/signer"
	"github.com/tweag/trustix/sth"
	"github.com/tweag/trustix/storage"
	"github.com/tweag/trustix/transport"
	"time"
)

type FlagConfig struct {
	StateDirectory string
}

type TrustixCore struct {
	store  storage.TrustixStorage
	signer signer.TrustixSigner

	mapRoot  []byte
	logRoot  []byte
	treeSize int
}

func (s *TrustixCore) Query(key []byte) (*schema.MapEntry, error) {
	var buf []byte

	err := s.store.View(func(txn storage.Transaction) error {
		tree := smt.ImportSparseMerkleTree(newMapStore(txn), sha256.New(), s.mapRoot)

		// TODO: Log verification (but optional?)

		v, err := tree.Get(key)
		if err != nil {
			return err
		}
		buf = v

		// Check inclusion proof
		proof, err := tree.Prove(key)
		if err != nil {
			return err
		}

		if !smt.VerifyProof(proof, tree.Root(), key, v, sha256.New()) {
			return fmt.Errorf("Proof verification failed")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	e := &schema.MapEntry{}
	err = proto.Unmarshal(buf, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (s *TrustixCore) Get(bucket []byte, key []byte) ([]byte, error) {
	var buf []byte

	err := s.store.View(func(txn storage.Transaction) error {
		v, err := txn.Get(bucket, key)
		if err != nil {
			return err
		}
		buf = v

		return nil
	})
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *TrustixCore) Submit(key []byte, value []byte) error {
	return s.store.Update(func(txn storage.Transaction) error {

		// The sparse merkle tree
		log.Debug("Creating sparse merkle tree from persisted data")
		smTree := smt.ImportSparseMerkleTree(newMapStore(txn), sha256.New(), s.mapRoot)

		// Get the old value and check it against new submitted value
		log.Debug("Checking if newly submitted value is already set")
		oldValue, err := smTree.Get(key)
		if err != nil {
			return err
		}
		if len(oldValue) > 0 {
			return fmt.Errorf("'%s' already exists in log", hex.EncodeToString(key))
		}

		// The append-only log
		log.WithField("size", s.treeSize).Debug("Creating log tree from persisted data")
		vLog, err := vlog.NewVerifiableLog(txn, s.treeSize)
		if err != nil {
			return err
		}

		// Append value to both verifiable log & sparse indexed tree
		log.Debug("Appending value to log")
		vLog.Append(value)

		entry, err := proto.Marshal(&schema.MapEntry{
			Value: value,
			Index: uint64(vLog.Size() - 1),
		})
		if err != nil {
			return err
		}

		smTree.Update(key, entry)

		sth, err := sth.SignHead(smTree, vLog, s.signer)
		if err != nil {
			return err
		}

		log.Debug("Signing tree heads")
		smhBytes, err := sth.Marshal()
		if err != nil {
			return err
		}

		mapRoot, err := sth.UnmarshalSMHRoot()
		if err != nil {
			return err
		}

		logRoot, err := sth.UnmarshalSTHRoot()
		if err != nil {
			return err
		}

		log.Debug("Setting new signed tree heads")
		err = txn.Set([]byte("META"), []byte("HEAD"), smhBytes)
		if err != nil {
			return err
		}

		s.mapRoot = mapRoot
		s.logRoot = logRoot
		s.treeSize = vLog.Size()

		return nil
	})
}

func (s *TrustixCore) updateRoot() error {
	return s.store.View(func(txn storage.Transaction) error {
		log.Debug("Updating tree root")

		oldHead, err := txn.Get([]byte("META"), []byte("HEAD"))
		if err != nil {
			return err
		} else {
			oldSMH, err := sth.NewSMHFromJSON(oldHead)
			if err != nil {
				return err
			}
			sthRootBytes, err := oldSMH.UnmarshalSTHRoot()
			if err != nil {
				return err
			}

			// Verify signed map head
			smhRootBytes, err := oldSMH.UnmarshalSMHRoot()
			if err != nil {
				return err
			}

			err = oldSMH.Verify(s.signer)
			if err != nil {
				return err
			}

			s.mapRoot = smhRootBytes
			s.logRoot = sthRootBytes
			s.treeSize = oldSMH.LogSth.Size

		}

		return nil
	})
}

func CoreFromConfig(conf *config.LogConfig, flags *FlagConfig) (*TrustixCore, error) {

	sig, err := signer.FromConfig(conf.Signer)
	if err != nil {
		return nil, err
	}

	if conf.Mode == "trustix-log" {
		if !sig.CanSign() {
			return nil, fmt.Errorf("Cannot sign using the current configuration, aborting.")
		}
	}

	var store storage.TrustixStorage
	switch conf.Mode {
	case "trustix-log":
		store, err = storage.FromConfig(conf.Name, flags.StateDirectory, conf.Storage)
		if err != nil {
			return nil, err
		}
	case "trustix-follower":
		store, err = transport.NewGRPCTransport(conf.Transport.GRPC, sig)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Mode '%s' unhandled", conf.Mode)
	}

	core := &TrustixCore{
		store:  store,
		signer: sig,
		// TODO: Log root
	}

	err = store.View(func(txn storage.Transaction) error {
		_, err := txn.Get([]byte("META"), []byte("HEAD"))
		if err != nil {

			// No STH yet, new tree
			// TODO: Create a completely separate command for new tree, no magic should happen at startup
			if err == storage.ObjectNotFoundError {
				hasher := sha256.New()
				tree := smt.NewSparseMerkleTree(newMapStore(txn), hasher)
				core.mapRoot = tree.Root()
				return nil
			}
			return err
		} else {
			// TODO: Implement local cache and set to old values so we can verify consistency between last known good HEAD
			// and the newest HEAD
			return core.updateRoot()
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	switch conf.Mode {
	case "trustix-follower":
		go func() {
			for {
				// TODO: This is just an arbitrary interval for testing, not particularly intelligent
				time.Sleep(10 * time.Second)
				err := core.updateRoot()
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	default:
	}

	return core, nil
}
