package transport

import (
	"context"
	"fmt"
	"github.com/tweag/trustix/config"
	pb "github.com/tweag/trustix/proto"
	"github.com/tweag/trustix/storage"
	"google.golang.org/grpc"
	"time"
)

type grpcTransport struct {
	c pb.TrustixKVClient
}

type grpcTxn struct {
	c pb.TrustixKVClient
}

func (t *grpcTxn) Get(key []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := t.c.Get(ctx, &pb.KVRequest{
		Key: key,
	})
	if err != nil {
		return nil, err
	}

	return r.Value, nil
}

func (t *grpcTxn) Set(key []byte, value []byte) error {
	return fmt.Errorf("Cannot Set on remote")
}

func NewGRPCTransport(t *config.GRPCTransportConfig) (*grpcTransport, error) {
	conn, err := grpc.Dial(t.Remote, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	c := pb.NewTrustixKVClient(conn)

	return &grpcTransport{
		c: c,
	}, nil
}

func (g *grpcTransport) Close() {
}

func (s *grpcTransport) runTX(fn func(storage.Transaction) error) error {
	t := &grpcTxn{
		c: s.c,
	}

	err := fn(t)
	if err != nil {
		return err
	}

	return err
}

func (s *grpcTransport) View(fn func(storage.Transaction) error) error {
	return s.runTX(fn)
}

func (s *grpcTransport) Update(fn func(storage.Transaction) error) error {
	return s.runTX(fn)
}
