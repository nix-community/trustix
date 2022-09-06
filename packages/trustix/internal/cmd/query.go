// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
	"github.com/nix-community/trustix/packages/trustix/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var queryKeyHex string

var queryCommand = &cobra.Command{
	Use:   "query",
	Short: "Query values from the log",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Verify input params
		{
			if logID == "" {
				return fmt.Errorf("Missing log ID")
			}

			if queryKeyHex == "" {
				return fmt.Errorf("Missing key parameter")
			}

		}

		inputBytes, err := hex.DecodeString(queryKeyHex)
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

		log.Debug("Requesting log head")
		sth, err := c.LogAPI.GetHead(ctx, &api.LogHeadRequest{
			LogID: &logID,
		})
		if err != nil {
			log.Fatalf("could not get STH: %v", err)
		}

		log.WithFields(log.Fields{
			"key": queryKeyHex,
		}).Debug("Requesting output mapping for")
		r, err := c.LogAPI.GetMapValue(ctx, &api.GetMapValueRequest{
			LogID:   &logID,
			Key:     inputBytes,
			MapRoot: sth.MapRoot,
		})
		if err != nil {
			log.Fatalf("could not query: %v", err)
		}

		mapEntry := &schema.MapEntry{}
		err = json.Unmarshal(r.Value, mapEntry)
		if err != nil {
			log.Fatalf("Could not unmarshal value")
		}

		fmt.Printf("Output digest: %s\n", hex.EncodeToString(mapEntry.Digest))

		return nil
	},
}

func initQuery() {
	queryCommand.Flags().StringVar(&queryKeyHex, "key", "", "Key in hex encoding")
}
