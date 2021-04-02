// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix/client"
)

var flushCommand = &cobra.Command{
	Use:   "flush",
	Short: "Flush submissions and write new tree head",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := checkLogID(); err != nil {
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

		_, err = c.Flush(ctx, &api.FlushRequest{
			LogID: &logID,
		})
		if err != nil {
			log.Fatalf("could not flush: %v", err)
		}

		return nil
	},
}
