// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package eval

import (
	"context"
	"github.com/nix-community/trustix/packages/trustix-nix-r13y/internal/lib"
)

type Evaluator func(context.Context, *EvalConfig) (chan *lib.Result[*EvalResult], error)
