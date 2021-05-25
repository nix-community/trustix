// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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
