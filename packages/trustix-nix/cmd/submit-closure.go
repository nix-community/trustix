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
