package filecache

import (
	"sync"
	"os"
)

const (
	FileMode = 0644
	DirLength = 4
)

// FileCache is the FileCache struct
type FileCache struct {
	mu  		sync.RWMutex
	path   		string
	dirMode    	os.FileMode
}

// Init initiates the FileCache
func Init(path string, dirMode os.FileMode) (fc FileCache, err error) {
	err = MakeDir(path, dirMode)
	if err != nil {
		return FileCache{}, err
	}
	return FileCache{
		//mu: nil,
		path: path,
		dirMode: dirMode,
	}, nil
}
