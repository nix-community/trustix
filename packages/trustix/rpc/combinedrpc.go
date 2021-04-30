// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package rpc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/lazyledger/smt"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	pb "github.com/tweag/trustix/packages/trustix-proto/proto"
	rpc "github.com/tweag/trustix/packages/trustix-proto/proto"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	tapi "github.com/tweag/trustix/packages/trustix/api"
	"github.com/tweag/trustix/packages/trustix/decider"
	pub "github.com/tweag/trustix/packages/trustix/publisher"
	"github.com/tweag/trustix/packages/trustix/rpc/auth"
	"github.com/tweag/trustix/packages/trustix/storage"
	storageapi "github.com/tweag/trustix/packages/trustix/storage/api"
)

type TrustixCombinedRPCServer struct {
	pb.UnimplementedTrustixCombinedRPCServer
	logs       *tapi.TrustixLogMap
	decider    decider.LogDecider
	store      storage.TrustixStorage
	publishers *pub.PublisherMap
	signerMeta *SignerMetaMap
}

func NewTrustixCombinedRPCServer(store storage.TrustixStorage, logs *tapi.TrustixLogMap, publishers *pub.PublisherMap, signerMeta *SignerMetaMap, decider decider.LogDecider) *TrustixCombinedRPCServer {
	rpc := &TrustixCombinedRPCServer{
		store:      store,
		logs:       logs,
		decider:    decider,
		publishers: publishers,
		signerMeta: signerMeta,
	}

	return rpc
}

func parseProof(p *api.SparseCompactMerkleProof) smt.SparseCompactMerkleProof {
	return smt.SparseCompactMerkleProof{
		SideNodes:             p.SideNodes,
		NonMembershipLeafData: p.NonMembershipLeafData,
		BitMask:               p.BitMask,
		NumSideNodes:          int(*p.NumSideNodes),
	}
}

func (l *TrustixCombinedRPCServer) GetLogEntries(ctx context.Context, in *pb.GetLogEntriesRequestNamed) (*api.LogEntriesResponse, error) {
	log, err := l.logs.Get(*in.LogID)
	if err != nil {
		return nil, err
	}

	return log.GetLogEntries(ctx, &api.GetLogEntriesRequest{
		LogID:  in.LogID,
		Start:  in.Start,
		Finish: in.Finish,
	})
}

func (l *TrustixCombinedRPCServer) getSTHMap() (map[string]*schema.STH, error) {
	m := make(map[string]*schema.STH)

	err := l.store.View(func(txn storage.Transaction) error {
		logAPI := storageapi.NewStorageAPI(txn)

		for _, logID := range l.logs.IDs() {
			sth, err := logAPI.GetSTH(logID)
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

func (l *TrustixCombinedRPCServer) Get(ctx context.Context, in *pb.KeyRequest) (*pb.EntriesResponse, error) {
	responses := make(map[string]*schema.MapEntry)

	var wg sync.WaitGroup
	var mux sync.Mutex

	hexInput := hex.EncodeToString(in.Key)
	log.WithField("key", hexInput).Info("Received Get request")

	sthMap, err := l.getSTHMap()
	if err != nil {
		return nil, err
	}

	for logID, l := range l.logs.Map() {
		logID := logID
		l := l

		wg.Add(1)

		go func() {
			defer wg.Done()

			log.WithFields(log.Fields{
				"key":   hexInput,
				"logID": logID,
			}).Info("Querying log")

			sth, ok := sthMap[logID]
			if !ok {
				log.Error(fmt.Errorf("could not get STH for logID: %s", logID))
				return
			}

			resp, err := l.GetMapValue(ctx, &api.GetMapValueRequest{
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

func (l *TrustixCombinedRPCServer) GetStream(srv pb.TrustixCombinedRPC_GetStreamServer) error {

	ctx := context.Background()
	var wg sync.WaitGroup
	jobChan := make(chan *pb.KeyRequest)
	errChan := make(chan error)

	numWorkers := 20
	for i := 0; i <= numWorkers; i++ {
		i := i
		go func() {
			log.WithField("count", i).Debug("Creating Get stream worker")

			for in := range jobChan {
				resp, err := l.Get(ctx, in)
				if err != nil {
					wg.Done()
					errChan <- err
					return
				}

				err = srv.Send(resp)
				if err != nil {
					wg.Done()
					errChan <- err
					return
				}

				wg.Done()
			}
		}()
	}

	go func() {
		for {
			in, err := srv.Recv()
			if err != nil {
				if err == io.EOF {
					close(jobChan)
					wg.Wait()
					close(errChan)
					break
				}
				errChan <- err
				return
			}
			wg.Add(1)
			jobChan <- in
		}
	}()

	for err := range errChan {
		return err
	}

	return nil
}

func (l *TrustixCombinedRPCServer) Decide(ctx context.Context, in *pb.KeyRequest) (*pb.DecisionResponse, error) {

	hexInput := hex.EncodeToString(in.Key)
	log.WithField("key", hexInput).Info("Received Decide request")

	var wg sync.WaitGroup
	var mux sync.Mutex
	var inputs []*decider.LogDeciderInput
	var misses []string

	logMap := l.logs.Map()

	sthMap, err := l.getSTHMap()
	if err != nil {
		return nil, err
	}

	for logID, l := range logMap {
		logID := logID
		l := l

		wg.Add(1)

		go func() {
			defer wg.Done()

			log.WithFields(log.Fields{
				"key":   hexInput,
				"logID": logID,
			}).Info("Querying log")

			sth, ok := sthMap[logID]
			if !ok {
				log.Error(fmt.Errorf("could not get STH for logID: %s", logID))
				return
			}
			resp, err := l.GetMapValue(ctx, &api.GetMapValueRequest{
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
				resp, err := logMap[id].GetValue(ctx, &api.ValueRequest{
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

func (l *TrustixCombinedRPCServer) DecideStream(srv pb.TrustixCombinedRPC_DecideStreamServer) error {

	ctx := context.Background()
	var wg sync.WaitGroup
	jobChan := make(chan *pb.KeyRequest)
	errChan := make(chan error)

	numWorkers := 20
	for i := 0; i <= numWorkers; i++ {
		i := i
		go func() {
			log.WithField("count", i).Debug("Creating Decide stream worker")

			for in := range jobChan {
				resp, err := l.Decide(ctx, in)
				if err != nil {
					wg.Done()
					errChan <- err
					return
				}

				err = srv.Send(resp)
				if err != nil {
					wg.Done()
					errChan <- err
					return
				}

				wg.Done()
			}
		}()
	}

	go func() {
		for {
			in, err := srv.Recv()
			if err != nil {
				if err == io.EOF {
					close(jobChan)
					wg.Wait()
					close(errChan)
					break
				}
				errChan <- err
				return
			}
			wg.Add(1)
			jobChan <- in
		}
	}()

	for err := range errChan {
		return err
	}

	return nil
}

func (l *TrustixCombinedRPCServer) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {

	log.Info("Received Value request")

	for name, logApi := range l.logs.Map() {
		resp, err := logApi.GetValue(ctx, in)
		if err != nil {
			log.WithFields(log.Fields{
				"logID": name,
			}).Info("Error receiving value")
		} else {
			return resp, nil
		}
	}

	return nil, fmt.Errorf("Value could not be retreived")
}

func (l *TrustixCombinedRPCServer) Submit(ctx context.Context, req *rpc.SubmitRequest) (*rpc.SubmitResponse, error) {
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

func (l *TrustixCombinedRPCServer) Flush(ctx context.Context, req *rpc.FlushRequest) (*rpc.FlushResponse, error) {
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

func (l *TrustixCombinedRPCServer) Logs(ctx context.Context, in *pb.LogsRequest) (*pb.LogsResponse, error) {
	logs := []*pb.Log{}

	sthMap, err := l.getSTHMap()
	if err != nil {
		return nil, err
	}

	// TODO: Don't hard-code meta, take from config
	var meta = map[string]string{
		"upstream": "https://cache.nixos.org",
	}

	for logID, _ := range l.logs.Map() {
		sth, ok := sthMap[logID]
		if !ok {
			return nil, fmt.Errorf("Missing STH for logID '%s'", logID)
		}

		signer, err := l.signerMeta.Get(logID)
		if err != nil {
			return nil, fmt.Errorf("Missing signer meta for log '%s': %v", logID, err)
		}

		logs = append(logs, &pb.Log{
			LogID:  &logID,
			STH:    sth,
			Meta:   meta,
			Signer: signer,
		})
	}

	return &pb.LogsResponse{
		Logs: logs,
	}, nil
}
