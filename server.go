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
	log.Println("[server] Creating")
	return &server{
		client: client,
		logID:  logID,
	}
}

func (s *server) put(r *Request) (*Response, error) {
	log.Println("[server:put] Entered")

	// Marshal a Thing (actually just 'name' which is a string) into []byte
	// Eventually we'll marshal a more interesting data structure
	leafValue, err := r.thing.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	// Marshal an Extra (again)
	extraData, err := r.extra.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	leaf := &trillian.LogLeaf{
		LeafValue: leafValue,
		ExtraData: extraData,
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
		log.Printf("[server:put] %s already Exists", leafValue)
	}

	return &Response{
		status: "ok",
	}, nil
}

func (s *server) get(r *Request) (*Response, error) {
	log.Println("[server:get] Entered")

	// Marshal a Thing (actually just 'name' which is a string) into []byte
	// Eventually we'll marshal a more interesting data structure
	leafValue, err := r.thing.Marshal()
	if err != nil {
		log.Fatal(err)
	}

	// Trillian uses its own (rfc6962) hasher
	hasher := rfc6962.DefaultHasher
	leafHash := hasher.HashLeaf(leafValue)
	// Output the hashed value (conventionally hex is used)
	log.Printf("[server:get] hash: %x\n", leafHash)

	// Create the request
	rqst := &trillian.GetLeavesByHashRequest{
		LogId:    s.logID,
		LeafHash: [][]byte{leafHash},
	}

	// Submit the request to the Trillian Log Server
	resp, err := s.client.GetLeavesByHash(context.Background(), rqst)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over the responses; there should be 0 or 1
	for i, logLeaf := range resp.GetLeaves() {
		leafValue := logLeaf.GetLeafValue()
		log.Printf("[server:get] %d: %s", i, leafValue)
	}

	return &Response{
		status: "ok",
	}, nil
}
