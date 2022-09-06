// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package config

import (
	"fmt"
)

func missingField(field string) error {
	return fmt.Errorf("Required field '%s' is missing", field)
}
