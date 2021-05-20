// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package storage

import (
	"path"
	"sync"

	badger "github.com/dgraph-io/badger/v2"
)

type nativeStorage struct {
	db       *badger.DB
	txnMutex *sync.RWMutex
}

type nativeTxn struct {
	txn *badger.Txn
}

func createCompoundNativeKey(bucket *Bucket, key []byte) []byte {
	cKey := []byte(bucket.Join())
	cKey = append(cKey, 0x2f)
	cKey = append(cKey, key...)
	return cKey
}

func (t *nativeTxn) Get(bucket *Bucket, key []byte) ([]byte, error) {
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

func (t *nativeTxn) Set(bucket *Bucket, key []byte, value []byte) error {
	return t.txn.Set(createCompoundNativeKey(bucket, key), value)
}

func (t *nativeTxn) Delete(bucket *Bucket, key []byte) error {
	return t.txn.Delete(createCompoundNativeKey(bucket, key))
}

func NewNativeStorage(stateDirectory string) (*nativeStorage, error) {
	path := path.Join(stateDirectory, "trustix.db")

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
