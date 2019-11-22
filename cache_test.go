package cache

import (
	"log"
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
	checkErr(e)
}

// TestFileCache_Get ...
func TestFileCache_Get(t *testing.T) {
	bys, e := c.Get("abc")
	checkErr(e)
	t.Log(string(bys))
}

// TestFileCache_Delete ...
func TestFileCache_Delete(t *testing.T) {
	e := c.Delete("abc")
	checkErr(e)
}

// TestFileCache_Has ...
func TestFileCache_Has(t *testing.T) {
	b, e := c.Has("abc")
	checkErr(e)
	t.Log(b)
}

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
