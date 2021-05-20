// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package log

import (
	"fmt"

	proto "github.com/golang/protobuf/proto"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	"github.com/tweag/trustix/packages/trustix/internal/storage"
)

type LogStorage struct {
	txn *storage.BucketTransaction
}

func NewLogStorage(txn *storage.BucketTransaction) *LogStorage {
	return &LogStorage{
		txn: txn,
	}
}

func (s *LogStorage) Get(level int, idx uint64) (*schema.LogLeaf, error) {
	key := []byte(fmt.Sprintf("%d-%d", level, idx))

	v, err := s.txn.Get(key)
	if err != nil {
		return nil, err
	}

	l, err := LeafFromBytes(v)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (s *LogStorage) Append(treeSize uint64, level int, leaf *schema.LogLeaf) error {
	v, err := proto.Marshal(leaf)
	if err != nil {
		return err
	}

	idx := levelSize(treeSize, level) - 1

	key := []byte(fmt.Sprintf("%d-%d", level, idx))

	err = s.txn.Set(key, v)
	if err != nil {
		return err
	}

	return nil

}
