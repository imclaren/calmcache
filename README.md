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

## Low memory put and get

calmcache allows putting bytes into the cache using an io.Reader, and getting bytes from the cache by providing access to an os.File or io.Writer.  For example:
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
	OK, err = Delete(bucket, key)
	if err != nil {
		return err
	}
	if !OK {
		return fmt.Errorf("delete key error")
	}
}
```
## sqlite database access

The sqlite database can be queried directly.  For example:
```
func getAllInBucket(cachePath, bucket string) ([]cacheitem.Item, error) {
	c, err := calmcache.Open(cachePath)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	
	c.RLock()
	defer c.RUnlock()
	sqlString := "SELECT * FROM cache WHERE bucket = ? ORDER BY key ASC"
	var items []cacheitem.Item
	err := c.DB.Select(&items, c.DB.Rebind(sqlString), bucket)
	if err != nil {
		return nil, err
	}
	return items, nil
}

for _, i := range getAllInBucket("/path/to/cachePath", "mybucket") {
	fmt.Println(i.Key, i.Size, i.CreatedAt, i.UpdatedAt)
}

```
Once open, calmcache is designed be accessed concurrently.

Calmcache has user accessible sync.RWMutexes at the top level (e.g. c.Lock() and c.Unlock()) and at the database and filecache levels (e.g. c.DB.Lock() and c.FC.Lock())
