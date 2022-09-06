// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"fmt"

	pb "github.com/nix-community/trustix/packages/trustix-proto/rpc"
	"github.com/nix-community/trustix/packages/trustix/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var flushCommand = &cobra.Command{
	Use:   "flush",
	Short: "Flush submissions and write new tree head",
	RunE: func(cmd *cobra.Command, args []string) error {

		if logID == "" {
			return fmt.Errorf("Missing log ID")
		}

		c, err := client.CreateClient(dialAddress)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer c.Close()

		ctx, cancel := client.CreateContext(timeout)
		defer cancel()

		_, err = c.LogRPC.Flush(ctx, &pb.FlushRequest{
			LogID: &logID,
		})
		if err != nil {
			log.Fatalf("could not flush: %v", err)
		}

		return nil
	},
}

func initFlush() {
}
