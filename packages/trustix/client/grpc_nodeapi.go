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

	api "github.com/nix-community/trustix/packages/trustix-proto/api"
	// schema "github.com/nix-community/trustix/packages/trustix-proto/schema"
	"google.golang.org/grpc"
)

type nodeAPIGRPCClient struct {
	client api.NodeAPIClient
}

func newNodeAPIGRPCClient(conn *grpc.ClientConn) *nodeAPIGRPCClient {
	c := api.NewNodeAPIClient(conn)
	return &nodeAPIGRPCClient{
		client: c,
	}
}

func (c *nodeAPIGRPCClient) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {
	return c.client.GetValue(ctx, in)
}

func (c *nodeAPIGRPCClient) Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error) {
	return c.client.Logs(ctx, in)
}
