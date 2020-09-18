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

var inputHashHex string
var outputHashHex string

var submitCommand = &cobra.Command{
	Use:   "submit",
	Short: "Submit hashes for inclusion in the log",
	RunE: func(cmd *cobra.Command, args []string) error {
		if inputHashHex == "" || outputHashHex == "" {
			return fmt.Errorf("Missing input/output hash")
		}

		inputBytes, err := hex.DecodeString(inputHashHex)
		if err != nil {
			log.Fatal(err)
		}

		outputBytes, err := hex.DecodeString(outputHashHex)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := grpc.Dial(dialAddress, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewTrustixRPCClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.SubmitMapping(ctx, &pb.SubmitRequest{
			OutputHash: outputBytes,
			InputHash:  inputBytes,
		})
		if err != nil {
			log.Fatalf("could not submit: %v", err)
		}

		fmt.Println(r.Status)

		return nil
	},
}

func initSubmit() {
	submitCommand.Flags().StringVar(&inputHashHex, "input-hash", "", "Input hash in hex encoding")
	submitCommand.Flags().StringVar(&outputHashHex, "output-hash", "", "Output hash in hex encoding")
}
