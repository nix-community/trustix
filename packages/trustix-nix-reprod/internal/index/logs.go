// Copyright (C) 2022 adisbladis
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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
