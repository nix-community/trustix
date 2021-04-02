// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package api

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	proto "github.com/golang/protobuf/proto"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	"github.com/tweag/trustix/packages/trustix/storage"
)

type StorageAPI struct {
	txn storage.Transaction
}

func NewStorageAPI(txn storage.Transaction) *StorageAPI {
	return &StorageAPI{
		txn: txn,
	}
}

func (s *StorageAPI) SetSTH(logID string, sth *schema.STH) error {
	buf, err := proto.Marshal(sth)
	if err != nil {
		return err
	}

	return s.txn.Set([]byte(logID), []byte("HEAD"), buf)
}

func (s *StorageAPI) GetSTH(logID string) (*schema.STH, error) {
	var buf []byte
	{
		v, err := s.txn.Get([]byte(logID), []byte("HEAD"))
		if err != nil {
			return nil, err
		}
		buf = v
	}
	if len(buf) == 0 {
		return nil, storage.ObjectNotFoundError
	}

	sth := &schema.STH{}
	err := proto.Unmarshal(buf, sth)
	if err != nil {
		return nil, err
	}

	return sth, nil
}

func (s *StorageAPI) SetSMTValue(logID string, key []byte, value []byte) error {
	return s.txn.Set([]byte("SMT"), key, value)
}

func (s *StorageAPI) GetSMTValue(logID string, key []byte) ([]byte, error) {
	return s.txn.Get([]byte("SMT"), key)
}

func (s *StorageAPI) GetQueueMeta(logID string) (*schema.SubmitQueue, error) {
	q := &schema.SubmitQueue{}

	qBytes, err := s.txn.Get([]byte("QUEUE"), []byte("META"))
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

func (s *StorageAPI) SetQueueMeta(logID string, q *schema.SubmitQueue) error {
	qBytes, err := proto.Marshal(q)
	if err != nil {
		return err
	}

	err = s.txn.Set([]byte("QUEUE"), []byte("META"), qBytes)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageAPI) WriteQueueItem(logID string, itemID int, item *api.KeyValuePair) error {
	itemBytes, err := proto.Marshal(item)
	if err != nil {
		return err
	}

	err = s.txn.Set([]byte("QUEUE"), []byte(fmt.Sprintf("%d", itemID)), itemBytes)
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageAPI) PopQueueItem(logID string, itemID int) (*api.KeyValuePair, error) {
	itemBytes, err := s.txn.Get([]byte("QUEUE"), []byte(fmt.Sprintf("%d", itemID)))
	if err != nil {
		return nil, err
	}

	item := &api.KeyValuePair{}
	err = proto.Unmarshal(itemBytes, item)
	if err != nil {
		return nil, err
	}

	err = s.txn.Delete([]byte("QUEUE"), []byte(fmt.Sprintf("%d", itemID)))
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *StorageAPI) SetCAValue(value []byte) error {
	digest := sha256.Sum256(value)
	return s.txn.Set([]byte("VALUES"), digest[:], value)
}

func (s *StorageAPI) GetCAValue(digest []byte) ([]byte, error) {
	value, err := s.txn.Get([]byte("VALUES"), digest)
	if err != nil {
		return nil, err
	}

	digest2 := sha256.Sum256(value)
	if !bytes.Equal(digest, digest2[:]) {
		return nil, fmt.Errorf("Digest mismatch")
	}

	return value, nil
}
