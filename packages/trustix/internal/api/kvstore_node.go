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

	"github.com/nix-community/trustix/packages/go-lib/set"
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
	// Default to returning all logs
	logs := kv.logMeta

	// If any protocols are passed filter out other protocols
	if in.Protocols != nil && len(in.Protocols) > 0 {
		logs = []*api.Log{}

		protocolSet := set.NewSet[string]()
		for _, p := range in.Protocols {
			protocolSet.Add(p)
		}

		for _, log := range kv.logMeta {
			if protocolSet.Has(*log.Protocol) {
				logs = append(logs, log)
			}
		}
	}

	return &api.LogsResponse{
		Logs: logs,
	}, nil
}
