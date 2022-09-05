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

	"github.com/nix-community/trustix/packages/trustix-proto/api"
	pb "github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var submitKeyHex string
var submitValueHex string

var submitCommand = &cobra.Command{
	Use:   "submit",
	Short: "Submit values for inclusion in a log",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Verify input params
		{

			if logID == "" {
				return fmt.Errorf("Missing log ID")
			}

			if submitKeyHex == "" {
				return fmt.Errorf("Missing key parameter")
			}

			if submitValueHex == "" {
				return fmt.Errorf("Missing value parameter")
			}

		}

		inputBytes, err := hex.DecodeString(submitKeyHex)
		if err != nil {
			log.Fatal(err)
		}

		outputBytes, err := hex.DecodeString(submitValueHex)
		if err != nil {
			log.Fatal(err)
		}

		c, err := client.CreateClient(dialAddress)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer c.Close()

		ctx, cancel := client.CreateContext(timeout)
		defer cancel()

		log.WithFields(log.Fields{
			"key":   submitKeyHex,
			"value": submitValueHex,
		}).Debug("Submitting mapping")

		r, err := c.LogRPC.Submit(ctx, &pb.SubmitRequest{
			LogID: &logID,
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
	submitCommand.Flags().StringVar(&submitKeyHex, "key", "", "Key in hex encoding")
	submitCommand.Flags().StringVar(&submitValueHex, "value", "", "Value in hex encoding")
}
