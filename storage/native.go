package storage

import (
	"github.com/tweag/trustix/config"
	bolt "go.etcd.io/bbolt"
	"path"
)

type nativeStorage struct {
	db *bolt.DB
}

type nativeTxn struct {
	txn *bolt.Tx
}

func (t *nativeTxn) Get(bucket []byte, key []byte) ([]byte, error) {
	b := t.txn.Bucket(bucket)
	if b == nil {
		return nil, ObjectNotFoundError
	}

	val := b.Get(key)
	if val == nil {
		return nil, ObjectNotFoundError
	}

	return val, nil
}

func (t *nativeTxn) Set(bucket []byte, key []byte, value []byte) error {
	b, err := t.txn.CreateBucketIfNotExists(bucket)
	if err != nil {
		return err
	}

	return b.Put(key, value)
}

func NativeStorageFromConfig(name string, stateDirectory string, conf *config.NativeStorageConfig) (*nativeStorage, error) {
	path := path.Join(stateDirectory, name+".db")

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &nativeStorage{
		db: db,
	}, nil
}

func (s *nativeStorage) runTX(readWrite bool, fn func(Transaction) error) error {
	txn, err := s.db.Begin(readWrite)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	t := &nativeTxn{
		txn: txn,
	}
	err = fn(t)
	if err != nil {
		return err
	} else {
		if readWrite {
			return txn.Commit()
		}
	}

	return err
}

func (s *nativeStorage) View(fn func(Transaction) error) error {
	return s.runTX(false, fn)
}

func (s *nativeStorage) Update(fn func(Transaction) error) error {
	return s.runTX(true, fn)
}

func (s *nativeStorage) Close() {
	s.db.Close()
}
