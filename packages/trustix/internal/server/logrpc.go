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
	"github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix-proto/rpc/rpcconnect"
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
	"github.com/nix-community/trustix/packages/trustix/internal/pool"
	pub "github.com/nix-community/trustix/packages/trustix/internal/publisher"
	"github.com/nix-community/trustix/packages/trustix/internal/rpc/auth"
	"github.com/nix-community/trustix/packages/trustix/internal/storage"
)

type LogRPCServer struct {
	rpcconnect.UnimplementedLogRPCHandler

	publishers *pub.PublisherMap
	clients    *pool.ClientPool
	store      storage.Storage
	rootBucket *storage.Bucket
}

func NewLogRPCServer(
	store storage.Storage,
	rootBucket *storage.Bucket,
	clients *pool.ClientPool,
	publishers *pub.PublisherMap,
) *LogRPCServer {
	return &LogRPCServer{
		store:      store,
		clients:    clients,
		publishers: publishers,
		rootBucket: rootBucket,
	}
}

func (l *LogRPCServer) GetHead(ctx context.Context, req *connect.Request[api.LogHeadRequest]) (*connect.Response[schema.LogHead], error) {
	msg := req.Msg

	logID := *msg.LogID
	var head *schema.LogHead
	err := l.store.View(func(txn storage.Transaction) error {
		var err error
		head, err = storage.GetLogHead(l.rootBucket.Cd(logID).Txn(txn))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(head), nil
}

func (l *LogRPCServer) GetLogEntries(ctx context.Context, req *connect.Request[api.GetLogEntriesRequest]) (*connect.Response[api.LogEntriesResponse], error) {
	msg := req.Msg

	client, err := l.clients.Get(*msg.LogID)
	if err != nil {
		return nil, err
	}

	logEntriesResp, err := client.LogAPI.GetLogEntries(ctx, &api.GetLogEntriesRequest{
		LogID:  msg.LogID,
		Start:  msg.Start,
		Finish: msg.Finish,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(logEntriesResp), nil
}

func (l *LogRPCServer) Submit(ctx context.Context, req *connect.Request[rpc.SubmitRequest]) (*connect.Response[rpc.SubmitResponse], error) {
	msg := req.Msg

	err := auth.CanWrite(ctx)
	if err != nil {
		return nil, err
	}

	q, err := l.publishers.Get(*msg.LogID)
	if err != nil {
		return nil, err
	}

	submitResp, err := q.Submit(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(submitResp), nil
}

func (l *LogRPCServer) Flush(ctx context.Context, req *connect.Request[rpc.FlushRequest]) (*connect.Response[rpc.FlushResponse], error) {
	msg := req.Msg

	err := auth.CanWrite(ctx)
	if err != nil {
		return nil, err
	}

	q, err := l.publishers.Get(*msg.LogID)
	if err != nil {
		return nil, err
	}

	flushResp, err := q.Flush(ctx, msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(flushResp), nil
}
