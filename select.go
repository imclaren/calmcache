package calmcache

import (
	"fmt"
	"io"
	"os"
	"io/ioutil"

	"github.com/imclaren/calmcache/filecache"
)

// Exists checks if an items exists in the cache
func (c *Cache) Exists(bucket, key string) (exists bool, err error) {
	c.RLock()
	defer c.RUnlock()

	return c.exists(bucket, key)
}

func (c *Cache) exists(bucket, key string) (exists bool, err error) {
	i, err := c.DB.GetItem(bucket, key)
	if err != nil {
		return false, err
	}
	return i != nil, err
}

// AllKeys returs the keys for all items in a bucket
func (c *Cache) AllKeys(bucket string) (allKeys []string, err error) {
	c.RLock()
	defer c.RUnlock()

	allKeys = []string{}
	items, err := c.DB.GetAllInBucket(bucket)
	if err != nil {
		return nil, err
	}
	for _, i := range items {
		allKeys = append(allKeys, i.Key)
	}
	return allKeys, err
}

// Get gets the cached item bytes
// Use GetPathAndLock / GetPathUnLock or GetToWriter instead to avoid holding the bytes in memory
func (c *Cache) Get(bucket, key string) (value []byte, err error) {
	c.RLock()
	defer c.RUnlock()

	OK, fullPath, _, err := c.getPath(bucket, key)
	if err != nil {
		return nil, err
	}
	if !OK {
		return nil, nil
	}
	return ioutil.ReadFile(fullPath)
}

// GetPathAndLock gets the path of the cached file to read.
// Note that the cache will lock until GetPathUnLock is called
func (c *Cache) GetPathAndLock(bucket, key string) (OK bool, fullPath string, size int64, err error) {
	c.RLock()
	//defer GetPathUnlock()

	OK, fullPath, size, err = c.getPath(bucket, key)
	if err != nil {
		c.GetPathUnlock()
	}
	return OK, fullPath, size, err
}

// GetPathUnlock unlocks the cache.
// Use this when finished with the item accessed using GetPathAndLock 
func (c *Cache) GetPathUnlock() {
	c.RUnlock()
}

func (c *Cache) getPath(bucket, key string) (OK bool, fullPath string, size int64, err error) {
	i, err := c.DB.GetItem(bucket, key)
	if err != nil {
		return false, "", 0, err
	}
	if i == nil {
		return false, "", 0, nil
	}
	err = c.DB.UpdateAccessCount(bucket, key)
	if err != nil {
		return false, "", 0, err
	}
	fullPath, err = c.FC.FilePath(bucket, key, true)
	if err != nil {
		return false, "", 0, err
	}
	return true, fullPath, i.Size, nil
}

// GetToWriter gets the cached item bytes as an io.Writer
func (c *Cache) GetToWriter(bucket, key string, w io.Writer) (OK bool, err error) {
	OK, fullPath, size, err := c.GetPathAndLock(bucket, key)
	if err != nil {
		return false, fmt.Errorf("cache GetToWriter GetPathAndLock error: %s %s %s", bucket, key, err.Error())
	}
	defer c.GetPathUnlock()
	if !OK {
		return false, nil
	}
	file, err := os.OpenFile(fullPath, os.O_RDONLY, filecache.FileMode)
	if err != nil {
		return false, fmt.Errorf("cache GetToWriter OpenFile error: %s %s %s %v", bucket, key, fullPath, err)
	}
	defer file.Close()
    n, err := io.Copy(w, file)
    if err != nil {
    	return false, fmt.Errorf("cache GetToWriter io.Copy error: %s %s %s", bucket, key, err.Error())
    }
	if n < size {
		return false, io.ErrShortWrite
	}
	return true, nil
}
