package calmcache

import (
	"path/filepath"
	"sync"
	"context"

	"github.com/imclaren/calmcache/dbcache"
	"github.com/imclaren/calmcache/filecache"
)

const (
	DBName = "cache.db"
	FCName = "files"
)

// Cache is the calmcache struct
type Cache struct {
	sync.RWMutex
	ctx 				context.Context 
	cancel 				context.CancelFunc
	Path 				string
	DBPath 				string
	DB 					*dbcache.DB
	FCPath				string
	FC 					*filecache.FileCache
}

// Open opens the and initiates the cache. 
// Note that this is not thread safe.  Use Cache.Open for thread safe openining of the Cache.
func Open(path string) (c Cache, err error) {
	DBPath := filepath.Join(path, DBName)
	FCPath := filepath.Join(path, FCName)
	path, dirMode, err := filecache.MakeCacheDir(path)
	if err != nil {
		return Cache{}, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	DB, err := dbcache.Init(DBPath, ctx, cancel)
	if err != nil {
		return Cache{}, err
	}
	FC, err := filecache.Init(FCPath, dirMode)
	if err != nil {
		return Cache{}, err
	}
	return Cache{
		//mu: nil,
		ctx: ctx,
		cancel: cancel,
		Path: path,
		DBPath: DBPath,
		DB: &DB,
		FCPath: FCPath,
		FC: &FC,
	}, nil
}

// Open opens the cache in a thread safe manner
func (c *Cache) Open() error {
	c.Lock()
	defer c.Unlock()

	DB, err := dbcache.Open(c.DBPath, c.ctx, c.cancel)
	if err != nil {
		return err
	}
	c.DB = &DB
	return nil
}

// Close closes the cache
func (c *Cache) Close() error {
	c.Lock()
	defer c.Unlock()

	return c.DB.Close()
}

