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
