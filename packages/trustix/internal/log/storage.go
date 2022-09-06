// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package log

import (
	"fmt"

	"github.com/nix-community/trustix/packages/trustix-proto/schema"
	"github.com/nix-community/trustix/packages/trustix/internal/storage"
	"google.golang.org/protobuf/proto"
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
