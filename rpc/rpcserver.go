package rpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/tweag/trustix/core"
	pb "github.com/tweag/trustix/proto"
)

type TrustixRPCServer struct {
	pb.UnimplementedTrustixRPCServer
	core *core.TrustixCore
}

func NewTrustixRPCServer(core *core.TrustixCore) *TrustixRPCServer {
	return &TrustixRPCServer{core: core}
}

func (s *TrustixRPCServer) SubmitMapping(ctx context.Context, in *pb.SubmitRequest) (*pb.SubmitResponse, error) {
	fmt.Println(fmt.Sprintf("Received input hash %s -> %s", hex.EncodeToString(in.InputHash), hex.EncodeToString(in.OutputHash)))

	err := s.core.Submit(in.InputHash, in.OutputHash)
	if err != nil {
		return nil, err
	}

	return &pb.SubmitResponse{
		Status: pb.SubmitResponse_OK,
	}, nil
}

func (s *TrustixRPCServer) QueryMapping(ctx context.Context, in *pb.QueryRequest) (*pb.QueryResponse, error) {
	fmt.Println(fmt.Sprintf("Received input hash query for %s", hex.EncodeToString(in.InputHash)))

	h, err := s.core.Query(in.InputHash)
	if err != nil {
		return nil, err
	}

	return &pb.QueryResponse{
		OutputHash: h,
	}, nil
}
