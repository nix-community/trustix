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

package cmd

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/url"
	"time"
)

func createClientConn() (*grpc.ClientConn, error) {

	u, err := url.Parse(dialAddress)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "unix" {
		return nil, fmt.Errorf("Only UNIX sockets are supported in the CLI for now")
	}

	sockPath := u.Host + u.Path

	log.WithFields(log.Fields{
		"address": dialAddress,
	}).Debug("Dialing gRPC")

	return grpc.Dial(
		sockPath,
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			unix_addr, err := net.ResolveUnixAddr("unix", addr)
			if err != nil {
				return nil, err
			}
			return net.DialUnix("unix", nil, unix_addr)
		}),
	)
}

// Create a context with the default timeout set
func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}
