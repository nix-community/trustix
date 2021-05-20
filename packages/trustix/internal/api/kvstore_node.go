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
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix/interfaces"
	"github.com/tweag/trustix/packages/trustix/internal/storage"
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
