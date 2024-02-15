package storage

// Storage is a interface :
type Storage interface {
	Set(key, value []byte) error
	Get(key []byte) []byte
	Delete(key []byte) error
	Close() error
}
