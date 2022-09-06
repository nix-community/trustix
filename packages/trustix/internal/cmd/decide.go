// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"encoding/hex"
	"fmt"

	pb "github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var decideKeyHex string
var decideProtocol string

var decideCommand = &cobra.Command{
	Use:   "decide",
	Short: "Decide on output from configured logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if decideKeyHex == "" {
			return fmt.Errorf("Missing key param")
		}

		if decideProtocol == "" {
			return fmt.Errorf("Missing protocol parameter")
		}

		inputBytes, err := hex.DecodeString(decideKeyHex)
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
			"key": decideKeyHex,
		}).Debug("Requesting output mappings for")

		r, err := c.RpcAPI.Decide(ctx, &pb.DecideRequest{
			Key:      inputBytes,
			Protocol: &decideProtocol,
		})
		if err != nil {
			log.Fatalf("could not decide: %v", err)
		}

		for _, miss := range r.Misses {
			fmt.Printf("Did not find digest in log '%s'\n", miss)
		}

		for _, unmatched := range r.Mismatches {
			fmt.Printf("Found mismatched digest '%s' in log '%s'\n", hex.EncodeToString(unmatched.Digest), *unmatched.LogID)
		}

		if r.Decision != nil {
			fmt.Printf("Decided on output digest '%s' with confidence %d\n", hex.EncodeToString(r.Decision.Digest), *r.Decision.Confidence)
		}

		return nil
	},
}

func initDecide() {
	decideCommand.Flags().StringVar(&decideKeyHex, "key", "", "Key in hex encoding")
	decideCommand.Flags().StringVar(&decideProtocol, "protocol", "", "Protocol ID")
}
