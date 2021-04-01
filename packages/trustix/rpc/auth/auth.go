// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package auth

import (
	"context"
	"fmt"
	"os/user"
	"strconv"
	"syscall"

	"google.golang.org/grpc/peer"
)

type AuthInfo struct {
	Ucred syscall.Ucred
}

func (AuthInfo) AuthType() string {
	return "ucred"
}

func AuthInfoFromContext(ctx context.Context) (*AuthInfo, bool) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return nil, false
	}
	info, ok := pr.AuthInfo.(AuthInfo)
	if !ok {
		return nil, false
	}
	return &info, true
}

func CanWrite(ctx context.Context) error {

	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %v", err)
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return fmt.Errorf("failed to get current user uid: %v", err)
	}

	pr, ok := peer.FromContext(ctx)
	if !ok {
		return fmt.Errorf("Could not get peer from context")
	}

	if pr.Addr.Network() != "unix" {
		return fmt.Errorf("Write only allowed over UNIX socket")
	}

	info, ok := pr.AuthInfo.(AuthInfo)
	if !ok {
		return fmt.Errorf("Could not get peer creds for socket")
	}

	// Deny connection from other than root and self
	if info.Ucred.Uid == 0 || info.Ucred.Uid == uint32(uid) {
		return nil
	}

	return fmt.Errorf("Denied peer creds")
}
