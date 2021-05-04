// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package storage

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/tweag/trustix/packages/trustix-proto/schema"
	"github.com/tweag/trustix/packages/trustix/constants"
)

// TODO: I don't like this living here but I also don't have a better option

func GetSTH(txn *BucketTransaction) (*schema.STH, error) {
	var buf []byte
	{
		v, err := txn.Get([]byte(constants.HeadBlob))
		if err != nil {
			return nil, err
		}
		buf = v
	}
	if len(buf) == 0 {
		return nil, ObjectNotFoundError
	}

	sth := &schema.STH{}
	err := proto.Unmarshal(buf, sth)
	if err != nil {
		return nil, err
	}

	return sth, nil
}
