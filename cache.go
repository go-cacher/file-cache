package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/gocacher/cacher"
)

// ErrNotFound ...
var ErrNotFound = errors.New("data not found")

// FileCache ...
type FileCache struct {
	path string
}

type cacheData struct {
	Val []byte
	TTL *time.Time
}

// DefaultCachePath ...
var DefaultCachePath = "cache"

func init() {
	cacher.Register(&FileCache{})
}

// New ...
func New() cacher.Cacher {
	s, e := filepath.Abs(DefaultCachePath)
	if e != nil {
		panic(e)
	}
	_ = os.MkdirAll(s, 0755)
	return &FileCache{path: DefaultCachePath}
}

func read(path, file string) ([]byte, error) {
	target := filepath.Join(path, file)
	_, e := os.Stat(path)
	if e != nil {
		return nil, e
	}
	open, e := os.Open(target)
	if e != nil {
		return nil, e
	}
	bytes, e := ioutil.ReadAll(open)
	if e != nil {
		return nil, e
	}
	return bytes, nil
}

func write(path, file string, data []byte) error {
	target := filepath.Join(path, file)
	_, e := os.Stat(path)
	if e != nil {
		return e
	}
	return ioutil.WriteFile(target, data, os.ModePerm)
}

func exist(path, file string) bool {
	target := filepath.Join(path, file)
	_, e := os.Stat(target)
	if e == nil {
		return true
	}
	return false
}

func remove(path, file string) error {
	target := filepath.Join(path, file)
	_, e := os.Stat(target)
	if e != nil {
		return e
	}
	return os.Remove(target)
}

// Get ...
func (f FileCache) Get(key string) ([]byte, error) {
	v, err := read(f.path, key)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%s:%+v", key, err)
	}
	var d cacheData
	err = json.Unmarshal(v, &d)
	if err != nil {
		return nil, fmt.Errorf("%s:%+v", key, err)
	}
	return d.Val, nil

}

// GetD ...
func (f FileCache) GetD(key string, v []byte) []byte {
	if ret, err := f.Get(key); err == nil {
		return ret
	}
	return v
}

// Set ...
func (f *FileCache) Set(key string, val []byte) error {
	bytes, e := json.Marshal(&cacheData{
		Val: val,
		TTL: nil,
	})
	if e != nil {
		return fmt.Errorf("%s:%+v", key, e)
	}
	return write(f.path, key, bytes)
}

// SetWithTTL ...
func (f *FileCache) SetWithTTL(key string, val []byte, ttl int64) error {
	t := time.Now().Add(time.Duration(ttl))
	bytes, e := json.Marshal(&cacheData{
		Val: val,
		TTL: &t,
	})
	if e != nil {
		return fmt.Errorf("%s:%+v", key, e)
	}
	return write(f.path, key, bytes)
}

// Has ...
func (f *FileCache) Has(key string) (bool, error) {
	return exist(f.path, key), nil
}

// Delete ...
func (f *FileCache) Delete(key string) error {
	if err := remove(f.path, key); err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return fmt.Errorf("%s:%+v", key, err)
	}
	return nil
}

// Clear ...
func (f *FileCache) Clear() error {
	return os.RemoveAll(f.path)
}

// GetMultiple ...
func (f *FileCache) GetMultiple(keys ...string) (map[string][]byte, error) {
	vals := make(map[string][]byte, len(keys))
	for _, key := range keys {
		if ret, e := f.Get(key); e == nil {
			vals[key] = ret
		}
		return nil, fmt.Errorf("%s:%+v", key, ErrNotFound)
	}
	return vals, nil
}

// SetMultiple ...
func (f *FileCache) SetMultiple(values map[string][]byte) error {
	for k, v := range values {
		e := f.Set(k, v)
		if e != nil {
			return fmt.Errorf("%s:%+v", k, e)
		}
	}
	return nil
}

// DeleteMultiple ...
func (f *FileCache) DeleteMultiple(keys ...string) error {
	for _, key := range keys {
		e := f.Delete(key)
		if e != nil {
			return fmt.Errorf("%s:%+v", key, e)
		}
	}
	return nil
}
