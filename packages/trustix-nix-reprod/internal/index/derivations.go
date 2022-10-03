// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package index

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/nix-community/trustix/packages/go-lib/executor"
	"github.com/nix-community/trustix/packages/go-lib/safemap"
	"github.com/nix-community/trustix/packages/go-lib/set"
	idb "github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/db"
	drvparse "github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/derivation"
	"github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/eval"
)

// Arbitrary large number of derivations to cache
const drvCacheSize = 30_000

// Sentinel values returned when indexing a derivation with errors or filtered
const (
	errorID = -1
)

func IndexEval(ctx context.Context, db *sql.DB) error {
	evalConfig := eval.NewConfig()
	evalConfig.Expr = "./pkgs.nix"

	// Indexing impl
	// commitSha := "c4c79f09a599717dfd57134cdd3c6e387a764f63"
	commitSha := "9c5efb63754024dd4026dceb6f3525934009fea9"
	maxWorkers := 15

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error creating db transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			panic(err)
		}
	}()

	queries := idb.New(db)
	qtx := queries.WithTx(tx)

	// Create the evaluation in the database
	dbEval, err := qtx.GetEval(ctx, commitSha)
	if err == nil {
		fmt.Println(fmt.Sprintf("eval '%s' already indexed", commitSha))
		return nil
	} else {
		if err == sql.ErrNoRows {
			dbEval, err = qtx.CreateEval(ctx, commitSha)
		}

		if err != nil {
			return fmt.Errorf("error retreiving eval: %w", err)
		}
	}

	results, err := eval.Eval(ctx, evalConfig)
	if err != nil {
		return fmt.Errorf("error initialising eval: %w", err)
	}

	drvParser, err := drvparse.NewCachedDrvParser(drvCacheSize)
	if err != nil {
		return fmt.Errorf("error creating cached derivation parser: %w", err)
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
				return errorID, fmt.Errorf("error getting map value: %w", err)
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
				return errorID, fmt.Errorf("error getting derivation id: %w", err)
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
				return errorID, fmt.Errorf("error getting derivation: %w", err)
			}

			// Create the derivation in the DB
			dbDrv, err = qtx.CreateDerivation(ctx, idb.CreateDerivationParams{
				Drv:    drvPath,
				System: drv.Platform,
			})
			if err != nil {
				return errorID, fmt.Errorf("error creating derivation: %w", err)
			}

			drvDBIDs.Set(drvPath, dbDrv.ID)

			// Index that this derivation was a part of this evaluation
			err = qtx.CreateDerivationEval(ctx, idb.CreateDerivationEvalParams{
				Drv:  dbDrv.ID,
				Eval: dbEval.ID,
			})
			if err != nil {
				return errorID, fmt.Errorf("error creating derivationeval: %w", err)
			}
		}

		// Direct dependencies
		refsDirect := set.NewSet[string]()
		for inputDrv := range drv.InputDerivations {
			refsDirect.Add(inputDrv)
		}

		// insert a self-reference
		refsDirect.Add(drvPath)

		// All dependencies (recursive, flattened)
		refsAll := refsDirect.Copy()

		for inputDrv := range drv.InputDerivations {
			// Recursively index drvs
			if !refs.Has(inputDrv) {
				_, err := indexDrv(inputDrv)
				if err != nil {
					return errorID, fmt.Errorf("error indexing ref derivation: %w", err)
				}
			}

			// If the input _still_ doesn't exist it means it's a fixed-output
			// and should be filtered out
			if refs.Has(inputDrv) {
				inputRefs, err := refs.Get(inputDrv)
				if err != nil {
					return errorID, fmt.Errorf("error getting references: %w", err)
				}

				refsAll.Update(inputRefs)
			} else {
				refsDirect.Remove(inputDrv)
				refsAll.Remove(inputDrv)
			}
		}

		refs.Set(drvPath, refsDirect)

		// Create derivation outputs
		for output, pathInfo := range drv.Outputs {
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

				err = qtx.CreateDerivationRefDirect(ctx, idb.CreateDerivationRefDirectParams{
					ReferrerID: dbDrv.ID,
					DrvID:      dbID,
				})
				if err != nil {
					return errorID, fmt.Errorf("error creating direct derivation ref: %w", err)
				}
			}

			// Create relation for all recursive references
			for _, ref := range refsAll.Values() {
				dbID, err := getDrvID(ref)
				if err != nil {
					return errorID, err
				}

				err = qtx.CreateDerivationRefRecursive(ctx, idb.CreateDerivationRefRecursiveParams{
					ReferrerID: dbDrv.ID,
					DrvID:      dbID,
				})
				if err != nil {
					return errorID, fmt.Errorf("error creating recursive derivation ref: %w", err)
				}
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
				return fmt.Errorf("error indexing derivation %s: %w", result.DrvPath, err)
			}

			// Add mapping from attribute to derivation
			if result.Attr != "" {
				err = qtx.CreateDerivationAttr(ctx, idb.CreateDerivationAttrParams{
					Attr:         result.Attr,
					DerivationID: drvID,
				})
				if err != nil {
					return fmt.Errorf("error creating attr reference for drv %s (%d): %w", result.DrvPath, drvID, err)
				}
			}

			fmt.Println(result.Attr, drvID)

			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	err = e.Wait()
	if err != nil {
		return fmt.Errorf("error in indexing: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
