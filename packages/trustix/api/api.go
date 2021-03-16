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

	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
)

type TrustixLogAPI interface {
	// GetSTH - Get a signed tree head
	GetSTH(context.Context, *api.STHRequest) (*schema.STH, error)

	// Get log consistency proof
	GetLogConsistencyProof(context.Context, *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error)

	// Get leaf audit proof
	GetLogAuditProof(context.Context, *api.GetLogAuditProofRequest) (*api.ProofResponse, error)

	// Get log entries (batched)
	GetLogEntries(context.Context, *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error)

	// Get value from the map
	GetMapValue(context.Context, *api.GetMapValueRequest) (*api.MapValueResponse, error)

	// Submit a value to the log
	Submit(context.Context, *api.SubmitRequest) (*api.SubmitResponse, error)

	// Flush the queue and write new heads
	Flush(context.Context, *api.FlushRequest) (*api.FlushResponse, error)

	// Get content-addressed computation outputs
	GetValue(context.Context, *api.ValueRequest) (*api.ValueResponse, error)

	// Get map head log consistency proof
	GetMHLogConsistencyProof(context.Context, *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error)

	// Get map head log leaf audit proof
	GetMHLogAuditProof(context.Context, *api.GetLogAuditProofRequest) (*api.ProofResponse, error)

	// Get map head log entries (batched)
	GetMHLogEntries(context.Context, *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error)
}
