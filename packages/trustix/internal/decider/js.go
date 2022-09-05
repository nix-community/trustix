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

	"github.com/dop251/goja"
)

type jsDecider struct {
	fn string
}

func NewJSDecider(function string) (LogDecider, error) {
	// Wrap function
	function = "(" + function + ")"

	// Sanity check input
	vm := goja.New()
	value, err := vm.RunString(function)
	if err != nil {
		return nil, err
	}

	_, ok := goja.AssertFunction(value)
	if !ok {
		return nil, fmt.Errorf("Script '%s' does not evaluate to a function", function)
	}

	return &jsDecider{
		fn: function,
	}, nil
}

func (j *jsDecider) Name() string {
	return "javascript"
}

func (j *jsDecider) Decide(inputs []*DeciderInput) (*DeciderOutput, error) {
	vm := goja.New()
	value, err := vm.RunString(j.fn)
	if err != nil {
		return nil, err
	}

	objects := make([]*goja.Object, len(inputs))
	for i, input := range inputs {
		obj := vm.NewObject()

		err = obj.Set("LogID", input.LogID)
		if err != nil {
			return nil, fmt.Errorf("error setting VM object value: %w", err)
		}

		err = obj.Set("Value", input.Value)
		if err != nil {
			return nil, fmt.Errorf("error setting VM object value: %w", err)
		}

		objects[i] = obj
	}

	arr := vm.NewArray(objects)

	fn, _ := goja.AssertFunction(value)
	ret, err := fn(goja.Undefined(), arr)
	if err != nil {
		return nil, err
	}

	return &DeciderOutput{
		Value:      ret.String(),
		Confidence: 1,
	}, nil
}
