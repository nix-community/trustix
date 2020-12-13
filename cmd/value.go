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
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/api"
)

var valueDigest string

var getValueCommand = &cobra.Command{
	Use:   "get-value",
	Short: "Get computed value based on content digest",
	RunE: func(cmd *cobra.Command, args []string) error {
		if valueDigest == "" {
			return fmt.Errorf("Missing input/output hash")
		}

		digest, err := hex.DecodeString(valueDigest)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := createClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := api.NewTrustixLogAPIClient(conn)
		ctx, cancel := createContext()
		defer cancel()

		resp, err := c.GetValue(ctx, &api.ValueRequest{
			Digest: digest,
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = os.Stdout.Write(resp.Value)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func initGetValue() {
	getValueCommand.Flags().StringVar(&valueDigest, "digest", "", "Value content digest")
}
