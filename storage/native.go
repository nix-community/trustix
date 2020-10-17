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
	"github.com/tweag/trustix/config"
	bolt "go.etcd.io/bbolt"
	"path"
)

type nativeStorage struct {
	db *bolt.DB
}

type nativeTxn struct {
	txn *bolt.Tx
}

func (t *nativeTxn) Get(bucket []byte, key []byte) ([]byte, error) {
	b := t.txn.Bucket(bucket)
	if b == nil {
		return nil, ObjectNotFoundError
	}

	val := b.Get(key)
	if val == nil {
		return nil, ObjectNotFoundError
	}

	return val, nil
}

func (t *nativeTxn) Set(bucket []byte, key []byte, value []byte) error {
	b, err := t.txn.CreateBucketIfNotExists(bucket)
	if err != nil {
		return err
	}

	err = b.Put(key, value)
	if err != nil {
		return err
	}

	return nil
}

func NativeStorageFromConfig(name string, stateDirectory string, conf *config.NativeStorageConfig) (*nativeStorage, error) {
	path := path.Join(stateDirectory, name+".db")

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &nativeStorage{
		db: db,
	}, nil
}

func (s *nativeStorage) runTX(readWrite bool, fn func(Transaction) error) error {
	txn, err := s.db.Begin(readWrite)
	if err != nil {
		return err
	}
	defer txn.Rollback()

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
