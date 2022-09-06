// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var once sync.Once

var dialAddress string

var logID string

var rootCmd = &cobra.Command{
	Use:   "trustix-nix",
	Short: "Trustix nix integration",
	Long:  `Trustix nix integration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func initCommands() {
	trustixSock := os.Getenv("TRUSTIX_RPC")
	if trustixSock == "" {
		tmpDir := "/tmp"
		trustixSock = fmt.Sprintf("unix://%s", filepath.Join(tmpDir, "trustix.sock"))
	}

	rootCmd.PersistentFlags().StringVar(&dialAddress, "address", trustixSock, "Connect to address")
	rootCmd.PersistentFlags().StringVar(&logID, "log-id", "", "Log ID")

	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stderr)

	rootCmd.AddCommand(nixHookCommand)

	rootCmd.AddCommand(binaryCacheCommand)
	initBinaryCache()

	rootCmd.AddCommand(submitClosureCommand)
}

func Execute() {
	once.Do(initCommands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
