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

type TrustixAPIGRPCClient struct {
	client TrustixLogAPIClient
}

func NewTrustixAPIGRPCClient(conn *grpc.ClientConn) (*TrustixAPIGRPCClient, error) {
	c := NewTrustixLogAPIClient(conn)
	return &TrustixAPIGRPCClient{
		client: c,
	}, nil
}

func (s *TrustixAPIGRPCClient) GetSTH(req *STHRequest) (*schema.STH, error) {
	return s.client.GetSTH(context.Background(), req)
}

func (s *TrustixAPIGRPCClient) GetLogConsistencyProof(req *GetLogConsistencyProofRequest) (*ProofResponse, error) {
	return s.client.GetLogConsistencyProof(context.Background(), req)
}

func (s *TrustixAPIGRPCClient) GetLogAuditProof(req *GetLogAuditProofRequest) (*ProofResponse, error) {
	return s.client.GetLogAuditProof(context.Background(), req)
}

func (s *TrustixAPIGRPCClient) GetLogEntries(req *GetLogEntriesRequest) (*LogEntriesResponse, error) {
	return s.client.GetLogEntries(context.Background(), req)
}

func (s *TrustixAPIGRPCClient) GetMapValue(req *GetMapValueRequest) (*MapValueResponse, error) {
	return s.client.GetMapValue(context.Background(), req)
}

func (s *TrustixAPIGRPCClient) Submit(req *SubmitRequest) (*SubmitResponse, error) {
	return s.client.Submit(context.Background(), req)
}
