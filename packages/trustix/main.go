// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package main // import "github.com/nix-community/trustix/packages/trustix/internal"

import (
	"runtime"

	"github.com/nix-community/trustix/packages/trustix/internal/cmd"
)

func main() {
	runtime.SetBlockProfileRate(100)
	runtime.GOMAXPROCS(128)

	cmd.Execute()
}
