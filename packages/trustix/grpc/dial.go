// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func Dial(address string) (*grpc.ClientConn, error) {

	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	var conn *grpc.ClientConn

	switch u.Scheme {
	case "grpc+unix", "unix":
		sockPath := u.Host + u.Path

		conn, err = grpc.Dial(
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
		if err != nil {
			return nil, err
		}
	case "grpc+https":
		creds := credentials.NewTLS(&tls.Config{})
		conn, err = grpc.Dial(address, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, err
		}
	case "grpc+http":
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("URL scheme '%s' not supported", u.Scheme)

	}

	return conn, nil
}
