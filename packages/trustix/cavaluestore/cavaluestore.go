// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cavaluestore

import (
	"crypto/sha256"

	"github.com/tweag/trustix/packages/trustix/storage"
)

const Bucket = "VALUES"

func Get(txn storage.Transaction, digest []byte) ([]byte, error) {
	return txn.Get([]byte(Bucket), digest)
}

func Set(txn storage.Transaction, value []byte) error {
	digest := sha256.Sum256(value)
	return txn.Set([]byte(Bucket), digest[:], value)
}
