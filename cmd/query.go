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
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/api"
	"github.com/tweag/trustix/schema"
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

		conn, err := createClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := api.NewTrustixLogAPIClient(conn)
		ctx, cancel := createContext()
		defer cancel()

		log.Debug("Requesting STH")
		sth, err := c.GetSTH(ctx, &api.STHRequest{})
		if err != nil {
			log.Fatalf("could not get STH: %v", err)
		}

		log.WithFields(log.Fields{
			"inputHash": inputHashHex,
		}).Debug("Requesting output mapping for")
		r, err := c.GetMapValue(ctx, &api.GetMapValueRequest{
			Key:     inputBytes,
			MapRoot: sth.MapRoot,
		})
		if err != nil {
			log.Fatalf("could not query: %v", err)
		}

		mapEntry := &schema.MapEntry{}
		err = json.Unmarshal(r.Value, mapEntry)
		if err != nil {
			log.Fatalf("Could not unmarshal value")
		}

		fmt.Println(fmt.Sprintf("Output hash: %s", hex.EncodeToString(mapEntry.Value)))

		return nil
	},
}

func initQuery() {
	queryCommand.Flags().StringVar(&inputHashHex, "input-hash", "", "Input hash in hex encoding")
}
