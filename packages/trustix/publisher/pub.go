// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package publisher

import (
	"bytes"
	"context"
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	// "time"

	"github.com/lazyledger/smt"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	rpc "github.com/tweag/trustix/packages/trustix-proto/proto"
	schema "github.com/tweag/trustix/packages/trustix-proto/schema"
	vlog "github.com/tweag/trustix/packages/trustix/log"
	sthsig "github.com/tweag/trustix/packages/trustix/sth"
	"github.com/tweag/trustix/packages/trustix/storage"
	storageapi "github.com/tweag/trustix/packages/trustix/storage/api"
)

func minUint64(x, y uint64) uint64 {
	if x > y {
		return y
	}
	return x
}

type Publisher struct {
	queueMux  *sync.Mutex
	store     storage.TrustixStorage
	signer    crypto.Signer
	submitMux *sync.Mutex
	logID     string
	// sth       *schema.STH
	closeChan chan interface{}
}

func NewPublisher(logID string, store storage.TrustixStorage, signer crypto.Signer) *Publisher {

	qm := &Publisher{
		store:     store,
		signer:    signer,
		logID:     logID,
		queueMux:  &sync.Mutex{},
		submitMux: &sync.Mutex{},
		closeChan: make(chan interface{}),
	}

	// go func() {
	// 	// TODO: Figure out a better method than hard coded sleep
	// 	duration := time.Second * 5

	// 	timeout := time.NewTimer(duration)
	// 	defer timeout.Stop()

	// 	for {
	// 		timeout.Reset(duration)
	// 		select {
	// 		case _ = <-qm.closeChan:
	// 			return
	// 		case <-timeout.C:
	// 			q, err := qm.submitBatch()
	// 			if err != nil {
	// 				log.Errorf("Unable to submit batch: %v", err)
	// 				continue
	// 			}

	// 			if *q.Min >= *q.Max {
	// 				time.Sleep(duration)
	// 			}
	// 		}
	// 	}
	// }()

	// go func() {
	// 	// TODO: Figure out a better method than hard coded sleep
	// 	duration := time.Second * 5
	// 	for {
	// 		q, err := qm.submitBatch()
	// 		if err != nil {
	// 			log.Errorf("Unable to submit batch: %v", err)
	// 			time.Sleep(duration)
	// 			continue
	// 		}

	// 		if *q.Min >= *q.Max {
	// 			time.Sleep(duration)
	// 		}
	// 	}
	// }()

	return qm
}

func (qm *Publisher) Submit(ctx context.Context, req *rpc.SubmitRequest) (*rpc.SubmitResponse, error) {
	if *req.LogID != qm.logID {
		return nil, fmt.Errorf("Log ID mismatch")
	}

	qm.queueMux.Lock()
	defer qm.queueMux.Unlock()

	err := qm.store.Update(func(txn storage.Transaction) error {

		storageAPI := storageapi.NewStorageAPI(txn)

		// Get the current state of the queue
		q, err := storageAPI.GetQueueMeta(qm.logID)
		if err != nil {
			return err
		}

		// Write each item to the DB while updating queue state
		for _, pair := range req.Items {
			itemId := *q.Max
			err = storageAPI.WriteQueueItem(qm.logID, int(itemId), pair)
			if err != nil {
				return err
			}
			next := itemId + 1
			q.Max = &next
		}

		// Write queue state
		storageAPI.SetQueueMeta(qm.logID, q)

		return nil
	})
	if err != nil {
		return nil, err
	}

	status := rpc.SubmitResponse_OK
	return &rpc.SubmitResponse{
		Status: &status,
	}, nil
}

func (qm *Publisher) Flush(ctx context.Context, in *rpc.FlushRequest) (*rpc.FlushResponse, error) {
	if *in.LogID != qm.logID {
		return nil, fmt.Errorf("Log ID mismatch")
	}

	for {
		q, err := qm.submitBatch()
		if err != nil {
			return nil, err
		}

		if *q.Min >= *q.Max {
			return &rpc.FlushResponse{}, nil
		}
	}
}

func (qm *Publisher) writeItems(txn storage.Transaction, items []*api.KeyValuePair) error {
	qm.submitMux.Lock()
	defer qm.submitMux.Unlock()

	if qm.signer == nil {
		return fmt.Errorf("Signing is disabled")
	}

	storageAPI := storageapi.NewStorageAPI(txn)

	sth, err := storageAPI.GetSTH(qm.logID)
	if err != nil {
		return err
	}

	// The append-only log
	log.WithField("size", *sth.TreeSize).Debug("Creating log tree from persisted data")
	vLog, err := vlog.NewVerifiableLog("log", txn, *sth.TreeSize)
	if err != nil {
		return err
	}

	// The sparse merkle tree
	log.Debug("Creating sparse merkle tree from persisted data")
	smTree := smt.ImportSparseMerkleTree(storageAPI.MapStore(qm.logID), sha256.New(), sth.MapRoot)

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
		err = storageapi.NewStorageAPI(txn).SetCAValue(pair.Value)
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

	log.Debug("Signing tree heads")
	sth, err = sthsig.SignHead(vLog, smTree, vMapLog, qm.signer)
	if err != nil {
		return err
	}

	log.WithField("size", *sth.TreeSize).Debug("Setting new signed tree heads")
	err = storageapi.NewStorageAPI(txn).SetSTH(qm.logID, sth)
	if err != nil {
		return err
	}

	return nil
}

func (qm *Publisher) submitBatch() (*schema.SubmitQueue, error) {
	qm.queueMux.Lock()
	defer qm.queueMux.Unlock()

	q := &schema.SubmitQueue{}

	err := qm.store.Update(func(txn storage.Transaction) error {
		var err error
		storageAPI := storageapi.NewStorageAPI(txn)

		// Get the current state of the queue
		q, err = storageAPI.GetQueueMeta(qm.logID)
		if err != nil {
			return err
		}

		maxBatchSize := uint64(500)
		max := minUint64(*q.Max, *q.Min+maxBatchSize)
		min := *q.Min
		if min >= max {
			return nil
		}

		items := []*api.KeyValuePair{}

		for itemId := min; itemId < max; itemId++ {
			q.Min = &itemId

			item, err := storageAPI.PopQueueItem(qm.logID, int(itemId))
			if err != nil {
				log.Error(fmt.Errorf("Error popping queue item '%d': %v", itemId, err))
				continue
			}

			items = append(items, item)
		}

		if len(items) == 0 {
			return nil
		}
		err = qm.writeItems(txn, items)
		if err != nil {
			return err
		}

		err = storageAPI.SetQueueMeta(qm.logID, q)
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

func (qm *Publisher) Close() {
	qm.closeChan <- nil
}
