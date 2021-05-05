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
	"github.com/tweag/trustix/packages/trustix-proto/api"
	pb "github.com/tweag/trustix/packages/trustix-proto/rpc"
	"github.com/tweag/trustix/packages/trustix/client"
)

var keyHex string
var valueHex string

var submitCommand = &cobra.Command{
	Use:   "submit",
	Short: "Submit hashes for inclusion in the log",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Verify input params
		{

			if logID == "" {
				return fmt.Errorf("Missing log ID")
			}

			if keyHex == "" {
				return fmt.Errorf("Missing key parameter")
			}

			if valueHex == "" {
				return fmt.Errorf("Missing value parameter")
			}

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

		c := pb.NewTrustixRPCClient(conn)

		log.WithFields(log.Fields{
			"key":   keyHex,
			"value": valueHex,
		}).Debug("Submitting mapping")

		r, err := c.Submit(ctx, &pb.SubmitRequest{
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
	submitCommand.Flags().StringVar(&keyHex, "input-hash", "", "Input hash in hex encoding")
	submitCommand.Flags().StringVar(&valueHex, "output-hash", "", "Output hash in hex encoding")
}
