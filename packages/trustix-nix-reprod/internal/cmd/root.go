// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/adrg/xdg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var once sync.Once

// Dial address to trustix daemon
var dialAddress string
var stateDirectory string

var rootCmd = &cobra.Command{
	Use:   "trustix-nix-reprod",
	Short: "Trustix Nix build reproducibility dashboard",
	Long:  `Trustix nix build reproducibility dashboard`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func initCommands() {

	// Trustix daemon RPC flag
	{
		trustixSock := os.Getenv("TRUSTIX_RPC")
		if trustixSock == "" {
			tmpDir := "/tmp"
			trustixSock = fmt.Sprintf("unix://%s", filepath.Join(tmpDir, "trustix.sock"))
		}

		rootCmd.PersistentFlags().StringVar(&dialAddress, "address", trustixSock, "Connect to address")
	}

	// State directory
	{
		defaultStateDir := filepath.Join(xdg.DataHome, "trustix-nix-reprod")
		rootCmd.PersistentFlags().StringVar(&stateDirectory, "state", defaultStateDir, "State directory")
	}

	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stderr)

	rootCmd.AddCommand(indexEvalCommand)
	rootCmd.AddCommand(indexLogsCommand)
	rootCmd.AddCommand(serveCommand)
	// rootCmd.AddCommand(queryLogsCommand)

	initServe()
}

func Execute() {
	once.Do(initCommands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
