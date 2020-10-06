package rpc

import (
	"context"
	"encoding/hex"
	"fmt"
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
	fmt.Println(fmt.Sprintf("Received KV request for %s", hex.EncodeToString(in.Key)))

	v, err := s.core.Get(in.Bucket, in.Key)
	if err != nil {
		return nil, err
	}

	return &pb.KVResponse{
		Value: v,
	}, nil
}
