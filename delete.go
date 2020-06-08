package calmcache

// DeleteCache deletes the cache
func (c *Cache) DeleteCache() (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.DB.Mutex.Lock()
	defer c.DB.Mutex.Unlock()

	return c.FC.DeleteCache()
}

// DeleteBucket deletes the bucket
func (c *Cache) DeleteBucket(bucket string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.DB.DeleteBucket(bucket)
	if err != nil {
		return err
	}
	return c.FC.DeleteBucket(bucket)
}

// Delete deletes an item from a bucket
func (c *Cache) Delete(bucket, key string) (OK bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.delete(bucket, key)
}

func (c *Cache) delete(bucket, key string) (OK bool, err error) {
	exists, err := c.exists(bucket, key)
	err = c.DB.Delete(bucket, key)
	if err != nil {
		return false, err
	}
	if !exists {
		return true, nil
	}
	err = c.DB.Delete(bucket, key)
	if err != nil {
		return false, err
	}
	err = c.FC.Delete(bucket, key)
	if err != nil {
		return false, err
	}
	return true, nil
}

