package dbcache

// UpdateAccessCount updates the access count of an item
func (db *DB) UpdateAccessCount(bucket, key string) (err error) {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    var accessCount int64
    sqlString := "SELECT access_count FROM cache WHERE bucket = ? AND key = ?"
	err = tx.QueryRow(db.Rebind(sqlString), bucket, key).Scan(&accessCount)
	if err != nil {
		return err
	}
	sqlString = "UPDATE cache SET access_count = ? WHERE bucket = ? AND key = ?"
	_, err = tx.Exec(db.Rebind(sqlString), accessCount+1, bucket, key)
	if err != nil {
		return err
	}
    return tx.Commit()
}