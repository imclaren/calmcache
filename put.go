package calmcache

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/imclaren/calmcache/cacheitem"
	"github.com/imclaren/calmcache/filecache"
)

// Put puts the contents of a byte slice in a bucket.
// Use PutWithFile or PutWithReader instead to avoid holding the bytes in memory
// If the bucket already contains a value for the key, OK returns false and the existing value is not overwritten.
func (c *Cache) Put(bucket, key string, value []byte) (OK bool, err error) {
	c.Lock()
	defer c.Unlock()

	if key == "" {
		fmt.Errorf("cache error: empty key provided")
	}
	i, err := c.DB.GetItem(bucket, key)
	if err != nil {
		return false, err
	}
	if i != nil {
		err = c.DB.UpdateAccessCount(bucket, key)
		if err != nil {
			return false, err
		}
		return false, nil
	}
	fullPath, err := c.FC.FilePath(bucket, key, true)
	if err != nil {
		return false, err
	}
	err = c.FC.WriteBytesToFile(fullPath, value)
	if err != nil {
		return false, err
	}
	err = c.DB.Insert(cacheitem.New(bucket, key, int64(len(value)), 0, time.Time{}))
	if err != nil {
		return false, err
	}
	return true, nil
}

// PutWithFile puts the contents of a file at the provided path in a bucket
// If the bucket already contains a value for the key, OK returns false and the existing value is not overwritten.
func (c *Cache) PutWithFile(bucket, key string, fullPath string) (OK bool, err error) {
	c.Lock()
	defer c.Unlock()

	file, err := os.OpenFile(fullPath, os.O_RDONLY, filecache.FileMode)
	if err != nil {
		return false, fmt.Errorf("cache PutWithFile open file error: %s", err.Error())
	}
	defer file.Close()
	fi, err := os.Stat(fullPath)
	if err != nil {
		return false, err
	}
	return c.putWithReader(bucket, key, file, fi.Size())
}

// PutWithReader puts the contents of an io.Reader in a bucket
// If the bucket already contains a value for the key, OK returns false and the existing value is not overwritten.
func (c *Cache) PutWithReader(bucket, key string, r io.Reader, size int64) (OK bool, err error) {
	c.Lock()
	defer c.Unlock()

	return c.putWithReader(bucket, key, r, size)
}

func (c *Cache) putWithReader(bucket, key string, r io.Reader, size int64) (OK bool, err error) {
	if key == "" {
		return false, fmt.Errorf("cache error: empty key provided")
	}
	i, err := c.DB.GetItem(bucket, key)
	if err != nil {
		return false, err
	}
	if i != nil {
		err = c.DB.UpdateAccessCount(bucket, key)
		if err != nil {
			return false, err
		}
		return false, nil
	}
	fullPath, err := c.FC.FilePath(bucket, key, true)
	if err != nil {
		return false, err
	}
	err = c.FC.WriteReaderToFile(fullPath, r, size)
	if err != nil {
		return false, err
	}
	err = c.DB.Insert(cacheitem.New(bucket, key, size, 0, time.Time{}))
	if err != nil {
		return false, err
	}
	return true, nil
}


