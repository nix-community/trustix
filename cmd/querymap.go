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
	"github.com/tweag/trustix/client"
	pb "github.com/tweag/trustix/proto"
)

var queryMap = &cobra.Command{
	Use:   "queryMap",
	Short: "Query hashes from the logs (multiple)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if keyHex == "" {
			return fmt.Errorf("Missing input/output hash")
		}

		inputBytes, err := hex.DecodeString(keyHex)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := client.CreateClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewTrustixCombinedRPCClient(conn)

		ctx, cancel := client.CreateContext(timeout)
		defer cancel()

		log.WithFields(log.Fields{
			"key": keyHex,
		}).Debug("Requesting output mappings for")
		r, err := c.Get(ctx, &pb.KeyRequest{
			Key: inputBytes,
		})
		if err != nil {
			log.Fatalf("could not query: %v", err)
		}

		for name, h := range r.Entries {
			s := hex.EncodeToString(h.Digest)
			if err != nil {
				return err
			}

			fmt.Println(fmt.Sprintf("%s: %s", name, s))
		}

		return nil
	},
}

func initMapQuery() {
	queryMap.Flags().StringVar(&keyHex, "input-hash", "", "Input hash in hex encoding")
}
