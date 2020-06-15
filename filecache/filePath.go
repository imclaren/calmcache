package filecache

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FilePath gets the file path of the cached file, and creates the required subfolders if necessary
func (fc *FileCache) FilePath(bucket, key string, makeDir bool) (filePath string, err error) {
	fc.Lock()
	defer fc.Unlock()

	subDirs, err := fc.subDirs(bucket, key, DirLength, true)
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(subDirs[len(subDirs)-1], key)
	return fullPath, nil
}

func (fc *FileCache) subDirs(bucket, key string, chunkSize int, makeDirs bool) (subDirs []string, err error) {
	baseDir := filepath.Ext(key) 
	baseName := strings.TrimSuffix(key, baseDir)
	if baseDir == "" {
		baseDir = "other"
	}
	baseDir = strings.TrimPrefix(baseDir, ".")
	
	// Get bucket dir
	currentPath := filepath.Join(fc.path, bucket)
	if makeDirs {
		err = fc.makeBucketDir(bucket)
		if err != nil {
			return nil, err
		}
	}

	// Get subdirs, and create subdirs if required
	subDirs = []string{}
	for _, d := range append([]string{baseDir}, stringChunks(baseName, chunkSize)...) {
		currentPath = filepath.Join(currentPath, d)
		if makeDirs {
			_, err = os.Stat(currentPath)
			if err != nil {
				if !os.IsNotExist(err) {
					return nil, err
				} 
				os.Mkdir(currentPath, fc.dirMode)
			}
		}
		subDirs = append(subDirs, currentPath)
	}
	return subDirs, nil
}

func (fc *FileCache) makeBucketDir(bucket string) error {
	cachePath := filepath.Dir(fc.path)
	_, err := os.Stat(cachePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} 
		return fmt.Errorf("cache path does not exist: %s", cachePath)
	}
	_, err = os.Stat(fc.path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} 
		return fmt.Errorf("file cache path does not exist: %s", fc.path)
	}
	bucketPath := filepath.Join(fc.path, bucket)
	_, err = os.Stat(bucketPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} 
		os.Mkdir(bucketPath, fc.dirMode)
	}
	return nil
}

// stringChunks was updated from code at https://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty
func dirIsEmpty(fullPath string) (bool, error) {
    f, err := os.Open(fullPath)
    if err != nil {
        return false, err
    }
    defer f.Close()
    _, err = f.Readdirnames(1) // Or f.Readdir(1)
    if err == io.EOF {
        return true, nil
    }
    return false, err
}

// stringChunks was updated from code at https://stackoverflow.com/questions/25686109/split-string-by-length-in-golang
func stringChunks(s string, chunkSize int) []string {
    if chunkSize >= len(s) {
        return []string{s}
    }
    var chunks []string
    chunk := make([]rune, chunkSize)
    len := 0
    for _, r := range s {
        chunk[len] = r
        len++
        if len == chunkSize {
            chunks = append(chunks, string(chunk))
            len = 0
        }
    }
    if len > 0 {
        chunks = append(chunks, string(chunk[:len]))
    }
    return chunks
}