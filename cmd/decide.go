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
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pb "github.com/tweag/trustix/proto"
)

var decideCommand = &cobra.Command{
	Use:   "decide",
	Short: "Decide output hash from the logs (multiple)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if inputHashHex == "" {
			return fmt.Errorf("Missing input/output hash")
		}

		inputBytes, err := hex.DecodeString(inputHashHex)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := createClientConn()
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewTrustixLogClient(conn)

		ctx, cancel := createContext()
		defer cancel()

		log.WithFields(log.Fields{
			"inputHash": inputHashHex,
		}).Debug("Requesting output mappings for")
		r, err := c.Decide(ctx, &pb.HashRequest{
			InputHash: inputBytes,
		})
		if err != nil {
			log.Fatalf("could not query: %v", err)
		}

		for _, miss := range r.Misses {
			fmt.Println(fmt.Sprintf("Did not find hash in log '%s'", miss))
		}

		for _, unmatched := range r.Mismatches {
			fmt.Println(fmt.Sprintf("Found mismatched hash '%s' in log '%s'", hex.EncodeToString(unmatched.OutputHash), unmatched.LogName))
		}

		if r.Decision != nil {
			fmt.Println(fmt.Sprintf("Decided on output hash '%s' with confidence %d", hex.EncodeToString(r.Decision.OutputHash), r.Decision.Confidence))
		}

		return nil
	},
}

func initDecide() {
	decideCommand.Flags().StringVar(&inputHashHex, "input-hash", "", "Input hash in hex encoding")
}
