# calmcache
calmcache is a low memory golang disk cache ([GoDoc](https://godoc.org/github.com/imclaren/calmcache)).  The filesystem cache is managed using a sqlite database.

## Example
```
cachePath := "/path/to/cachePath"
bucket := "mybucket"
key := "mykey"
value := []byte{'a', 'b', 'c'}

c, err := Open(cachePath)
if err != nil {
	log.Fatal(err)
}
OK, err := c.Put(bucket, key, value)
if err != nil {
	log.Fatal(err)
}
if !OK {
	log.Fatal(fmt.Errorf("key already contains a value - delete the key first if you want to overwrite the value"))
}
b, err := c.Get(bucket, key)
if err != nil {
	log.Fatal(err)
}
if b == nil {
	log.Fatal(fmt.Errorf("key does not exist"))
}
if string(value) != string(b) {
	log.Fatal(fmt.Errorf("returned value error"))
}
```
## Getting and putting bytes without saving in memory
calmcache writes and reads directly to and from files, and does not require values to be saved in memory.  For example:

```
func putAndGetBytes(cachePath, bucket, key string, value []byte) {
	c, err := calmcache.Open(cachePath)
	if err != nil {
		return err
	}
	defer c.Close()

	// Put
	OK, err := c.PutWithReader(bucket, key, bytes.NewReader(value), int64(len(value)))
	if err != nil {
		return err
	}
	if !OK {
		log.Fatal(fmt.Errorf("key already contains a value - delete the key first if you want to overwrite the value"))
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
	if string(b) != string(value) {
		return fmt.Errorf("returned value error")
	}
}
```
