// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package index

import (
	"context"
	"database/sql"
	"fmt"

	idb "github.com/nix-community/trustix/packages/trustix-nix-reprod/internal/db"
	"github.com/nix-community/trustix/packages/trustix-proto/api"
	"github.com/nix-community/trustix/packages/trustix/client"
	"golang.org/x/exp/constraints"
)

var logProtocols = []string{}

func max[T constraints.Ordered](x T, y T) T {
	if x > y {
		return x
	}
	return y
}

func IndexLogs(ctx context.Context, db *sql.DB, client *client.Client) error {
	logsResp, err := client.NodeAPI.Logs(ctx, &api.LogsRequest{
		Protocols: logProtocols,
	})
	if err != nil {
		return fmt.Errorf("error getting logs list: %w", err)
	}

	tx, err := db.Begin()
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

	// Map stringly logID to a Log instance
	logMap := make(map[string]idb.Log)

	// create non existing logs
	for _, log := range logsResp.Logs {
		dbLog, err := qtx.CreateLog(ctx, *log.LogID)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("error creating log '%s' in db: %s", log.LogID, err)
		}

		if err == sql.ErrNoRows {
			dbLog, err = qtx.GetLog(ctx, *log.LogID)

			if err != nil {
				return fmt.Errorf("error getting log '%s' in db: %s", log.LogID, err)
			}
		}

		logMap[*log.LogID] = dbLog
	}

	// index any logs that has updates
	for _, log := range logsResp.Logs {
		// Get the log head
		logHead, err := client.LogAPI.GetHead(ctx, &api.LogHeadRequest{
			LogID: log.LogID,
		})
		if err != nil {
			return fmt.Errorf("error getting logs list: %w", err)
		}

		dbLog := logMap[*log.LogID]

		if uint64(dbLog.TreeSize) >= *logHead.TreeSize {
			continue // return nil
		}
	}

	return tx.Commit()
}
