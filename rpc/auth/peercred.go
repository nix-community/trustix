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
package auth

import (
	"errors"
	"fmt"
	"net"
	"syscall"

	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
)

type SoPeercred struct{}

func (creds *SoPeercred) ClientHandshake(context.Context, string, net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return nil, nil, fmt.Errorf("Not implemented")
}

func (creds *SoPeercred) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {

	log.Debug("Checking peer credential for socket")

	if conn.LocalAddr().Network() != "unix" {
		return conn, nil, nil
	}

	errLogger := log.WithFields(log.Fields{
		"socket":  conn.LocalAddr().String(),
		"network": conn.LocalAddr().Network(),
	})

	uconn, ok := conn.(*net.UnixConn)
	if !ok {
		errLogger.Error("Not a UNIX socket")
		return nil, nil, errors.New("Not a UNIX socket")
	}

	file, err := uconn.File()
	if err != nil {
		errLogger.Error("Could not get UNIX file")
		return nil, nil, err
	}
	defer file.Close()

	cred, err := syscall.GetsockoptUcred(
		int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_PEERCRED,
	)
	if err != nil {
		errLogger.Error("Failed to get SO_PEERCRED")
		return nil, nil, err
	}

	return conn, AuthInfo{Ucred: *cred}, nil
}

func (creds *SoPeercred) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{
		SecurityProtocol: "peercred",
		ServerName:       "localhost",
	}
}

func (creds *SoPeercred) Clone() credentials.TransportCredentials {
	return creds
}

func (creds *SoPeercred) OverrideServerName(string) error {
	return fmt.Errorf("Not implemented")
}
