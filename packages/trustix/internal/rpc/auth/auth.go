// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package auth

import (
	"context"
)

func CanWrite(ctx context.Context) error {
	// TODO: Introduce a concept of an auth token to write
	return nil
}
