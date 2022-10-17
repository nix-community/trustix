// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package main // import "github.com/nix-community/trustix/packages/trustix-nix-r13y"

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/cmd"
)

func main() {
	cmd.Execute()
}
