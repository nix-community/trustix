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
	"github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix-proto/rpc/rpcconnect"
	schema "github.com/nix-community/trustix/packages/trustix-proto/schema"
)

type logRPCConnectClient struct {
	client rpcconnect.LogRPCClient
}

func newLogRPCConnectClient(client *http.Client, baseURL string) *logRPCConnectClient {
	c := rpcconnect.NewLogRPCClient(client, baseURL)
	return &logRPCConnectClient{
		client: c,
	}
}

func (c *logRPCConnectClient) GetHead(ctx context.Context, req *api.LogHeadRequest) (*schema.LogHead, error) {
	resp, err := c.client.GetHead(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *logRPCConnectClient) GetLogEntries(ctx context.Context, in *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	resp, err := c.client.GetLogEntries(ctx, connect.NewRequest(in))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *logRPCConnectClient) Submit(ctx context.Context, in *rpc.SubmitRequest) (*rpc.SubmitResponse, error) {
	resp, err := c.client.Submit(ctx, connect.NewRequest(in))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}

func (c *logRPCConnectClient) Flush(ctx context.Context, in *rpc.FlushRequest) (*rpc.FlushResponse, error) {
	resp, err := c.client.Flush(ctx, connect.NewRequest(in))
	if err != nil {
		return nil, err
	}

	return resp.Msg, nil
}
