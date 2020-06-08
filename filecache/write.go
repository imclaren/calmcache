package filecache

import (
	"fmt"
	"io"
	"os"

	"github.com/imclaren/fs"
)

// WriteBytesToFile writes the contents of []byte to a file at filepath
// Use WriteReaderToFile instead to avoid holding the bytes in memory
func (fc *FileCache) WriteBytesToFile(fullPath string, b []byte) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fullPath, err := fs.RealPath(fullPath)
	if err != nil {
		return fmt.Errorf("cache RealPath error: %s %v", fullPath, err)
	}
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, FileMode)
	if err != nil {
		return fmt.Errorf("cache OpenFile error: %s %v", fullPath, err)
	}
	defer file.Close()
	n, err := file.Write(b)
	if err != nil {
		return err
	}
	if n < len(b) {
		return io.ErrShortWrite
	}
	err = file.Sync()
	if err != nil {
		return fmt.Errorf("cache file.Sync error: %s %v", fullPath, err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("cache file.Close error: %s %v", fullPath, err)
	}
	return nil
}

// WriteReaderToFile streams the contents of an io.Reader to a file at filepath
func (fc *FileCache) WriteReaderToFile(fullPath string, r io.Reader, size int64) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fullPath, err := fs.RealPath(fullPath)
	if err != nil {
		return fmt.Errorf("cache RealPath error: %s %v", fullPath, err)
	}
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, FileMode)
	if err != nil {
		return fmt.Errorf("cache OpenFile error: %s %v", fullPath, err)
	}
	defer file.Close()
    n, err := io.Copy(file, r)
    if err != nil {
    	return err
    }
	if n < size {
		return io.ErrShortWrite
	}
	err = file.Sync()
	if err != nil {
		return fmt.Errorf("cache file.Sync error: %s %v", fullPath, err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("cache file.Close error: %s %v", fullPath, err)
	}
	return nil
}


