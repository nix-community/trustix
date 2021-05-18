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
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix/client"
)

var valueDigest string

var getValueCommand = &cobra.Command{
	Use:   "get-value",
	Short: "Get computed value based on content digest",
	RunE: func(cmd *cobra.Command, args []string) error {
		if valueDigest == "" {
			return fmt.Errorf("Missing value digest parameter")
		}

		digest, err := hex.DecodeString(valueDigest)
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

		resp, err := c.NodeAPI.GetValue(ctx, &api.ValueRequest{
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
