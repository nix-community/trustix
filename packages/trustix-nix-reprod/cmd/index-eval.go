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
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/go-lib/executor"
	"github.com/tweag/trustix/packages/go-lib/safemap"
	"github.com/tweag/trustix/packages/go-lib/set"
	idb "github.com/tweag/trustix/packages/trustix-nix-reprod/db"
	drvparse "github.com/tweag/trustix/packages/trustix-nix-reprod/derivation"
	"github.com/tweag/trustix/packages/trustix-nix-reprod/eval"
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

		evalConfig := eval.NewConfig()
		evalConfig.Expr = "./pkgs.nix"

		ctx := context.Background()

		db, err := sql.Open(sqlDialect, "/home/adisbladis/foo.sqlite3?_journal_mode=WAL")
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

		drvParser, err := drvparse.NewCachedDrvParser(drvCacheSize)
		if err != nil {
			panic(err)
		}

		// Map drv to it's direct references for later re-use
		refs := safemap.NewMap[string, *set.Set[string]]()

		// Map drv paths to DB ids so we can avoid queries in the hot indexing path
		drvDBIDs := safemap.NewMap[string, int64]()

		alreadyIndexed := set.NewSafeSet[string]()

		// indexDrv is somewhat racy but we can work around that by getting
		// a value in a loop with a timeout
		getDrvID := func(drvPath string) (dbID int64, err error) {
			for i := 0; i < 10_000; i++ {
				dbID, err = drvDBIDs.Get(drvPath)
				if err == nil {
					return dbID, nil
				}

				if err != nil && !errors.Is(err, safemap.ErrNotExist) {
					return errorID, err
				}

				time.Sleep(5 * time.Millisecond)
			}

			return -1, fmt.Errorf("Couldnt get derivation id for derivation path '%s': %w", drvPath, err)
		}

		// Index a derivation including it's dependencies
		var indexDrv func(string) (int64, error)
		indexDrv = func(drvPath string) (int64, error) {
			// No-op if already indexed, populate map early to act as a lock per drvPath
			if alreadyIndexed.Has(drvPath) {
				dbID, err := getDrvID(drvPath)
				if err != nil {
					return errorID, err
				}

				return dbID, nil
			} else {
				alreadyIndexed.Add(drvPath)
			}

			drv, err := drvParser.ReadPath(drvPath)
			if err != nil {
				return errorID, fmt.Errorf("Error reading '%s': %w", drvPath, err)
			}

			var dbDrv idb.Derivation
			{
				// Check if the derivation is already indexed
				dbDrv, err = qtx.GetDerivation(ctx, drvPath)
				if err == nil {
					drvDBIDs.Set(drvPath, dbDrv.ID)
					return dbDrv.ID, nil
				} else if err != sql.ErrNoRows {
					return errorID, err
				}

				// Create the derivation in the DB
				dbDrv, err = qtx.CreateDerivation(ctx, idb.CreateDerivationParams{
					Drv:    drvPath,
					System: drv.Platform,
				})
				if err != nil {
					return errorID, err
				}

				drvDBIDs.Set(drvPath, dbDrv.ID)
			}

			// Direct dependencies
			refsDirect := set.NewSet[string]()
			for inputDrv, _ := range drv.InputDerivations {
				refsDirect.Add(inputDrv)
			}

			// All dependencies (recursive, flattened)
			refsAll := refsDirect.Copy()

			for inputDrv, _ := range drv.InputDerivations {
				// Recursively index drvs
				if !refs.Has(inputDrv) {
					_, err := indexDrv(inputDrv)
					if err != nil {
						return errorID, err
					}
				}

				// If the input _still_ doesn't exist it means it's a fixed-output
				// and should be filtered out
				if refs.Has(inputDrv) {
					inputRefs, err := refs.Get(inputDrv)
					if err != nil {
						return errorID, err
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
					return fixedOutputID, nil
				}
			}

			refs.Set(drvPath, refsDirect)

			// Create derivation outputs
			for output, pathInfo := range drv.Outputs {
				_, err := qtx.GetDerivationOutput(ctx, idb.GetDerivationOutputParams{
					DerivationID: dbDrv.ID,
					StorePath:    pathInfo.Path,
				})
				if err == nil {
					continue
				} else if err != sql.ErrNoRows {
					return errorID, fmt.Errorf("Error fetching derivation output: %w", err)
				}

				err = qtx.CreateDerivationOutput(ctx, idb.CreateDerivationOutputParams{
					Output:       output,
					StorePath:    pathInfo.Path,
					DerivationID: dbDrv.ID,
				})
				if err != nil {
					return errorID, fmt.Errorf("Error creating derivation output: %w", err)
				}
			}

			// Create relations to referenced derivations
			{
				// Create relation for direct references
				for _, ref := range refsDirect.Values() {
					dbID, err := getDrvID(ref)
					if err != nil {
						return errorID, err
					}

					qtx.CreateDerivationRefDirect(ctx, idb.CreateDerivationRefDirectParams{
						ReferrerID: dbDrv.ID,
						DrvID:      dbID,
					})
				}

				// Create relation for all recursive references
				for _, ref := range refsDirect.Values() {
					dbID, err := getDrvID(ref)
					if err != nil {
						return errorID, err
					}

					qtx.CreateDerivationRefRecursive(ctx, idb.CreateDerivationRefRecursiveParams{
						ReferrerID: dbDrv.ID,
						DrvID:      dbID,
					})
				}
			}

			return dbDrv.ID, nil
		}

		e := executor.NewLimitedParallellExecutor(maxWorkers)

		for wrappedResult := range results {
			result, err := wrappedResult.Unwrap()
			if err != nil {
				panic(err)
			}

			if result.Error != "" || result.DrvPath == "" {
				continue
			}

			// Index the derivation + attribute mappings
			err = e.Add(func() error {
				// Index the derivation
				drvID, err := indexDrv(result.DrvPath)
				if err != nil {
					return err
				}

				// Don't index fixed outputs
				if drvID == fixedOutputID {
					return nil
				}

				// Add mapping from attribute to derivation
				if result.Attr != "" {
					fmt.Println(drvID, result.Attr)
				}

				return nil
			})
			if err != nil {
				panic(err)
			}
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
