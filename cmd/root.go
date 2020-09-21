package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/core"
	pb "github.com/tweag/trustix/proto"
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

type pbServer struct {
	pb.UnimplementedTrustixRPCServer
	core *core.TrustixCore
}

func (s *pbServer) SubmitMapping(ctx context.Context, in *pb.SubmitRequest) (*pb.SubmitResponse, error) {
	fmt.Println(fmt.Sprintf("Received input hash %s -> %s", hex.EncodeToString(in.InputHash), hex.EncodeToString(in.OutputHash)))

	err := s.core.Submit(in.InputHash, in.OutputHash)
	if err != nil {
		return nil, err
	}

	return &pb.SubmitResponse{
		Status: pb.SubmitResponse_OK,
	}, nil
}

func (s *pbServer) QueryMapping(ctx context.Context, in *pb.QueryRequest) (*pb.QueryResponse, error) {
	fmt.Println(fmt.Sprintf("Received input hash query for %s", hex.EncodeToString(in.InputHash)))

	h, err := s.core.Query(in.InputHash)
	if err != nil {
		return nil, err
	}

	return &pb.QueryResponse{
		OutputHash: h,
	}, nil
}

type kvServer struct {
	pb.UnimplementedTrustixKVServer
	core *core.TrustixCore
}

func (s *kvServer) GetKey(ctx context.Context, in *pb.KVRequest) (*pb.KVResponse, error) {
	fmt.Println(fmt.Sprintf("Received KV request for %s", hex.EncodeToString(in.Key)))

	v, err := s.core.Get(in.Key)
	if err != nil {
		return nil, err
	}

	return &pb.KVResponse{
		Value: v,
	}, nil
}

type logServer struct {
	pb.UnimplementedTrustixLogServer
	m map[string]*core.TrustixCore
}

func (l *logServer) HashMap(ctx context.Context, in *pb.HashRequest) (*pb.HashMapResponse, error) {

	responses := make(map[string][]byte)

	for name, l := range l.m {
		h, err := l.Query(in.InputHash)
		if err != nil {
			continue
		}
		responses[name] = h
	}

	return &pb.HashMapResponse{
		Hashes: responses,
	}, nil

}

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
				pb.RegisterTrustixRPCServer(s, &pbServer{
					core: c,
				})
				pb.RegisterTrustixKVServer(s, &kvServer{
					core: c,
				})
			}
		}

		pb.RegisterTrustixLogServer(s, &logServer{
			m: logMap,
		})

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
