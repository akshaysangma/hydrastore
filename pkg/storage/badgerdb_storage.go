package storage

import "github.com/dgraph-io/badger/v4"

// BadgerDBStorage represents a key-value store backed by BadgerDB.
type BadgerDBStorage struct {
	db *badger.DB
}

// NewBadgerDBStorage create a new instance of BadgerDB backed storage instance.
func NewBadgerDBStorage(path string) (*BadgerDBStorage, error) {
	opts := badger.DefaultOptions(path)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &BadgerDBStorage{
		db: db,
	}, nil
}

// Set stores a key-value pair in the BadgerDB storage.
func (s *BadgerDBStorage) Set(key, value []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
	return err
}

// Get retrieves a value for a given key from BadgerDB storage.
func (s *BadgerDBStorage) Get(key []byte) ([]byte, error) {
	var value []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	})
	return value, err
}

// Delete removes a key from the BadgerDB storage.
func (s *BadgerDBStorage) Delete(key []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
	return err
}

// Close the underlying BadgerDB storage.
func (s *BadgerDBStorage) Close() error {
	return s.db.Close()
}
