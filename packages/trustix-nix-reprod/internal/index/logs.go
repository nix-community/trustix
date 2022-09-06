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
)

func IndexLogs(ctx context.Context, db *sql.DB, client *client.Client) error {
	logsResp, err := client.NodeAPI.Logs(ctx, &api.LogsRequest{})
	if err != nil {
		return fmt.Errorf("error getting logs list: %w", err)
	}

	fmt.Println(logsResp)

	if true {
		return nil
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

	fmt.Println(qtx)

	return nil
}
