package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDump(t *testing.T) {
	storage, _ := NewStorage("in-memory")
	storage.Store("https://example.com/", []string{"https://example.com/about-us/"})
	sitemap, err := storage.Dump()
	assert.Nil(t, err)
	assert.NotEmpty(t, sitemap)
}
