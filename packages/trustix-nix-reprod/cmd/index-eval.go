// Copyright (C) 2022 adisbladis
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/trustix-nix-reprod/index"
)

const sqlDialect = "sqlite"

// Arbitrary large number of derivations to cache
const drvCacheSize = 30_000

// Sentinel values returned when indexing a derivation with errors or filtered
const (
	errorID       = -1
	fixedOutputID = -2
)

var indexEvalCommand = &cobra.Command{
	Use:   "index-eval",
	Short: "Index evaluation",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		db, err := sql.Open(sqlDialect, "/home/adisbladis/foo.sqlite3?_journal_mode=WAL")
		if err != nil {
			return fmt.Errorf("error opening database: %w", err)
		}

		err = migrate(db, sqlDialect)
		if err != nil {
			panic(err)
		}

		err = index.IndexEval(ctx, db)
		if err != nil {
			panic(err)
		}

		return nil
	},
}
