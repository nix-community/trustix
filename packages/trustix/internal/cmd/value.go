// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

		interceptors, err := getAuthInterceptors()
		if err != nil {
			log.Fatalf("Could not get auth interceptor: %v", err)
		}

		c, err := client.CreateClient(dialAddress, interceptors)
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
