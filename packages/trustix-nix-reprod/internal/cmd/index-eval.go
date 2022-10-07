// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/hydra"
	"github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/index"
	"github.com/spf13/cobra"
)

var indexEvalCommand = &cobra.Command{
	Use:   "index-eval",
	Short: "Index evaluation",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		dbs, err := setupDatabases(stateDirectory)
		if err != nil {
			return fmt.Errorf("error opening database: %w", err)
		}

		timestamp := time.Now().UTC()
		revision := "cabcec1477db472a0272a909fd88588ec3afc2d3"
		channel := "nixos-unstable"

		var nixpkgs string
		{
			githubOrg := "NixOS"
			githubRepo := "nixpkgs"

			u, err := url.Parse("https://github.com")
			if err != nil {
				panic(err)
			}

			u.Path = path.Join(githubOrg, githubRepo, "archive", revision+".tar.gz")

			nixpkgs = u.String()
		}

		evalMeta := &hydra.HydraEval{
			ID:        1234,
			Timestamp: 1234,
			EvalInputs: map[string]*hydra.JobsetEvalInput{
				"nixpkgs": &hydra.JobsetEvalInput{
					Revision: revision,
				},
			},
		}

		err = index.IndexEval(ctx, dbs.dbRW, nixpkgs, channel, timestamp, evalMeta)
		if err != nil {
			panic(err)
		}

		return nil
	},
}
