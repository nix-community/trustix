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

		for _, logConfig := range config.Logs {
			c, err := core.CoreFromConfig(logConfig, flagConfig)
			if err != nil {
				log.Fatal(err)
			}

			pb.RegisterTrustixRPCServer(s, &pbServer{
				core: c,
			})
		}

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
}

func Execute() {
	once.Do(initCommands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
