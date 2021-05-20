// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	"github.com/tweag/trustix/packages/trustix/client"
)

var queryCommand = &cobra.Command{
	Use:   "query",
	Short: "Query values from the log",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Verify input params
		{

			if logID == "" {
				return fmt.Errorf("Missing log ID")
			}

			if keyHex == "" {
				return fmt.Errorf("Missing key parameter")
			}

		}

		inputBytes, err := hex.DecodeString(keyHex)
		if err != nil {
			log.Fatal(err)
		}

		c, err := client.CreateClientConn(dialAddress)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer c.Close()

		ctx, cancel := client.CreateContext(timeout)
		defer cancel()

		log.Debug("Requesting log head")
		sth, err := c.LogAPI.GetHead(ctx, &api.LogHeadRequest{
			LogID: &logID,
		})
		if err != nil {
			log.Fatalf("could not get STH: %v", err)
		}

		log.WithFields(log.Fields{
			"key": keyHex,
		}).Debug("Requesting output mapping for")
		r, err := c.LogAPI.GetMapValue(ctx, &api.GetMapValueRequest{
			LogID:   &logID,
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

		fmt.Println(fmt.Sprintf("Output digest: %s", hex.EncodeToString(mapEntry.Digest)))

		return nil
	},
}

func initQuery() {
	queryCommand.Flags().StringVar(&keyHex, "key", "", "Key in hex encoding")
}
