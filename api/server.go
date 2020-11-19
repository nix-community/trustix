// MIT License
//
// Copyright (c) 2020 Tweag IO
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package api

import (
	"context"

	"github.com/tweag/trustix/rpc/auth"
	"github.com/tweag/trustix/schema"
)

// TrustixAPIServer wraps kvStoreLogApi and turns it into a gRPC implementation
type TrustixAPIServer struct {
	UnimplementedTrustixLogAPIServer
	impl TrustixLogAPI
}

func NewTrustixAPIServer(impl TrustixLogAPI) (*TrustixAPIServer, error) {
	return &TrustixAPIServer{
		impl: impl,
	}, nil
}

func (s *TrustixAPIServer) GetSTH(ctx context.Context, req *STHRequest) (*schema.STH, error) {
	return s.impl.GetSTH(ctx, req)
}

func (s *TrustixAPIServer) GetLogConsistencyProof(ctx context.Context, req *GetLogConsistencyProofRequest) (*ProofResponse, error) {
	return s.impl.GetLogConsistencyProof(ctx, req)
}

func (s *TrustixAPIServer) GetLogAuditProof(ctx context.Context, req *GetLogAuditProofRequest) (*ProofResponse, error) {
	return s.impl.GetLogAuditProof(ctx, req)
}

func (s *TrustixAPIServer) GetLogEntries(ctx context.Context, req *GetLogEntriesRequest) (*LogEntriesResponse, error) {
	return s.impl.GetLogEntries(ctx, req)
}

func (s *TrustixAPIServer) GetMapValue(ctx context.Context, req *GetMapValueRequest) (*MapValueResponse, error) {
	return s.impl.GetMapValue(ctx, req)
}

func (s *TrustixAPIServer) Submit(ctx context.Context, req *SubmitRequest) (*SubmitResponse, error) {

	err := auth.CanWrite(ctx)
	if err != nil {
		return nil, err
	}

	return s.impl.Submit(ctx, req)
}

func (s *TrustixAPIServer) Flush(ctx context.Context, req *FlushRequest) (*FlushResponse, error) {
	err := auth.CanWrite(ctx)
	if err != nil {
		return nil, err
	}

	return s.impl.Flush(ctx, req)
}
