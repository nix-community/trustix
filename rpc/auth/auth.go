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
	"context"
	"fmt"
	"google.golang.org/grpc/peer"
	"os/user"
	"strconv"
	"syscall"
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
