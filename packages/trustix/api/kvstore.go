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
	"bytes"
	"context"
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/lazyledger/smt"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	ca "github.com/tweag/trustix/packages/trustix/cavaluestore"
	vlog "github.com/tweag/trustix/packages/trustix/log"
	sthsig "github.com/tweag/trustix/packages/trustix/sth"
	"github.com/tweag/trustix/packages/trustix/storage"
)

type kvStoreLogApi struct {
	store     storage.TrustixStorage
	signer    crypto.Signer
	sth       *schema.STH
	queueMux  *sync.Mutex
	submitMux *sync.Mutex
}

func minUint64(x, y uint64) uint64 {
	if x > y {
		return y
	}
	return x
}

// NewKVStoreAPI - Returns an instance of the log API for an authoritive log implemented on top
// of a key/value store
//
// This is the underlying implementation used by all other abstractions
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

		vLog, err := vlog.NewVerifiableLog("log", txn, 0)
		if err != nil {
			return err
		}

		smTree := smt.NewSparseMerkleTree(newMapStore(txn), sha256.New())

		vMapLog, err := vlog.NewVerifiableLog("maplog", txn, 0)
		if err != nil {
			return err
		}

		log.Debug("Signing STH for empty tree")
		sth, err = sthsig.SignHead(vLog, smTree, vMapLog, signer)
		if err != nil {
			return err
		}

		smhBytes, err := proto.Marshal(sth)
		if err != nil {
			return err
		}

		log.WithField("size", *sth.TreeSize).Debug("Setting STH for empty tree")
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
		store:     store,
		signer:    signer,
		sth:       sth,
		queueMux:  &sync.Mutex{},
		submitMux: &sync.Mutex{},
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

	go func() {
		// TODO: Figure out a better method than hard coded sleep
		duration := time.Second * 5
		for {
			q, err := api.submitBatch()
			if err != nil {
				log.Error(err)
				time.Sleep(duration)
				continue
			}

			if *q.Min >= *q.Max {
				time.Sleep(duration)
			}
		}
	}()

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

func (kv *kvStoreLogApi) GetSTH(ctx context.Context, req *api.STHRequest) (*schema.STH, error) {
	return kv.sth, nil
}

func (kv *kvStoreLogApi) GetLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (resp *api.ProofResponse, err error) {
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

func (kv *kvStoreLogApi) GetLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (resp *api.ProofResponse, err error) {
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

func (kv *kvStoreLogApi) GetLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (resp *api.LogEntriesResponse, err error) {
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

func (kv *kvStoreLogApi) GetMapValue(ctx context.Context, req *api.GetMapValueRequest) (*api.MapValueResponse, error) {

	resp := &api.MapValueResponse{}

	err := kv.store.View(func(txn storage.Transaction) error {
		tree := smt.ImportSparseMerkleTree(newMapStore(txn), sha256.New(), req.MapRoot)

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

func (kv *kvStoreLogApi) Submit(ctx context.Context, req *api.SubmitRequest) (*api.SubmitResponse, error) {
	kv.queueMux.Lock()
	defer kv.queueMux.Unlock()

	err := kv.store.Update(func(txn storage.Transaction) error {

		// Get the current state of the queue
		q := &schema.SubmitQueue{}
		qBytes, err := txn.Get([]byte("QUEUE"), []byte("META"))
		if err != nil && err == storage.ObjectNotFoundError {
			min := uint64(0)
			max := uint64(0)
			q.Min = &min
			q.Max = &max
		} else if err != nil {
			return err
		} else {
			err = proto.Unmarshal(qBytes, q)
			if err != nil {
				return err
			}
		}

		// Write each item to the DB while updating queue state
		for _, pair := range req.Items {
			itemBytes, err := proto.Marshal(pair)
			if err != nil {
				return err
			}

			itemId := *q.Max
			err = txn.Set([]byte("QUEUE"), []byte(fmt.Sprintf("%d", itemId)), itemBytes)
			if err != nil {
				return err
			}

			next := itemId + 1
			q.Max = &next
		}

		// Write queue state
		qBytes, err = proto.Marshal(q)
		if err != nil {
			return err
		}
		err = txn.Set([]byte("QUEUE"), []byte("META"), qBytes)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	status := api.SubmitResponse_OK
	return &api.SubmitResponse{
		Status: &status,
	}, nil
}

func (kv *kvStoreLogApi) Flush(ctx context.Context, in *api.FlushRequest) (*api.FlushResponse, error) {
	for {
		q, err := kv.submitBatch()
		if err != nil {
			return nil, err
		}

		if *q.Min >= *q.Max {
			return &api.FlushResponse{}, nil
		}
	}
}

func (kv *kvStoreLogApi) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {
	var value []byte
	err := kv.store.View(func(txn storage.Transaction) error {
		v, err := ca.Get(txn, in.Digest)
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

func (kv *kvStoreLogApi) submitBatch() (*schema.SubmitQueue, error) {
	kv.queueMux.Lock()
	defer kv.queueMux.Unlock()

	q := &schema.SubmitQueue{}

	err := kv.store.Update(func(txn storage.Transaction) error {

		// Get the current state of the queue
		qBytes, err := txn.Get([]byte("QUEUE"), []byte("META"))
		if err != nil && err == storage.ObjectNotFoundError {
			min := uint64(0)
			max := uint64(0)
			q.Min = &min
			q.Max = &max
			return nil
		} else if err != nil {
			return err
		}
		err = proto.Unmarshal(qBytes, q)
		if err != nil {
			return err
		}

		maxBatchSize := uint64(500)

		items := []*api.KeyValuePair{}
		max := minUint64(*q.Max-1, *q.Min+maxBatchSize)
		for itemId := *q.Min; itemId <= max; itemId++ {
			q.Min = &itemId

			itemBytes, err := txn.Get([]byte("QUEUE"), []byte(fmt.Sprintf("%d", itemId)))
			if err != nil {
				log.Error(err)
				continue
			}

			item := &api.KeyValuePair{}
			err = proto.Unmarshal(itemBytes, item)
			if err != nil {
				log.Error(err)
				continue
			}

			err = txn.Delete([]byte("QUEUE"), []byte(fmt.Sprintf("%d", itemId)))
			if err != nil {
				log.Error(err)
				continue
			}

			items = append(items, item)
		}

		if len(items) == 0 {
			return nil
		}
		err = kv.writeItems(txn, items)
		if err != nil {
			return err
		}

		// Write queue state
		qBytes, err = proto.Marshal(q)
		if err != nil {
			return err
		}
		err = txn.Set([]byte("QUEUE"), []byte("META"), qBytes)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (kv *kvStoreLogApi) writeItems(txn storage.Transaction, items []*api.KeyValuePair) error {
	kv.submitMux.Lock()
	defer kv.submitMux.Unlock()

	if kv.signer == nil {
		return fmt.Errorf("Signing is disabled")
	}

	sth := kv.sth

	// The append-only log
	log.WithField("size", *sth.TreeSize).Debug("Creating log tree from persisted data")
	vLog, err := vlog.NewVerifiableLog("log", txn, *sth.TreeSize)
	if err != nil {
		return err
	}

	// The sparse merkle tree
	log.Debug("Creating sparse merkle tree from persisted data")
	smTree := smt.ImportSparseMerkleTree(newMapStore(txn), sha256.New(), sth.MapRoot)

	// The append-only log tracking published map heads
	log.WithField("size", *sth.MHTreeSize).Debug("Creating log tree from persisted data")
	vMapLog, err := vlog.NewVerifiableLog("maplog", txn, *sth.MHTreeSize)
	if err != nil {
		return err
	}

	wrote := false

	for _, pair := range items {

		// Get the old value and check it against new submitted value
		log.Debug("Checking if newly submitted value is already set")
		oldValue, err := smTree.Get(pair.Key)
		if err != nil {
			return err
		}
		if len(oldValue) > 0 {
			oldEntry := &schema.MapEntry{}
			err = json.Unmarshal(oldValue, oldEntry)
			if err != nil {
				return err
			}

			valueDigest := sha256.Sum256(pair.Value)
			if bytes.Equal(oldEntry.Digest, valueDigest[:]) {
				continue
			}

			log.WithFields(log.Fields{
				"key":   hex.EncodeToString(pair.Key),
				"value": hex.EncodeToString(pair.Value),
			}).Error("Already exists in log with a different value")
			continue
		}

		wrote = true

		// Add value to content-addressed value store
		err = ca.Set(txn, pair.Value)
		if err != nil {
			return err
		}

		// Append value to both verifiable log & sparse indexed tree
		log.Debug("Appending value to log")
		leaf, err := vLog.AppendKV(pair.Key, pair.Value)
		if err != nil {
			return err
		}

		vLogSize := uint64(vLog.Size() - 1)
		entry, err := json.Marshal(&schema.MapEntry{
			Digest: leaf.ValueDigest,
			Index:  &vLogSize,
		})
		if err != nil {
			return err
		}

		smTree.Update(pair.Key, entry)

	}

	if !wrote {
		log.WithField("size", *sth.TreeSize).Debug("Nothing written, skipping head signatures")
		return nil
	}

	sth, err = sthsig.SignHead(vLog, smTree, vMapLog, kv.signer)
	if err != nil {
		return err
	}

	log.Debug("Signing tree heads")
	smhBytes, err := proto.Marshal(sth)
	if err != nil {
		return err
	}

	log.WithField("size", *sth.TreeSize).Debug("Setting new signed tree heads")
	err = txn.Set([]byte("META"), []byte("HEAD"), smhBytes)
	if err != nil {
		return err
	}

	kv.sth = sth

	return nil
}

func (kv *kvStoreLogApi) GetMHLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (resp *api.ProofResponse, err error) {
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

func (kv *kvStoreLogApi) GetMHLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (resp *api.ProofResponse, err error) {
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

func (kv *kvStoreLogApi) GetMHLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (resp *api.LogEntriesResponse, err error) {
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
