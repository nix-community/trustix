package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/config"
	"github.com/tweag/trustix/core"
	pb "github.com/tweag/trustix/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"
)

var once sync.Once
var configPath string

var dialAddress string

type pbServer struct {
	pb.UnimplementedTrustixServer
	core *core.TrustixCore
}

func (s *pbServer) SubmitMapping(ctx context.Context, in *pb.SubmitRequest) (*pb.SubmitReply, error) {
	fmt.Println(fmt.Sprintf("Received input hash %s", in.InputHash))

	err := s.core.Submit(in.InputHash, in.OutputHash)
	if err != nil {
		return nil, err
	}

	return &pb.SubmitReply{
		Status: pb.SubmitReply_OK,
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

		for _, logConfig := range config.Logs {
			c, err := core.CoreFromConfig(logConfig)
			if err != nil {
				log.Fatal(err)
			}

			pb.RegisterTrustixServer(s, &pbServer{
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

	rootCmd.AddCommand(generateKeyCmd)
	initGenerate()

	rootCmd.AddCommand(submitCommand)
	initSubmit()
}

func Execute() {
	once.Do(initCommands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
