// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package server

import (
	"context"

	connect "connectrpc.com/connect"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/api/apiconnect"
	"github.com/nix-community/trustix/packages/trustix/interfaces"
)

// NodeAPIServer wraps NodeAPI and turns it into a gRPC implementation
type NodeAPIServer struct {
	apiconnect.UnimplementedNodeAPIHandler
	node interfaces.NodeAPI
}

func NewNodeAPIServer(node interfaces.NodeAPI) *NodeAPIServer {
	return &NodeAPIServer{
		node: node,
	}
}

func (s *NodeAPIServer) Logs(ctx context.Context, req *connect.Request[api.LogsRequest]) (*connect.Response[api.LogsResponse], error) {
	msg := req.Msg

	resp, err := s.node.Logs(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (s *NodeAPIServer) GetValue(ctx context.Context, req *connect.Request[api.ValueRequest]) (*connect.Response[api.ValueResponse], error) {
	msg := req.Msg

	resp, err := s.node.GetValue(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}
