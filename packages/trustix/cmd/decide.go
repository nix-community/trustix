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
	pb "github.com/tweag/trustix/packages/trustix-proto/rpc"
	"github.com/tweag/trustix/packages/trustix/client"
)

var decideCommand = &cobra.Command{
	Use:   "decide",
	Short: "Decide on output from configured logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if keyHex == "" {
			return fmt.Errorf("Missing key param")
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

		log.WithFields(log.Fields{
			"key": keyHex,
		}).Debug("Requesting output mappings for")

		r, err := c.RpcAPI.Decide(ctx, &pb.KeyRequest{
			Key: inputBytes,
		})
		if err != nil {
			log.Fatalf("could not decide: %v", err)
		}

		for _, miss := range r.Misses {
			fmt.Println(fmt.Sprintf("Did not find hash in log '%s'", miss))
		}

		for _, unmatched := range r.Mismatches {
			fmt.Println(fmt.Sprintf("Found mismatched digest '%s' in log '%s'", hex.EncodeToString(unmatched.Digest), *unmatched.LogID))
		}

		if r.Decision != nil {
			fmt.Println(fmt.Sprintf("Decided on output digest '%s' with confidence %d", hex.EncodeToString(r.Decision.Digest), *r.Decision.Confidence))
		}

		return nil
	},
}

func initDecide() {
	decideCommand.Flags().StringVar(&keyHex, "key", "", "Key in hex encoding")
}
