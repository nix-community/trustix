package rpc

import (
	"context"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/tweag/trustix/core"
	pb "github.com/tweag/trustix/proto"
)

type TrustixKVServer struct {
	pb.UnimplementedTrustixKVServer
	core *core.TrustixCore
}

func NewTrustixKVServer(core *core.TrustixCore) *TrustixKVServer {
	return &TrustixKVServer{core: core}
}

func (s *TrustixKVServer) Get(ctx context.Context, in *pb.KVRequest) (*pb.KVResponse, error) {
	log.WithField("key", hex.EncodeToString(in.Key)).Info("Received KV request")

	v, err := s.core.Get(in.Bucket, in.Key)
	if err != nil {
		return nil, err
	}

	return &pb.KVResponse{
		Value: v,
	}, nil
}
