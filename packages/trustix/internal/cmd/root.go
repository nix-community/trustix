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

var initOnce sync.Once

var (
	dialAddress string
	logID       string
	timeout     int
)

var rootCmd = &cobra.Command{
	Use:   "trustix",
	Short: "Trustix",
	Long:  `Trustix`,
}

func initCommands() {

	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 20, "Timeout in seconds")

	rootCmd.PersistentFlags().StringVar(&logID, "log-id", "", "Log ID")

	trustixSock := os.Getenv("TRUSTIX_RPC")
	if trustixSock == "" {
		tmpDir := "/tmp"
		trustixSock = fmt.Sprintf("unix://%s", filepath.Join(tmpDir, "trustix.sock"))
	}
	rootCmd.PersistentFlags().StringVar(&dialAddress, "address", trustixSock, "Connect to address")

	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stderr)

	rootCmd.AddCommand(daemonCmd)
	initDaemon()

	rootCmd.AddCommand(generateKeyCmd)
	initGenerate()

	rootCmd.AddCommand(printLogIdCmd)
	initPrintLogId()

	rootCmd.AddCommand(submitCommand)
	initSubmit()

	rootCmd.AddCommand(queryCommand)
	initQuery()

	rootCmd.AddCommand(getValueCommand)
	initGetValue()

	rootCmd.AddCommand(decideCommand)
	initDecide()

	rootCmd.AddCommand(flushCommand)
	initFlush()

	rootCmd.AddCommand(docCommand)
	initDoc()

	rootCmd.AddCommand(generateTokenCmd)
	initGenerateToken()
}

func Execute() {
	initOnce.Do(initCommands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
