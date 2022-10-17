// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bakins/logrus-middleware"
	"github.com/coreos/go-systemd/activation"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/config"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/cron"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/index"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/server"
	apiconnect "github.com/nix-community/trustix/packages/trustix-nix-r13y/reprod-api/reprod_apiconnect"
	tclient "github.com/nix-community/trustix/packages/trustix/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var serveListenAddresses []string
var serveConfig string

var serveCommand = &cobra.Command{
	Use:   "serve",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		if serveConfig == "" {
			panic("Missing config path parameter")
		}

		conf, err := config.NewConfigFromFile(serveConfig)
		if err != nil {
			panic(err)
		}

		// config options
		logIndexCronInterval := time.Second * time.Duration(conf.Cron.LogInterval)
		evalIndexCronInterval := time.Second * time.Duration(conf.Cron.EvalInterval)

		dbs, err := setupDatabases(stateDirectory)
		if err != nil {
			panic(fmt.Errorf("error opening database: %w", err))
		}

		client, err := tclient.CreateClient(dialAddress)
		if err != nil {
			panic(err)
		}

		os.Unsetenv("NIX_PATH") // Prevents eval from inheriting NIX_PATH

		// Start indexing logs
		{
			log.WithFields(log.Fields{
				"interval": logIndexCronInterval,
			}).Info("Starting log index cron")

			logIndexCron := cron.NewSingletonCronJob("log_index", logIndexCronInterval, func(ctx context.Context) {
				log.Info("Triggering log index cron job")

				err = index.IndexLogs(ctx, dbs.dbRW, client)
				if err != nil {
					panic(err)
				}

				log.Info("Done executing log index cron job")
			})
			defer logIndexCron.Close()
		}

		// Start indexing evaluations
		{
			log.WithFields(log.Fields{
				"interval": evalIndexCronInterval,
			}).Info("Starting evaluation index cron")

			evalIndexCron := cron.NewSingletonCronJob("eval_index", evalIndexCronInterval, func(ctx context.Context) {
				indexedEvals := 0

				for channel, channelConfig := range conf.Channels.Hydra {
					l := log.WithFields(log.Fields{
						"channel": channel,
					})

					l.Info("indexing channel")

					n, err := index.IndexHydraJobset(ctx, dbs.dbRW, channel, channelConfig)
					if err != nil {
						l.WithFields(log.Fields{
							"erroro": err,
						}).Error("error indexing channel")
					}

					indexedEvals += n
				}

				log.WithFields(log.Fields{
					"num_evals": indexedEvals,
				}).Info("done executing evaluation index cron job")
			})
			defer evalIndexCron.Close()
		}

		apiServer := server.NewAPIServer(dbs.dbRO, dbs.cacheDbRW, dbs.cacheDbRO, client, conf.Lognames, conf.Attrs)

		errChan := make(chan error)

		createServer := func(lis net.Listener) *http.Server {
			mux := http.NewServeMux()

			mux.Handle(apiconnect.NewReproducibilityAPIHandler(apiServer))

			l := logrusmiddleware.Middleware{
				Name:   "trustix-nix-r13y",
				Logger: log.New(),
			}

			loggedHandler := l.Handler(h2c.NewHandler(mux, &http2.Server{}), "/")

			server := &http.Server{Handler: loggedHandler}

			go func() {
				err := server.Serve(lis)
				if err != nil {
					errChan <- fmt.Errorf("failed to serve: %v", err)
				}
			}()

			return server
		}

		var servers []*http.Server

		// Systemd socket activation
		listeners, err := activation.Listeners()
		if err != nil {
			log.Fatalf("cannot retrieve listeners: %s", err)
		}
		for _, lis := range listeners {
			log.WithFields(log.Fields{
				"address": lis.Addr(),
			}).Info("Using socket activated listener")

			servers = append(servers, createServer(lis))
		}

		for _, addr := range serveListenAddresses {
			u, err := url.Parse(addr)
			if err != nil {
				log.Fatalf("Could not parse url: %v", err)
			}

			family := ""
			host := ""

			switch u.Scheme {
			case "unix":
				family = "unix"
				host = u.Host + u.Path
			case "http":
				family = "tcp"
				host = u.Host
			default:
				log.Fatalf("Socket with scheme '%s' unsupported", u.Scheme)
			}

			lis, err := net.Listen(family, host)
			if err != nil {
				log.Fatalf("failed to listen: %v", err)
			}

			log.WithFields(log.Fields{
				"address": addr,
			}).Info("Listening to address")

			servers = append(servers, createServer(lis))
		}

		if len(servers) <= 0 {
			log.Fatal("No listeners configured!")
		}

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-quit
			wg := new(sync.WaitGroup)

			log.Info("Received shutdown signal, closing down server gracefully")

			for _, server := range servers {
				server := server
				wg.Add(1)
				go func() {
					defer wg.Done()
					server.Close()
				}()
			}

			wg.Wait()

			log.Info("Done closing down servers")

			close(errChan)
		}()

		for err := range errChan {
			log.Fatal(err)
		}
	},
}

func initServe() {
	serveCommand.Flags().StringSliceVar(&serveListenAddresses, "listen", []string{}, "Listen to address")
	serveCommand.Flags().StringVar(&serveConfig, "config", "", "Path to config.toml/json")
}
