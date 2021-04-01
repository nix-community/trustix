// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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
