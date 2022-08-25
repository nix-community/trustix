// Copyright (C) 2022 Tweag IO
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
	"github.com/tweag/trustix/packages/go-lib/executor"
	"github.com/tweag/trustix/packages/go-lib/safemap"
	"github.com/tweag/trustix/packages/go-lib/set"
	idb "github.com/tweag/trustix/packages/trustix-nix-reprod/db"
	drvparse "github.com/tweag/trustix/packages/trustix-nix-reprod/derivation"
	"github.com/tweag/trustix/packages/trustix-nix-reprod/eval"
)

const sqlDialect = "sqlite"

var indexEvalCommand = &cobra.Command{
	Use:   "index-eval",
	Short: "Index evaluation",
	RunE: func(cmd *cobra.Command, args []string) error {

		evalConfig := eval.NewConfig()
		evalConfig.Expr = "./pkgs.nix"

		ctx := context.Background()

		db, err := sql.Open(sqlDialect, "./foo.sqlite3")
		if err != nil {
			return err
		}

		err = migrate(db, sqlDialect)
		if err != nil {
			panic(err)
		}

		// Indexing impl
		commitSha := "c4c79f09a599717dfd57134cdd3c6e387a764f63"
		maxWorkers := 15

		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		queries := idb.New(db)
		qtx := queries.WithTx(tx)

		// Create the evaluation in the database
		_, err = qtx.GetEval(ctx, commitSha)
		if err != nil {
			if err == sql.ErrNoRows {
				_, err = qtx.CreateEval(ctx, commitSha)
			}

			if err != nil {
				panic(err)
			}
		}

		results, err := eval.Eval(ctx, evalConfig)
		if err != nil {
			panic(err)
		}

		drvParser, err := drvparse.NewCachedDrvParser()
		if err != nil {
			panic(err)
		}

		// Map drv to it's direct references for later re-use
		refs := safemap.NewMap[string, *set.Set[string]]()

		alreadyIndexed := set.NewSet[string]()

		// Index a derivation including it's dependencies
		var indexDrv func(string) error
		indexDrv = func(drvPath string) error {
			// No-op if already indexed
			if alreadyIndexed.Has(drvPath) {
				return nil
			} else {
				alreadyIndexed.Add(drvPath)
			}

			drv, err := drvParser.ReadPath(drvPath)
			if err != nil {
				return fmt.Errorf("Error reading '%s': %w", drvPath, err)
			}

			// Direct dependencies
			refsDirect := set.NewSet[string]()
			for inputDrv, _ := range drv.InputDerivations {
				refsDirect.Add(inputDrv)
			}

			// All dependencies (recursive, flattened)
			refsAll := refsDirect.Copy()

			for inputDrv, _ := range drv.InputDerivations {
				if !refs.Has(inputDrv) {
					err := indexDrv(inputDrv)
					if err != nil {
						return err
					}
				}

				// If the input _still_ doesn't exist it means it's a fixed-output
				// and should be filtered out
				if refs.Has(inputDrv) {
					inputRefs, err := refs.Get(inputDrv)
					if err != nil {
						return err
					}

					refsAll.Update(inputRefs)
				} else {
					refsDirect.Remove(inputDrv)
					refsAll.Remove(inputDrv)
				}
			}

			// Filter fixed outputs
			for _, output := range drv.Outputs {
				if output.HashAlgorithm != "" {
					return nil
				}
			}

			refs.Set(drvPath, refsDirect)

			// Create the derivation in the database
			dbDrv, err := qtx.GetDerivation(ctx, drvPath)
			if err != nil {
				if err == sql.ErrNoRows {
					dbDrv, err = qtx.CreateDerivation(ctx, idb.CreateDerivationParams{
						Drv:    drvPath,
						System: drv.Platform,
					})
				}

				if err != nil {
					panic(err)
				}
			}

			fmt.Println(dbDrv)

			return nil
		}

		e := executor.NewLimitedParallellExecutor(maxWorkers)

		for wrappedResult := range results {
			result, err := wrappedResult.Unwrap()
			if err != nil {
				panic(err)
			}

			e.Add(func() error {
				return indexDrv(result.DrvPath)
			})
		}

		err = e.Wait()
		if err != nil {
			panic(err)
		}

		err = tx.Commit()
		if err != nil {
			panic(err)
		}

		return nil
	},
}
