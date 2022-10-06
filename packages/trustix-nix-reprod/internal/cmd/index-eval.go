// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"context"
	"fmt"

	"github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/index"
	"github.com/spf13/cobra"
)

var indexEvalCommand = &cobra.Command{
	Use:   "index-eval",
	Short: "Index evaluation",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		db, err := setupDB(stateDirectory)
		if err != nil {
			return fmt.Errorf("error opening database: %w", err)
		}

		err = index.IndexEval(ctx, db)
		if err != nil {
			panic(err)
		}

		return nil
	},
}
