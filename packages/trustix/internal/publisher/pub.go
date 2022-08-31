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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	proto "github.com/golang/protobuf/proto"
	"github.com/lazyledger/smt"
	log "github.com/sirupsen/logrus"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	rpc "github.com/nix-community/trustix/packages/trustix-proto/rpc"
	schema "github.com/nix-community/trustix/packages/trustix-proto/schema"
	"github.com/nix-community/trustix/packages/trustix/internal/constants"
	vlog "github.com/nix-community/trustix/packages/trustix/internal/log"
	"github.com/nix-community/trustix/packages/trustix/internal/protocols"
	sthsig "github.com/nix-community/trustix/packages/trustix/internal/sth"
	"github.com/nix-community/trustix/packages/trustix/internal/storage"
)

func minUint64(x, y uint64) uint64 {
	if x > y {
		return y
	}
	return x
}

func hashSum(pd *protocols.ProtocolDescriptor, b []byte) []byte {
	h := pd.NewHash()
	h.Write(b)
	return h.Sum(nil)
}

type Publisher struct {
	queueMux *sync.Mutex
	store    storage.Storage

	// Storage buckets
	queueBucket  *storage.Bucket
	vLogBucket   *storage.Bucket
	mapBucket    *storage.Bucket
	mapLogBucket *storage.Bucket
	caBucket     *storage.Bucket // Content-addressed
	logBucket    *storage.Bucket // Root-level bucket for log

	pd *protocols.ProtocolDescriptor

	signer    crypto.Signer
	submitMux *sync.Mutex
	logID     string
	// closeChan chan interface{}
}

func NewPublisher(logID string, store storage.Storage, caBucket *storage.Bucket, logBucket *storage.Bucket, signer crypto.Signer, pd *protocols.ProtocolDescriptor) (*Publisher, error) {

	qm := &Publisher{
		store:     store,
		signer:    signer,
		logID:     logID,
		queueMux:  &sync.Mutex{},
		submitMux: &sync.Mutex{},
		// closeChan: make(chan interface{}),

		// Storage buckets
		logBucket:    logBucket,
		queueBucket:  logBucket.Cd(constants.QueueBucket),
		vLogBucket:   logBucket.Cd(constants.VLogBucket),
		mapBucket:    logBucket.Cd(constants.MapBucket),
		mapLogBucket: logBucket.Cd(constants.VMapLogBucket),
		caBucket:     caBucket,

		pd: pd,
	}

	// Ensure STH for an empty tree
	err := store.Update(func(txn storage.Transaction) error {

		logBucketTxn := qm.logBucket.Txn(txn)

		sth, err := storage.GetLogHead(logBucketTxn)
		if err == nil {
			return nil
		}
		if err != storage.ObjectNotFoundError {
			return err
		}

		vLog, err := vlog.NewVerifiableLog(qm.pd.NewHash, qm.vLogBucket.Txn(txn), 0)
		if err != nil {
			return err
		}

		smTree := smt.NewSparseMerkleTree(qm.mapBucket.Txn(txn), pd.NewHash())

		vMapLog, err := vlog.NewVerifiableLog(qm.pd.NewHash, qm.mapLogBucket.Txn(txn), 0)
		if err != nil {
			return err
		}

		log.Debug("Signing STH for empty tree")
		sth, err = sthsig.SignHead(vLog, smTree, vMapLog, signer, qm.pd)
		if err != nil {
			return err
		}

		log.WithField("size", *sth.TreeSize).Debug("Setting STH for empty tree")
		{
			buf, err := proto.Marshal(sth)
			if err != nil {
				return err
			}

			return logBucketTxn.Set([]byte(constants.HeadBlob), buf)
		}
	})
	if err != nil {
		return nil, err
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

	return qm, nil
}

func (qm *Publisher) Submit(ctx context.Context, req *rpc.SubmitRequest) (*rpc.SubmitResponse, error) {
	if *req.LogID != qm.logID {
		return nil, fmt.Errorf("Log ID mismatch")
	}

	qm.queueMux.Lock()
	defer qm.queueMux.Unlock()

	err := qm.store.Update(func(txn storage.Transaction) error {
		var err error

		queueBucketTxn := qm.queueBucket.Txn(txn)

		// Get the current state of the queue
		q, err := qm.getQueueMeta(queueBucketTxn)
		if err != nil {
			return err
		}

		// Write each item to the DB while updating queue state
		for _, pair := range req.Items {
			itemId := *q.Max
			err = qm.writeQueueItem(queueBucketTxn, int(itemId), pair)
			if err != nil {
				return err
			}
			next := itemId + 1
			q.Max = &next
		}

		// Write queue state
		err = qm.setQueueMeta(queueBucketTxn, q)
		if err != nil {
			return err
		}

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

	logBucketTxn := qm.logBucket.Txn(txn)

	sth, err := storage.GetLogHead(logBucketTxn)
	if err != nil {
		return err
	}

	// The append-only log
	vLogBucketTxn := qm.vLogBucket.Txn(txn)
	log.WithField("size", *sth.TreeSize).Debug("Creating log tree from persisted data")
	vLog, err := vlog.NewVerifiableLog(qm.pd.NewHash, vLogBucketTxn, *sth.TreeSize)
	if err != nil {
		return err
	}

	// The sparse merkle tree
	log.Debug("Creating sparse merkle tree from persisted data")
	mapBucketTxn := qm.mapBucket.Txn(txn)
	smTree := smt.ImportSparseMerkleTree(mapBucketTxn, qm.pd.NewHash(), sth.MapRoot)

	// The append-only log tracking published map heads
	vMapLogBucketTxn := qm.mapLogBucket.Txn(txn)
	log.WithField("size", *sth.MHTreeSize).Debug("Creating log tree from persisted data")
	vMapLog, err := vlog.NewVerifiableLog(qm.pd.NewHash, vMapLogBucketTxn, *sth.MHTreeSize)
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

			valueDigest := hashSum(qm.pd, pair.Value)
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
		{
			digest := hashSum(qm.pd, pair.Value)
			err = qm.caBucket.Txn(txn).Set(digest[:], pair.Value)
			if err != nil {
				return err
			}
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
	sth, err = sthsig.SignHead(vLog, smTree, vMapLog, qm.signer, qm.pd)
	if err != nil {
		return err
	}

	log.WithField("size", *sth.TreeSize).Debug("Setting new signed tree heads")
	{
		buf, err := proto.Marshal(sth)
		if err != nil {
			return err
		}

		return logBucketTxn.Set([]byte(constants.HeadBlob), buf)
	}
}

func (qm *Publisher) setQueueMeta(txn *storage.BucketTransaction, q *schema.SubmitQueue) error {
	qBytes, err := proto.Marshal(q)
	if err != nil {
		return err
	}

	err = txn.Set([]byte(constants.QueueMetaBlob), qBytes)
	if err != nil {
		return err
	}

	return nil
}

func (qm *Publisher) writeQueueItem(txn *storage.BucketTransaction, itemID int, item *api.KeyValuePair) error {
	itemBytes, err := proto.Marshal(item)
	if err != nil {
		return err
	}

	err = txn.Set([]byte(fmt.Sprintf("%d", itemID)), itemBytes)
	if err != nil {
		return err
	}

	return nil
}

func (qm *Publisher) popQueueItem(txn *storage.BucketTransaction, itemID int) (*api.KeyValuePair, error) {
	itemBytes, err := txn.Get([]byte(fmt.Sprintf("%d", itemID)))
	if err != nil {
		return nil, err
	}

	item := &api.KeyValuePair{}
	err = proto.Unmarshal(itemBytes, item)
	if err != nil {
		return nil, err
	}

	err = txn.Delete([]byte(fmt.Sprintf("%d", itemID)))
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (qm *Publisher) getQueueMeta(txn *storage.BucketTransaction) (*schema.SubmitQueue, error) {
	q := &schema.SubmitQueue{}

	qBytes, err := txn.Get([]byte(constants.QueueMetaBlob))
	if err != nil && err == storage.ObjectNotFoundError {
		min := uint64(0)
		max := uint64(0)
		q.Min = &min
		q.Max = &max
		return q, nil
	} else if err != nil {
		return nil, err
	}

	err = proto.Unmarshal(qBytes, q)
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (qm *Publisher) submitBatch() (*schema.SubmitQueue, error) {
	qm.queueMux.Lock()
	defer qm.queueMux.Unlock()

	q := &schema.SubmitQueue{}

	err := qm.store.Update(func(txn storage.Transaction) error {
		var err error

		queueBucketTxn := qm.queueBucket.Txn(txn)

		// Get the current state of the queue
		q, err = qm.getQueueMeta(queueBucketTxn)
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

			item, err := qm.popQueueItem(queueBucketTxn, int(itemId))
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

		err = qm.setQueueMeta(queueBucketTxn, q)
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
	// qm.closeChan <- nil
}
