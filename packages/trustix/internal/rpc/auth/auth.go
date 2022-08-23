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

	"google.golang.org/grpc/peer"
)

type AuthInfo struct{}

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
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return fmt.Errorf("Could not get peer from context")
	}

	if pr.Addr.Network() != "unix" {
		return fmt.Errorf("Write only allowed over UNIX socket")
	}

	return fmt.Errorf("Denied peer creds")
}
