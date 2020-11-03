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

package api

import (
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	"github.com/lazyledger/smt"
	log "github.com/sirupsen/logrus"
	vlog "github.com/tweag/trustix/log"
	"github.com/tweag/trustix/schema"
	sthsig "github.com/tweag/trustix/sth"
	"github.com/tweag/trustix/storage"
)

type kvStoreLogApi struct {
	store  storage.TrustixStorage
	signer crypto.Signer
	sth    *schema.STH
}

func NewKVStoreAPI(store storage.TrustixStorage, signer crypto.Signer) (TrustixLogAPI, error) {

	var sth *schema.STH

	// Create an empty initial log if it doesn't exist already
	err := store.Update(func(txn storage.Transaction) error {
		_, err := txn.Get([]byte("META"), []byte("HEAD"))
		if err == nil {
			return nil
		}
		if err != storage.ObjectNotFoundError {
			return err
		}

		smTree := smt.NewSparseMerkleTree(newMapStore(txn), sha256.New())
		vLog, err := vlog.NewVerifiableLog(txn, 0)
		if err != nil {
			return err
		}

		log.Debug("Signing STH for empty tree")
		sth, err = sthsig.SignHead(smTree, vLog, signer)
		if err != nil {
			return err
		}

		smhBytes, err := proto.Marshal(sth)
		if err != nil {
			return err
		}

		log.WithField("size", sth.TreeSize).Debug("Setting STH for empty tree")
		err = txn.Set([]byte("META"), []byte("HEAD"), smhBytes)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	api := &kvStoreLogApi{
		store:  store,
		signer: signer,
		sth:    sth,
	}

	if api.sth == nil {
		err := store.View(func(txn storage.Transaction) error {
			sth, err := api.getSTH(txn)
			if err != nil {
				return err
			}
			api.sth = sth
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if api.sth == nil {
		return nil, fmt.Errorf("Could not find STH")
	}

	return api, nil
}

func (kv *kvStoreLogApi) getSTH(txn storage.Transaction) (*schema.STH, error) {
	sth := &schema.STH{}
	var err error

	sthBytes, err := txn.Get([]byte("META"), []byte("HEAD"))
	if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(sthBytes, sth)
	if err != nil {
		return nil, err
	}

	return sth, nil
}

func (kv *kvStoreLogApi) GetSTH(req *STHRequest) (*schema.STH, error) {
	return kv.sth, nil
}

func (kv *kvStoreLogApi) GetLogConsistencyProof(req *GetLogConsistencyProofRequest) (resp *ProofResponse, err error) {

	err = kv.store.View(func(txn storage.Transaction) error {
		vLog, err := vlog.NewVerifiableLog(txn, req.SecondSize)
		if err != nil {
			return err
		}

		proof, err := vLog.ConsistencyProof(req.FirstSize, req.SecondSize)
		if err != nil {
			return err
		}

		resp.Proof = proof

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (kv *kvStoreLogApi) GetLogAuditProof(req *GetLogAuditProofRequest) (resp *ProofResponse, err error) {

	err = kv.store.View(func(txn storage.Transaction) error {
		vLog, err := vlog.NewVerifiableLog(txn, req.TreeSize)
		if err != nil {
			return err
		}

		proof, err := vLog.AuditProof(req.Index, req.TreeSize)
		if err != nil {
			return err
		}

		resp.Proof = proof

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (kv *kvStoreLogApi) GetLogEntries(req *GetLogEntriesRequest) (*LogEntriesResponse, error) {

	resp := &LogEntriesResponse{
		Leaves: []*schema.LogLeaf{},
	}

	err := kv.store.View(func(txn storage.Transaction) error {
		logStorage := vlog.NewLogStorage(txn)

		for i := int(req.Start); i <= int(req.Finish); i++ {
			leaf, err := logStorage.Get(0, uint64(i))
			if err != nil {
				return err
			}
			resp.Leaves = append(resp.Leaves, leaf)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (kv *kvStoreLogApi) GetMapValue(req *GetMapValueRequest) (*MapValueResponse, error) {

	resp := &MapValueResponse{}

	err := kv.store.View(func(txn storage.Transaction) error {
		tree := smt.ImportSparseMerkleTree(newMapStore(txn), sha256.New(), req.MapRoot)

		v, err := tree.Get(req.Key)
		if err != nil {
			return err
		}

		proof, err := tree.ProveCompact(req.Key)
		if err != nil {
			return err
		}

		resp.Value = v
		resp.Proof = &SparseCompactMerkleProof{
			SideNodes:             proof.SideNodes,
			NonMembershipLeafData: proof.NonMembershipLeafData,
			BitMask:               proof.BitMask,
			NumSideNodes:          uint64(proof.NumSideNodes),
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (kv *kvStoreLogApi) Submit(req *SubmitRequest) (*SubmitResponse, error) {

	if kv.signer == nil {
		return nil, fmt.Errorf("Signing is disabled")
	}

	sth := kv.sth

	err := kv.store.Update(func(txn storage.Transaction) error {
		// The sparse merkle tree
		log.Debug("Creating sparse merkle tree from persisted data")
		smTree := smt.ImportSparseMerkleTree(newMapStore(txn), sha256.New(), sth.MapRoot)

		// The append-only log
		log.WithField("size", sth.TreeSize).Debug("Creating log tree from persisted data")
		vLog, err := vlog.NewVerifiableLog(txn, sth.TreeSize)
		if err != nil {
			return err
		}

		for _, pair := range req.Items {

			// Get the old value and check it against new submitted value
			log.Debug("Checking if newly submitted value is already set")
			oldValue, err := smTree.Get(pair.Key)
			if err != nil {
				return err
			}
			if len(oldValue) > 0 {
				return fmt.Errorf("'%s' already exists in log", hex.EncodeToString(pair.Key))
			}

			// Append value to both verifiable log & sparse indexed tree
			log.Debug("Appending value to log")
			err = vLog.Append(pair.Value)
			if err != nil {
				return err
			}

			entry, err := proto.Marshal(&schema.MapEntry{
				Value: pair.Value,
				Index: uint64(vLog.Size() - 1),
			})
			if err != nil {
				return err
			}

			smTree.Update(pair.Key, entry)

		}

		sth, err = sthsig.SignHead(smTree, vLog, kv.signer)
		if err != nil {
			return err
		}

		log.Debug("Signing tree heads")
		smhBytes, err := proto.Marshal(sth)
		if err != nil {
			return err
		}

		log.WithField("size", sth.TreeSize).Debug("Setting new signed tree heads")
		err = txn.Set([]byte("META"), []byte("HEAD"), smhBytes)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	kv.sth = sth

	return &SubmitResponse{
		Status: SubmitResponse_OK,
	}, nil
}
