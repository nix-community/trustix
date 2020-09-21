package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/spf13/cobra"
	pb "github.com/tweag/trustix/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

var queryMap = &cobra.Command{
	Use:   "queryMap",
	Short: "Query hashes from the logs (multiple)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if inputHashHex == "" {
			return fmt.Errorf("Missing input/output hash")
		}

		inputBytes, err := hex.DecodeString(inputHashHex)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := grpc.Dial(dialAddress, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewTrustixLogClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.HashMap(ctx, &pb.HashRequest{
			InputHash: inputBytes,
		})
		if err != nil {
			log.Fatalf("could not submit: %v", err)
		}

		for name, h := range r.Hashes {
			s := hex.EncodeToString(h)
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("%s: %s", name, s))
		}

		return nil
	},
}

func initMapQuery() {
	queryMap.Flags().StringVar(&inputHashHex, "input-hash", "", "Input hash in hex encoding")
}
