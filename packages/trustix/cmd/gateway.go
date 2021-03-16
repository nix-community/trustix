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
