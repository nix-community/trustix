package log

import (
	"fmt"
	"github.com/tweag/trustix/storage"
)

type logStorage struct {
	txn storage.Transaction
}

func (s *logStorage) LevelSize(level int) int {
	n, err := s.txn.Size([]byte(fmt.Sprintf("log-%d")))
	if err != nil {
		if err == storage.ObjectNotFoundError {
			return 1
		}
		panic(err)
	}

	return n
}

func (s *logStorage) Size() int {
	n, err := s.txn.Size([]byte("log-root"))
	if err != nil {
		if err == storage.ObjectNotFoundError {
			return 1
		}
		panic(err)
	}
	return n
}

func (s *logStorage) Get(level int, idx int) *Leaf {
	bucket := []byte(fmt.Sprintf("log-%d", level))
	key := []byte(fmt.Sprintf("%d", idx))

	v, err := s.txn.Get(bucket, key)
	if err != nil {
		panic(err)
	}

	l, err := LeafFromBytes(v)
	if err != nil {
		panic(err)
	}

	return l
}

func (s *logStorage) Append(level int, leaf *Leaf) {
	if s.Size() == level {
		// "Grow" root level by adding another "level"
		err := s.txn.Set([]byte("log-root"), []byte(fmt.Sprintf("%d", s.Size())), []byte(""))
		if err != nil {
			panic(err)
		}
	}

	v, err := leaf.Marshal()
	if err != nil {
		panic(err)
	}

	idx := s.LevelSize(level) - 1

	bucket := []byte(fmt.Sprintf("log-%d", level))
	key := []byte(fmt.Sprintf("%d", idx))

	err = s.txn.Set(bucket, key, v)
	if err != nil {
		panic(err)
	}

}
