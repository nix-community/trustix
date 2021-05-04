// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package storage

type Transaction interface {
	Get(bucket []byte, key []byte) ([]byte, error)
	Set(bucket []byte, key []byte, value []byte) error
	Delete(bucket []byte, key []byte) error
}

type Storage interface {
	Close()

	// View - Start a read-only transaction
	View(func(txn Transaction) error) error

	// Update - Start a read-write transaction
	Update(func(txn Transaction) error) error
}
