// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package interfaces

import (
	"context"

	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
)

type RpcAPI interface {
	Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error)
	Decide(ctx context.Context, in *rpc.DecideRequest) (*rpc.DecisionResponse, error)
	GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error)
}

type LogRPC interface {
	GetHead(context.Context, *api.LogHeadRequest) (*schema.LogHead, error)
	GetLogEntries(ctx context.Context, in *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error)
	Submit(ctx context.Context, in *rpc.SubmitRequest) (*rpc.SubmitResponse, error)
	Flush(ctx context.Context, in *rpc.FlushRequest) (*rpc.FlushResponse, error)
}
