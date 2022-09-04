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

	connect "github.com/bufbuild/connect-go"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/api/apiconnect"
)

type nodeAPIConnectClient struct {
	client apiconnect.NodeAPIClient
}

func newNodeAPIConnectClient(client *http.Client, baseURL string) *nodeAPIConnectClient {
	c := apiconnect.NewNodeAPIClient(client, baseURL)
	return &nodeAPIConnectClient{
		client: c,
	}
}

func (c *nodeAPIConnectClient) GetValue(ctx context.Context, in *api.ValueRequest) (*api.ValueResponse, error) {
	resp, err := c.client.GetValue(ctx, connect.NewRequest(in))
	if err != nil {
		return nil, err
	}

	return resp.Msg, err
}

func (c *nodeAPIConnectClient) Logs(ctx context.Context, in *api.LogsRequest) (*api.LogsResponse, error) {
	resp, err := c.client.Logs(ctx, connect.NewRequest(in))
	if err != nil {
		return nil, err
	}

	return resp.Msg, err
}
