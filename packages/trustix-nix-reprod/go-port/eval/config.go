// Copyright (C) 2022 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package eval

import (
	"fmt"
	"github.com/pbnjay/memory"
	"strconv"
)

type EvalConfig struct {
	Expr   string
	Arg    map[string]string
	ArgStr map[string]string

	Flake string

	Workers       int
	MaxMemorySize int

	ResultBufferSize int
}

func NewConfig() *EvalConfig {
	workers := 4
	maxMemorySize := 0

	// Max a rough estimate at a reasonable default evaluator memory usage
	// Each individual worker can still go above the threshold, it's merely the point
	// where they are restarted.
	// This attempts to give each worker about equal ram.
	freeMem := memory.FreeMemory()
	if freeMem > 0 {
		maxMemorySize = int(float64(freeMem/uint64(workers))*0.9) / 1_000_000
	}

	return &EvalConfig{
		Arg:              make(map[string]string),
		ArgStr:           make(map[string]string),
		Workers:          workers,
		MaxMemorySize:    maxMemorySize,
		ResultBufferSize: 1024,
	}
}

func (c *EvalConfig) AddArg(arg string, value string) {
	delete(c.Arg, arg)
	delete(c.ArgStr, arg)
	c.Arg[arg] = value
}

func (c *EvalConfig) AddArgStr(arg string, value string) {
	delete(c.Arg, arg)
	delete(c.ArgStr, arg)
	c.ArgStr[arg] = value
}

func (c *EvalConfig) toArgs() ([]string, error) {

	args := []string{
		"--quiet",
	}

	if c.Flake != "" {
		args = append(args, "--flake", c.Flake)
	}

	for arg, value := range c.Arg {
		args = append(args, "--arg", arg, value)
	}
	for arg, value := range c.ArgStr {
		args = append(args, "--argstr", arg, value)
	}

	if c.Workers > 0 {
		args = append(args, "--workers", strconv.Itoa(c.Workers))
	}

	if c.MaxMemorySize > 0 {
		args = append(args, "--max-memory-size", strconv.Itoa(c.MaxMemorySize))
	}

	if c.Expr == "" {
		return nil, fmt.Errorf("Missing expression to evaluate")
	}
	args = append(args, c.Expr)

	return args, nil
}
