// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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
}

func Execute() {
	initOnce.Do(initCommands)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
