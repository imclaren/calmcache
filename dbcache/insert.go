package dbcache

// Insert inserts an item in the database
func (db *DB) Insert(i Item) error {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	sqlString := "INSERT INTO cache (bucket, key, size, access_count, expires_at) VALUES (?,?,?,?,?)"
	_, err := db.Exec(db.Rebind(sqlString),
		i.Bucket,
		i.Key,
		i.Size,
		i.AccessCount,
		i.ExpiresAt,
	)
	return err
}