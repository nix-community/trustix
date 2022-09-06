// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix/interfaces"
	"github.com/nix-community/trustix/packages/trustix/internal/storage"
)

type kvStoreNodeApi struct {
	store    storage.Storage
	caBucket *storage.Bucket
	logMeta  []*api.Log
}

func NewKVStoreNodeAPI(store storage.Storage, caBucket *storage.Bucket, logMeta []*api.Log) interfaces.NodeAPI {
	return &kvStoreNodeApi{
		store:    store,
		caBucket: caBucket,
		logMeta:  logMeta,
	}
}

func (kv *kvStoreNodeApi) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {
	var value []byte
	err := kv.store.View(func(txn storage.Transaction) error {
		bucketTxn := kv.caBucket.Txn(txn)

		var err error

		value, err = bucketTxn.Get(in.Digest)
		if err != nil {
			return err
		}

		digest2 := sha256.Sum256(value)
		if !bytes.Equal(in.Digest, digest2[:]) {
			return fmt.Errorf("Digest mismatch")
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &api.ValueResponse{
		Value: value,
	}, nil
}

func (kv *kvStoreNodeApi) Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error) {
	return &api.LogsResponse{
		Logs: kv.logMeta,
	}, nil
}
