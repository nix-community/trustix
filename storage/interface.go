package storage

type Transaction interface {
	Get(key []byte) ([]byte, error)
	Set(key []byte, value []byte) error
}

type TrustixStorage interface {
	Close()

	// View - Start a read-only transaction
	View(func(txn Transaction) error) error

	// Update - Start a read-write transaction
	Update(func(txn Transaction) error) error
}
