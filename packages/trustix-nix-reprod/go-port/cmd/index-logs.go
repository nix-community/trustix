// Copyright (C) 2022 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"github.com/spf13/cobra"
)

var indexLogsCommand = &cobra.Command{
	Use:   "index-logs",
	Short: "Index log build outputs (all known logs)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
