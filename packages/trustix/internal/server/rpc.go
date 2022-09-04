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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lazyledger/smt"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	pb "github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
	"github.com/nix-community/trustix/packages/trustix/client"
	"github.com/nix-community/trustix/packages/trustix/internal/decider"
	"github.com/nix-community/trustix/packages/trustix/internal/pool"
	"github.com/nix-community/trustix/packages/trustix/internal/protocols"
	pub "github.com/nix-community/trustix/packages/trustix/internal/publisher"
	"github.com/nix-community/trustix/packages/trustix/internal/storage"
	log "github.com/sirupsen/logrus"
)

type RPCServer struct {
	pb.UnimplementedRPCApiServer
	decider    map[string]decider.LogDecider
	store      storage.Storage
	publishers *pub.PublisherMap
	rootBucket *storage.Bucket
	logs       []*api.Log
	clients    *pool.ClientPool
}

func NewRPCServer(
	store storage.Storage,
	rootBucket *storage.Bucket,
	clients *pool.ClientPool,
	publishers *pub.PublisherMap,
	logs []*api.Log,
	decider map[string]decider.LogDecider,
) *RPCServer {
	return &RPCServer{
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

func (l *RPCServer) GetLogEntries(ctx context.Context, in *api.GetLogEntriesRequest) (*api.LogEntriesResponse, error) {

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

func (l *RPCServer) Decide(ctx context.Context, in *pb.DecideRequest) (*pb.DecisionResponse, error) {

	hexInput := hex.EncodeToString(in.Key)
	log.WithField("key", hexInput).Info("Received Decide request")

	pd, err := protocols.Get(*in.Protocol)
	if err != nil {
		return nil, fmt.Errorf("Unknown protocol: %w", err)
	}

	deciderEngine, ok := l.decider[pd.ID]
	if !ok {
		return nil, fmt.Errorf("No decider configured for protocol '%s'", pd.ID)
	}

	var wg sync.WaitGroup
	var mux sync.Mutex
	var inputs []*decider.DeciderInput
	var misses []string

	logs := l.logs

	sthMap := make(map[string]*schema.LogHead)
	{
		err := l.store.View(func(txn storage.Transaction) error {
			for _, log := range logs {
				logID := *log.LogID
				head, err := getLogHead(l.rootBucket, txn, logID)
				if err != nil {
					return err
				}
				sthMap[logID] = head
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
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

	for _, logMeta := range logs {
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

			valid := smt.VerifyCompactProof(parseProof(resp.Proof), sth.MapRoot, in.Key, resp.Value, pd.NewHash())
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

			inputs = append(inputs, &decider.DeciderInput{
				LogID: logID,
				Value: hex.EncodeToString(mapEntry.Digest),
			})

		}()
	}

	wg.Wait()

	resp := &pb.DecisionResponse{
		Misses: misses,
	}

	var decision *decider.DeciderOutput
	if len(inputs) > 0 {

		var err error
		decision, err = deciderEngine.Decide(inputs)
		if err != nil {
			return nil, err
		}

		outputHash, err := hex.DecodeString(decision.Value)
		if err != nil {
			return nil, err
		}

		confidence := int32(decision.Confidence)
		if decision.Value != "" {
			logIDs := []string{}
			for _, input := range inputs {
				if input.Value == decision.Value {
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
		if input.Value == decision.Value {
			continue
		}

		digest, err := hex.DecodeString(input.Value)
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

func (l *RPCServer) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {

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

func (l *RPCServer) Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error) {
	return &api.LogsResponse{
		Logs: l.logs,
	}, nil
}
