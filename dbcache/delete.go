package dbcache

// Delete deletes an item from the database
func (db *DB) Delete(bucket, key string) error {
	db.Lock()
	defer db.Unlock()

	sqlString := "DELETE FROM cache WHERE bucket = ? AND key = ?"
	_, err := db.Exec(db.Rebind(sqlString), bucket, key)
	return err
}

// DeleteBucket deletes all of the items in the bucket
func (db *DB) DeleteBucket(bucket string) error {
	db.Lock()
	defer db.Unlock()

	sqlString := "DELETE FROM cache WHERE bucket = ?"
	_, err := db.Exec(db.Rebind(sqlString), bucket)
	return err
}

// DeleteAll deletes all of the items from the cache
func (db *DB) DeleteAll() error {
	db.Lock()
	defer db.Unlock()

	sqlString := "DELETE FROM cache"
	_, err := db.Exec(db.Rebind(sqlString))
	return err
}