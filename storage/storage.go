package storage

import "fmt"

var storageEngines map[string]bool

const inMemoryStorage = "in-memory"

func init() {
	storageEngines = make(map[string]bool)
	storageEngines[inMemoryStorage] = true
}

// Storage abstracts different implementation for the crawler results store.
type Storage interface {
	Store(u string, l []string) error
	Dump() (string, error)
}

// NewStorage returns a `Storage` interface given the Storage `kind` or `error` if it is not supported.
func NewStorage(kind string) (Storage, error) {
	if !storageEngines[kind] {
		return nil, fmt.Errorf("storage engine %s not supported", kind)
	}

	var storage Storage
	var err error

	switch kind {
	case inMemoryStorage:
		im := &inMemory{
			db: make(map[string][]string),
		}
		storage = im
	}
	return storage, err
}
