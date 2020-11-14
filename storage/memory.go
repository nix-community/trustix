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

func (t *memoryTxn) Get(bucket []byte, key []byte) ([]byte, error) {
	val, err := t.txn.First("record", "id", base64.StdEncoding.EncodeToString(createCompoundNativeKey(bucket, key)))
	if err != nil {
		return nil, err
	}

	if val == nil {
		return nil, ObjectNotFoundError
	}

	record, ok := val.(*memdbRecord)
	if !ok {
		return nil, fmt.Errorf("Could not convert value to record")
	}

	return base64.StdEncoding.DecodeString(record.Value)
}

func (t *memoryTxn) Set(bucket []byte, key []byte, value []byte) error {
	return t.txn.Insert("record", &memdbRecord{
		Key:   base64.StdEncoding.EncodeToString(createCompoundNativeKey(bucket, key)),
		Value: base64.StdEncoding.EncodeToString(value),
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
