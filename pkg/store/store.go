package store

type Store interface {
	Open() error
	Close() error
}
