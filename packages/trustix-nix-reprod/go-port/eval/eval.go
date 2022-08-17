// Copyright (C) 2022 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package eval

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"

	"github.com/buger/jsonparser"
	"github.com/tweag/trustix/packages/trustix-nix-reprod/lib"
)

type EvalResult struct {
	Attr     string
	AttrPath []string
	Error    string
	DrvPath  string
	Name     string
	System   string
	Outputs  map[string]string
}

var jsonPaths = [][]string{
	{"attr"},
	{"attrPath"},
	{"error"},
	{"drvPath"},
	{"name"},
	{"system"},
	{"outputs"},
}

func Eval(ctx context.Context, config *EvalConfig) (chan *lib.Result[*EvalResult], error) {

	args, err := config.toArgs()
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "nix-eval-jobs", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	resultChan := make(chan *lib.Result[*EvalResult], config.ResultBufferSize)

	alive := func() bool {
		return cmd.Process.Signal(syscall.Signal(0)) == nil
	}

	// Wait until both handling stdout stream and process is done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Wait for process
	go func() {
		defer wg.Done()
		stderrMux := sync.Mutex{}

		var stderrBytes []byte

		// Gather stderr for error message
		go func() {
			stderrMux.Lock()
			defer stderrMux.Unlock()
			var err error

			stderrBytes, err = io.ReadAll(stderr)
			if err != nil {

				err = fmt.Errorf("Error reading stderr: %w", err)
				resultChan <- lib.NewResult[*EvalResult](nil, err)
			}
		}()

		err := cmd.Wait()
		if err != nil {
			stderrMux.Lock()
			err = fmt.Errorf("Nix eval error '%w' with stderr: %s", err, string(stderrBytes))
			resultChan <- lib.NewResult[*EvalResult](nil, err)
			return
		}
	}()

	// Handle stdout stream
	go func() {
		defer wg.Done()

		die := func(err error) {
			resultChan <- lib.NewResult[*EvalResult](nil, err)

			// Kill process, disregarding the error on purpose
			// since it might already be dead
			cmd.Process.Kill() // nolint:errcheck
		}

		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
			if !alive() {
				return
			}

			s := scanner.Bytes()

			r := &EvalResult{
				Outputs:  make(map[string]string),
				AttrPath: []string{},
			}

			jsonparser.EachKey(s, func(idx int, value []byte, vt jsonparser.ValueType, keyErr error) {
				if keyErr != nil {
					die(keyErr)
					return
				}

				switch idx {

				case 0: // attr
					if vt != jsonparser.String {
						die(fmt.Errorf("'attr' not of type string"))
					}
					r.Attr = string(value)

				case 1: // attrPath
					if vt != jsonparser.Array {
						die(fmt.Errorf("'attrPath' not of type array"))
						return
					}

					var err error
					_, arrayEachErr := jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, arrErr error) {
						if arrErr != nil {
							err = fmt.Errorf("error in array index '%d': %w", offset, arrErr)
							return
						}

						if dataType != jsonparser.String {
							err = fmt.Errorf("attrPath member at index '%d' not of type string", offset)
							return
						}

						r.AttrPath = append(r.AttrPath, string(value))
					})
					if err != nil {
						die(err)
						return
					}
					if arrayEachErr != nil {
						die(err)
						return
					}

				case 2: // error
					if vt != jsonparser.String {
						die(fmt.Errorf("'error' not of type string"))
						return
					}
					r.Error = string(value)

				case 3: // drvPath
					if vt != jsonparser.String {
						die(fmt.Errorf("'drvPath' not of type string"))
						return
					}
					r.DrvPath = string(value)

				case 4: // name
					if vt != jsonparser.String {
						die(fmt.Errorf("'name' not of type string"))
						return
					}
					r.Name = string(value)

				case 5: // system
					if vt != jsonparser.String {
						die(fmt.Errorf("'system' not of type string"))
						return
					}
					r.System = string(value)

				case 6: // outputs
					if vt != jsonparser.Object {
						die(fmt.Errorf("'outputs' not of type object"))
						return
					}

					err := jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
						if dataType != jsonparser.String {
							return fmt.Errorf("'outputs' key '%s' not of type string", string(key))
						}
						r.Outputs[string(key)] = string(value)

						return nil
					})
					if err != nil {
						die(err)
						return
					}

				default:
					panic(fmt.Errorf("Unhandled index '%d' with accessor '%s'", idx, jsonPaths[idx]))
				}
			}, jsonPaths...)

			resultChan <- lib.NewResult(r, nil)
		}
	}()

	return resultChan, nil
}
