package seen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	cache, err := NewCache("in-memory")
	assert.Nil(t, err)
	assert.NotNil(t, cache)
	assert.IsType(t, new(inMemory), cache)
	assert.Implements(t, new(Cache), cache)

	notImplementedCache, err := NewCache("not-implemented")
	assert.EqualError(t, err, fmt.Sprintf("cache engine %s not supported", "not-implemented"))
	assert.Nil(t, notImplementedCache)
}
