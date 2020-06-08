package calmcache

import (
	"fmt"
	"time"
)

// PruneToSize prunes the bucket to targetSize (by last accessed time)
func (c *Cache) PruneToSize(bucket string, targetSize int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	bucketSize, err := c.DB.BucketSize(bucket)
	if err != nil {
		return err
	}
	if bucketSize == int64(0) || bucketSize <= targetSize {
		return nil
	}
	for bucketSize > targetSize {
		i, err := c.DB.GetOldestInBucket(bucket)
		if err != nil {
			return err
		}
		if i == nil {
			return nil
		}
		OK, err := c.delete(i.Bucket, i.Key)
	   	if err != nil {
	    	return err
	    }
	    if !OK {
	    	return fmt.Errorf("pruneToSize bucket (%s) delete error for key: %s", bucket, i.Key)
	    }
	    bucketSize = bucketSize-i.Size
	}
	return nil
}

// PruneOlderThan prunes the bucket of all items with an access time that is earlier than the time.Duration provided
func (c *Cache) PruneOlderThan(bucket string, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	items, err := c.DB.AllInBucketOlderThan(bucket, d)
	if err != nil {
		return err
	}
	for _, i := range items {
		OK, err := c.delete(bucket, i.Key)
	   	if err != nil {
	    	return err
	    }
	    if !OK {
	    	return fmt.Errorf("PruneOlderThan bucket (%s) error for key: %s", bucket, i.Key)
	    }	
	}
	return nil
}


