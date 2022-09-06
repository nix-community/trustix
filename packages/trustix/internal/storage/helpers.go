// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package storage

import (
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
	"github.com/nix-community/trustix/packages/trustix/internal/constants"
	"google.golang.org/protobuf/proto"
)

// TODO: I don't like this living here but I also don't have a better option

func GetLogHead(txn *BucketTransaction) (*schema.LogHead, error) {
	var buf []byte
	{
		v, err := txn.Get([]byte(constants.HeadBlob))
		if err != nil {
			return nil, err
		}
		buf = v
	}
	if len(buf) == 0 {
		return nil, objectNotFoundError([]byte(constants.HeadBlob))
	}

	sth := &schema.LogHead{}
	err := proto.Unmarshal(buf, sth)
	if err != nil {
		return nil, err
	}

	return sth, nil
}
