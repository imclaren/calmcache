# calmcache
calmcache is a golang disk cache ([GoDoc](https://godoc.org/github.com/imclaren/calmcache)).  The filesystem cache is managed via a sqlite database.

calmcache reads and writes directly to files, and does not require values to be saved in memory.  For example:

```
func putAndGetBytes(key string, value []byte) {
	c, err := calmcache.Open(cachePath)
	if err != nil {
		return err
	}
	defer c.Close()

	// Put
	_, err = c.PutWithReader(bucket, key, bytes.NewReader(value), int64(len(value)))
	if err != nil {
		return err
	}

	// Get
	OK, fullPath, _, err := c.GetPathAndLock(bucket, key) 
	if err != nil {
		return err
	}
	defer c.GetPathUnlock()
	if !OK {
		return fmt.Errorf("key does not exist")
	}
	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	if b != value {
		return fmt.Errorf("value returned error")
	}
}
```
