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
	pb "github.com/tweag/trustix/packages/trustix-proto/rpc"
	rpc "github.com/tweag/trustix/packages/trustix-proto/rpc"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	"github.com/tweag/trustix/packages/trustix/internal/pool"
	pub "github.com/tweag/trustix/packages/trustix/internal/publisher"
	"github.com/tweag/trustix/packages/trustix/internal/rpc/auth"
	"github.com/tweag/trustix/packages/trustix/internal/storage"
)

type LogRPCServer struct {
	pb.UnimplementedLogRPCServer

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

func (l *LogRPCServer) GetHead(ctx context.Context, req *api.LogHeadRequest) (*schema.LogHead, error) {
	logID := *req.LogID
	var head *schema.LogHead
	err := l.store.View(func(txn storage.Transaction) error {
		var err error
		head, err = storage.GetLogHead(l.rootBucket.Cd(logID).Txn(txn))
		if err != nil {
			return err
		}

		return nil
	})
	return head, err
}

func (l *LogRPCServer) GetLogEntries(ctx context.Context, in *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {

	client, err := l.clients.Get(*in.LogID)
	if err != nil {
		return nil, err
	}

	return client.LogAPI.GetLogEntries(ctx, &api.GetLogEntriesRequest{
		LogID:  in.LogID,
		Start:  in.Start,
		Finish: in.Finish,
	})
}

func (l *LogRPCServer) Submit(ctx context.Context, req *rpc.SubmitRequest) (*rpc.SubmitResponse, error) {
	err := auth.CanWrite(ctx)
	if err != nil {
		return nil, err
	}

	q, err := l.publishers.Get(*req.LogID)
	if err != nil {
		return nil, err
	}

	return q.Submit(ctx, req)
}

func (l *LogRPCServer) Flush(ctx context.Context, req *rpc.FlushRequest) (*rpc.FlushResponse, error) {
	err := auth.CanWrite(ctx)
	if err != nil {
		return nil, err
	}

	q, err := l.publishers.Get(*req.LogID)
	if err != nil {
		return nil, err
	}

	return q.Flush(ctx, req)
}
