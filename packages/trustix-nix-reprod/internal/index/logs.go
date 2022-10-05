// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package index

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/nix-community/go-nix/pkg/nixpath"
	idb "github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/db"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix-proto/protocols"
	"github.com/nix-community/trustix/packages/trustix/client"
	logger "github.com/sirupsen/logrus"
)

// an arbitrary chunk size that's big enough
const logChunkSize = 500

// Filter only the Nix protocol
var logProtocols = []string{
	(func() string {
		pd, err := protocols.Get("nix")
		if err != nil {
			panic(err)
		}
		return pd.ID
	})(),
}

// index a single chunk of a log
func indexLogChunk(ctx context.Context, client *client.Client, log *api.Log, dbLog idb.Log, db *sql.DB, start uint64, finish uint64) error {
	if start >= finish {
		return nil
	}

	logger.WithFields(logger.Fields{
		"logID":  *log.LogID,
		"start":  start,
		"finish": finish,
	}).Debug("indexing log chunk")

	resp, err := client.LogAPI.GetLogEntries(ctx, &api.GetLogEntriesRequest{
		LogID:  log.LogID,
		Start:  &start,
		Finish: &finish,
	})
	if err != nil {
		return fmt.Errorf("error getting log entries: %w", err)
	}

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

	// create an entry for each leaf
	for _, leaf := range resp.Leaves {
		storePath := string(leaf.Key)

		if !strings.HasPrefix(storePath, nixpath.StoreDir) {
			continue
		}

		_, err := qtx.CreateDerivationOutputResult(
			ctx,
			idb.CreateDerivationOutputResultParams{
				OutputHash: base64.URLEncoding.EncodeToString(leaf.ValueDigest),
				StorePath:  string(leaf.Key),
				LogID:      dbLog.ID,
			},
		)

		if err != nil {
			return fmt.Errorf("could not create derivation output result: %w", err)
		}
	}

	// set new tree size
	{
		treeSize := int64(finish) + 1

		logger.WithFields(logger.Fields{
			"logID":    *log.LogID,
			"dbLogID":  dbLog.ID,
			"treeSize": treeSize,
		}).Debug("setting new tree size")

		err = qtx.SetTreeSize(ctx, idb.SetTreeSizeParams{
			TreeSize: treeSize,
			ID:       dbLog.ID,
		})
		if err != nil {
			return fmt.Errorf("error updating tree size: %w", err)
		}

		dbLog.TreeSize = treeSize
	}
	return tx.Commit()
}

// index full log
func indexLog(ctx context.Context, log *api.Log, dbLog idb.Log, client *client.Client, db *sql.DB) error {
	// Get the log head
	logHead, err := client.LogAPI.GetHead(ctx, &api.LogHeadRequest{
		LogID: log.LogID,
	})
	if err != nil {
		return fmt.Errorf("error getting log head: %w", err)
	}

	start := uint64(dbLog.TreeSize)
	if start > 0 {
		start += 1
	}
	finish := *logHead.TreeSize - 1

	if start >= finish {
		logger.WithFields(logger.Fields{
			"logID":  *log.LogID,
			"start":  start - 1,
			"finish": finish,
		}).Debug("log already up to date")

		return err
	}

	logger.WithFields(logger.Fields{
		"logID":        *log.LogID,
		"treeSizeDiff": finish - start,
		"start":        start,
		"finish":       finish,
	}).Debug("indexing log")

	// calculate request chunk boundaries
	chunks := []uint64{}
	{
		// chunk requests
		for i := uint64(dbLog.TreeSize); i <= *logHead.TreeSize; i += logChunkSize {
			chunks = append(chunks, i)
		}

		if chunks[len(chunks)-1] != finish {
			chunks = append(chunks, finish)
		}
	}

	// get log entries and index built outputs
	for _, finish := range chunks {
		err := indexLogChunk(ctx, client, log, dbLog, db, start, finish)
		if err != nil {
			return err
		}
	}

	return nil
}

func createLogs(ctx context.Context, db *sql.DB, logsResp *api.LogsResponse) (map[string]idb.Log, error) {
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

	// Map stringly logID to a Log instance
	logMap := make(map[string]idb.Log)

	logger.Info("indexing logs")

	// create non existing logs
	for _, log := range logsResp.Logs {
		logger.WithFields(logger.Fields{
			"logID": *log.LogID,
		}).Debug("trying to get log from database")

		dbLog, err := qtx.GetLog(ctx, *log.LogID)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("error getting log '%s' in db: %s", *log.LogID, err)
		} else if err == sql.ErrNoRows {
			logger.WithFields(logger.Fields{
				"logID": *log.LogID,
			}).Debug("log not found, creating database entry")

			dbLog, err = qtx.CreateLog(ctx, *log.LogID)
			if err != nil {
				return nil, fmt.Errorf("error creating log '%s' in db: %s", *log.LogID, err)
			}
		}

		logMap[*log.LogID] = dbLog
	}

	logger.Info("finished indexing logs")

	return logMap, tx.Commit()
}

func IndexLogs(ctx context.Context, db *sql.DB, client *client.Client) error {
	logsResp, err := client.NodeAPI.Logs(ctx, &api.LogsRequest{
		Protocols: logProtocols,
	})
	if err != nil {
		return fmt.Errorf("error getting logs list: %w", err)
	}

	// create logs that don't exist yet
	logMap, err := createLogs(ctx, db, logsResp)
	if err != nil {
		return nil
	}

	// index any logs that has updates
	for _, log := range logsResp.Logs {
		dbLog, ok := logMap[*log.LogID]
		if !ok {
			panic(fmt.Sprintf("expected to find log with id '%s'", *log.LogID))
		}

		err := indexLog(ctx, log, dbLog, client, db)
		if err != nil {
			return err
		}
	}

	return nil
}
