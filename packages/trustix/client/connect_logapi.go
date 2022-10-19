// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package client

import (
	"context"
	"net/http"

	connect "github.com/bufbuild/connect-go"
	api "github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/api/apiconnect"
	schema "github.com/nix-community/trustix/packages/trustix-proto/schema"
)

type logAPIConnectClient struct {
	client apiconnect.LogAPIClient
}

func newLogAPIConnectClient(client *http.Client, baseURL string, options ...connect.ClientOption) *logAPIConnectClient {
	c := apiconnect.NewLogAPIClient(client, baseURL, options...)
	return &logAPIConnectClient{
		client: c,
	}
}

func (c *logAPIConnectClient) GetHead(ctx context.Context, req *api.LogHeadRequest) (*schema.LogHead, error) {
	resp, err := c.client.GetHead(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *logAPIConnectClient) GetLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error) {
	resp, err := c.client.GetLogConsistencyProof(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *logAPIConnectClient) GetLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (*api.ProofResponse, error) {
	resp, err := c.client.GetLogAuditProof(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *logAPIConnectClient) GetLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	resp, err := c.client.GetLogEntries(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *logAPIConnectClient) GetMapValue(ctx context.Context, req *api.GetMapValueRequest) (*api.MapValueResponse, error) {
	resp, err := c.client.GetMapValue(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *logAPIConnectClient) GetMHLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error) {
	resp, err := c.client.GetMHLogConsistencyProof(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *logAPIConnectClient) GetMHLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (*api.ProofResponse, error) {
	resp, err := c.client.GetMHLogAuditProof(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *logAPIConnectClient) GetMHLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	resp, err := c.client.GetMHLogEntries(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}
