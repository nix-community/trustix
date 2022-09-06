// Copyright (c) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: MIT

//go:build tools

package tools

import (
	_ "github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
