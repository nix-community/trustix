// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package sql

import (
	"embed"
)

//go:embed schema/*.sql
var SchemaFS embed.FS
