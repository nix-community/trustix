// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package api

import (
	"context"

	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	ca "github.com/tweag/trustix/packages/trustix/cavaluestore"
	"github.com/tweag/trustix/packages/trustix/rpc/auth"
	"github.com/tweag/trustix/packages/trustix/storage"
)

// TrustixAPIServer wraps kvStoreLogApi and turns it into a gRPC implementation
// It is also responsible for extracting the relevant log implementation for calls that require routing
type TrustixAPIServer struct {
	api.UnimplementedTrustixLogAPIServer
	store  storage.TrustixStorage
	logMap *TrustixLogMap
}

func NewTrustixAPIServer(logMap *TrustixLogMap, store storage.TrustixStorage) (*TrustixAPIServer, error) {
	return &TrustixAPIServer{
		logMap: logMap,
	}, nil
}

func (s *TrustixAPIServer) GetSTH(ctx context.Context, req *api.STHRequest) (*schema.STH, error) {
	log, err := s.logMap.Get(*req.LogID)
	if err != nil {
		return nil, err
	}

	return log.GetSTH(ctx, req)
}

func (s *TrustixAPIServer) GetLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error) {
	log, err := s.logMap.Get(*req.LogID)
	if err != nil {
		return nil, err
	}
	return log.GetLogConsistencyProof(ctx, req)
}

func (s *TrustixAPIServer) GetLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (*api.ProofResponse, error) {
	log, err := s.logMap.Get(*req.LogID)
	if err != nil {
		return nil, err
	}
	return log.GetLogAuditProof(ctx, req)
}

func (s *TrustixAPIServer) GetLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	log, err := s.logMap.Get(*req.LogID)
	if err != nil {
		return nil, err
	}
	return log.GetLogEntries(ctx, req)
}

func (s *TrustixAPIServer) GetMapValue(ctx context.Context, req *api.GetMapValueRequest) (*api.MapValueResponse, error) {
	log, err := s.logMap.Get(*req.LogID)
	if err != nil {
		return nil, err
	}
	return log.GetMapValue(ctx, req)
}

func (s *TrustixAPIServer) Submit(ctx context.Context, req *api.SubmitRequest) (*api.SubmitResponse, error) {

	err := auth.CanWrite(ctx)
	if err != nil {
		return nil, err
	}

	log, err := s.logMap.Get(*req.LogID)
	if err != nil {
		return nil, err
	}

	return log.Submit(ctx, req)
}

func (s *TrustixAPIServer) Flush(ctx context.Context, req *api.FlushRequest) (*api.FlushResponse, error) {
	err := auth.CanWrite(ctx)
	if err != nil {
		return nil, err
	}

	log, err := s.logMap.Get(*req.LogID)
	if err != nil {
		return nil, err
	}

	return log.Flush(ctx, req)
}

func (s *TrustixAPIServer) GetValue(ctx context.Context, req *api.ValueRequest) (*api.ValueResponse, error) {
	var value []byte
	err := s.store.View(func(txn storage.Transaction) error {
		v, err := ca.Get(txn, req.Digest)
		value = v
		return err

	})
	if err != nil {
		return nil, err
	}

	return &api.ValueResponse{
		Value: value,
	}, nil
}

func (s *TrustixAPIServer) GetMHLogConsistencyProof(ctx context.Context, req *api.GetLogConsistencyProofRequest) (*api.ProofResponse, error) {
	log, err := s.logMap.Get(*req.LogID)
	if err != nil {
		return nil, err
	}
	return log.GetMHLogConsistencyProof(ctx, req)
}

func (s *TrustixAPIServer) GetMHLogAuditProof(ctx context.Context, req *api.GetLogAuditProofRequest) (*api.ProofResponse, error) {
	log, err := s.logMap.Get(*req.LogID)
	if err != nil {
		return nil, err
	}
	return log.GetMHLogAuditProof(ctx, req)
}

func (s *TrustixAPIServer) GetMHLogEntries(ctx context.Context, req *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {
	log, err := s.logMap.Get(*req.LogID)
	if err != nil {
		return nil, err
	}
	return log.GetMHLogEntries(ctx, req)
}
