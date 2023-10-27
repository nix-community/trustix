// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	connect "connectrpc.com/connect"
	"github.com/nix-community/trustix/packages/go-lib/executor"
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

	attrsByChannel map[string][]string

	diffExecutor     *future.KeyedFutures[*pb.DiffResponse]
	downloadExecutor *future.KeyedFutures[*refcount.RefCountedValue[*narDownload]]
}

func NewAPIServer(db *sql.DB, cacheDB *sql.DB, cacheDBRO *sql.DB, client *client.Client, logNames map[string]string, attrsByChannel map[string][]string) *APIServer {
	return &APIServer{
		db:               db,
		client:           client,
		cacheDbRw:        cacheDB,
		cacheDbRo:        cacheDBRO,
		logNames:         logNames,
		attrsByChannel:   attrsByChannel,
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

	resp, err := s.attrReproducibilityTimeSeries(ctx, tx, req.Msg)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (s *APIServer) AttrReproducibilityTimeSeriesGroupedbyChannel(ctx context.Context, req *connect.Request[pb.AttrReproducibilityTimeSeriesGroupedbyChannelRequest]) (*connect.Response[pb.AttrReproducibilityTimeSeriesGroupedbyChannelResponse], error) {
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

	var start, stop int64
	{
		now := time.Now().UTC()
		stop = now.Unix()
		start = now.Add(-time.Hour * 24 * 90).Unix()
	}

	resp := &pb.AttrReproducibilityTimeSeriesGroupedbyChannelResponse{
		Channels: make(map[string]*pb.AttrReproducibilityTimeSeriesGroupedbyChannelResponse_Channel),
	}

	e := executor.NewParallellExecutor()

	for channel, attrs := range s.attrsByChannel {
		channel := channel
		attrs := attrs

		c := &pb.AttrReproducibilityTimeSeriesGroupedbyChannelResponse_Channel{
			Attrs: make(map[string]*pb.AttrReproducibilityTimeSeriesResponse),
		}
		resp.Channels[channel] = c

		var mux sync.Mutex

		for _, attr := range attrs {
			attr := attr

			err := e.Add(func() error {
				r, err := s.attrReproducibilityTimeSeries(ctx, tx, &pb.AttrReproducibilityTimeSeriesRequest{
					Attr:    attr,
					Start:   start,
					Stop:    stop,
					Channel: channel,
				})
				if err != nil {
					return fmt.Errorf("error getting reproducibility time series for attr '%s', channel '%s': %w", attr, channel, err)
				}

				mux.Lock()
				c.Attrs[attr] = r
				mux.Unlock()

				return nil
			})
			if err != nil {
				return nil, err
			}

		}
	}

	err = e.Wait()
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(resp), nil
}

func (s *APIServer) attrReproducibilityTimeSeries(ctx context.Context, tx *sql.Tx, msg *pb.AttrReproducibilityTimeSeriesRequest) (*pb.AttrReproducibilityTimeSeriesResponse, error) {
	attr := msg.Attr
	start := time.Unix(msg.Start, 0)
	stop := time.Unix(msg.Stop, 0)

	queries := idb.New(s.db)
	qtx := queries.WithTx(tx)

	rows, err := qtx.GetDerivationReproducibilityTimeSeriesByAttr(ctx, idb.GetDerivationReproducibilityTimeSeriesByAttrParams{
		Attr:        attr,
		Timestamp:   start,
		Timestamp_2: stop,
		Channel:     msg.Channel,
	})
	if err != nil {
		return nil, fmt.Errorf("error retreiving time series rows for attr '%s': %w", attr, err)
	}

	resp := &pb.AttrReproducibilityTimeSeriesResponse{
		Points: []*pb.AttrReproducibilityTimeSeriesPoint{},
	}

	numReproduced := 0
	numDerivations := 0
	points := 0

	for i, row := range rows {
		numDerivations++

		if row.OutputHashCount == 1 && row.ResultCount >= 2 {
			numReproduced++
		}

		// Return rows grouped by their derivation
		if i+1 == len(rows) || rows[i+1].Drv != row.Drv {
			pctReproduced := 100 / float32(numDerivations) * float32(numReproduced)

			resp.Points = append(resp.Points, &pb.AttrReproducibilityTimeSeriesPoint{
				EvalID:        row.EvalID,
				EvalTimestamp: row.EvalTimestamp.Unix(),
				DrvPath:       row.Drv,
				PctReproduced: pctReproduced,
			})

			numDerivations = 0
			numReproduced = 0

			points++

			resp.PctReproduced += pctReproduced
		}
	}

	if len(rows) > 0 {
		resp.PctReproduced = resp.PctReproduced / float32(points)
	}

	return resp, nil
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
