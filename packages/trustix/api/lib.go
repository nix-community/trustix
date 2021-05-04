// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package api

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/tweag/trustix/packages/trustix/storage"
)

func getCAValue(txn *storage.BucketTransaction, digest []byte) ([]byte, error) {
	value, err := txn.Get(digest)
	if err != nil {
		return nil, err
	}

	digest2 := sha256.Sum256(value)
	if !bytes.Equal(digest, digest2[:]) {
		return nil, fmt.Errorf("Digest mismatch")
	}

	return value, nil
}
