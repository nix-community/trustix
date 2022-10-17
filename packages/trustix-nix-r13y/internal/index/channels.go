// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package index

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/config"
	idb "github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/db"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/hydra"
	log "github.com/sirupsen/logrus"
)

func IndexHydraJobset(ctx context.Context, db *sql.DB, channel string, jobsetConfig *config.HydraJobset) (int, error) {
	l := log.WithFields(log.Fields{
		"channel": channel,
	})

	// Hydra
	{

		l = l.WithFields(log.Fields{
			"channelType": "hydra",
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

		evalResp, err := hydra.GetEvaluations(jobsetConfig.BaseURL, jobsetConfig.Project, jobsetConfig.Jobset)
		if err != nil {
			return 0, fmt.Errorf("error getting response from Hydra at '%s': %w", jobsetConfig.BaseURL, err)
		}

		// Create a list of evaluations to index
		evals := []*hydra.HydraEval{}
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

		for _, evalMeta := range evals {
			timestamp := time.Unix(evalMeta.Timestamp, 0)

			nixPath, err := evalMeta.NixPath()
			if err != nil {
				return 0, fmt.Errorf("error getting nix path: %w", err)
			}

			err = IndexEval(ctx, db, nixPath, channel, timestamp, evalMeta)
			if err != nil {
				return 0, fmt.Errorf("error indexing evaluation as a part of channel '%s': %w", channel, err)
			}
		}

		return len(evals), nil
	}
}
