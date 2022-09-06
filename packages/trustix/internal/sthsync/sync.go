// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package sthsync

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	apipb "github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
	"github.com/nix-community/trustix/packages/trustix/internal/constants"
	vlog "github.com/nix-community/trustix/packages/trustix/internal/log"
	"github.com/nix-community/trustix/packages/trustix/internal/pool"
	"github.com/nix-community/trustix/packages/trustix/internal/protocols"
	"github.com/nix-community/trustix/packages/trustix/internal/signer"
	sthlib "github.com/nix-community/trustix/packages/trustix/internal/sth"
	"github.com/nix-community/trustix/packages/trustix/internal/storage"
	log "github.com/sirupsen/logrus"
	proto "google.golang.org/protobuf/proto"
)

type sthSyncer struct {
	store     storage.Storage
	logID     string
	closeChan chan interface{}
}

func NewSTHSyncer(
	logID string,
	store storage.Storage,
	logBucket *storage.Bucket,
	clients *pool.ClientPool,
	verifier signer.Verifier,
	pollInterval time.Duration,
	pd *protocols.ProtocolDescriptor,
) io.Closer {
	c := &sthSyncer{
		store:     store,
		logID:     logID,
		closeChan: make(chan interface{}),
	}

	updateSTH := func() error {

		var oldSTH *schema.LogHead
		err := store.View(func(txn storage.Transaction) error {
			var err error
			logBucketTxn := logBucket.Txn(txn)
			oldSTH, err = storage.GetLogHead(logBucketTxn)
			return err
		})
		if err != nil {
			if err != storage.ObjectNotFoundError {
				return err
			} else {
				// New tree, no local state yet
				size := uint64(0)
				oldSTH = &schema.LogHead{
					TreeSize: &size,
				}
			}
		}

		client, err := clients.Get(logID)
		if err != nil {
			return err
		}
		logapi := client.LogAPI

		sth, err := logapi.GetHead(context.Background(), &apipb.LogHeadRequest{
			LogID: &logID,
		})
		if err != nil {
			return err
		}

		newTreeSize := *sth.TreeSize
		oldTreeSize := *oldSTH.TreeSize

		if oldTreeSize > 0 {

			if newTreeSize < oldTreeSize {
				return fmt.Errorf("Refusing to go back in time")
			}

			if newTreeSize == oldTreeSize {

				if !bytes.Equal(sth.LogRoot, oldSTH.LogRoot) {
					return fmt.Errorf("Log root hash mismatch")
				}

				if !bytes.Equal(sth.MapRoot, oldSTH.MapRoot) {
					return fmt.Errorf("Map root hash mismatch")
				}

				if !bytes.Equal(sth.Signature, oldSTH.Signature) {
					return fmt.Errorf("Signature mismatch")
				}

				return nil // Old and new trees are the same
			}
		}

		valid := sthlib.VerifyLogHeadSig(verifier, sth, pd)
		if !valid {
			return fmt.Errorf("STH signature invalid")
		}

		resp, err := logapi.GetLogConsistencyProof(context.Background(), &apipb.GetLogConsistencyProofRequest{
			LogID:      &logID,
			FirstSize:  &oldTreeSize,
			SecondSize: &newTreeSize,
		})
		if err != nil {
			return err
		}

		valid = vlog.ValidConsistencyProof(
			pd.NewHash,
			oldSTH.LogRoot,
			sth.LogRoot,
			oldTreeSize,
			newTreeSize,
			resp.Proof,
		)
		if !valid {
			return fmt.Errorf("Consistency proof invalid")
		}

		err = store.Update(func(txn storage.Transaction) error {
			buf, err := proto.Marshal(sth)
			if err != nil {
				return err
			}

			return logBucket.Txn(txn).Set([]byte(constants.HeadBlob), buf)
		})
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"logID":       logID,
			"oldTreeSize": *oldSTH.TreeSize,
			"treeSize":    *sth.TreeSize,
		}).Info("Updated STH")

		return nil
	}

	go func() {
		run := func() {
			err := updateSTH()
			if err != nil {
				log.WithFields(log.Fields{
					"logID": logID,
					"error": err,
				}).Error("Could not update STH")
			}
		}

		run()

		timeout := time.NewTimer(pollInterval)
		defer timeout.Stop()

		for {
			timeout.Reset(pollInterval)
			select {
			case <-c.closeChan:
				timeout.Stop()
				return
			case <-timeout.C:
				run()
			}
		}
	}()

	return c
}

func (c *sthSyncer) Close() error {
	c.closeChan <- nil
	return nil
}
