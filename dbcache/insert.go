package dbcache

import "github.com/imclaren/calmcache/cacheitem"

// Insert inserts an item in the database
func (db *DB) Insert(i cacheitem.Item) error {
	db.Lock()
	defer db.Unlock()

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