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

// Common log operations

import (
	"context"
	"fmt"

	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	vlog "github.com/tweag/trustix/packages/trustix/log"
	"github.com/tweag/trustix/packages/trustix/storage"
)

func getLogConsistencyProof(prefix string, txn storage.Transaction, ctx context.Context, req *api.GetLogConsistencyProofRequest) (resp *api.ProofResponse, err error) {
	resp = &api.ProofResponse{}

	vLog, err := vlog.NewVerifiableLog(prefix, txn, *req.SecondSize)
	if err != nil {
		return nil, err
	}

	proof, err := vLog.ConsistencyProof(*req.FirstSize, *req.SecondSize)
	if err != nil {
		return nil, err
	}

	resp.Proof = proof

	return resp, nil
}

func getLogAuditProof(prefix string, txn storage.Transaction, ctx context.Context, req *api.GetLogAuditProofRequest) (resp *api.ProofResponse, err error) {

	vLog, err := vlog.NewVerifiableLog(prefix, txn, *req.TreeSize)
	if err != nil {
		return nil, err
	}

	proof, err := vLog.AuditProof(*req.Index, *req.TreeSize)
	if err != nil {
		return nil, err
	}

	resp.Proof = proof

	return resp, nil
}

func getLogEntries(prefix string, txn storage.Transaction, ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {

	resp := &api.LogEntriesResponse{
		Leaves: []*schema.LogLeaf{},
	}

	logStorage := vlog.NewLogStorage(prefix, txn)

	for i := int(*req.Start); i <= int(*req.Finish); i++ {
		leaf, err := logStorage.Get(0, uint64(i))
		if err != nil {
			return nil, err
		}
		resp.Leaves = append(resp.Leaves, leaf)
	}

	total := *req.Finish - *req.Start + 1
	if uint64(len(resp.Leaves)) != total {
		return nil, fmt.Errorf("%d != %d", len(resp.Leaves), total)
	}

	return resp, nil
}
