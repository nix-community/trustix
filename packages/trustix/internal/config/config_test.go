// Copyright (C) 2021 Tweag IO
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

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
