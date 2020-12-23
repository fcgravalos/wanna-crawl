package seen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeen(t *testing.T) {
	cache, _ := NewCache("in-memory")
	cache.Add("https://example.com/")
	assert.True(t, cache.Seen("https://example.com/"))
	assert.False(t, cache.Seen("https://example.com/about-us"))
}
