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

	api "github.com/tweag/trustix/packages/trustix-proto/api"
	schema "github.com/tweag/trustix/packages/trustix-proto/schema"
	"google.golang.org/grpc"
)

type logAPIGRPCClient struct {
	client api.LogAPIClient
}

func newLogAPIGRPCClient(conn *grpc.ClientConn) *logAPIGRPCClient {
	c := api.NewLogAPIClient(conn)
	return &logAPIGRPCClient{
		client: c,
	}
}

func (c *logAPIGRPCClient) GetHead(ctx context.Context, req *api.LogHeadRequest) (*schema.LogHead, error) {
	return c.client.GetHead(ctx, req)
}

func (c *logAPIGRPCClient) GetLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error) {
	return c.client.GetLogConsistencyProof(ctx, req)
}

func (c *logAPIGRPCClient) GetLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (*api.ProofResponse, error) {
	return c.client.GetLogAuditProof(ctx, req)
}

func (c *logAPIGRPCClient) GetLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	return c.client.GetLogEntries(ctx, req)
}

func (c *logAPIGRPCClient) GetMapValue(ctx context.Context, req *api.GetMapValueRequest) (*api.MapValueResponse, error) {
	return c.client.GetMapValue(ctx, req)
}

func (c *logAPIGRPCClient) GetMHLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error) {
	return c.client.GetMHLogConsistencyProof(ctx, req)
}

func (c *logAPIGRPCClient) GetMHLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (*api.ProofResponse, error) {
	return c.client.GetMHLogAuditProof(ctx, req)
}

func (c *logAPIGRPCClient) GetMHLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	return c.client.GetMHLogEntries(ctx, req)
}
