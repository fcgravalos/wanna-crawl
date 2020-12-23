package seen

import "fmt"

var cacheEngines map[string]bool

const inMemoryCache = "in-memory"

func init() {
	cacheEngines = make(map[string]bool)
	cacheEngines[inMemoryCache] = true
}

// Cache provides an interface to keep track of seen urls.
type Cache interface {
	Seen(url string) bool
	Add(url string) error
}

// NewCache returns the `kind` specific `Cache` implementation
func NewCache(kind string) (Cache, error) {
	if !cacheEngines[kind] {
		return nil, fmt.Errorf("cache engine %s not supported", kind)
	}

	var cache Cache
	var err error

	switch kind {
	case inMemoryCache:
		im := &inMemory{
			seen: make(map[string]bool),
		}
		cache = im
		err = nil
	}
	return cache, err
}
