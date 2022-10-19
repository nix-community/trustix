// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"os/exec"
	"strings"

	"github.com/nix-community/trustix/packages/trustix-proto/api"
	pb "github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var submitClosureCommand = &cobra.Command{
	Use:   "submit-closure",
	Short: "Submit an entire closure for inclusion in the log (development/testing ONLY)",
	Long: `Submit an entire closure for inclusion in the log.
           This is meant for development use ONLY as it will submit all packages, even substituted ones.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if logID == "" {
			log.Fatal("Missing log id parameter")
		}

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

			for key := range requisites {
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

		req := &pb.SubmitRequest{
			LogID: &logID,
			Items: items,
		}

		interceptors, err := getAuthInterceptors()
		if err != nil {
			log.Fatalf("Could not get auth interceptor: %w", err)
		}

		c, err := client.CreateClient(dialAddress, interceptors)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer c.Close()

		ctx, cancel := client.CreateContext(30)
		defer cancel()

		_, err = c.LogRPC.Submit(ctx, req)
		if err != nil {
			log.Fatalf("could not submit: %v", err)
		}

		return nil
	},
}
