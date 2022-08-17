//go:build tools

package main

import (
	_ "github.com/kyleconroy/sqlc/cmd/sqlc"
	_ "github.com/pressly/goose/cmd/goose"
)
