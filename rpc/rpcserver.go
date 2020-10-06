package rpc

import (
	"context"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
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

	log.WithFields(log.Fields{
		"inputHash":  hex.EncodeToString(in.InputHash),
		"outputHash": hex.EncodeToString(in.OutputHash),
	}).Info("Received input hash")

	err := s.core.Submit(in.InputHash, in.OutputHash)
	if err != nil {
		return nil, err
	}

	return &pb.SubmitResponse{
		Status: pb.SubmitResponse_OK,
	}, nil
}

func (s *TrustixRPCServer) QueryMapping(ctx context.Context, in *pb.QueryRequest) (*pb.QueryResponse, error) {

	log.WithFields(log.Fields{
		"inputHash": hex.EncodeToString(in.InputHash),
	}).Info("Received input hash query")

	h, err := s.core.Query(in.InputHash)
	if err != nil {
		return nil, err
	}

	return &pb.QueryResponse{
		OutputHash: h,
	}, nil
}
