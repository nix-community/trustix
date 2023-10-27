// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package client

import (
	"context"
	"net/http"

	connect "connectrpc.com/connect"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix-proto/rpc/rpcconnect"
)

type rpcAPIConnectClient struct {
	client rpcconnect.RPCApiClient
}

func newRpcAPIConnectClient(client *http.Client, baseURL string, options ...connect.ClientOption) *rpcAPIConnectClient {
	c := rpcconnect.NewRPCApiClient(client, baseURL, options...)
	return &rpcAPIConnectClient{
		client: c,
	}
}

func (c *rpcAPIConnectClient) Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error) {
	resp, err := c.client.Logs(ctx, connect.NewRequest(in))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *rpcAPIConnectClient) Decide(ctx context.Context, in *rpc.DecideRequest) (*rpc.DecisionResponse, error) {
	resp, err := c.client.Decide(ctx, connect.NewRequest(in))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *rpcAPIConnectClient) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {
	resp, err := c.client.GetValue(ctx, connect.NewRequest(in))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}
