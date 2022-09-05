// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

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
