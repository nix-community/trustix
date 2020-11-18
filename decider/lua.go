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
		state.PushString(in.LogName)
		state.SetField(-2, "LogName")
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

	state.Field(-1, "LogNames")
	if !state.IsTable(-1) {
		return nil, fmt.Errorf("LogNames is not of type table")
	}
	state.Length(-1)
	logNamesLen, ok := state.ToInteger(-1)
	state.Pop(1)
	if !ok {
		return nil, fmt.Errorf("OutputHash is not of type string")
	}
	for i := 0; i < logNamesLen; i++ {
		idx := i + 1 // Arrays start at 1
		state.RawGetInt(-1, idx)
		value, ok := state.ToString(-1)
		state.Pop(1)
		if !ok {
			return nil, fmt.Errorf("Member of OutputHash at index %d is not a string", idx)
		}
		ret.LogNames = append(ret.LogNames, value)
	}

	ret.Confidence = 1

	return ret, nil
}
