// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"net"
	"net/http"

	"github.com/coreos/go-systemd/activation"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix/client"
)

var gatewayCommand = &cobra.Command{
	Use:   "gateway",
	Short: "Trustix gateway translating REST calls to gRPC",
	RunE: func(cmd *cobra.Command, args []string) error {

		conn, err := client.CreateClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := client.CreateContext(timeout)
		defer cancel()

		mux := runtime.NewServeMux()
		err = api.RegisterTrustixLogAPIHandler(ctx, mux, conn)
		if err != nil {
			return err
		}

		var listeners []net.Listener

		{
			systemdListeners, err := activation.Listeners()
			if err != nil {
				panic(err)
			}

			for _, lis := range systemdListeners {
				log.WithFields(log.Fields{
					"address": lis.Addr(),
				}).Info("Using socket activated listener")

				listeners = append(listeners, lis)
			}
		}

		for _, addr := range listenAddresses {
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}

			log.WithFields(log.Fields{
				"address": addr,
			}).Info("Listening to address")

			listeners = append(listeners, lis)
		}

		if len(listeners) == 0 {
			log.Fatal("No listeners configured")
		}

		errChan := make(chan error)
		for _, listener := range listeners {
			go func(l net.Listener) {
				err := http.Serve(l, mux)
				if err != nil {
					errChan <- err
				}
			}(listener)
		}
		for err := range errChan {
			log.Fatalf("Error in HTTP handler: %v", err)
			panic(err)
		}

		return nil
	},
}
