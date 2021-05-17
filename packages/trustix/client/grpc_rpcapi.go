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
	"google.golang.org/grpc"
)

type rpcAPIGRPCClient struct {
	client rpc.RPCApiClient
}

func newRpcAPIGRPCClient(conn *grpc.ClientConn) *rpcAPIGRPCClient {
	c := rpc.NewRPCApiClient(conn)
	return &rpcAPIGRPCClient{
		client: c,
	}
}

func (c *rpcAPIGRPCClient) Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error) {
	return c.client.Logs(ctx, in)
}

func (c *rpcAPIGRPCClient) Decide(ctx context.Context, in *rpc.KeyRequest) (*rpc.DecisionResponse, error) {
	return c.client.Decide(ctx, in)
}

func (c *rpcAPIGRPCClient) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {
	return c.client.GetValue(ctx, in)
}
