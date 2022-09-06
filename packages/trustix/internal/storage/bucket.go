// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package storage

import (
	"strings"
)

// Bucket - Emulate a directory tree like storage interface
type Bucket struct {
	path []string
}

func (b *Bucket) Cd(path ...string) *Bucket {
	return &Bucket{
		path: append(b.path, path...),
	}
}

func (b *Bucket) Join() string {
	return strings.Join(b.path, "/")
}

func (b *Bucket) Txn(txn Transaction) *BucketTransaction {
	return NewBucketTransaction(b, txn)
}

// BucketTransaction - Compound type for txn with preapplied bucket
type BucketTransaction struct {
	bucket *Bucket
	txn    Transaction
}

func NewBucketTransaction(bucket *Bucket, txn Transaction) *BucketTransaction {
	return &BucketTransaction{
		bucket: bucket,
		txn:    txn,
	}
}

func (bt *BucketTransaction) Get(key []byte) ([]byte, error) {
	return bt.txn.Get(bt.bucket, key)
}

func (bt *BucketTransaction) Set(key []byte, value []byte) error {
	return bt.txn.Set(bt.bucket, key, value)
}

func (bt *BucketTransaction) Delete(key []byte) error {
	return bt.txn.Delete(bt.bucket, key)
}
