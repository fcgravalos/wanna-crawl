package storage

import (
	"encoding/json"
	"sync"
)

type inMemory struct {
	sync.RWMutex
	db map[string][]string
}

func (im *inMemory) Store(u string, l []string) error {
	im.Lock()
	im.db[u] = l
	im.Unlock()
	return nil
}

func (im *inMemory) Dump() (string, error) {
	im.RLock()
	db := im.db
	im.RUnlock()
	jsonData, err := json.MarshalIndent(db, "", "\t")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
