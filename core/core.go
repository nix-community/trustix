package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/lazyledger/smt"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/correlator"
	"github.com/tweag/trustix/signer"
	"github.com/tweag/trustix/sth"
	"github.com/tweag/trustix/storage"
	"github.com/tweag/trustix/transport"
	"hash"
	"time"
)

type FlagConfig struct {
	StateDirectory string
}

type TrustixCore struct {
	store      storage.TrustixStorage
	hasher     hash.Hash
	signer     signer.TrustixSigner
	correlator correlator.LogCorrelator
	root       []byte
}

func (s *TrustixCore) Query(key []byte) ([]byte, error) {
	var buf []byte

	err := s.store.View(func(txn storage.Transaction) error {
		tree := smt.ImportSparseMerkleTree(newMapStore(txn), s.hasher, s.root)

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

		if !smt.VerifyProof(proof, tree.Root(), key, v, s.hasher) {
			return fmt.Errorf("Proof verification failed")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *TrustixCore) Get(key []byte) ([]byte, error) {
	var buf []byte

	err := s.store.View(func(txn storage.Transaction) error {
		v, err := txn.Get(key)
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
		mapStore := newMapStore(txn)
		tree := smt.ImportSparseMerkleTree(mapStore, s.hasher, s.root)

		sthManager := sth.NewSTHManager(tree, s.signer)

		tree.Update(key, value)

		sth, err := sthManager.Sign()
		if err != nil {
			return err
		}

		s.root = tree.Root()

		return mapStore.Set([]byte("HEAD"), sth)
	})
}

func (s *TrustixCore) updateRoot() error {
	return s.store.View(func(txn storage.Transaction) error {
		mapStore := newMapStore(txn)
		tree := smt.ImportSparseMerkleTree(mapStore, s.hasher, s.root)

		oldHead, err := txn.Get([]byte("HEAD"))
		if err != nil {
			return err
		} else {
			oldSTH := &sth.STH{}
			err = oldSTH.FromJSON(oldHead)
			if err != nil {
				return err
			}

			rootBytes, err := oldSTH.UnmarshalRoot()
			if err != nil {
				return err
			}

			sigBytes, err := oldSTH.UnmarshalSignature()
			if err != nil {
				return err
			}

			if !s.signer.Verify(rootBytes, sigBytes) {
				return fmt.Errorf("STH signature verification failed")
			}

			if bytes.Compare(rootBytes, tree.Root()) != 0 {
				s.root = rootBytes
				fmt.Println("Updated root")
			}
		}

		return nil
	})
}

func CoreFromConfig(conf *config.LogConfig, flags *FlagConfig) (*TrustixCore, error) {

	hasher := sha256.New()

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
		store, err = transport.NewGRPCTransport(conf.Transport.GRPC)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("Mode '%s' unhandled", conf.Mode)
	}

	var tree *smt.SparseMerkleTree

	corr, err := correlator.NewMinimumPercentCorrelator(100)
	if err != nil {
		return nil, err
	}

	var root []byte
	err = store.View(func(txn storage.Transaction) error {
		mapStore := newMapStore(txn)

		oldHead, err := txn.Get([]byte("HEAD"))
		if err != nil {
			// No STH yet, new tree
			// TODO: Create a completely separate command for new tree, no magic should happen at startup
			if err == storage.ObjectNotFoundError {
				tree = smt.NewSparseMerkleTree(mapStore, hasher)
			} else {
				return err
			}
			root = tree.Root()
		} else {
			oldSTH := &sth.STH{}
			err = oldSTH.FromJSON(oldHead)
			if err != nil {
				return err
			}

			rootBytes, err := oldSTH.UnmarshalRoot()
			if err != nil {
				return err
			}
			root = rootBytes

			sigBytes, err := oldSTH.UnmarshalSignature()
			if err != nil {
				return err
			}

			if !sig.Verify(rootBytes, sigBytes) {
				return fmt.Errorf("STH signature verification failed")
			}

		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	core := &TrustixCore{
		store:      store,
		hasher:     hasher,
		signer:     sig,
		correlator: corr,
		root:       root,
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
