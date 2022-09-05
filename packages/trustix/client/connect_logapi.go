// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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

func newLogAPIConnectClient(client *http.Client, baseURL string) *logAPIConnectClient {
	c := apiconnect.NewLogAPIClient(client, baseURL)
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
