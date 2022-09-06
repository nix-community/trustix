// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

//go:build tools

package main

import (
	_ "github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go"
	_ "github.com/kyleconroy/sqlc/cmd/sqlc"
	_ "github.com/pressly/goose/cmd/goose"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
