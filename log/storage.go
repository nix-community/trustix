package log

import (
	"fmt"
)

type logStorage struct {
	// Use an in-memory map for now, this needs to be changed to persistent storage later
	m map[string][]byte
}

func (l *logStorage) Get(idx int) ([]byte, error) {
	key := fmt.Sprintf("log-%s", idx)
	v, ok := l.m[key]
	if !ok {
		return nil, fmt.Errorf("Could not find entry")
	}
	return v, nil
}

func (l *logStorage) Append(idx int, entry []byte) error {
	key := fmt.Sprintf("log-%s", idx)
	l.m[key] = entry
	return nil
}
