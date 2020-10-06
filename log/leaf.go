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

package log

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type Leaf struct {
	Digest []byte
	Value  []byte
}

func NewLeaf(digest []byte, value []byte) (*Leaf, error) {
	if len(digest) > 65535 || len(value) > 65535 {
		return nil, fmt.Errorf("Records over 65535 bytes is not supported")
	}

	return &Leaf{
		Digest: digest,
		Value:  value,
	}, nil
}

func LeafFromBytes(data []byte) (*Leaf, error) {
	l := &Leaf{}
	s := reflect.ValueOf(l).Elem()

	offset := uint16(0)

	for i := 0; i < s.NumField(); i++ {
		fLen := binary.LittleEndian.Uint16(data[offset : offset+2])
		newOffset := offset + 2 + fLen

		field := s.Field(i)
		field.Set(reflect.ValueOf(data[offset+2 : newOffset]))

		offset = newOffset
	}

	return l, nil

}

func (l *Leaf) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, value := range [][]byte{l.Digest, l.Value} {
		num := uint16(len(value))
		err := binary.Write(buf, binary.LittleEndian, num)
		if err != nil {
			return nil, fmt.Errorf("binary.Write failed:", err)
		}
		err = binary.Write(buf, binary.LittleEndian, value)
		if err != nil {
			return nil, fmt.Errorf("binary.Write failed:", err)
		}
	}

	return buf.Bytes(), nil
}
