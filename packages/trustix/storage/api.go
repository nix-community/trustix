// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package storage

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
)

type StorageAPI struct {
	txn Transaction
}

func NewStorageAPI(txn Transaction) *StorageAPI {
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
		return nil, ObjectNotFoundError
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
