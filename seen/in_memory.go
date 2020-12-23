package seen

import "sync"

// inMemory is a basic implementation of the seen cache using a map.
type inMemory struct {
	sync.RWMutex
	seen map[string]bool
}

//Seen returns `true` if `url` is in `seen` cache, `false` otherwise.
func (im *inMemory) Seen(url string) bool {
	im.RLock()
	s := im.seen[url]
	im.RUnlock()
	return s
}

//Add will insert a new `url` in `seen` cache
func (im *inMemory) Add(url string) error {
	im.Lock()
	im.seen[url] = true
	im.Unlock()
	return nil
}
