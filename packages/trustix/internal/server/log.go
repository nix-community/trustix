// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package server

import (
	"context"
	"fmt"

	connect "github.com/bufbuild/connect-go"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/api/apiconnect"
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
	"github.com/nix-community/trustix/packages/trustix/interfaces"
	"github.com/nix-community/trustix/packages/trustix/internal/pool"
)

// LogAPIServer wraps kvStoreLogApi and turns it into a gRPC implementation
type LogAPIServer struct {
	apiconnect.UnimplementedNodeAPIHandler

	clients *pool.ClientPool

	// Keep a set of published logs around to filter on, we don't want to leak what we subscribe to
	published map[string]struct{}
}

func NewLogAPIServer(logs []*api.Log, clients *pool.ClientPool) *LogAPIServer {
	published := make(map[string]struct{})
	for _, log := range logs {
		published[*log.LogID] = struct{}{}
	}

	return &LogAPIServer{
		published: published,
		clients:   clients,
	}
}

func (s *LogAPIServer) getClient(logID string) (interfaces.LogAPI, error) {
	_, ok := s.published[logID]
	if !ok {
		return nil, fmt.Errorf("Log with ID '%s' is not published", logID)
	}

	client, err := s.clients.Get(logID)
	if err != nil {
		return nil, err
	}

	return client.LogAPI, nil
}

func (s *LogAPIServer) GetHead(ctx context.Context, req *connect.Request[api.LogHeadRequest]) (*connect.Response[schema.LogHead], error) {
	msg := req.Msg

	log, err := s.getClient(*msg.LogID)
	if err != nil {
		return nil, err
	}

	resp, err := log.GetHead(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (s *LogAPIServer) GetLogConsistencyProof(ctx context.Context, req *connect.Request[api.GetLogConsistencyProofRequest]) (*connect.Response[api.ProofResponse], error) {
	msg := req.Msg

	log, err := s.getClient(*msg.LogID)
	if err != nil {
		return nil, err
	}

	resp, err := log.GetLogConsistencyProof(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (s *LogAPIServer) GetLogAuditProof(ctx context.Context, req *connect.Request[api.GetLogAuditProofRequest]) (*connect.Response[api.ProofResponse], error) {
	msg := req.Msg

	log, err := s.getClient(*msg.LogID)
	if err != nil {
		return nil, err
	}

	resp, err := log.GetLogAuditProof(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (s *LogAPIServer) GetLogEntries(ctx context.Context, req *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error) {
	msg := req.Msg

	log, err := s.getClient(*msg.LogID)
	if err != nil {
		return nil, err
	}

	resp, err := log.GetLogEntries(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (s *LogAPIServer) GetMapValue(ctx context.Context, req *connect.Request[api.GetMapValueRequest]) (*connect.Response[api.MapValueResponse], error) {
	msg := req.Msg

	log, err := s.getClient(*msg.LogID)
	if err != nil {
		return nil, err
	}

	resp, err := log.GetMapValue(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (s *LogAPIServer) GetMHLogConsistencyProof(ctx context.Context, req *connect.Request[api.GetLogConsistencyProofRequest]) (*connect.Response[api.ProofResponse], error) {
	msg := req.Msg

	log, err := s.getClient(*msg.LogID)
	if err != nil {
		return nil, err
	}

	resp, err := log.GetMHLogConsistencyProof(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (s *LogAPIServer) GetMHLogAuditProof(ctx context.Context, req *connect.Request[api.GetLogAuditProofRequest]) (*connect.Response[api.ProofResponse], error) {
	msg := req.Msg

	log, err := s.getClient(*msg.LogID)
	if err != nil {
		return nil, err
	}

	resp, err := log.GetMHLogAuditProof(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (s *LogAPIServer) GetMHLogEntries(ctx context.Context, req *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error) {
	msg := req.Msg

	log, err := s.getClient(*msg.LogID)
	if err != nil {
		return nil, err
	}

	resp, err := log.GetMHLogEntries(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}
