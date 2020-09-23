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

func (t *nativeTxn) Get(key []byte) ([]byte, error) {
	bucket := t.txn.Bucket([]byte("somebucket"))
	if bucket == nil {
		return nil, ObjectNotFoundError
	}

	val := bucket.Get(key)
	if val == nil {
		return nil, ObjectNotFoundError
	}

	return val, nil
}

func (t *nativeTxn) Set(key []byte, value []byte) error {
	bucket, err := t.txn.CreateBucketIfNotExists([]byte("somebucket"))
	if err != nil {
		return err
	}

	return bucket.Put(key, value)
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
