// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package storage

import (
	"errors"
	"fmt"
)

var ObjectNotFoundError = errors.New("could not find object")

// Factory function to create a nice error message that contains the key
func objectNotFoundError(key []byte) error {
	return fmt.Errorf("error retreiving object with key '%v': %w", key, ObjectNotFoundError)
}
