package filecache

import (
	"os"
	"path/filepath"
)

// DeleteCache deletes all items in the cache
func (fc *FileCache) DeleteCache() error {
	fc.Lock()
	defer fc.Unlock()

	cachePath := filepath.Dir(fc.path)
	return os.RemoveAll(cachePath)
}

// DeleteBucket deletes all the items in a bucket
func (fc *FileCache) DeleteBucket(bucket string) error {
	fc.Lock()
	defer fc.Unlock()

	bucketPath := filepath.Join(fc.path, bucket)
	return os.RemoveAll(bucketPath)
}

// Delete deletes an item
func (fc *FileCache) Delete(bucket, key string) error {
	fc.Lock()
	defer fc.Unlock()

	// Delete subDirs
	subDirs, err := fc.subDirs(bucket, key, DirLength, false)
	if err != nil {
		return err
	}

	// Delete file
	fullPath := filepath.Join(subDirs[len(subDirs)-1], key)
	err = os.Remove(fullPath)
	if err != nil {
		return err
	}

	// Delete empty subdirs
	for i := len(subDirs)-1; i >= 0; i-- {
		d := subDirs[i]
		isEmpty, err := dirIsEmpty(d)
		if err != nil {
			return err
		}
		if isEmpty {
			err = os.Remove(d)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
