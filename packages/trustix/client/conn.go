// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package client

import (
	"context"
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	tgrpc "github.com/tweag/trustix/packages/trustix/internal/grpc"
	"google.golang.org/grpc"
)

func CreateClientConn(address string) (*Client, error) {

	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"address": address,
	}).Debug("Dialing remote")

	var conn *grpc.ClientConn

	switch u.Scheme {
	case "grpc+unix", "unix", "grpc+https", "grpc+http":
		conn, err = tgrpc.Dial(address)
		if err != nil {
			return nil, fmt.Errorf("Error dialing grpc: %w", err)
		}
	default:
		return nil, fmt.Errorf("URL '%s' with scheme '%s' not supported", address, u.Scheme)

	}

	client := &Client{
		LogAPI:  newLogAPIGRPCClient(conn),
		RpcAPI:  newRpcAPIGRPCClient(conn),
		NodeAPI: newNodeAPIGRPCClient(conn),
		LogRPC:  newLogRPCGRPCClient(conn),

		close: conn.Close,
	}

	return client, nil
}

// Create a context with the default timeout set
func CreateContext(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
}
