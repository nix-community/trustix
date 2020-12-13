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
	"github.com/tweag/trustix/api"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/decider"
	pb "github.com/tweag/trustix/proto"
	"github.com/tweag/trustix/rpc"
	"github.com/tweag/trustix/rpc/auth"
	"github.com/tweag/trustix/schema"
	"github.com/tweag/trustix/signer"
	"github.com/tweag/trustix/sthmanager"
	"github.com/tweag/trustix/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var once sync.Once
var configPath string
var stateDirectory string

var listenAddress string
var dialAddress string

var rootCmd = &cobra.Command{
	Use:   "trustix",
	Short: "Trustix",
	Long:  `Trustix`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if configPath == "" {
			return fmt.Errorf("Missing config flag")
		}

		config, err := config.NewConfigFromFile(configPath)
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

		// Check if any names are non-unique
		seenNames := make(map[string]struct{})
		// The number of authoritive logs, can't exceed 1
		numLogs := 0
		for _, logConfig := range config.Logs {
			_, ok := seenNames[logConfig.Name]
			if ok {
				log.Fatalf("Found non-unique log name: %s", logConfig.Name)
			}
			seenNames[logConfig.Name] = struct{}{}

			if logConfig.Mode == "trustix-log" {
				numLogs += 1
				if numLogs > 1 {
					log.Fatal("More than 1 authoritive logs in the same instance is not supported.")
				}
			}
		}

		var logAPIServer api.TrustixLogAPIServer

		errChan := make(chan error)
		wg := new(sync.WaitGroup)

		go func() {
			wg.Wait()
			close(errChan)
		}()

		sthmgr := sthmanager.NewSTHManager()
		defer sthmgr.Close()
		sthstore, err := storage.NewNativeStorage("STH", stateDirectory)
		if err != nil {
			log.Fatal("Could not create local cache storage")
		}

		logMap := rpc.NewTrustixCombinedRPCServerMap()
		for _, logConfig := range config.Logs {
			logConfig := logConfig
			wg.Add(1)

			mkLog := func() error {
				log.WithFields(log.Fields{
					"name":   logConfig.Name,
					"pubkey": logConfig.Signer.PublicKey,
					"mode":   logConfig.Mode,
				}).Info("Adding log")

				switch logConfig.Mode {

				case "trustix-log":
					log.WithFields(log.Fields{
						"name": logConfig.Name,
						"mode": logConfig.Mode,
					}).Info("Adding authoritive log to gRPC")

					signerConfig := logConfig.Signer

					if signerConfig.Type == "" {
						fmt.Errorf("Missing signer config field 'type'.")
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

					var store storage.TrustixStorage
					switch logConfig.Storage.Type {
					case "native":
						store, err = storage.NewNativeStorage("log-"+logConfig.Name, stateDirectory)
					case "memory":
						store, err = storage.NewMemoryStorage()
					}
					if err != nil {
						return err
					}

					logAPI, err := api.NewKVStoreAPI(store, sig)
					if err != nil {
						return err
					}

					logAPIServer, err = api.NewTrustixAPIServer(logAPI)
					if err != nil {
						return err
					}

					logMap.Add(logConfig.Name, logAPI)
					sthmgr.Add(logConfig.Name, sthmanager.NewDummySTHCache(func() (*schema.STH, error) {
						return logAPI.GetSTH(context.Background(), new(api.STHRequest))
					}))

				case "trustix-follower":
					var verifier signer.TrustixVerifier

					signerConfig := logConfig.Signer
					switch signerConfig.Type {
					case "ed25519":
						verifier, err = signer.NewED25519Verifier(logConfig.Signer.PublicKey)
						if err != nil {
							return err
						}
					default:
						return fmt.Errorf("Verifier type '%s' is not supported.", signerConfig.Type)
					}

					conn, err := createClientConn(logConfig.Transport.GRPC.Remote, verifier.Public())
					if err != nil {
						return err
					}

					c, err := api.NewTrustixAPIGRPCClient(conn)
					if err != nil {
						return err
					}

					logMap.Add(logConfig.Name, c)

					sthCache, err := sthmanager.NewSTHCache(logConfig.Name, sthstore, c, verifier)
					if err != nil {
						return err
					}

					sthmgr.Add(logConfig.Name, sthCache)

				default:
					return fmt.Errorf("Mode '%s' could not be initialised for log name %s", logConfig.Mode, logConfig.Name)

				}

				return nil
			}

			go func() {
				defer wg.Done()

				err := mkLog()
				if err != nil {
					errChan <- fmt.Errorf("Got error in log initialisation for log name '%s': %v", logConfig.Name, err)
				}
			}()
		}

		for err := range errChan {
			if err != nil {
				return err
			}
		}
		wg.Wait()

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
				case "logname":
					decider, err := decider.NewLogNameDecider(deciderConfig.LogName.Name)
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

		errChan = make(chan error)

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
		for _, addr := range []string{listenAddress} {

			if addr == "" {
				continue
			}

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

	rootCmd.PersistentFlags().StringVar(&listenAddress, "listen", "", "Listen to address")

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
