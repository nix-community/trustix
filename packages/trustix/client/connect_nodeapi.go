// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package client

import (
	"context"
	"net/http"

	connect "github.com/bufbuild/connect-go"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/api/apiconnect"
)

type nodeAPIConnectClient struct {
	client apiconnect.NodeAPIClient
}

func newNodeAPIConnectClient(client *http.Client, baseURL string) *nodeAPIConnectClient {
	c := apiconnect.NewNodeAPIClient(client, baseURL)
	return &nodeAPIConnectClient{
		client: c,
	}
}

func (c *nodeAPIConnectClient) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {
	resp, err := c.client.GetValue(ctx, connect.NewRequest(in))
	if err != nil {
		return nil, err
	}

	return resp.Msg, err
}

func (c *nodeAPIConnectClient) Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error) {
	resp, err := c.client.Logs(ctx, connect.NewRequest(in))
	if err != nil {
		return nil, err
	}

	return resp.Msg, err
}
