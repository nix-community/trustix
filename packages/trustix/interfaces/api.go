// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package interfaces

import (
	"context"

	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
)

type NodeAPI interface {
	// List published logs
	Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error)

	// Get content-addressed computation outputs
	GetValue(context.Context, *api.ValueRequest) (*api.ValueResponse, error)
}

type LogAPI interface {

	// GetHead - Get a signed tree head
	GetHead(context.Context, *api.LogHeadRequest) (*schema.LogHead, error)

	// Get log consistency proof
	GetLogConsistencyProof(context.Context, *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error)

	// Get leaf audit proof
	GetLogAuditProof(context.Context, *api.GetLogAuditProofRequest) (*api.ProofResponse, error)

	// Get log entries (batched)
	GetLogEntries(context.Context, *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error)

	// Get value from the map
	GetMapValue(context.Context, *api.GetMapValueRequest) (*api.MapValueResponse, error)

	// Get map head log consistency proof
	GetMHLogConsistencyProof(context.Context, *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error)

	// Get map head log leaf audit proof
	GetMHLogAuditProof(context.Context, *api.GetLogAuditProofRequest) (*api.ProofResponse, error)

	// Get map head log entries (batched)
	GetMHLogEntries(context.Context, *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error)
}
