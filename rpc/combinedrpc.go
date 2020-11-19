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

package rpc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/lazyledger/smt"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/api"
	"github.com/tweag/trustix/decider"
	pb "github.com/tweag/trustix/proto"
	"github.com/tweag/trustix/schema"
	"github.com/tweag/trustix/sthmanager"
)

type TrustixCombinedRPCServer struct {
	pb.UnimplementedTrustixCombinedRPCServer
	logs       *TrustixCombinedRPCServerMap
	decider    decider.LogDecider
	sthmanager *sthmanager.STHManager
}

func NewTrustixCombinedRPCServer(sthmanager *sthmanager.STHManager, logs *TrustixCombinedRPCServerMap, decider decider.LogDecider) *TrustixCombinedRPCServer {
	return &TrustixCombinedRPCServer{
		logs:       logs,
		decider:    decider,
		sthmanager: sthmanager,
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

func (l *TrustixCombinedRPCServer) Get(ctx context.Context, in *pb.KeyRequest) (*pb.EntriesResponse, error) {
	responses := make(map[string]*schema.MapEntry)

	var wg sync.WaitGroup
	var mux sync.Mutex

	hexInput := hex.EncodeToString(in.Key)
	log.WithField("key", hexInput).Info("Received Get request")

	getSTH := l.sthmanager.Get

	for name, l := range l.logs.Map() {
		name := name
		l := l

		wg.Add(1)

		go func() {
			defer wg.Done()

			log.WithFields(log.Fields{
				"key":     hexInput,
				"logName": name,
			}).Info("Querying log")

			sth, err := getSTH(name)
			if err != nil {
				log.Error(fmt.Errorf("could not get STH: %v", err))
				return
			}

			resp, err := l.GetMapValue(ctx, &api.GetMapValueRequest{
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
			err = proto.Unmarshal(resp.Value, entry)
			if err != nil {
				log.Error("Could not unmarshal map value")
				return
			}

			mux.Lock()
			responses[name] = entry
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

	getSTH := l.sthmanager.Get

	for name, l := range l.logs.Map() {
		name := name
		l := l

		wg.Add(1)

		go func() {
			defer wg.Done()

			log.WithFields(log.Fields{
				"key":     hexInput,
				"logName": name,
			}).Info("Querying log")

			sth, err := getSTH(name)
			if err != nil {
				log.WithField("error", err).Error("Could not get STH")
				return
			}
			resp, err := l.GetMapValue(ctx, &api.GetMapValueRequest{
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
				misses = append(misses, name)
				return
			}

			if len(resp.Value) == 0 {
				misses = append(misses, name)
				return
			}

			inputs = append(inputs, &decider.LogDeciderInput{
				LogName:    name,
				OutputHash: hex.EncodeToString(mapEntry.Value),
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
			logNames := []string{}
			for _, input := range inputs {
				if input.OutputHash == decision.OutputHash {
					logNames = append(logNames, input.LogName)
				}
			}

			resp.Decision = &pb.LogValueDecision{
				LogNames:   logNames,
				Value:      outputHash,
				Confidence: &confidence,
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

		h, err := hex.DecodeString(input.OutputHash)
		if err != nil {
			return nil, err
		}

		resp.Mismatches = append(resp.Mismatches, &pb.LogValueResponse{
			LogName: &input.LogName,
			Value:   h,
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
