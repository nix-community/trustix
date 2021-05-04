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
	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	"github.com/tweag/trustix/packages/trustix/constants"
	"github.com/tweag/trustix/packages/trustix/storage"
)

type KvStoreLogApi struct {
	store storage.Storage

	logBucket    *storage.Bucket
	vLogBucket   *storage.Bucket
	caBucket     *storage.Bucket
	mapBucket    *storage.Bucket
	mapLogBucket *storage.Bucket

	signer crypto.Signer
	logID  string
}

// NewKVStoreAPI - Returns an instance of the log API for an authoritive log implemented on top
// of a key/value store
//
// This is the underlying implementation used by all other abstractions
func NewKVStoreAPI(logID string, store storage.Storage, caBucket *storage.Bucket, logBucket *storage.Bucket, signer crypto.Signer) (*KvStoreLogApi, error) {
	return &KvStoreLogApi{
		store:        store,
		logBucket:    logBucket,
		vLogBucket:   logBucket.Cd(constants.VLogBucket),
		caBucket:     caBucket,
		mapBucket:    logBucket.Cd(constants.MapBucket),
		mapLogBucket: logBucket.Cd(constants.VMapLogBucket),
		signer:       signer,
		logID:        logID,
	}, nil
}

func (kv *KvStoreLogApi) GetSTH(ctx context.Context, req *api.STHRequest) (*schema.STH, error) {
	var sth *schema.STH
	err := kv.store.View(func(txn storage.Transaction) error {
		var err error
		sth, err = storage.GetSTH(kv.logBucket.Txn(txn))
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
		var err error

		vLogBucketTxn := kv.vLogBucket.Txn(txn)

		resp, err = getLogConsistencyProof(vLogBucketTxn, ctx, req)
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
		var err error

		vLogBucketTxn := kv.vLogBucket.Txn(txn)

		resp, err = getLogAuditProof(vLogBucketTxn, ctx, req)
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
		var err error

		vLogBucketTxn := kv.vLogBucket.Txn(txn)

		resp, err = getLogEntries(vLogBucketTxn, ctx, req)
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
		mapBucketTxn := kv.mapBucket.Txn(txn)
		tree := smt.ImportSparseMerkleTree(mapBucketTxn, sha256.New(), req.MapRoot)

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
		v, err := getCAValue(kv.caBucket.Txn(txn), in.Digest)
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
		var err error

		vMapLogBucketTxn := kv.mapLogBucket.Txn(txn)

		resp, err = getLogConsistencyProof(vMapLogBucketTxn, ctx, req)
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
		var err error

		vMapLogBucketTxn := kv.mapLogBucket.Txn(txn)

		resp, err = getLogAuditProof(vMapLogBucketTxn, ctx, req)
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
		var err error

		vMapLogBucketTxn := kv.mapLogBucket.Txn(txn)

		resp, err = getLogEntries(vMapLogBucketTxn, ctx, req)
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
