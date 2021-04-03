// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package api

import (
	"context"
	"crypto"
	"crypto/sha256"
	"fmt"

	"github.com/lazyledger/smt"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	vlog "github.com/tweag/trustix/packages/trustix/log"
	sthsig "github.com/tweag/trustix/packages/trustix/sth"
	"github.com/tweag/trustix/packages/trustix/storage"
	storageapi "github.com/tweag/trustix/packages/trustix/storage/api"
)

type KvStoreLogApi struct {
	store  storage.TrustixStorage
	signer crypto.Signer
	logID  string
}

// NewKVStoreAPI - Returns an instance of the log API for an authoritive log implemented on top
// of a key/value store
//
// This is the underlying implementation used by all other abstractions
func NewKVStoreAPI(logID string, store storage.TrustixStorage, signer crypto.Signer) (*KvStoreLogApi, error) {

	var sth *schema.STH

	// Create an empty initial log if it doesn't exist already
	err := store.Update(func(txn storage.Transaction) error {
		storageAPI := storageapi.NewStorageAPI(txn)

		_, err := storageAPI.GetSTH(logID)
		if err == nil {
			return nil
		}
		if err != storage.ObjectNotFoundError {
			return err
		}

		vLog, err := vlog.NewVerifiableLog(logID, txn, 0)
		if err != nil {
			return err
		}

		smTree := smt.NewSparseMerkleTree(storageAPI.MapStore(logID), sha256.New())

		vMapLog, err := vlog.NewVerifiableLog("maplog", txn, 0)
		if err != nil {
			return err
		}

		log.Debug("Signing STH for empty tree")
		sth, err = sthsig.SignHead(vLog, smTree, vMapLog, signer)
		if err != nil {
			return err
		}

		log.WithField("size", *sth.TreeSize).Debug("Setting STH for empty tree")
		err = storageAPI.SetSTH(logID, sth)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &KvStoreLogApi{
		store:  store,
		signer: signer,
		logID:  logID,
	}, nil
}

func (kv *KvStoreLogApi) GetSTH(ctx context.Context, req *api.STHRequest) (*schema.STH, error) {
	var sth *schema.STH
	err := kv.store.View(func(txn storage.Transaction) error {
		var err error
		storageAPI := storageapi.NewStorageAPI(txn)
		sth, err = storageAPI.GetSTH(kv.logID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return sth, nil
}

func (kv *KvStoreLogApi) GetLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (resp *api.ProofResponse, err error) {
	resp = &api.ProofResponse{}
	err = kv.store.View(func(txn storage.Transaction) error {
		resp, err = getLogConsistencyProof("log", txn, ctx, req)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (kv *KvStoreLogApi) GetLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (resp *api.ProofResponse, err error) {
	resp = &api.ProofResponse{}
	err = kv.store.View(func(txn storage.Transaction) error {
		resp, err = getLogAuditProof("log", txn, ctx, req)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (kv *KvStoreLogApi) GetLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (resp *api.LogEntriesResponse, err error) {
	resp = &api.LogEntriesResponse{
		Leaves: []*schema.LogLeaf{},
	}

	err = kv.store.View(func(txn storage.Transaction) error {
		resp, err = getLogEntries("log", txn, ctx, req)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (kv *KvStoreLogApi) GetMapValue(ctx context.Context, req *api.GetMapValueRequest) (*api.MapValueResponse, error) {

	resp := &api.MapValueResponse{}

	err := kv.store.View(func(txn storage.Transaction) error {
		storageAPI := storageapi.NewStorageAPI(txn)

		tree := smt.ImportSparseMerkleTree(storageAPI.MapStore(kv.logID), sha256.New(), req.MapRoot)

		v, err := tree.Get(req.Key)
		if err != nil {
			return err
		}

		if len(v) == 0 {
			return fmt.Errorf("Map value not found")
		}

		proof, err := tree.ProveCompact(req.Key)
		if err != nil {
			return err
		}

		numSideNodes := uint64(proof.NumSideNodes)
		resp.Value = v

		resp.Proof = &api.SparseCompactMerkleProof{
			SideNodes:             proof.SideNodes,
			NonMembershipLeafData: proof.NonMembershipLeafData,
			BitMask:               proof.BitMask,
			NumSideNodes:          &numSideNodes,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (kv *KvStoreLogApi) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {
	var value []byte
	err := kv.store.View(func(txn storage.Transaction) error {
		v, err := storageapi.NewStorageAPI(txn).GetCAValue(in.Digest)
		value = v
		return err

	})
	if err != nil {
		return nil, err
	}

	return &api.ValueResponse{
		Value: value,
	}, nil
}

func (kv *KvStoreLogApi) GetMHLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (resp *api.ProofResponse, err error) {
	resp = &api.ProofResponse{}
	err = kv.store.View(func(txn storage.Transaction) error {
		resp, err = getLogConsistencyProof("log", txn, ctx, req)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (kv *KvStoreLogApi) GetMHLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (resp *api.ProofResponse, err error) {
	resp = &api.ProofResponse{}
	err = kv.store.View(func(txn storage.Transaction) error {
		resp, err = getLogAuditProof("log", txn, ctx, req)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (kv *KvStoreLogApi) GetMHLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (resp *api.LogEntriesResponse, err error) {
	resp = &api.LogEntriesResponse{
		Leaves: []*schema.LogLeaf{},
	}

	err = kv.store.View(func(txn storage.Transaction) error {
		resp, err = getLogEntries("log", txn, ctx, req)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
