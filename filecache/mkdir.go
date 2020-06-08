package filecache

import (
	"os"
	"path/filepath"
)

// MakeCacheDir makes the cache folder and returns the absolute path and os.FileMode of the parent directory of the cache
func MakeCacheDir(path string) (cachePath string, dirMode os.FileMode, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", dirMode, err
	}
	parentFi, err := os.Stat(filepath.Dir(absPath))
	if err != nil {
		return "", dirMode, err
	}
	_, err = os.Stat(absPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", dirMode, err
		} 
		os.Mkdir(absPath, parentFi.Mode())
	}
	cachePath, err = filepath.EvalSymlinks(absPath)
	if err != nil {
		return "", dirMode, err
	}
	return cachePath, parentFi.Mode(), nil
}

// MakeDir makes a dir at the provided path with the provided os.FileMode
func MakeDir(path string, dirMode os.FileMode) (err error) {
	_, err = os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} 
		os.Mkdir(path, dirMode)
	}
	return nil
}




