// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package nar

import (
	"fmt"
	"strconv"
	"strings"
)

type NarInfo struct {
	StorePath   string
	URL         string
	Compression string
	FileHash    string
	FileSize    uint64
	NarHash     string
	NarSize     uint64
	References  []string
	Deriver     string
	Sig         string
	System      string
	CA          string
}

func ParseNarInfo(input []byte) (*NarInfo, error) {

	n := &NarInfo{}

	parseUint := func(in string) (uint64, error) {
		return strconv.ParseUint(in, 10, 64)
	}

	var err error

	for _, line := range strings.Split(string(input), "\n") {
		tok := strings.SplitN(line, ": ", 2)
		if len(tok) <= 1 {
			continue
		}

		if len(tok) != 2 {
			return nil, fmt.Errorf("Unexpected number of tokens '%d' for value '%s'", len(tok), tok)
		}
		value := tok[1]

		switch tok[0] {
		case "StorePath":
			n.StorePath = value
		case "URL":
			n.URL = value
		case "Compression":
			n.Compression = value
		case "FileHash":
			n.FileHash = value
		case "FileSize":
			n.FileSize, err = parseUint(value)
			if err != nil {
				return nil, err
			}
		case "NarHash":
			n.NarHash = value
		case "NarSize":
			n.NarSize, err = parseUint(value)
			if err != nil {
				return nil, err
			}
		case "References":
			n.References = strings.Split(value, " ")
		case "Deriver":
			n.Deriver = value
		case "Sig":
			n.Sig = value
		case "System":
			n.System = value
		case "CA":
			n.CA = value
		default:
			return nil, fmt.Errorf("Unknown field: %s", tok[0])
		}
	}

	return n, nil
}
