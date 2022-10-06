package index

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/config"
	idb "github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/db"
	"github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/hydra"
)

func IndexChannel(ctx context.Context, db *sql.DB, channel string, channelConfig *config.Channel) error {
	switch channelConfig.Type {

	case "hydra":
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
			return fmt.Errorf("error retrieving latest hydra eval: %w", err)
		}

		// are we creating completely from scratch?
		isNew := errors.Is(err, sql.ErrNoRows)

		evalResp, err := hydra.GetEvaluations(channelConfig.Hydra.BaseURL, channelConfig.Hydra.Project, channelConfig.Hydra.Jobset)
		if err != nil {
			return fmt.Errorf("error getting response from Hydra at '%s': %w", channelConfig.Hydra.BaseURL, err)
		}

		// Create a list of evaluations to index
		evals := []*hydra.HydraEval{}
		if isNew {
			evals = append(evals, evalResp.Evals[0])
		} else {
			for {
				evalResp, err = evalResp.NextPage()
				if err != nil {
					if err == io.EOF {
						break
					}

					return fmt.Errorf("error getting response from Hydra at '%s': %w", channelConfig.Hydra.BaseURL, err)
				}

				for _, eval := range evalResp.Evals {
					if latestEval.HydraEvalID >= eval.ID {
						break
					}

					evals = append(evals)
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
				return fmt.Errorf("error getting nix path: %w", err)
			}

			err = IndexEval(ctx, db, nixPath, channel, timestamp, evalMeta)
			if err != nil {
				return fmt.Errorf("error indexing evaluation as a part of channel '%s': %w", channel, err)
			}
		}

		fmt.Println(evals)

	default:
		return fmt.Errorf("unhandled channel type: %s", channelConfig.Type)
	}

	return nil
}
