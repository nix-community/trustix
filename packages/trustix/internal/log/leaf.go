// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package log

import (
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
	"google.golang.org/protobuf/proto"
)

func NewLeaf(digest []byte, value []byte) (*schema.LogLeaf, error) {
	return &schema.LogLeaf{
		LeafDigest: digest,
	}, nil
}

func LeafFromBytes(data []byte) (*schema.LogLeaf, error) {
	l := &schema.LogLeaf{}
	err := proto.Unmarshal(data, l)
	if err != nil {
		return nil, err
	}
	return l, nil
}
