// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package nar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testInput = `
StorePath: /nix/store/byghkc1k0xmrrl2jk04lp0qipmpmz547-hello-2.10
URL: nar/06vv5hjmdyrklwsxxq5d4fnslkgfzpy3z3ri4s7a9fawi2d20ivb.nar.xz
Compression: xz
FileHash: sha256:06vv5hjmdyrklwsxxq5d4fnslkgfzpy3z3ri4s7a9fawi2d20ivb
FileSize: 41272
NarHash: sha256:1llxabk0xq0gc15yi6kkysfbvn5gzisj9dxk6g29sh5ncqx3if8y
NarSize: 206000
References: a6rnjp15qgp8a699dlffqj94hzy1nldg-glibc-2.32 byghkc1k0xmrrl2jk04lp0qipmpmz547-hello-2.10
Deriver: m0i10ghpcwhi2dml0dj6b437jjrh8ia3-hello-2.10.drv
Sig: cache.nixos.org-1:IHkSz9VMQC/KGYgah2Vr2ISz0uawXUKqm/yP4JtcaBkyLO13B3yD2k578ZrP3RyJIVyvMdn4KOjUxvCungdkDA==
`

func TestParseNarInfo(t *testing.T) {

	assert := assert.New(t)

	n, err := ParseNarInfo(testInput)
	assert.Nil(err)
	assert.NotNil(n)

	assert.Equal("/nix/store/byghkc1k0xmrrl2jk04lp0qipmpmz547-hello-2.10", *n.StorePath)
	assert.Equal("nar/06vv5hjmdyrklwsxxq5d4fnslkgfzpy3z3ri4s7a9fawi2d20ivb.nar.xz", *n.URL)
	assert.Equal("xz", *n.Compression)
	assert.Equal("sha256:06vv5hjmdyrklwsxxq5d4fnslkgfzpy3z3ri4s7a9fawi2d20ivb", *n.FileHash)
	assert.Equal(uint64(41272), *n.FileSize)
	assert.Equal("sha256:1llxabk0xq0gc15yi6kkysfbvn5gzisj9dxk6g29sh5ncqx3if8y", *n.NarHash)
	assert.Equal(uint64(206000), *n.NarSize)
	assert.Equal(2, len(n.References))

}
