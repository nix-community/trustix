// MIT License
//
// Copyright (c) 2020 Tweag IO
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

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

	assert.Equal(t, 1, len(config.Logs), "Exactly one log should have been returned")

	log := config.Logs[0]

	assert.Equal(t, "trustix-test1", log.Name, "Unexpected name returned")

	assert.Equal(t, "native", log.Storage.Type, "Unexpected storage type returned")

}
