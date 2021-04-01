// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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
