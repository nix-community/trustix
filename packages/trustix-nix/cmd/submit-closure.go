// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/trustix-proto/api"
	"github.com/tweag/trustix/packages/trustix/client"
)

var submitClosureCommand = &cobra.Command{
	Use:   "submit-closure",
	Short: "Submit an entire closur for inclusion in the log (development/testing ONLY)",
	Long: `Submit an entire closur for inclusion in the log.
           This is meant for development use ONLY as it will submit all packages, even substituted ones.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		storePaths := []string{}
		{
			requisites := make(map[string]struct{})
			for _, arg := range args {
				out, err := exec.Command("nix-store", "--query", "--requisites", arg).Output()
				if err != nil {
					log.Fatalf("Could not query requisites: %v", err)
				}
				for _, path := range strings.Split(string(out), "\n") {
					if path == "" {
						continue
					}
					requisites[path] = struct{}{}
				}
			}

			for key, _ := range requisites {
				storePaths = append(storePaths, key)
			}

			if len(storePaths) < 1 {
				log.Fatal("Store paths is empty, expected at least one path to submit")
			}
		}

		items := []*api.KeyValuePair{}

		for _, storePath := range storePaths {

			item, err := createKVPair(storePath)
			if err != nil {
				log.Fatal(err)
			}

			items = append(items, item)
		}

		req := &api.SubmitRequest{
			Items: items,
		}

		conn, err := client.CreateClientConn(dialAddress, nil)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		ctx, cancel := client.CreateContext(30)
		defer cancel()

		c := api.NewTrustixLogAPIClient(conn)
		_, err = c.Submit(ctx, req)
		if err != nil {
			log.Fatalf("could not submit: %v", err)
		}

		return nil
	},
}
