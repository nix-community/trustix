package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pb "github.com/tweag/trustix/proto"
	"google.golang.org/grpc"
	"time"
)

var queryCommand = &cobra.Command{
	Use:   "query",
	Short: "Query hashes from the log",
	RunE: func(cmd *cobra.Command, args []string) error {
		if inputHashHex == "" {
			return fmt.Errorf("Missing input/output hash")
		}

		inputBytes, err := hex.DecodeString(inputHashHex)
		if err != nil {
			log.Fatal(err)
		}

		log.WithFields(log.Fields{
			"address": dialAddress,
		}).Debug("Dialing gRPC")
		conn, err := grpc.Dial(dialAddress, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewTrustixRPCClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		log.WithFields(log.Fields{
			"inputHash": inputHashHex,
		}).Debug("Requesting output mapping for")
		r, err := c.QueryMapping(ctx, &pb.QueryRequest{
			InputHash: inputBytes,
		})
		if err != nil {
			log.Fatalf("could not submit: %v", err)
		}

		fmt.Println(fmt.Sprintf("Output hash: %s", hex.EncodeToString(r.OutputHash)))

		return nil
	},
}

func initQuery() {
	queryCommand.Flags().StringVar(&inputHashHex, "input-hash", "", "Input hash in hex encoding")
}
