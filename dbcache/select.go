package dbcache

import (
	"database/sql"
	"strings"
	"fmt"
	"time"
)

// GetItem returns a database item
func (db *DB) GetItem(bucket, key string) (i *Item, err error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	if bucket == "" || key == "" {
		return nil, fmt.Errorf("empty bucket (%s) or key (%s)", bucket, key)
	}
	sqlString := "SELECT * FROM cache WHERE bucket = ? AND key = ?"
	var newItem Item 
	err = db.QueryRowx(db.Rebind(sqlString), bucket, key).StructScan(&newItem)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &newItem, nil
}

// GetAllInBucket returns all of the database items in a bucket
func (db *DB) GetAllInBucket(bucket string) ([]Item, error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	sqlString := "SELECT * FROM cache WHERE bucket = ? ORDER BY key ASC"
	var items []Item
	err := db.Select(&items, db.Rebind(sqlString), bucket)
	return items, err
}

// GetOldestInBucket returns the oldest (i.e. last accessed) database item
func (db *DB) GetOldestInBucket(bucket string) (i *Item, err error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	sqlString := "SELECT * FROM cache WHERE bucket = ? ORDER BY updated_at ASC LIMIT 1"
	var newItem Item 
	err = db.QueryRowx(db.Rebind(sqlString), bucket).StructScan(&newItem)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &newItem, nil
}

// AllInBucketOlderThan returns all database items that are older than (i.e. last accessed before) the provided time.Duration
func (db *DB) AllInBucketOlderThan(bucket string, d time.Duration) (items []Item, err error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	targetTS := time.Now().Add(-d)
	sqlString := "SELECT * FROM cache WHERE bucket = ? AND updated_at < ? ORDER BY updated_at ASC"

	var newItems []Item
	err = db.Select(&items, db.Rebind(sqlString), bucket, targetTS)
	if err != nil {
		if err == sql.ErrNoRows {
			return []Item{}, nil
		}
		return nil, err
	}
	return newItems, nil
}

// All returns all database items in the cache
func (db *DB) All() ([]Item, error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	sqlString := "SELECT * FROM cache"
	var items []Item
	err := db.Select(&items, db.Rebind(sqlString))
	return items, err
}

// AllInBucketCount returns the number of items in a bucket
func (db *DB) AllInBucketCount(bucket string) (count int, err error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	sqlString := "SELECT COUNT(*) FROM cache WHERE bucket = ?"
	var c int
	err = db.Get(&c, db.Rebind(sqlString), bucket)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "converting NULL to int64 is unsupported"):
			return 0, nil
		case strings.Contains(err.Error(), "invalid syntax"):
			return 0, nil
		default:
		}
		return 0, err
	}
	return c, nil
}

// AllInBucketCount returns the number of items in the cache
func (db *DB) AllCount() (count int, err error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	sqlString := "SELECT COUNT(*) FROM cache"
	var c int
	err = db.Get(&c, db.Rebind(sqlString))
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "converting NULL to int64 is unsupported"):
			return 0, nil
		case strings.Contains(err.Error(), "invalid syntax"):
			return 0, nil
		default:
		}
		return 0, err
	}
	return c, nil
}

// BucketSize returns the total size (in bytes) of the items in a bucket
func (db *DB) BucketSize(bucket string) (size int64, err error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	sqlString := "SELECT SUM(size) FROM cache WHERE bucket = ?"
	var s int64
	err = db.Get(&s, db.Rebind(sqlString), bucket)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "converting NULL to int64 is unsupported"):
			return 0, nil
		case strings.Contains(err.Error(), "invalid syntax"):
			return 0, nil
		default:
		}
		return 0, err
	}
	return s, nil
}

// Size returns the total size (in bytes) of the items in the cache 
func (db *DB) Size() (size int64, err error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	sqlString := "SELECT SUM(size) FROM cache"
	var s int64
	err = db.Get(&s, db.Rebind(sqlString))
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "converting NULL to int64 is unsupported"):
			return 0, nil
		case strings.Contains(err.Error(), "invalid syntax"):
			return 0, nil
		default:
		}
		return 0, err
	}
	return s, nil
}
