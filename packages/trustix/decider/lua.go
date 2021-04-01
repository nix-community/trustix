// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package decider

import (
	"fmt"

	"github.com/Shopify/go-lua"
)

type luaDecider struct {
	state  *lua.State
	script string
}

const luaScriptWrapper = `
  DeciderFunction = (%s)
`

func NewLuaDecider(function string) (LogDecider, error) {
	state := lua.NewState()
	lua.OpenLibraries(state)

	script := fmt.Sprintf(luaScriptWrapper, function)
	if err := lua.DoString(state, script); err != nil {
		return nil, err
	}

	// Initial state is function call
	state.Global("DeciderFunction")

	return &luaDecider{
		state: state,
	}, nil
}

func (l *luaDecider) Name() string {
	return "lua"
}

func (l *luaDecider) Decide(inputs []*LogDeciderInput) (*LogDeciderOutput, error) {
	state := l.state

	// TODO: Recover from panic() in call

	// Set state
	state.Global("DeciderFunction")

	// Reset state to function call on exit
	defer state.Global("DeciderFunction")

	// Create corresponding []*LogDeciderInput
	state.NewTable()

	// Create corresponding *LogDeciderInput
	for i, in := range inputs {
		state.NewTable()
		state.PushString(in.LogID)
		state.SetField(-2, "LogID")
		state.PushString(in.OutputHash)
		state.SetField(-2, "OutputHash")

		idx := i + 1             // In Lua arrays start at 1...
		state.RawSetInt(-2, idx) // Append to list
	}

	// Call the function
	numArgs := 1
	numResults := 1
	state.Call(numArgs, numResults)

	if !state.IsTable(-1) {
		return nil, fmt.Errorf("Return of function is not a table")
	}

	// Translate lua table back to Go struct
	ret := &LogDeciderOutput{}

	state.Field(-1, "OutputHash")
	outputHash, ok := state.ToString(-1)
	state.Pop(1)
	if !ok {
		return nil, fmt.Errorf("OutputHash is not of type string")
	}
	ret.OutputHash = outputHash

	ret.Confidence = 1

	return ret, nil
}
