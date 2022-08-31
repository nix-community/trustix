// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package log

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/nix-community/trustix/packages/trustix-proto/schema"
)

func NewLeaf(digest []byte, value []byte) (*schema.LogLeaf, error) {
	return &schema.LogLeaf{
		LeafDigest: digest,
	}, nil
}

func LeafFromBytes(data []byte) (*schema.LogLeaf, error) {
	l := &schema.LogLeaf{}
	err := proto.Unmarshal(data, l)
	if err != nil {
		return nil, err
	}
	return l, nil
}
