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
	"net/http"
	"sync"
	"time"

	"github.com/nix-community/trustix/packages/unixtransport"
	log "github.com/sirupsen/logrus"
)

// Use a shared client connection pool
var initOnce sync.Once
var client *http.Client

func CreateClient(URL string) (*Client, error) {
	initOnce.Do(func() {
		t := &http.Transport{}
		unixtransport.Register(t)
		client = &http.Client{Transport: t}
	})

	log.WithFields(log.Fields{
		"address": URL,
	}).Debug("Creating client for remote")

	return &Client{
		LogAPI:  newLogAPIConnectClient(client, URL),
		RpcAPI:  newRpcAPIConnectClient(client, URL),
		NodeAPI: newNodeAPIConnectClient(client, URL),
		LogRPC:  newLogRPCConnectClient(client, URL),
	}, nil
}

// Create a context with the default timeout set
func CreateContext(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
}
