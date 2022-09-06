// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package storage

type Transaction interface {
	Get(bucket *Bucket, key []byte) ([]byte, error)
	Set(bucket *Bucket, key []byte, value []byte) error
	Delete(bucket *Bucket, key []byte) error
}

type Storage interface {
	Close()

	// View - Start a read-only transaction
	View(func(txn Transaction) error) error

	// Update - Start a read-write transaction
	Update(func(txn Transaction) error) error
}
