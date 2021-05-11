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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lazyledger/smt"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	pb "github.com/tweag/trustix/packages/trustix-proto/rpc"
	rpc "github.com/tweag/trustix/packages/trustix-proto/rpc"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	"github.com/tweag/trustix/packages/trustix/client"
	"github.com/tweag/trustix/packages/trustix/decider"
	pub "github.com/tweag/trustix/packages/trustix/publisher"
	"github.com/tweag/trustix/packages/trustix/rpc/auth"
	"github.com/tweag/trustix/packages/trustix/storage"
)

type TrustixRPCServer struct {
	pb.UnimplementedTrustixRPCServer
	decider    decider.LogDecider
	store      storage.Storage
	publishers *pub.PublisherMap
	rootBucket *storage.Bucket
	logs       []*api.Log
	clients    *client.ClientPool
}

func NewTrustixRPCServer(
	store storage.Storage,
	rootBucket *storage.Bucket,
	clients *client.ClientPool,
	publishers *pub.PublisherMap,
	logs []*api.Log,
	decider decider.LogDecider,
) *TrustixRPCServer {
	return &TrustixRPCServer{
		store:      store,
		decider:    decider,
		publishers: publishers,
		logs:       logs,
		rootBucket: rootBucket,
		clients:    clients,
	}
}

func parseProof(p *api.SparseCompactMerkleProof) smt.SparseCompactMerkleProof {
	return smt.SparseCompactMerkleProof{
		SideNodes:             p.SideNodes,
		NonMembershipLeafData: p.NonMembershipLeafData,
		BitMask:               p.BitMask,
		NumSideNodes:          int(*p.NumSideNodes),
	}
}

func (l *TrustixRPCServer) GetLogEntries(ctx context.Context, in *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {

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

func (l *TrustixRPCServer) getSTHMap() (map[string]*schema.STH, error) {
	m := make(map[string]*schema.STH)

	err := l.store.View(func(txn storage.Transaction) error {
		for _, log := range l.logs {
			logID := *log.LogID
			bucket := l.rootBucket.Cd(logID)
			sth, err := storage.GetSTH(bucket.Txn(txn))
			if err != nil {
				return err
			}
			m[logID] = sth
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (l *TrustixRPCServer) Get(ctx context.Context, in *pb.KeyRequest) (*pb.EntriesResponse, error) {
	responses := make(map[string]*schema.MapEntry)

	var wg sync.WaitGroup
	var mux sync.Mutex

	hexInput := hex.EncodeToString(in.Key)
	log.WithField("key", hexInput).Info("Received Get request")

	sthMap, err := l.getSTHMap()
	if err != nil {
		return nil, err
	}

	for _, logMeta := range l.logs {
		logMeta := logMeta

		wg.Add(1)

		go func() {
			defer wg.Done()

			logID := *logMeta.LogID

			log.WithFields(log.Fields{
				"key":   hexInput,
				"logID": logID,
			}).Info("Querying log")

			client, err := l.clients.Get(logID)
			if err != nil {
				log.Error(fmt.Errorf("could not get client for logID '%s': %v", logID, err))
				return
			}

			sth, ok := sthMap[logID]
			if !ok {
				log.Error(fmt.Errorf("could not get STH for logID: %s", logID))
				return
			}

			resp, err := client.LogAPI.GetMapValue(ctx, &api.GetMapValueRequest{
				LogID:   &logID,
				Key:     in.Key,
				MapRoot: sth.MapRoot,
			})
			if err != nil {
				log.Error(fmt.Sprintf("could not query: %v", err))
				return
			}

			valid := smt.VerifyCompactProof(parseProof(resp.Proof), sth.MapRoot, in.Key, resp.Value, sha256.New())
			if !valid {
				log.Error("SMT proof verification failed")
				return
			}

			entry := &schema.MapEntry{}
			err = json.Unmarshal(resp.Value, entry)
			if err != nil {
				log.Error("Could not unmarshal map value")
				return
			}

			mux.Lock()
			responses[logID] = entry
			mux.Unlock()
		}()
	}

	wg.Wait()

	return &pb.EntriesResponse{
		Key:     in.Key,
		Entries: responses,
	}, nil

}

func (l *TrustixRPCServer) Decide(ctx context.Context, in *pb.KeyRequest) (*pb.DecisionResponse, error) {

	hexInput := hex.EncodeToString(in.Key)
	log.WithField("key", hexInput).Info("Received Decide request")

	var wg sync.WaitGroup
	var mux sync.Mutex
	var inputs []*decider.LogDeciderInput
	var misses []string

	sthMap, err := l.getSTHMap()
	if err != nil {
		return nil, err
	}

	logIDs := []string{}
	clients := make(map[string]*client.Client)
	for logID := range sthMap {
		client, err := l.clients.Get(logID)
		if err != nil {
			log.Error(fmt.Errorf("could not get client for logID '%s': %v", logID, err))
			continue
		}
		clients[logID] = client
		logIDs = append(logIDs, logID)
	}

	for _, logMeta := range l.logs {
		logMeta := logMeta

		wg.Add(1)

		go func() {
			defer wg.Done()

			logID := *logMeta.LogID

			log.WithFields(log.Fields{
				"key":   hexInput,
				"logID": logID,
			}).Info("Querying log")

			sth, ok := sthMap[logID]
			if !ok {
				log.Error(fmt.Errorf("could not get STH for logID: %s", logID))
				return
			}
			resp, err := clients[logID].LogAPI.GetMapValue(ctx, &api.GetMapValueRequest{
				LogID:   &logID,
				Key:     in.Key,
				MapRoot: sth.MapRoot,
			})
			if err != nil {
				log.WithField("error", err).Error("Could not get query")
				return
			}

			valid := smt.VerifyCompactProof(parseProof(resp.Proof), sth.MapRoot, in.Key, resp.Value, sha256.New())
			if !valid {
				log.Error("SMT proof verification failed")
				return
			}

			mapEntry := &schema.MapEntry{}
			err = json.Unmarshal(resp.Value, mapEntry)
			if err != nil {
				log.Error("Could not unmarshal value")
				return
			}

			mux.Lock()
			defer mux.Unlock()

			if err != nil {
				log.Error(err)
				misses = append(misses, logID)
				return
			}

			if len(resp.Value) == 0 {
				misses = append(misses, logID)
				return
			}

			inputs = append(inputs, &decider.LogDeciderInput{
				LogID:      logID,
				OutputHash: hex.EncodeToString(mapEntry.Digest),
			})

		}()
	}

	wg.Wait()

	resp := &pb.DecisionResponse{
		Misses: misses,
	}

	var decision *decider.LogDeciderOutput
	if len(inputs) > 0 {

		var err error
		decision, err = l.decider.Decide(inputs)
		if err != nil {
			return nil, err
		}

		outputHash, err := hex.DecodeString(decision.OutputHash)
		if err != nil {
			return nil, err
		}

		confidence := int32(decision.Confidence)
		if decision.OutputHash != "" {
			logIDs := []string{}
			for _, input := range inputs {
				if input.OutputHash == decision.OutputHash {
					logIDs = append(logIDs, input.LogID)
				}
			}

			value := []byte{}
			for _, id := range logIDs {
				resp, err := clients[id].NodeAPI.GetValue(ctx, &api.ValueRequest{
					Digest: outputHash,
				})
				if err != nil {
					return nil, err
				}
				value = resp.Value
			}

			resp.Decision = &pb.LogValueDecision{
				LogIDs:     logIDs,
				Digest:     outputHash,
				Confidence: &confidence,
				Value:      value,
			}
		}
	} else {
		return nil, fmt.Errorf("No outputs found for input")
	}

	// inputMap := make(map[string][]byte)
	for _, input := range inputs {
		if input.OutputHash == decision.OutputHash {
			continue
		}

		digest, err := hex.DecodeString(input.OutputHash)
		if err != nil {
			return nil, err
		}

		resp.Mismatches = append(resp.Mismatches, &pb.LogValueResponse{
			LogID:  &input.LogID,
			Digest: digest,
		})
	}

	return resp, nil

}

func (l *TrustixRPCServer) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {

	log.Info("Received Value request")

	for _, logMeta := range l.logs {
		logID := *logMeta.LogID

		client, err := l.clients.Get(logID)
		if err != nil {
			log.WithFields(log.Fields{
				"logID": logID,
			}).Info("Could not find active client")
			continue
		}

		resp, err := client.NodeAPI.GetValue(ctx, in)
		if err != nil {
			log.WithFields(log.Fields{
				"logID": logID,
				"error": err,
			}).Info("Error receiving value")
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("Value could not be retreived")
}

func (l *TrustixRPCServer) Submit(ctx context.Context, req *rpc.SubmitRequest) (*rpc.SubmitResponse, error) {
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

func (l *TrustixRPCServer) Flush(ctx context.Context, req *rpc.FlushRequest) (*rpc.FlushResponse, error) {
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

func (l *TrustixRPCServer) Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error) {
	return &api.LogsResponse{
		Logs: l.logs,
	}, nil
}
