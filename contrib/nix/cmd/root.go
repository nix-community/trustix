// MIT License
//
// Copyright (c) 2020 Tweag IO
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

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
