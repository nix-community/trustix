// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package cmd

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"

	schema "github.com/nix-community/trustix/packages/trustix-nix-reprod/sql"
	cache_schema "github.com/nix-community/trustix/packages/trustix-nix-reprod/sql-cache"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
)

const sqlDialect = "sqlite3"
const dbConnectionString = "?_journal_mode=WAL"
const dbFilename = "db.sqlite3"
const cacheDBFilename = "cachedb.sqlite3"

type databases struct {
	dbRW      *sql.DB
	dbRO      *sql.DB
	cacheDbRW *sql.DB
	cacheDbRO *sql.DB
}

func setupDatabases(stateDirectory string) (*databases, error) {
	err := os.MkdirAll(stateDirectory, 0755)
	if err != nil {
		return nil, fmt.Errorf("error creating state directory: %w", err)
	}

	db, err := setupDB(stateDirectory, dbFilename, schema.SchemaFS, false)
	if err != nil {
		return nil, fmt.Errorf("error opening rw database: %w", err)
	}

	dbRO, err := setupDB(stateDirectory, dbFilename, schema.SchemaFS, true)
	if err != nil {
		return nil, fmt.Errorf("error opening ro database: %w", err)
	}

	cacheDB, err := setupDB(stateDirectory, cacheDBFilename, cache_schema.SchemaFS, false)
	if err != nil {
		return nil, fmt.Errorf("error opening rw cache database: %w", err)
	}

	cacheDBRO, err := setupDB(stateDirectory, cacheDBFilename, cache_schema.SchemaFS, true)
	if err != nil {
		return nil, fmt.Errorf("error opening ro cache database: %w", err)
	}

	return &databases{
		dbRW:      db,
		dbRO:      dbRO,
		cacheDbRW: cacheDB,
		cacheDbRO: cacheDBRO,
	}, nil
}

func migrateDB(db *sql.DB, fs embed.FS) error {
	goose.SetBaseFS(fs)

	if err := goose.SetDialect(sqlDialect); err != nil {
		return err
	}

	if err := goose.Up(db, "schema"); err != nil {
		return err
	}

	return nil
}

func setupDB(stateDirectory string, filename string, migrationFS embed.FS, readonly bool) (*sql.DB, error) {
	dbPath := "file:" + filepath.Join(stateDirectory, filename+dbConnectionString)
	if readonly {
		dbPath += "&mode=ro"
	} else {
		dbPath += "&mode=rwc"
	}

	l := log.WithFields(log.Fields{
		"path": dbPath,
	})

	l.Info("Opening database")

	db, err := sql.Open(sqlDialect, dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(1)

	if !readonly {
		l.Info("Migrating database")

		err = migrateDB(db, migrationFS)
		if err != nil {
			return nil, fmt.Errorf("error migrating database: %w", err)
		}
	}

	return db, nil
}
