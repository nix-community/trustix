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
