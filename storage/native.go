package storage

import (
	// "fmt"
	badger "github.com/dgraph-io/badger/v2"
	"github.com/tweag/trustix/config"
)

type nativeTxn struct {
	txn *badger.Txn
}

func (t *nativeTxn) Get(key []byte) ([]byte, error) {
	val, err := t.txn.Get(key)
	if err != nil {
		// Normalise error
		if err == badger.ErrKeyNotFound {
			return nil, ObjectNotFoundError
		}
		return nil, err
	}

	return val.ValueCopy(nil)
}

func (t *nativeTxn) Set(key []byte, value []byte) error {
	return t.txn.Set(key, value)
}

func newNativeTXN() *nativeTxn {
	return &nativeTxn{}
}

type NativeStorage struct {
	db *badger.DB
}

func NativeStorageFromConfig(conf *config.NativeStorageConfig) (*NativeStorage, error) {
	db, err := badger.Open(badger.DefaultOptions(conf.Path))
	if err != nil {
		return nil, err
	}

	return &NativeStorage{
		db: db,
	}, nil
}

func (s *NativeStorage) runTX(readWrite bool, fn func(Transaction) error) error {
	txn := s.db.NewTransaction(readWrite)
	if readWrite {
		defer txn.Discard()
	}

	t := &nativeTxn{
		txn: txn,
	}

	err := fn(t)
	if err != nil {
		return err
	} else {
		if readWrite {
			return txn.Commit()
		}
	}

	return err
}

func (s *NativeStorage) View(fn func(Transaction) error) error {
	return s.runTX(false, fn)
}

func (s *NativeStorage) Update(fn func(Transaction) error) error {
	return s.runTX(true, fn)
}

func (s *NativeStorage) Close() {
	s.db.Close()
}
