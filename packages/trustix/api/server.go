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
	"github.com/tweag/trustix/packages/trustix/rpc/auth"
)

// TrustixAPIServer wraps kvStoreLogApi and turns it into a gRPC implementation
type TrustixAPIServer struct {
	api.UnimplementedTrustixLogAPIServer
	impl TrustixLogAPI
}

func NewTrustixAPIServer(impl TrustixLogAPI) (*TrustixAPIServer, error) {
	return &TrustixAPIServer{
		impl: impl,
	}, nil
}

func (s *TrustixAPIServer) GetSTH(ctx context.Context, req *api.STHRequest) (*schema.STH, error) {
	return s.impl.GetSTH(ctx, req)
}

func (s *TrustixAPIServer) GetLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error) {
	return s.impl.GetLogConsistencyProof(ctx, req)
}

func (s *TrustixAPIServer) GetLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (*api.ProofResponse, error) {
	return s.impl.GetLogAuditProof(ctx, req)
}

func (s *TrustixAPIServer) GetLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	return s.impl.GetLogEntries(ctx, req)
}

func (s *TrustixAPIServer) GetMapValue(ctx context.Context, req *api.GetMapValueRequest) (*api.MapValueResponse, error) {
	return s.impl.GetMapValue(ctx, req)
}

func (s *TrustixAPIServer) Submit(ctx context.Context, req *api.SubmitRequest) (*api.SubmitResponse, error) {

	err := auth.CanWrite(ctx)
	if err != nil {
		return nil, err
	}

	return s.impl.Submit(ctx, req)
}

func (s *TrustixAPIServer) Flush(ctx context.Context, req *api.FlushRequest) (*api.FlushResponse, error) {
	err := auth.CanWrite(ctx)
	if err != nil {
		return nil, err
	}

	return s.impl.Flush(ctx, req)
}

func (s *TrustixAPIServer) GetValue(ctx context.Context, req *api.ValueRequest) (*api.ValueResponse, error) {
	return s.impl.GetValue(ctx, req)
}

func (s *TrustixAPIServer) GetMHLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error) {
	return s.impl.GetMHLogConsistencyProof(ctx, req)
}

func (s *TrustixAPIServer) GetMHLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (*api.ProofResponse, error) {
	return s.impl.GetMHLogAuditProof(ctx, req)
}

func (s *TrustixAPIServer) GetMHLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	return s.impl.GetMHLogEntries(ctx, req)
}
