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

	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix/interfaces"
)

// NodeAPIServer wraps NodeAPI and turns it into a gRPC implementation
type NodeAPIServer struct {
	api.UnimplementedNodeAPIServer
	node interfaces.NodeAPI
}

func NewNodeAPIServer(node interfaces.NodeAPI) *NodeAPIServer {
	return &NodeAPIServer{
		node: node,
	}
}

func (s *NodeAPIServer) Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error) {
	return s.node.Logs(ctx, in)
}

func (s *NodeAPIServer) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {
	return s.node.GetValue(ctx, in)
}
