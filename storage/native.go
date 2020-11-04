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
	badger "github.com/dgraph-io/badger/v2"
	"path"
	"sync"
)

type nativeStorage struct {
	db       *badger.DB
	txnMutex *sync.RWMutex
}

type nativeTxn struct {
	txn *badger.Txn
}

func createCompoundNativeKey(bucket []byte, key []byte) []byte {
	cKey := []byte{}
	cKey = append(cKey, bucket...)
	cKey = append(cKey, 0x2f)
	cKey = append(cKey, key...)
	return cKey
}

func (t *nativeTxn) Get(bucket []byte, key []byte) ([]byte, error) {
	val, err := t.txn.Get(createCompoundNativeKey(bucket, key))
	if err != nil {
		// Normalise error
		if err == badger.ErrKeyNotFound {
			return nil, ObjectNotFoundError
		}
		return nil, err
	}

	return val.ValueCopy(nil)
}

func (t *nativeTxn) Set(bucket []byte, key []byte, value []byte) error {
	return t.txn.Set(createCompoundNativeKey(bucket, key), value)
}

func NewNativeStorage(name string, stateDirectory string) (*nativeStorage, error) {
	path := path.Join(stateDirectory, name+".db")

	options := badger.DefaultOptions(path)
	options.Logger = nil

	db, err := badger.Open(options)
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

	txn := s.db.NewTransaction(readWrite)
	if readWrite {
		defer txn.Discard()
	}

	t := &nativeTxn{
		txn: txn,
	}

	err := fn(t)
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
