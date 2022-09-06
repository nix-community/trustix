// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"database/sql"

	schema "github.com/nix-community/trustix/packages/trustix-nix-reprod/sql"
	"github.com/pressly/goose/v3"
)

func migrate(db *sql.DB, dialect string) error {
	goose.SetBaseFS(schema.SchemaFS)

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db, "schema"); err != nil {
		return err
	}

	return nil
}
