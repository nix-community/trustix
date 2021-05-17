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

	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/rpc"
	schema "github.com/tweag/trustix/packages/trustix-proto/schema"
	"google.golang.org/grpc"
)

type logRPCGRPCClient struct {
	client rpc.LogRPCClient
}

func newLogRPCGRPCClient(conn *grpc.ClientConn) *logRPCGRPCClient {
	c := rpc.NewLogRPCClient(conn)
	return &logRPCGRPCClient{
		client: c,
	}
}

func (c *logRPCGRPCClient) GetHead(ctx context.Context, req *api.LogHeadRequest) (*schema.LogHead, error) {
	return c.client.GetHead(ctx, req)
}

func (c *logRPCGRPCClient) GetLogEntries(ctx context.Context, in *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	return c.client.GetLogEntries(ctx, in)
}

func (c *logRPCGRPCClient) Submit(ctx context.Context, in *rpc.SubmitRequest) (*rpc.SubmitResponse, error) {
	return c.client.Submit(ctx, in)
}

func (c *logRPCGRPCClient) Flush(ctx context.Context, in *rpc.FlushRequest) (*rpc.FlushResponse, error) {
	return c.client.Flush(ctx, in)
}
