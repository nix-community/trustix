// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docFormat string
var docOutputDir string

var docCommand = &cobra.Command{
	Use:    "__doc",
	Short:  "Generate documentation",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		{
			if _, err := os.Stat(docOutputDir); !os.IsNotExist(err) {
				os.RemoveAll(docOutputDir)
			}

			err := os.Mkdir(docOutputDir, 0755)
			if err != nil && err != os.ErrExist {
				log.Fatal(err)
			}
		}

		switch docFormat {

		case "markdown":
			err := doc.GenMarkdownTree(rootCmd, docOutputDir)
			if err != nil {
				log.Fatal(err)
			}

		default:
			log.Fatalf("Unhandled doc format: %s", docFormat)
		}

		return nil
	},
}

func initDoc() {
	docCommand.Flags().StringVar(&docFormat, "format", "markdown", "Output format")
	docCommand.Flags().StringVar(&docOutputDir, "out", "doc", "Output directory")
}
