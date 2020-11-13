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

package sthmanager

import (
	"bytes"
	"fmt"
	proto "github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/api"
	vlog "github.com/tweag/trustix/log"
	"github.com/tweag/trustix/schema"
	"github.com/tweag/trustix/signer"
	sthlib "github.com/tweag/trustix/sth"
	"github.com/tweag/trustix/storage"
	"time"
)

type STHCache interface {
	Get() (*schema.STH, error)
	Close()
}

type sthCache struct {
	store     storage.TrustixStorage
	logName   string
	closeChan chan interface{}
}

func NewSTHCache(logName string, store storage.TrustixStorage, logapi api.TrustixLogAPI, verifier signer.TrustixVerifier) STHCache {
	c := &sthCache{
		store:     store,
		logName:   logName,
		closeChan: make(chan interface{}),
	}

	updateSTH := func() error {
		oldSTH, err := c.Get()
		if err != nil {
			return err
		}

		sth, err := logapi.GetSTH(new(api.STHRequest))
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

		resp, err := logapi.GetLogConsistencyProof(&api.GetLogConsistencyProofRequest{
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

		err = c.Set(sth)

		log.WithFields(log.Fields{
			"logName":     logName,
			"oldTreeSize": oldSTH.TreeSize,
			"treeSize":    sth.TreeSize,
		}).Info("Updated STH")

		return err
	}

	go func() {
		run := func() {
			err := updateSTH()
			if err != nil {
				log.WithFields(log.Fields{
					"logName": logName,
					"error":   err,
				}).Error("Could not update STH")
			}
		}

		run()

		// TODO: Make timeout configurable (& manually triggerable)
		for {
			select {
			case _ = <-c.closeChan:
				return
			case <-time.After(10 * time.Second):
				run()
			}
		}
	}()

	return c
}

func (c *sthCache) Set(sth *schema.STH) error {
	buf, err := proto.Marshal(sth)
	if err != nil {
		return err
	}

	return c.store.Update(func(txn storage.Transaction) error {
		return txn.Set([]byte(c.logName), []byte("HEAD"), buf)
	})
}

func (c *sthCache) Get() (*schema.STH, error) {
	var buf []byte
	err := c.store.View(func(txn storage.Transaction) error {
		v, err := txn.Get([]byte(c.logName), []byte("HEAD"))
		buf = v
		return err
	})

	sth := &schema.STH{}
	err = proto.Unmarshal(buf, sth)
	if err != nil {
		return nil, err
	}

	return sth, nil
}

func (c *sthCache) Close() {
	c.closeChan <- nil
}
