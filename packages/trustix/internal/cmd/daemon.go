// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"crypto"
	"fmt"
	"math"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/activation"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	pb "github.com/tweag/trustix/packages/trustix-proto/rpc"
	"github.com/tweag/trustix/packages/trustix/client"
	tapi "github.com/tweag/trustix/packages/trustix/internal/api"
	conf "github.com/tweag/trustix/packages/trustix/internal/config"
	"github.com/tweag/trustix/packages/trustix/internal/constants"
	"github.com/tweag/trustix/packages/trustix/internal/decider"
	"github.com/tweag/trustix/packages/trustix/internal/lib"
	"github.com/tweag/trustix/packages/trustix/internal/pool"
	"github.com/tweag/trustix/packages/trustix/internal/protocols"
	pub "github.com/tweag/trustix/packages/trustix/internal/publisher"
	"github.com/tweag/trustix/packages/trustix/internal/rpc/auth"
	"github.com/tweag/trustix/packages/trustix/internal/server"
	"github.com/tweag/trustix/packages/trustix/internal/signer"
	"github.com/tweag/trustix/packages/trustix/internal/sthsync"
	"github.com/tweag/trustix/packages/trustix/internal/storage"
	"google.golang.org/grpc"
)

var daemonListenAddresses []string
var daemonConfigPath string
var daemonStateDirectory string
var daemonPollInterval float64

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Trustix daemon",
	RunE: func(cmd *cobra.Command, args []string) error {

		if daemonConfigPath == "" {
			return fmt.Errorf("Missing config flag")
		}

		config, err := conf.NewConfigFromFile(daemonConfigPath)
		if err != nil {
			log.Fatal(err)
		}

		log.WithFields(log.Fields{
			"directory": daemonStateDirectory,
		}).Info("Creating state directory")
		err = os.MkdirAll(daemonStateDirectory, 0700)
		if err != nil {
			log.Fatalf("Could not create state directory: %s", daemonStateDirectory)
		}

		var store storage.Storage
		{
			switch config.Storage.Type {
			case "native":
				store, err = storage.NewNativeStorage(daemonStateDirectory)
			case "memory":
				store, err = storage.NewMemoryStorage()
			}
			if err != nil {
				log.Fatalf("Could not initialise store: %v", err)
			}
		}

		signers := make(map[string]crypto.Signer)
		{
			for name, signerConfig := range config.Signers {
				var sig crypto.Signer

				log.WithFields(log.Fields{
					"type": signerConfig.Type,
					"name": name,
				}).Info("Creating signer")

				switch signerConfig.Type {
				case "ed25519":
					sig, err = signer.NewED25519Signer(signerConfig.ED25519.PrivateKeyPath)
					if err != nil {
						return err
					}
				default:
					return fmt.Errorf("Signer type '%s' is not supported.", signerConfig.Type)
				}

				signers[name] = sig
			}
		}

		// These APIs are static and fully controlled by the configuration file
		logs := []*api.Log{}
		logsPublished := []*api.Log{}
		{
			for _, pubConf := range config.Publishers {
				pd, err := protocols.Get(pubConf.Protocol)
				if err != nil {
					return err
				}

				logMode := api.Log_LogModes(0)

				logID, err := pubConf.PublicKey.LogID(pd, logMode)
				if err != nil {
					log.Fatal(err)
				}

				signer, err := pubConf.PublicKey.Signer()
				if err != nil {
					log.Fatal(err)
				}

				log := &api.Log{
					LogID:    &logID,
					Meta:     pubConf.GetMeta(),
					Signer:   signer,
					Protocol: &pd.ID,
					Mode:     logMode.Enum(), // Hard-coded for now
				}

				logs = append(logs, log)
				logsPublished = append(logsPublished, log)
			}

			for _, subConf := range config.Subscribers {
				pd, err := protocols.Get(subConf.Protocol)
				if err != nil {
					return err
				}

				logMode := api.Log_LogModes(0)

				logID, err := subConf.PublicKey.LogID(pd, logMode)
				if err != nil {
					log.Fatal(err)
				}

				signer, err := subConf.PublicKey.Signer()
				if err != nil {
					log.Fatal(err)
				}

				log := &api.Log{
					LogID:    &logID,
					Meta:     subConf.GetMeta(),
					Signer:   signer,
					Protocol: &pd.ID,
					Mode:     logMode.Enum(),
				}

				logs = append(logs, log)
			}
		}

		clientPool := pool.NewClientPool()
		defer clientPool.Close()

		for _, remote := range config.Remotes {
			remote := remote
			go func() {

				pc, err := clientPool.Dial(remote)
				if err != nil {
					log.WithFields(log.Fields{
						"remote": remote,
					}).Error("Couldn't dial remote")
					return
				}

				pc.Activate()
			}()
		}

		rootBucket := &storage.Bucket{}
		caValueBucket := rootBucket.Cd(constants.CaValueBucket)

		nodeAPI := tapi.NewKVStoreNodeAPI(store, caValueBucket, logsPublished)
		nodeAPIServer := server.NewNodeAPIServer(nodeAPI)

		headSyncCloser := lib.NewMultiCloser()
		defer headSyncCloser.Close()

		pubMap := pub.NewPublisherMap()
		defer pubMap.Close()

		{

			logInitExecutor := lib.NewParallellExecutor()

			for _, subscriberConfig := range config.Subscribers {
				subConf := subscriberConfig
				logInitExecutor.Add(func() error {
					pubBytes, err := subConf.PublicKey.Decode()
					if err != nil {
						return err
					}

					pd, err := protocols.Get(subConf.Protocol)
					if err != nil {
						return err
					}

					logMode := api.Log_LogModes(0)

					logID, err := subConf.PublicKey.LogID(pd, logMode)
					if err != nil {
						return err
					}

					log.WithFields(log.Fields{
						"id":     logID,
						"pubkey": subConf.PublicKey.Pub,
					}).Info("Adding log subscriber")

					var verifier signer.Verifier
					{
						switch subConf.PublicKey.Type {
						case "ed25519":
							verifier, err = signer.NewED25519Verifier(pubBytes)
							if err != nil {
								return err
							}
						default:
							return fmt.Errorf("Verifier type '%s' is not supported.", subConf.PublicKey.Type)
						}
					}

					pollDuration := time.Millisecond * time.Duration(math.Round(daemonPollInterval/1000))
					headSyncCloser.Add(sthsync.NewSTHSyncer(logID, store, rootBucket.Cd(logID), clientPool, verifier, pollDuration, pd))

					return nil
				})

			}

			for i, publisherConfig := range config.Publishers {
				i := i
				publisherConfig := publisherConfig
				logInitExecutor.Add(func() error {
					logID := *logsPublished[i].LogID

					log.WithFields(log.Fields{
						"id":     logID,
						"pubkey": publisherConfig.PublicKey.Pub,
					}).Info("Adding log")

					pd, err := protocols.Get(publisherConfig.Protocol)
					if err != nil {
						return err
					}

					logAPI, err := tapi.NewKVStoreLogAPI(logID, store, rootBucket.Cd(logID), pd)
					if err != nil {
						return err
					}

					sig, ok := signers[publisherConfig.Signer]
					if !ok {
						return fmt.Errorf("Missing signer '%s'", publisherConfig.Signer)
					}

					publisher, err := pub.NewPublisher(logID, store, caValueBucket, rootBucket.Cd(logID), sig, pd)
					if err != nil {
						return err
					}

					if err = pubMap.Set(logID, publisher); err != nil {
						return err
					}

					pc, err := clientPool.Add(&client.Client{
						NodeAPI: nodeAPI,
						LogAPI:  logAPI,
						Type:    client.LocalClientType,
					}, []string{logID})
					if err != nil {
						return err
					}
					pc.Activate()

					return nil
				})

			}

			err = logInitExecutor.Wait()
			if err != nil {
				return err
			}

		}

		logAPIServer := server.NewLogAPIServer(logsPublished, clientPool)

		deciders := make(map[string]decider.LogDecider)
		{
			for protocol, deciderConfigs := range config.Deciders {
				current := []decider.LogDecider{}
				for _, deciderConfig := range deciderConfigs {
					switch deciderConfig.Engine {
					case "javascript":
						decider, err := decider.NewJSDecider(deciderConfig.JS.Function)
						if err != nil {
							return err
						}
						current = append(current, decider)
					case "percentage":
						decider, err := decider.NewMinimumPercentDecider(deciderConfig.Percentage.Minimum)
						if err != nil {
							return err
						}
						current = append(current, decider)
					case "logid":
						decider, err := decider.NewLogIDDecider(deciderConfig.LogID.ID)
						if err != nil {
							return err
						}
						current = append(current, decider)
					default:
						return fmt.Errorf("No such engine: %s", deciderConfig.Engine)
					}
				}

				pd, err := protocols.Get(protocol)
				if err != nil {
					return err
				}

				deciders[pd.ID] = decider.NewAggDecider(current...)
			}
		}

		logRpcServer := server.NewLogRPCServer(store, rootBucket, clientPool, pubMap)
		rpcServer := server.NewRPCServer(store, rootBucket, clientPool, pubMap, logs, deciders)

		log.Debug("Creating gRPC servers")

		errChan := make(chan error)

		createServer := func(lis net.Listener) (s *grpc.Server) {
			_, isUnix := lis.(*net.UnixListener)

			if isUnix {
				s = grpc.NewServer(
					grpc.Creds(&auth.SoPeercred{}), // Attach SO_PEERCRED auth to UNIX sockets
				)

				pb.RegisterLogRPCServer(s, logRpcServer)
				pb.RegisterRPCApiServer(s, rpcServer)

			} else {
				s = grpc.NewServer()
			}

			api.RegisterLogAPIServer(s, logAPIServer)
			api.RegisterNodeAPIServer(s, nodeAPIServer)

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

			servers = append(servers, createServer(lis))
		}

		// Create sockets
		for _, addr := range daemonListenAddresses {
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

func initDaemon() {

	homeDir, _ := os.UserHomeDir()
	defaultStateDir := path.Join(homeDir, ".local/share/trustix")
	daemonCmd.Flags().StringVar(&daemonStateDirectory, "state", defaultStateDir, "State directory")

	// Default poll every 30 minutes
	daemonCmd.Flags().Float64Var(&daemonPollInterval, "interval", 60*30, "Log poll interval in seconds")

	daemonCmd.Flags().StringSliceVar(&daemonListenAddresses, "listen", []string{}, "Listen to address")

	daemonCmd.Flags().StringVar(&daemonConfigPath, "config", "", "Path to config.toml")
}
