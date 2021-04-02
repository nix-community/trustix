// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package sthsync

import (
	"bytes"
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	apipb "github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	"github.com/tweag/trustix/packages/trustix/api"
	vlog "github.com/tweag/trustix/packages/trustix/log"
	"github.com/tweag/trustix/packages/trustix/signer"
	sthlib "github.com/tweag/trustix/packages/trustix/sth"
	"github.com/tweag/trustix/packages/trustix/storage"
)

type sthSyncer struct {
	store     storage.TrustixStorage
	logID     string
	closeChan chan interface{}
}

func NewSTHSyncer(logID string, store storage.TrustixStorage, logapi api.TrustixLogAPI, verifier signer.TrustixVerifier) STHSyncer {
	c := &sthSyncer{
		store:     store,
		logID:     logID,
		closeChan: make(chan interface{}),
	}

	updateSTH := func() error {

		var oldSTH *schema.STH
		err := store.View(func(txn storage.Transaction) error {
			storageAPI := storage.NewStorageAPI(txn)
			var err error
			oldSTH, err = storageAPI.GetSTH(logID)
			return err
		})

		if err != nil {
			if err != storage.ObjectNotFoundError {
				return err
			} else {
				// New tree, no local state yet
				size := uint64(0)
				oldSTH = &schema.STH{
					TreeSize: &size,
				}
			}
		}

		sth, err := logapi.GetSTH(context.Background(), new(apipb.STHRequest))
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

		valid := sthlib.VerifySTHSig(verifier, sth)
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
			storageAPI := storage.NewStorageAPI(txn)
			return storageAPI.SetSTH(logID, sth)
		})

		log.WithFields(log.Fields{
			"logID":       logID,
			"oldTreeSize": *oldSTH.TreeSize,
			"treeSize":    *sth.TreeSize,
		}).Info("Updated STH")

		return err
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

		duration := time.Second * 10
		timeout := time.NewTimer(duration)
		defer timeout.Stop()

		// TODO: Make timeout configurable (& manually triggerable)
		for {
			timeout.Reset(duration)
			select {
			case _ = <-c.closeChan:
				return
			case <-timeout.C:
				run()
			}
		}
	}()

	return c
}

func (c *sthSyncer) Close() {
	c.closeChan <- nil
}
