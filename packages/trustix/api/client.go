// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package api

import (
	"context"

	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	"google.golang.org/grpc"
)

// TrustixAPIGRPCClient is a gRPC based implementation of TrustixLogAPI wrapping a gRPC ClientConn
type TrustixAPIGRPCClient struct {
	client api.TrustixLogAPIClient
}

func NewTrustixAPIGRPCClient(conn *grpc.ClientConn) (*TrustixAPIGRPCClient, error) {
	c := api.NewTrustixLogAPIClient(conn)
	return &TrustixAPIGRPCClient{
		client: c,
	}, nil
}

func (s *TrustixAPIGRPCClient) GetSTH(ctx context.Context, req *api.STHRequest) (*schema.STH, error) {
	return s.client.GetSTH(ctx, req)
}

func (s *TrustixAPIGRPCClient) GetLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error) {
	return s.client.GetLogConsistencyProof(ctx, req)
}

func (s *TrustixAPIGRPCClient) GetLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (*api.ProofResponse, error) {
	return s.client.GetLogAuditProof(ctx, req)
}

func (s *TrustixAPIGRPCClient) GetLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	return s.client.GetLogEntries(ctx, req)
}

func (s *TrustixAPIGRPCClient) GetMapValue(ctx context.Context, req *api.GetMapValueRequest) (*api.MapValueResponse, error) {
	return s.client.GetMapValue(ctx, req)
}

func (s *TrustixAPIGRPCClient) Submit(ctx context.Context, req *api.SubmitRequest) (*api.SubmitResponse, error) {
	return s.client.Submit(ctx, req)
}

func (s *TrustixAPIGRPCClient) Flush(ctx context.Context, in *api.FlushRequest) (*api.FlushResponse, error) {
	return s.client.Flush(ctx, in)
}

func (s *TrustixAPIGRPCClient) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {
	return s.client.GetValue(ctx, in)
}

func (s *TrustixAPIGRPCClient) GetMHLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error) {
	return s.client.GetMHLogConsistencyProof(ctx, req)
}

func (s *TrustixAPIGRPCClient) GetMHLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (*api.ProofResponse, error) {
	return s.client.GetMHLogAuditProof(ctx, req)
}

func (s *TrustixAPIGRPCClient) GetMHLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	return s.client.GetMHLogEntries(ctx, req)
}
