package file_cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var NotFoundError = errors.New("data not found")

type FileCache struct {
}

type cacheData struct {
	Val []byte
	TTL *time.Time
}

var DefaultPath = "cache"

func read(file string) ([]byte, error) {
	target := filepath.Join(DefaultPath, file)
	_, e := os.Stat(DefaultPath)
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

func write(file string, data []byte) error {
	target := filepath.Join(DefaultPath, file)
	_, e := os.Stat(DefaultPath)
	if e != nil {
		return e
	}
	return ioutil.WriteFile(target, data, os.ModePerm)
}

func exist(file string) bool {
	target := filepath.Join(DefaultPath, file)
	_, e := os.Stat(target)
	if e == nil {
		return true
	}
	return false
}

func remove(file string) error {
	target := filepath.Join(DefaultPath, file)
	_, e := os.Stat(target)
	if e != nil {
		return e
	}
	return os.Remove(target)
}

func (f FileCache) Get(key string) ([]byte, error) {
	v, err := read(key)
	if err != nil {
		return nil, fmt.Errorf("%s:%+v", key, err)
	}
	var d cacheData
	err = json.Unmarshal(v, &d)
	if err != nil {
		return nil, fmt.Errorf("%s:%+v", key, err)
	}
	return d.Val, nil

}

func (f FileCache) GetD(key string, v []byte) []byte {
	if ret, err := f.Get(key); err == nil {
		return ret
	}
	return v
}

func (f *FileCache) Set(key string, val []byte) error {
	bytes, e := json.Marshal(&cacheData{
		Val: val,
		TTL: nil,
	})
	if e != nil {
		return fmt.Errorf("%s:%+v", key, e)
	}
	return write(key, bytes)
}

func (f *FileCache) SetWithTTL(key string, val []byte, ttl int64) error {
	t := time.Now().Add(time.Duration(ttl))
	bytes, e := json.Marshal(&cacheData{
		Val: val,
		TTL: &t,
	})
	if e != nil {
		return fmt.Errorf("%s:%+v", key, e)
	}
	return write(key, bytes)
}

func (f *FileCache) Has(key string) (bool, error) {
	return exist(key), nil
}

func (f *FileCache) Delete(key string) error {
	if err := remove(key); err != nil {
		return fmt.Errorf("%s:%+v", key, err)
	}
	return nil
}

func (f *FileCache) Clear() error {
	if err := os.RemoveAll(DefaultPath); err != nil {
		return err
	}
	return nil
}

func (f *FileCache) GetMultiple(keys ...string) (map[string][]byte, error) {
	vals := make(map[string][]byte, len(keys))
	for _, key := range keys {
		if ret, e := f.Get(key); e == nil {
			vals[key] = ret
		}
		return nil, fmt.Errorf("%s:%+v", key, NotFoundError)
	}
	return vals, nil
}

func (f *FileCache) SetMultiple(values map[string][]byte) error {
	for k, v := range values {
		e := f.Set(k, v)
		if e != nil {
			return fmt.Errorf("%s:%+v", k, e)
		}
	}
	return nil
}

func (f *FileCache) DeleteMultiple(keys ...string) error {
	for _, key := range keys {
		e := f.Delete(key)
		if e != nil {
			return fmt.Errorf("%s:%+v", key, e)
		}
	}
	return nil
}
