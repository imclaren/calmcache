# calmcache
calmcache is a low memory golang disk cache ([GoDoc](https://godoc.org/github.com/imclaren/calmcache)).  The filesystem cache is managed using a sqlite database.

## Example

```
import (
	"github.com/imclaren/calmcache"
	"log"
)

cachePath := "/path/to/cachePath"
bucket := "mybucket"
key := "mykey"
value := []byte{'a', 'b', 'c'}

c, err := calmcache.Open(cachePath)
if err != nil {
	log.Fatal(err)
}
OK, err := c.Put(bucket, key, value)
if err != nil {
	log.Fatal(err)
}
if !OK {
	log.Fatal("key already contains a value - delete the key first if you want to overwrite the value")
}
b, err := c.Get(bucket, key)
if err != nil {
	log.Fatal(err)
}
if b == nil {
	log.Fatal("key does not exist")
}
if string(value) != string(b) {
	log.Fatal("returned value error")
}
OK, err = c.Delete(bucket, key)
if err != nil {
	log.Fatal(err)
}
if !OK {
	log.Fatal("delete key error")
}
```

## Low memory put and get

calmcache allows putting bytes into the cache using an io.Reader, and getting bytes from the cache by providing access to an os.File or io.Writer.  For example:
```
import (
	"github.com/imclaren/calmcache"
	"fmt"
)

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
		return fmt.Errorf("key already contains a value - delete the key first if you want to overwrite the value")
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
## sqlite database access

The sqlite database can be queried directly.  For example:
```
import (
	"github.com/imclaren/calmcache"
	"github.com/imclaren/calmcache/cacheitem"
	"github.com/imclaren/calmcache/dbcache"
	"database/sql"
	"fmt"
)

// getNewestInBucket returns the newest (i.e. most recently accessed) database item
func getNewestInBucket(db *dbcache.DB, bucket string) (i *cacheitem.Item, err error) {
	db.RLock()
	defer db.RUnlock()

	sqlString := "SELECT * FROM cache WHERE bucket = ? ORDER BY updated_at DESC LIMIT 1"
	var newItem cacheitem.Item 
	err = db.QueryRowx(db.Rebind(sqlString), bucket).StructScan(&newItem)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &newItem, nil
}

c, err := calmcache.Open(cachePath)
if err != nil {
	return err
}
defer c.Close()
item, err := getNewestInBucket(&c.DB, bucket)
if err != nil {
	return err
}
fmt.Println(i.Key, i.Size, i.CreatedAt, i.UpdatedAt)
```
Once open, calmcache is designed be accessed concurrently.

Calmcache has user accessible sync.RWMutexes at the top level (e.g. c.Lock() and c.Unlock()) and at the database and filecache levels (e.g. c.DB.Lock() and c.FC.Lock())
