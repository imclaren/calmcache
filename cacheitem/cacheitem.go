package cacheitem

import (
	"time"
)

// Item is a cache item
type Item struct {
	Id              int
	Bucket 			string 		`db:"bucket"`
	Key 			string 		`db:"key"`
	Size 			int64 		`db:"size"`
	AccessCount 	int64  		`db:"access_count"`
	ExpiresAt 		time.Time  	`db:"expires_at"`
	CreatedAt       time.Time 	`db:"created_at"`
	UpdatedAt       time.Time 	`db:"updated_at"`
}

// New returns a new cache item
func New(bucket, key string, size, accessCount int64, expiresAt time.Time) Item {
	return Item{
		//Id              int
		Bucket: 		bucket,
		Key: 			key,
		Size: 			size,
		AccessCount: 	accessCount,
		ExpiresAt: 		expiresAt,
		//CreatedAt       time.Time 	`db:"created_at"`
		//UpdatedAt       time.Time 	`db:"updated_at"`
	}
}