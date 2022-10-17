// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package hydra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/config"
	idb "github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/db"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/eval"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/index"
	log "github.com/sirupsen/logrus"
)

const requestPoolSize = 15

func indexHydraJobset(ctx context.Context, db *sql.DB, channel string, baseURL string, hydraJobset *HydraJobset, hydraEval *HydraEval) error {
	timestamp := time.Unix(hydraEval.Timestamp, 0)

	nixpath, err := hydraEval.NixPath()
	if err != nil {
		return fmt.Errorf("error getting NIX_PATH: %w", err)
	}

	evalConfig := eval.NewConfig()
	evalConfig.NixPath = nixpath.String()
	evalConfig.ForceRecurse = true

	// Resolve the evaluation expression to the full store path and ensure it's not trying to acceess files outside of NIX_PATH
	{
		storePath, ok := nixpath[hydraJobset.NixExprInput]
		if !ok {
			return fmt.Errorf("Nix path missing expr input '%s'", hydraJobset.NixExprInput)
		}

		exprPath := path.Join(storePath, hydraJobset.NixExprPath)
		absExprPath, err := filepath.Abs(exprPath)
		if err != nil {
			return fmt.Errorf("error getting absolute path for '%s': %w", exprPath, err)
		}

		if !strings.HasPrefix(absExprPath, storePath) {
			return fmt.Errorf("Nix expression path '%s' outside of prefix '%s', possible injection attempt", absExprPath, storePath)
		}

		evalConfig.ExprPath = exprPath
	}

	// Add jobset inputs as parameters
	for k, v := range nixpath {
		evalConfig.AddArgStr(k, v)
	}

	evalResults, err := eval.Eval(ctx, evalConfig)
	if err != nil {
		return fmt.Errorf("error initialising eval: %w", err)
	}

	// For simplicity gather all builds up-front
	evalAttrs := []*index.EvalAttribute{}
	for wrappedResult := range evalResults {
		result, err := wrappedResult.Unwrap()
		if err != nil {
			return fmt.Errorf("eval result error: %w", err)
		}

		if result.Error != "" || result.DrvPath == "" {
			continue
		}

		evalAttrs = append(evalAttrs, &index.EvalAttribute{
			DrvPath: result.DrvPath,
			Attr:    result.Attr,
		})
	}

	return index.IndexEval(ctx, db, channel, timestamp, evalAttrs, func(ctx context.Context, dbEval idb.Evaluation, qtx *idb.Queries) error {
		var revision string
		{
			for _, input := range hydraEval.EvalInputs {
				if input.Type != "git" {
					continue
				}

				if input.Revision != "" {
					revision = input.Revision
					break
				}
			}

			if revision == "" {
				return fmt.Errorf("No revision could be extracted from hydra inputs")
			}
		}

		_, err := qtx.CreateHydraEval(ctx, idb.CreateHydraEvalParams{
			Evaluation:  dbEval.ID,
			HydraEvalID: hydraEval.ID,
			Revision:    revision,
		})

		return err
	})
}

func IndexHydraJobset(ctx context.Context, db *sql.DB, channel string, jobsetConfig *config.HydraJobset) (int, error) {
	l := log.WithFields(log.Fields{
		"channel": channel,
		"import":  "hydra",
	})

	l.Info("getting latest evaluation from database")

	latestEval, err := func() (idb.Hydraevaluation, error) {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return idb.Hydraevaluation{}, fmt.Errorf("error creating db transaction: %w", err)
		}

		defer func() {
			err := tx.Rollback()
			if err != nil && err != sql.ErrTxDone {
				panic(err)
			}
		}()

		queries := idb.New(db)
		qtx := queries.WithTx(tx)

		return qtx.GetLatestHydraEval(ctx, channel)
	}()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("error retrieving latest hydra eval: %w", err)
	}

	// are we creating completely from scratch?
	isNew := errors.Is(err, sql.ErrNoRows)

	l.Info("getting evaluations from hydra API")

	evalResp, err := GetEvaluations(jobsetConfig.BaseURL, jobsetConfig.Project, jobsetConfig.Jobset)
	if err != nil {
		return 0, fmt.Errorf("error getting response from Hydra at '%s': %w", jobsetConfig.BaseURL, err)
	}

	// Create a list of evaluations to index
	evals := []*HydraEval{}
	if isNew {
		evals = append(evals, evalResp.Evals[0])
	} else {
	idloop:
		for {
			for _, eval := range evalResp.Evals {
				if latestEval.HydraEvalID >= eval.ID {
					break idloop
				}

				evals = append(evals, eval)
			}

			evalResp, err = evalResp.NextPage()
			if err != nil {
				if err == io.EOF {
					break
				}

				return 0, fmt.Errorf("error getting response from Hydra at '%s': %w", jobsetConfig.BaseURL, err)
			}

		}
	}

	// Reverse the list so we create older missing evaluations first
	for i, j := 0, len(evals)-1; i < j; i, j = i+1, j-1 {
		evals[i], evals[j] = evals[j], evals[i]
	}

	l.Info("getting jobset from hydra API")
	hydraJobset, err := GetJobset(jobsetConfig.BaseURL, jobsetConfig.Project, jobsetConfig.Jobset)
	if err != nil {
		return 0, fmt.Errorf("error getting response from Hydra at '%s': %w", jobsetConfig.BaseURL, err)
	}

	for _, evalMeta := range evals {
		err = indexHydraJobset(ctx, db, channel, jobsetConfig.BaseURL, hydraJobset, evalMeta)
		if err != nil {
			return 0, fmt.Errorf("error indexing evaluation as a part of channel '%s': %w", channel, err)
		}
	}

	return len(evals), nil
}
