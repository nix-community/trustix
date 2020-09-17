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
var listen string

type pbServer struct {
	pb.UnimplementedTrustixServer
}

func (s *pbServer) SubmitMapping(ctx context.Context, in *pb.SubmitRequest) (*pb.SubmitReply, error) {

	// in.InputHash
	// in.OutputHash
	fmt.Println("?")

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

		lis, err := net.Listen("tcp", listen)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()

		for _, logConfig := range config.Logs {
			c, err := core.CoreFromConfig(logConfig)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(c)
		}

		pb.RegisterTrustixServer(s, &pbServer{})

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		return nil
	},
}

func initCommands() {
	rootCmd.Flags().StringVar(&configPath, "config", "", "Path to config.toml")
	rootCmd.Flags().StringVar(&listen, "listen", ":8080", "Path to config.toml")

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
