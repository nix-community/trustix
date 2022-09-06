// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package storage

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/hashicorp/go-memdb"
)

type memdbRecord struct {
	Key   string
	Value string
}

type memoryStorage struct {
	txnMutex *sync.RWMutex
	db       *memdb.MemDB
}

type memoryTxn struct {
	txn *memdb.Txn
}

func (t *memoryTxn) createKey(bucket *Bucket, key []byte) []byte {
	cKey := []byte(bucket.Join())
	cKey = append(cKey, 0x2f)
	cKey = append(cKey, key...)
	return cKey
}

func (t *memoryTxn) Get(bucket *Bucket, key []byte) ([]byte, error) {
	val, err := t.txn.First("record", "id", base64.StdEncoding.EncodeToString(t.createKey(bucket, key)))
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, objectNotFoundError(key)
	}

	record, ok := val.(*memdbRecord)
	if !ok {
		return nil, fmt.Errorf("Could not convert value to record")
	}

	return base64.StdEncoding.DecodeString(record.Value)
}

func (t *memoryTxn) Set(bucket *Bucket, key []byte, value []byte) error {
	return t.txn.Insert("record", &memdbRecord{
		Key:   base64.StdEncoding.EncodeToString(t.createKey(bucket, key)),
		Value: base64.StdEncoding.EncodeToString(value),
	})
}

func (t *memoryTxn) Delete(bucket *Bucket, key []byte) error {
	return t.txn.Delete("record", &memdbRecord{
		Key: base64.StdEncoding.EncodeToString(t.createKey(bucket, key)),
	})
}

func NewMemoryStorage() (*memoryStorage, error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"record": &memdb.TableSchema{
				Name: "record",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Key"},
					},
					"value": &memdb.IndexSchema{
						Name:    "value",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Value"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	return &memoryStorage{
		txnMutex: &sync.RWMutex{},
		db:       db,
	}, nil
}

func (s *memoryStorage) runTX(readWrite bool, fn func(Transaction) error) error {
	if readWrite {
		s.txnMutex.Lock()
		defer s.txnMutex.Unlock()
	} else {
		s.txnMutex.RLock()
		defer s.txnMutex.RUnlock()
	}

	txn := s.db.Txn(readWrite)
	if readWrite {
		defer txn.Abort()
	}

	t := &memoryTxn{
		txn: txn,
	}

	err := fn(t)
	if err != nil {
		return err
	} else {
		if readWrite {
			txn.Commit()
			return nil
		}
	}

	return err
}

func (s *memoryStorage) View(fn func(Transaction) error) error {
	return s.runTX(false, fn)
}

func (s *memoryStorage) Update(fn func(Transaction) error) error {
	return s.runTX(true, fn)
}

func (s *memoryStorage) Close() {

}
