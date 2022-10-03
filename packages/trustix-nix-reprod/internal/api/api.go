package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	idb "github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/db"
	pb "github.com/nix-community/trustix/packages/trustix-nix-reprod/reprod-api"
)

const nullJSONGroupObjectString = "{:null}"

type respDerivation = pb.DerivationReproducibilityResponse_Derivation
type respDerivationOutput = pb.DerivationReproducibilityResponse_DerivationOutput
type respDerivationOutputHash = pb.DerivationReproducibilityResponse_DerivationOutputHash

func getFirstMapKey(m map[string][]int) string {
	for k := range m {
		return k
	}

	panic("No key found")
}

func toInt64(s []int) []int64 {
	x := make([]int64, len(s))

	for i, v := range s {
		x[i] = int64(v)
	}

	return x

}

func GetDerivationReproducibility(ctx context.Context, db *sql.DB, drvPath string) (*pb.DerivationReproducibilityResponse, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating db transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			panic(err)
		}
	}()

	queries := idb.New(db)
	qtx := queries.WithTx(tx)

	rows, err := qtx.GetDerivationReproducibility(ctx, drvPath)
	if err != nil {
		return nil, fmt.Errorf("error retrieving rows: %w", err)
	}

	resp := &pb.DerivationReproducibilityResponse{
		MissingPaths:      make(map[string]*respDerivation),
		ReproducedPaths:   make(map[string]*respDerivation),
		UnknownPaths:      make(map[string]*respDerivation),
		UnreproducedPaths: make(map[string]*respDerivation),
	}

	appendOutput := func(drvs map[string]*respDerivation, row idb.GetDerivationReproducibilityRow, outputHashes map[string][]int) {
		drv, ok := drvs[row.Drv]
		if !ok {
			drv = &respDerivation{
				Outputs: make(map[string]*respDerivationOutput),
			}

			drvs[row.Drv] = drv
		}

		drvOutput, ok := drv.Outputs[row.Drv]
		if !ok {
			drvOutput = &respDerivationOutput{
				Output:       row.Output,
				StorePath:    row.StorePath,
				OutputHashes: make(map[string]*respDerivationOutputHash, len(outputHashes)),
			}

			for outputHash, logIDs := range outputHashes {
				out := &respDerivationOutputHash{
					LogIDs: toInt64(logIDs),
				}

				drvOutput.OutputHashes[outputHash] = out
			}

			drv.Outputs[row.Drv] = drvOutput
		}
	}

	for _, row := range rows {
		// Decode output hashes from aggregate JSON object from SQLite
		outputHashes := make(map[string][]int)
		{
			outputHashesString := row.OutputResults.(string)

			if outputHashesString != nullJSONGroupObjectString {
				outputHashesObj := make(map[int]string)

				err = json.Unmarshal([]byte(outputHashesString), &outputHashesObj)
				if err != nil {
					return nil, fmt.Errorf("couldnt decode JSON result object: %w", err)
				}

				for logID, outputHash := range outputHashesObj {
					outputHashes[outputHash] = append(outputHashes[outputHash], logID)
				}
			}
		}

		// fmt.Println(row)

		if len(outputHashes) < 1 {
			appendOutput(resp.MissingPaths, row, outputHashes)
		} else if len(outputHashes) == 1 && len(outputHashes[getFirstMapKey(outputHashes)]) > 1 {
			appendOutput(resp.ReproducedPaths, row, outputHashes)
		} else if len(outputHashes) == 1 {
			appendOutput(resp.UnknownPaths, row, outputHashes)
		} else if len(outputHashes) > 1 {
			appendOutput(resp.UnreproducedPaths, row, outputHashes)
		} else {
			panic("logic error")
		}
	}

	return resp, nil
}

func GetDerivationReproducibilityTimeSeriesByAttr(ctx context.Context, db *sql.DB) error {
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

	fmt.Println(qtx)

	return nil
}
