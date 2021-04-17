// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const logIDLength = 64

// LogID - Generate a deterministic ID based on known facts
func LogID(keyType string, publicKey []byte) string {
	h := sha256.New()

	h.Write([]byte(keyType))
	h.Write([]byte(":"))

	h.Write(publicKey)
	h.Write([]byte(":"))

	return hex.EncodeToString(h.Sum(nil))
}

func ValidLogID(logID string) error {
	if logID == "" {
		return fmt.Errorf("Empty logID")
	}

	if len(logID) != logIDLength {
		return fmt.Errorf("LogID length not correct: %d != %d", len(logID), logIDLength)
	}

	return nil
}
