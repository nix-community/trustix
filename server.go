package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/trillian"
	"github.com/google/trillian/merkle/rfc6962"
	"google.golang.org/grpc/codes"
)

type server struct {
	client trillian.TrillianLogClient
	logID  int64
}

func newServer(client trillian.TrillianLogClient, logID int64) *server {
	return &server{
		client: client,
		logID:  logID,
	}
}

func (s *server) put(r *Request) (*Response, error) {

	leaf := &trillian.LogLeaf{
		LeafIdentityHash: r.input.IdentityHash(),
		LeafValue:        r.input.OutputHash(),
	}

	rqst := &trillian.QueueLeafRequest{
		LogId: s.logID,
		Leaf:  leaf,
	}
	resp, err := s.client.QueueLeaf(context.Background(), rqst)
	if err != nil {
		log.Fatal(err)
	}

	c := codes.Code(resp.QueuedLeaf.GetStatus().GetCode())
	if c != codes.OK && c != codes.AlreadyExists {
		return &Response{}, fmt.Errorf("[server:put] Bad status: %v", resp.QueuedLeaf.GetStatus())
	}
	if c == codes.OK {
		log.Println("[server:put] ok")
	} else if c == codes.AlreadyExists {
		return nil, fmt.Errorf("Input hash %s already exists", string(r.input.IdentityHashString()))
	}

	return &Response{
		status: "ok",
	}, nil
}

func (s *server) get(r *Request) (*Response, error) {
	log.Println("[server:get] Entered")

	leafValue := r.input.OutputHash()

	hasher := rfc6962.DefaultHasher
	leafHash := hasher.HashLeaf(leafValue)
	log.Printf("[server:get] hash: %x\n", leafHash)

	req := &trillian.GetLeavesByHashRequest{
		LogId: s.logID,
		// LeafIdentityHash: r.input.IdentityHash(),
		LeafHash: [][]byte{leafHash},
	}

	resp, err := s.client.GetLeavesByHash(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	for i, logLeaf := range resp.GetLeaves() {
		leafValue := logLeaf.GetLeafValue()
		extraData := logLeaf.GetExtraData()
		log.Printf("[server:get] %d: %s", i, extraData)
		log.Printf("[server:get] %d: %s", i, leafValue)
	}

	return &Response{
		status: "ok",
	}, nil
}
