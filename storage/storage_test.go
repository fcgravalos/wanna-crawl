package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	db, err := NewStorage("in-memory")
	assert.Nil(t, err)
	assert.NotNil(t, db)
	assert.IsType(t, new(inMemory), db)
	assert.Implements(t, new(Storage), db)

	notImplementedStorage, err := NewStorage("not-implemented")
	assert.EqualError(t, err, fmt.Sprintf("storage engine %s not supported", "not-implemented"))
	assert.Nil(t, notImplementedStorage)
}
