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
	trustixSock := os.Getenv("TRUSTIX_SOCK")
	if trustixSock == "" {
		tmpDir := "/tmp"
		trustixSock = filepath.Join(tmpDir, "trustix.sock")
	}
	trustixSock = fmt.Sprintf("unix://%s", trustixSock)

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
