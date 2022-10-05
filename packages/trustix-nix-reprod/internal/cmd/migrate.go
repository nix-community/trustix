// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"database/sql"
	"fmt"
	"path/filepath"

	schema "github.com/nix-community/trustix/packages/trustix-nix-reprod/sql"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
)

const sqlDialect = "sqlite"

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

func setupDB(stateDirectory string) (*sql.DB, error) {
	dbPath := filepath.Join(stateDirectory, "db.sqlite3?_journal_mode=WAL")

	l := log.WithFields(log.Fields{
		"path": dbPath,
	})

	l.Info("Opening database")

	db, err := sql.Open(sqlDialect, dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	l.Info("Migrating database")

	err = migrate(db, sqlDialect)
	if err != nil {
		return nil, fmt.Errorf("error migrating database: %w", err)
	}

	return db, nil
}
