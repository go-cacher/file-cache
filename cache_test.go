package cache

import (
	"testing"

	"github.com/gocacher/cacher"
)

var c cacher.Cacher

func init() {
	c = New()
}

// TestFileCache_Set ...
func TestFileCache_Set(t *testing.T) {
	e := c.Set("abc", []byte("123"))
	if e != nil {
		t.Fatal(e)
	}
}

// TestFileCache_Get ...
func TestFileCache_Get(t *testing.T) {
	bys, e := c.Get("abc")
	if e != nil {
		t.Fatal(e)
	}
	t.Log(string(bys))
}
