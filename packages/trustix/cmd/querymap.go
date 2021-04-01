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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	pb "github.com/tweag/trustix/packages/trustix-proto/proto"
	"github.com/tweag/trustix/packages/trustix/client"
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
