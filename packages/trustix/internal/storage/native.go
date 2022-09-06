// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package storage

import (
	"path"
	"sync"

	bolt "go.etcd.io/bbolt"
)

type nativeStorage struct {
	db       *bolt.DB
	txnMutex *sync.RWMutex
}

type nativeTxn struct {
	txn *bolt.Tx
}

func (t *nativeTxn) Get(bucket *Bucket, key []byte) ([]byte, error) {
	// TODO: Better bucket scheme
	b := t.txn.Bucket([]byte(bucket.Join()))
	if b == nil {
		return nil, ObjectNotFoundError
	}

	val := b.Get(key)
	if val == nil {
		return nil, ObjectNotFoundError
	}

	return val, nil
}

func (t *nativeTxn) Set(bucket *Bucket, key []byte, value []byte) error {
	// TODO: Better bucket scheme
	b, err := t.txn.CreateBucketIfNotExists([]byte(bucket.Join()))
	if err != nil {
		return err
	}

	err = b.Put(key, value)
	if err != nil {
		return err
	}

	return nil
}

func (t *nativeTxn) Delete(bucket *Bucket, key []byte) error {
	// TODO: Better bucket scheme
	b, err := t.txn.CreateBucketIfNotExists([]byte(bucket.Join()))
	if err != nil {
		return err
	}

	return b.Delete(key)
}

func NewNativeStorage(stateDirectory string) (*nativeStorage, error) {
	path := path.Join(stateDirectory, "trustix.db")

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &nativeStorage{
		db:       db,
		txnMutex: &sync.RWMutex{},
	}, nil

}

func (s *nativeStorage) runTX(readWrite bool, fn func(Transaction) error) error {
	if readWrite {
		s.txnMutex.Lock()
		defer s.txnMutex.Unlock()
	} else {
		s.txnMutex.RLock()
		defer s.txnMutex.RUnlock()
	}

	txn, err := s.db.Begin(readWrite)
	if err != nil {
		return err
	}
	defer func() {
		err := txn.Rollback()
		if err != nil && err != bolt.ErrTxClosed {
			panic(err)
		}
	}()

	t := &nativeTxn{
		txn: txn,
	}

	err = fn(t)
	if err != nil {
		return err
	} else {
		if readWrite {
			return txn.Commit()
		}
	}

	return err
}

func (s *nativeStorage) View(fn func(Transaction) error) error {
	return s.runTX(false, fn)
}

func (s *nativeStorage) Update(fn func(Transaction) error) error {
	return s.runTX(true, fn)
}

func (s *nativeStorage) Close() {
	s.db.Close()
}
