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

		conn, err := createClientConn()
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := createContext()
		defer cancel()

		c := pb.NewTrustixRPCClient(conn)

		log.WithFields(log.Fields{
			"inputHash":  inputHashHex,
			"outputHash": outputHashHex,
		}).Debug("Submitting mapping")
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
