// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"context"
	"crypto"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/coreos/go-systemd/activation"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	pb "github.com/tweag/trustix/packages/trustix-proto/proto"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	tapi "github.com/tweag/trustix/packages/trustix/api"
	"github.com/tweag/trustix/packages/trustix/client"
	conf "github.com/tweag/trustix/packages/trustix/config"
	"github.com/tweag/trustix/packages/trustix/decider"
	"github.com/tweag/trustix/packages/trustix/lib"
	"github.com/tweag/trustix/packages/trustix/rpc"
	"github.com/tweag/trustix/packages/trustix/rpc/auth"
	"github.com/tweag/trustix/packages/trustix/signer"
	"github.com/tweag/trustix/packages/trustix/sthmanager"
	"github.com/tweag/trustix/packages/trustix/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var once sync.Once
var configPath string
var stateDirectory string

var listenAddresses []string
var dialAddress string

var logID string

var timeout int

var rootCmd = &cobra.Command{
	Use:   "trustix",
	Short: "Trustix",
	Long:  `Trustix`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if configPath == "" {
			return fmt.Errorf("Missing config flag")
		}

		config, err := conf.NewConfigFromFile(configPath)
		if err != nil {
			log.Fatal(err)
		}

		log.WithFields(log.Fields{
			"directory": stateDirectory,
		}).Info("Creating state directory")
		err = os.MkdirAll(stateDirectory, 0700)
		if err != nil {
			log.Fatalf("Could not create state directory: %s", stateDirectory)
		}

		var store storage.TrustixStorage
		{
			switch config.Storage.Type {
			case "native":
				store, err = storage.NewNativeStorage(stateDirectory)
			case "memory":
				store, err = storage.NewMemoryStorage()
			}
			if err != nil {
				log.Fatalf("Could not initialise store: %v", err)
			}
		}

		sthmgr := sthmanager.NewSTHManager()
		defer sthmgr.Close()
		sthstore := store // TODO: Create a more narrow Store view

		logMap := tapi.NewTrustixLogMap()
		{

			errChan := make(chan error)
			wg := new(sync.WaitGroup)

			go func() {
				wg.Wait()
				close(errChan)
			}()

			// Wrap to handle errors and propagate to channel
			wrapLogInit := func(fn func() error) {
				err := fn()
				if err != nil {
					errChan <- err
				}
			}

			for _, subscriberConfig := range config.Subscribers {
				subConf := subscriberConfig
				wg.Add(1)

				go wrapLogInit(func() error {
					defer wg.Done()

					pubBytes, err := signer.Decode(subConf.Signer.PublicKey)
					if err != nil {
						return err
					}

					logID := lib.LogID(subConf.Signer.Type, pubBytes)

					log.WithFields(log.Fields{
						"id":     logID,
						"pubkey": subConf.Signer.PublicKey,
						// "mode":   subConf.Mode,
					}).Info("Adding log subscriber")

					var verifier signer.TrustixVerifier
					{
						switch subConf.Signer.Type {
						case "ed25519":
							verifier, err = signer.NewED25519Verifier(pubBytes)
							if err != nil {
								return err
							}
						default:
							return fmt.Errorf("Verifier type '%s' is not supported.", subConf.Signer.Type)
						}
					}

					conn, err := client.CreateClientConn(subConf.Transport.GRPC.Remote, verifier.Public())
					if err != nil {
						return err
					}

					c, err := tapi.NewTrustixAPIGRPCClient(conn)
					if err != nil {
						return err
					}

					logMap.Add(logID, c)

					sthCache, err := sthmanager.NewSTHCache(logID, sthstore, c, verifier)
					if err != nil {
						return err
					}

					sthmgr.Set(logID, sthCache)

					return nil
				})

			}

			for _, logConfig := range config.Logs {
				logConfig := logConfig
				wg.Add(1)

				go wrapLogInit(func() error {
					defer wg.Done()

					var logID string
					{
						pubBytes, err := signer.Decode(logConfig.Signer.PublicKey)
						if err != nil {
							return err
						}

						logID = lib.LogID(logConfig.Signer.Type, pubBytes)
					}

					log.WithFields(log.Fields{
						"id":     logID,
						"pubkey": logConfig.Signer.PublicKey,
						"mode":   logConfig.Mode,
					}).Info("Adding log")

					// TODO: Define a logger with fields already applied

					log.WithFields(log.Fields{
						"id":   logID,
						"mode": logConfig.Mode,
					}).Info("Adding authoritive log to gRPC")

					signerConfig := logConfig.Signer

					if signerConfig.Type == "" {
						return fmt.Errorf("Missing signer config field 'type'.")
					}

					var sig crypto.Signer

					log.WithField("type", signerConfig.Type).Info("Creating signer")
					switch signerConfig.Type {
					case "ed25519":
						sig, err = signer.NewED25519Signer(signerConfig.ED25519.PrivateKeyPath)
						if err != nil {
							return err
						}
					default:
						return fmt.Errorf("Signer type '%s' is not supported.", signerConfig.Type)
					}

					logAPI, err := tapi.NewKVStoreAPI(store, sig)
					if err != nil {
						return err
					}

					logMap.Add(logID, logAPI)
					// TODO: Allow to fail on startup
					sthmgr.Set(logID, sthmanager.NewDummySTHCache(func() (*schema.STH, error) {
						return logAPI.GetSTH(context.Background(), &api.STHRequest{
							LogID: &logID,
						})
					}))

					return nil
				})

			}

			for err := range errChan {
				if err != nil {
					return err
				}
			}
			wg.Wait()

		}

		logAPIServer, err := tapi.NewTrustixAPIServer(logMap, store)
		if err != nil {
			return err
		}

		decider, err := func() (decider.LogDecider, error) {
			deciders := []decider.LogDecider{}
			for _, deciderConfig := range config.Deciders {
				switch deciderConfig.Engine {
				case "lua":
					decider, err := decider.NewLuaDecider(deciderConfig.Lua.Script)
					if err != nil {
						return nil, err
					}
					deciders = append(deciders, decider)
				case "percentage":
					decider, err := decider.NewMinimumPercentDecider(deciderConfig.Percentage.Minimum)
					if err != nil {
						return nil, err
					}
					deciders = append(deciders, decider)
				case "logid":
					decider, err := decider.NewLogIDDecider(deciderConfig.LogID.ID)
					if err != nil {
						return nil, err
					}
					deciders = append(deciders, decider)
				default:
					return nil, fmt.Errorf("No such engine: %s", deciderConfig.Engine)
				}
			}
			return decider.NewAggDecider(deciders...), nil
		}()
		if err != nil {
			return fmt.Errorf("Error creating decision engine: %v", err)
		}

		logServer := rpc.NewTrustixCombinedRPCServer(sthmgr, logMap, decider)

		log.Debug("Creating gRPC servers")

		errChan := make(chan error)

		createServer := func(lis net.Listener, insecure bool) (s *grpc.Server) {
			_, isUnix := lis.(*net.UnixListener)

			if isUnix {
				s = grpc.NewServer(
					grpc.Creds(&auth.SoPeercred{}), // Attach SO_PEERCRED auth to UNIX sockets
				)

				pb.RegisterTrustixCombinedRPCServer(s, logServer)

			} else {

				if insecure {
					s = grpc.NewServer()
				} else {
					cert, err := generateCert()
					if err != nil {
						log.Fatalf("Could not create cert")
					}

					config := &tls.Config{
						Certificates: []tls.Certificate{*cert},
						ClientAuth:   tls.NoClientCert,
					}

					s = grpc.NewServer(grpc.Creds(credentials.NewTLS(config)))
				}

			}

			if logAPIServer != nil {
				api.RegisterTrustixLogAPIServer(s, logAPIServer)
			}

			go func() {
				err := s.Serve(lis)
				if err != nil {
					errChan <- fmt.Errorf("failed to serve: %v", err)
				}
			}()

			return s
		}

		var servers []*grpc.Server

		// Systemd socket activation
		listeners, err := activation.Listeners()
		if err != nil {
			log.Fatalf("cannot retrieve listeners: %s", err)
		}
		for _, lis := range listeners {
			log.WithFields(log.Fields{
				"address": lis.Addr(),
			}).Info("Using socket activated listener")

			servers = append(servers, createServer(lis, false))
		}

		// Create sockets
		for _, addr := range listenAddresses {
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
			case "https":
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

			servers = append(servers, createServer(lis, u.Scheme == "http"))
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
				wg.Add(1)
				go func() {
					defer wg.Done()
					server.GracefulStop()
				}()
			}

			wg.Wait()

			log.Info("Done closing down servers")

			close(errChan)
		}()

		for err := range errChan {
			log.Fatal(err)
		}

		return nil
	},
}

func initCommands() {
	rootCmd.Flags().StringVar(&configPath, "config", "", "Path to config.toml")

	rootCmd.PersistentFlags().StringSliceVar(&listenAddresses, "listen", []string{}, "Listen to address")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 20, "Timeout in seconds")

	rootCmd.PersistentFlags().StringVar(&logID, "log-id", "", "Log ID")

	trustixSock := os.Getenv("TRUSTIX_SOCK")
	if trustixSock == "" {
		tmpDir := "/tmp"
		trustixSock = filepath.Join(tmpDir, "trustix.sock")
	}
	trustixSock = fmt.Sprintf("unix://%s", trustixSock)

	rootCmd.PersistentFlags().StringVar(&dialAddress, "address", trustixSock, "Connect to address")

	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stderr)

	homeDir, _ := os.UserHomeDir()
	defaultStateDir := path.Join(homeDir, ".local/share/trustix")
	rootCmd.PersistentFlags().StringVar(&stateDirectory, "state", defaultStateDir, "State directory")

	rootCmd.AddCommand(generateKeyCmd)
	initGenerate()

	rootCmd.AddCommand(submitCommand)
	initSubmit()

	rootCmd.AddCommand(queryCommand)
	initQuery()

	rootCmd.AddCommand(getValueCommand)
	initGetValue()

	rootCmd.AddCommand(queryMap)
	initMapQuery()

	rootCmd.AddCommand(decideCommand)
	initDecide()

	rootCmd.AddCommand(flushCommand)

	rootCmd.AddCommand(exportCommand)
	initExport()

	rootCmd.AddCommand(gatewayCommand)
}

func Execute() {
	once.Do(initCommands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
