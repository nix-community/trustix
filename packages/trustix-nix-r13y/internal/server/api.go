// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	connect "github.com/bufbuild/connect-go"
	"github.com/nix-community/trustix/packages/go-lib/set"
	idb "github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/db"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/future"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/refcount"
	pb "github.com/nix-community/trustix/packages/trustix-nix-r13y/reprod-api"
	apiconnect "github.com/nix-community/trustix/packages/trustix-nix-r13y/reprod-api/reprod_apiconnect"
	"github.com/nix-community/trustix/packages/trustix/client"
)

const nullJSONGroupObjectString = "{:null}"

type respDerivation = pb.DerivationReproducibilityResponse_Derivation
type respDerivationOutput = pb.DerivationReproducibilityResponse_DerivationOutput
type respDerivationOutputHash = pb.DerivationReproducibilityResponse_DerivationOutputHash

type APIServer struct {
	apiconnect.UnimplementedReproducibilityAPIHandler

	client *client.Client

	db        *sql.DB
	cacheDbRo *sql.DB
	cacheDbRw *sql.DB

	logNames map[string]string

	diffExecutor     *future.KeyedFutures[*pb.DiffResponse]
	downloadExecutor *future.KeyedFutures[*refcount.RefCountedValue[*narDownload]]
}

func NewAPIServer(db *sql.DB, cacheDB *sql.DB, cacheDBRO *sql.DB, client *client.Client, logNames map[string]string) *APIServer {
	return &APIServer{
		db:               db,
		client:           client,
		cacheDbRw:        cacheDB,
		cacheDbRo:        cacheDBRO,
		logNames:         logNames,
		diffExecutor:     future.NewKeyedFutures[*pb.DiffResponse](),
		downloadExecutor: future.NewKeyedFutures[*refcount.RefCountedValue[*narDownload]](),
	}
}

func getFirstMapKey(m map[string][]string) string {
	for k := range m {
		return k
	}

	panic("No key found")
}

func (s *APIServer) DerivationReproducibility(ctx context.Context, req *connect.Request[pb.DerivationReproducibilityRequest]) (*connect.Response[pb.DerivationReproducibilityResponse], error) {
	msg := req.Msg
	drvPath := msg.DrvPath

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating db transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			panic(err)
		}
	}()

	queries := idb.New(s.db)
	qtx := queries.WithTx(tx)

	rows, err := qtx.GetDerivationReproducibility(ctx, drvPath)
	if err != nil {
		return nil, fmt.Errorf("error retrieving rows: %w", err)
	}

	resp := &pb.DerivationReproducibilityResponse{
		DrvPath:           drvPath,
		MissingPaths:      make(map[string]*respDerivation),
		ReproducedPaths:   make(map[string]*respDerivation),
		UnknownPaths:      make(map[string]*respDerivation),
		UnreproducedPaths: make(map[string]*respDerivation),
		Logs:              make(map[string]*pb.Log),
	}

	logIDSet := set.NewSet[string]()

	appendOutput := func(drvs map[string]*respDerivation, row idb.GetDerivationReproducibilityRow, outputHashes map[string][]string) {
		drv, ok := drvs[row.Drv]
		if !ok {
			drv = &respDerivation{
				Outputs: make(map[string]*respDerivationOutput),
			}

			drvs[row.Drv] = drv
		}

		_, ok = drv.Outputs[row.Output]
		if !ok {
			drvOutput := &respDerivationOutput{
				Output:       row.Output,
				StorePath:    row.StorePath,
				OutputHashes: make(map[string]*respDerivationOutputHash, len(outputHashes)),
			}

			for outputHash, logIDs := range outputHashes {
				out := &respDerivationOutputHash{
					LogIDs: logIDs,
				}

				drvOutput.OutputHashes[outputHash] = out

				// Collect all log ids used in response so we can aggregate it on the response object later
				for _, logID := range logIDs {
					logIDSet.Add(logID)
				}
			}

			drv.Outputs[row.Output] = drvOutput
		}
	}

	for _, row := range rows {
		// Decode output hashes from aggregate JSON object from SQLite
		outputHashes := make(map[string][]string)
		{
			outputHashesString := row.OutputResults.(string)

			if outputHashesString != nullJSONGroupObjectString {
				outputHashesObj := make(map[string]string)

				err = json.Unmarshal([]byte(outputHashesString), &outputHashesObj)
				if err != nil {
					return nil, fmt.Errorf("couldnt decode JSON result object: %w", err)
				}

				for logID, outputHash := range outputHashesObj {
					outputHashes[outputHash] = append(outputHashes[outputHash], logID)
				}
			}
		}

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

	for _, logID := range logIDSet.Values() {
		// Is there a human friendly name configured?
		// If there is set that as the name otherwise fall back to the raw ID
		name, ok := s.logNames[logID]
		if !ok {
			name = logID
		}

		resp.Logs[logID] = &pb.Log{
			LogID: logID,
			Name:  name,
		}
	}

	return connect.NewResponse(resp), nil
}

func (s *APIServer) AttrReproducibilityTimeSeries(ctx context.Context, req *connect.Request[pb.AttrReproducibilityTimeSeriesRequest]) (*connect.Response[pb.AttrReproducibilityTimeSeriesResponse], error) {
	msg := req.Msg

	attr := msg.Attr
	start := time.Unix(msg.Start, 0)
	stop := time.Unix(msg.Stop, 0)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating db transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			panic(err)
		}
	}()

	queries := idb.New(s.db)
	qtx := queries.WithTx(tx)

	rows, err := qtx.GetDerivationReproducibilityTimeSeriesByAttr(ctx, idb.GetDerivationReproducibilityTimeSeriesByAttrParams{
		Attr:        attr,
		Timestamp:   start,
		Timestamp_2: stop,
		Channel:     msg.Channel,
	})
	if err != nil {
		return nil, fmt.Errorf("error retreiving time series rows: %w", err)
	}

	resp := &pb.AttrReproducibilityTimeSeriesResponse{
		Points: make([]*pb.AttrReproducibilityTimeSeriesPoint, len(rows)),
	}

	for i, row := range rows {
		// out of all built paths, how many were reproduced
		var pctReproduced float32
		if row.OutputHashCount > 0 {
			pctReproduced = (100 / float32(row.OutputHashCount)) * float32(row.StorePathCount)
		} else {
			pctReproduced = 0.0
		}

		// out of the total amount of paths, how many were reproduced
		pctReproducedCum := 100 / float32(row.OutputCount) * (pctReproduced / 100 * float32(row.StorePathCount))

		resp.Points[i] = &pb.AttrReproducibilityTimeSeriesPoint{
			EvalID:        row.EvalID,
			EvalTimestamp: row.EvalTimestamp.Unix(),
			DrvPath:       row.Drv,
			PctReproduced: pctReproducedCum,
		}

		resp.PctReproduced += pctReproducedCum
	}

	if len(rows) > 0 {
		resp.PctReproduced = resp.PctReproduced / float32(len(rows))
	}

	return connect.NewResponse(resp), nil
}

func (s *APIServer) SuggestAttribute(ctx context.Context, req *connect.Request[pb.SuggestsAttributeRequest]) (*connect.Response[pb.SuggestAttributeResponse], error) {
	msg := req.Msg
	attrPrefix := msg.AttrPrefix

	if len(attrPrefix) < 3 {
		return nil, fmt.Errorf("attribute prefix '%s' is too short (minimum 3)", attrPrefix)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating db transaction: %w", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			panic(err)
		}
	}()

	queries := idb.New(s.db)
	qtx := queries.WithTx(tx)

	suggestions, err := qtx.SuggestAttribute(ctx, attrPrefix+"%")
	if err != nil {
		return nil, fmt.Errorf("error retreiving suggested attributes: %w", err)
	}

	resp := &pb.SuggestAttributeResponse{
		Attrs: make([]string, len(suggestions)),
	}

	copy(resp.Attrs, suggestions)

	return connect.NewResponse(resp), nil
}
