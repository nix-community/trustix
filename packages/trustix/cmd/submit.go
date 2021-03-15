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
	"github.com/tweag/trustix/packages/trustix/api"
	"github.com/tweag/trustix/packages/trustix/client"
)

var keyHex string
var valueHex string

var submitCommand = &cobra.Command{
	Use:   "submit",
	Short: "Submit hashes for inclusion in the log",
	RunE: func(cmd *cobra.Command, args []string) error {
		if keyHex == "" || valueHex == "" {
			return fmt.Errorf("Missing input/output hash")
		}

		inputBytes, err := hex.DecodeString(keyHex)
		if err != nil {
			log.Fatal(err)
		}

		outputBytes, err := hex.DecodeString(valueHex)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := client.CreateClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := client.CreateContext(timeout)
		defer cancel()

		c := api.NewTrustixLogAPIClient(conn)

		log.WithFields(log.Fields{
			"key":   keyHex,
			"value": valueHex,
		}).Debug("Submitting mapping")

		r, err := c.Submit(ctx, &api.SubmitRequest{
			Items: []*api.KeyValuePair{
				&api.KeyValuePair{
					Key:   inputBytes,
					Value: outputBytes,
				},
			},
		})
		if err != nil {
			log.Fatalf("could not submit: %v", err)
		}

		fmt.Println(r.Status)

		return nil
	},
}

func initSubmit() {
	submitCommand.Flags().StringVar(&keyHex, "input-hash", "", "Input hash in hex encoding")
	submitCommand.Flags().StringVar(&valueHex, "output-hash", "", "Output hash in hex encoding")
}
