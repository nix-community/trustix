// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fixturePath(name string) string {
	wd, _ := os.Getwd()
	return path.Join(wd, "fixtures", name)
}

func TestPathNotExists(t *testing.T) {
	_, err := NewConfigFromFile(fixturePath("notexist.toml"))
	assert.NotNil(t, err)
}

func TestConfig(t *testing.T) {
	f := fixturePath("example.toml")

	config, err := NewConfigFromFile(f)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(config.Publishers), "Exactly one publisher should have been returned")

	// log := config.Publishers[0]

	// assert.Equal(t, "trustix-test1", log.Name, "Unexpected name returned")

	// assert.Equal(t, "native", config.Storage.Type, "Unexpected storage type returned")

}
