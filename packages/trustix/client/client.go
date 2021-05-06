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
	"crypto"
	"crypto/ed25519"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func CreateClientConn(address string, pubKey crypto.PublicKey) (*grpc.ClientConn, error) {

	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "unix" {
		return nil, fmt.Errorf("Only UNIX sockets are supported in the CLI for now")
	}

	log.WithFields(log.Fields{
		"address": address,
	}).Debug("Dialing gRPC")

	var conn *grpc.ClientConn

	switch u.Scheme {
	case "https":
		config := &tls.Config{
			InsecureSkipVerify: true,
			VerifyConnection: func(state tls.ConnectionState) error {
				if len(state.PeerCertificates) != 1 {
					return fmt.Errorf("Dont know how to handle %d certs", len(state.PeerCertificates))
				}

				cert := state.PeerCertificates[0]

				edPub, ok := cert.PublicKey.(ed25519.PublicKey)
				if !ok {
					return fmt.Errorf("Key not ed25519")
				}

				if !edPub.Equal(pubKey) {
					return fmt.Errorf("Expected key mismatch")
				}

				err := cert.CheckSignature(cert.SignatureAlgorithm, cert.RawTBSCertificate, cert.Signature)
				if err != nil {
					fmt.Errorf("Signature check failed %d certs", len(state.PeerCertificates))
					return err
				}

				return nil
			},
		}

		creds := credentials.NewTLS(config)

		conn, err = grpc.Dial(address, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, err
		}

	case "unix":
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
	default:
		return nil, fmt.Errorf("URL scheme '%s' not supported", u.Scheme)

	}

	return conn, nil
}

// Create a context with the default timeout set
func CreateContext(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
}
