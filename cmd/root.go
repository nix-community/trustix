package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/core"
	pb "github.com/tweag/trustix/proto"
	"github.com/tweag/trustix/rpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"path"
	"sync"
)

var once sync.Once
var configPath string
var stateDirectory string

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

		lis, err := net.Listen("tcp", dialAddress)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()

		err = os.MkdirAll(stateDirectory, 0700)
		if err != nil {
			log.Fatalf("Could not create state directory: %v")
		}

		flagConfig := &core.FlagConfig{
			StateDirectory: stateDirectory,
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
					log.Fatal("More than 1 authoritive logs is not supported.")
				}
			}
		}

		logMap := make(map[string]*core.TrustixCore)
		for _, logConfig := range config.Logs {
			c, err := core.CoreFromConfig(logConfig, flagConfig)
			if err != nil {
				log.Fatal(err)
			}

			logMap[logConfig.Name] = c

			if logConfig.Mode == "trustix-log" {
				// Authoritive APIs
				pb.RegisterTrustixRPCServer(s, rpc.NewTrustixRPCServer(c))
				pb.RegisterTrustixKVServer(s, rpc.NewTrustixKVServer(c))
			}
		}

		pb.RegisterTrustixLogServer(s, rpc.NewTrustixLogServer(logMap))

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		return nil
	},
}

func initCommands() {
	rootCmd.Flags().StringVar(&configPath, "config", "", "Path to config.toml")

	rootCmd.PersistentFlags().StringVar(&dialAddress, "address", ":8080", "Path to config.toml")

	homeDir, _ := os.UserHomeDir()
	defaultStateDir := path.Join(homeDir, ".local/share/trustix")
	rootCmd.PersistentFlags().StringVar(&stateDirectory, "state", defaultStateDir, "State directory")

	rootCmd.AddCommand(generateKeyCmd)
	initGenerate()

	rootCmd.AddCommand(submitCommand)
	initSubmit()

	rootCmd.AddCommand(queryCommand)
	initQuery()

	rootCmd.AddCommand(queryMap)
	initMapQuery()
}

func Execute() {
	once.Do(initCommands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
