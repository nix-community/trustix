// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package client

import (
	"github.com/nix-community/trustix/packages/trustix/interfaces"
)

const (
	LocalClientType = iota
	GRPCClientType
)

type Client struct {
	NodeAPI interfaces.NodeAPI
	LogAPI  interfaces.LogAPI
	RpcAPI  interfaces.RpcAPI
	LogRPC  interfaces.LogRPC

	Type int

	close func() error
}

func (c *Client) Close() error {
	if c.close != nil {
		return c.Close()
	}

	return nil
}
