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
	"fmt"

	proto "github.com/golang/protobuf/proto"
	"github.com/tweag/trustix/schema"
	"github.com/tweag/trustix/storage"
)

type LogStorage struct {
	txn storage.Transaction
}

func NewLogStorage(txn storage.Transaction) *LogStorage {
	return &LogStorage{
		txn: txn,
	}
}

func (s *LogStorage) Get(level int, idx uint64) (*schema.LogLeaf, error) {
	bucket := []byte(fmt.Sprintf("log-%d", level))
	key := []byte(fmt.Sprintf("%d", idx))

	v, err := s.txn.Get(bucket, key)
	if err != nil {
		return nil, err
	}

	l, err := LeafFromBytes(v)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (s *LogStorage) Append(treeSize uint64, level int, leaf *schema.LogLeaf) error {
	v, err := proto.Marshal(leaf)
	if err != nil {
		return err
	}

	idx := levelSize(treeSize, level) - 1

	bucket := []byte(fmt.Sprintf("log-%d", level))
	key := []byte(fmt.Sprintf("%d", idx))

	err = s.txn.Set(bucket, key, v)
	if err != nil {
		return err
	}

	return nil

}
