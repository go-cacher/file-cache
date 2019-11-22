package easy

import (
	"github.com/gocacher/cacher"
	cache "github.com/gocacher/file-cache"
)

func init() {
	cacher.Register(cache.New())
}
