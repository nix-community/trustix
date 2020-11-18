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

	"github.com/tweag/trustix/schema"
	"google.golang.org/grpc"
)

// TrustixAPIGRPCClient is a gRPC based implementation of TrustixLogAPI wrapping a gRPC ClientConn
type TrustixAPIGRPCClient struct {
	client TrustixLogAPIClient
}

func NewTrustixAPIGRPCClient(conn *grpc.ClientConn) (*TrustixAPIGRPCClient, error) {
	c := NewTrustixLogAPIClient(conn)
	return &TrustixAPIGRPCClient{
		client: c,
	}, nil
}

func (s *TrustixAPIGRPCClient) GetSTH(ctx context.Context, req *STHRequest) (*schema.STH, error) {
	return s.client.GetSTH(ctx, req)
}

func (s *TrustixAPIGRPCClient) GetLogConsistencyProof(ctx context.Context, req *GetLogConsistencyProofRequest) (*ProofResponse, error) {
	return s.client.GetLogConsistencyProof(ctx, req)
}

func (s *TrustixAPIGRPCClient) GetLogAuditProof(ctx context.Context, req *GetLogAuditProofRequest) (*ProofResponse, error) {
	return s.client.GetLogAuditProof(ctx, req)
}

func (s *TrustixAPIGRPCClient) GetLogEntries(ctx context.Context, req *GetLogEntriesRequest) (*LogEntriesResponse, error) {
	return s.client.GetLogEntries(ctx, req)
}

func (s *TrustixAPIGRPCClient) GetMapValue(ctx context.Context, req *GetMapValueRequest) (*MapValueResponse, error) {
	return s.client.GetMapValue(ctx, req)
}

func (s *TrustixAPIGRPCClient) Submit(ctx context.Context, req *SubmitRequest) (*SubmitResponse, error) {
	return s.client.Submit(ctx, req)
}
