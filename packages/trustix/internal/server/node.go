// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package server

import (
	"context"

	connect "github.com/bufbuild/connect-go"
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
