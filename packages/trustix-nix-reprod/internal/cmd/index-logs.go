// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/index"
	tclient "github.com/nix-community/trustix/packages/trustix/client"
	"github.com/spf13/cobra"
)

var indexLogsCommand = &cobra.Command{
	Use:   "index-logs",
	Short: "Index log build outputs (all known logs)",
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

		{
			client, err := tclient.CreateClient(dialAddress)
			if err != nil {
				panic(err)
			}

			err = index.IndexLogs(ctx, db, client)
			if err != nil {
				panic(err)
			}
		}

		return nil
	},
}
