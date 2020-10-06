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
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/core"
	pb "github.com/tweag/trustix/proto"
	"sync"
)

type TrustixLogServer struct {
	pb.UnimplementedTrustixLogServer
	logMap map[string]*core.TrustixCore
}

func NewTrustixLogServer(logMap map[string]*core.TrustixCore) *TrustixLogServer {
	return &TrustixLogServer{logMap: logMap}
}

func (l *TrustixLogServer) HashMap(ctx context.Context, in *pb.HashRequest) (*pb.HashMapResponse, error) {
	responses := make(map[string][]byte)

	var wg sync.WaitGroup
	var mux sync.Mutex

	hexInput := hex.EncodeToString(in.InputHash)
	log.WithField("inputHash", hexInput).Info("Received HashMap request")

	for name, l := range l.logMap {
		// Create copies for goroutine
		name := name
		l := l

		wg.Add(1)

		go func() {
			defer wg.Done()

			log.WithFields(log.Fields{
				"inputHash": hexInput,
				"logName":   name,
			}).Info("Querying log")

			h, err := l.Query(in.InputHash)
			if err != nil {
				fmt.Println(err)
			}

			mux.Lock()
			responses[name] = h
			mux.Unlock()
		}()
	}

	wg.Wait()

	return &pb.HashMapResponse{
		Hashes: responses,
	}, nil

}
