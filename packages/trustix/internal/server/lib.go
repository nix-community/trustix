// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package server

import (
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
	"github.com/nix-community/trustix/packages/trustix/internal/storage"
)

func getLogHead(rootBucket *storage.Bucket, txn storage.Transaction, logID string) (*schema.LogHead, error) {
	bucket := rootBucket.Cd(logID)
	return storage.GetLogHead(bucket.Txn(txn))
}
