package cmd

import (
	"database/sql"

	"github.com/pressly/goose/v3"
	schema "github.com/tweag/trustix/packages/trustix-nix-reprod/sql"
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
