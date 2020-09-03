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
