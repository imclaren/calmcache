package calmcache

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"strconv"
)

const (
	Kilobyte = int64(1024)
	Megabyte = Kilobyte * int64(1024)
	Gigabyte = Megabyte * int64(1024)
	Terabyte = Gigabyte * int64(1024)
)

func createDb(t *testing.T) (Cache, func()) {
	// First, create a temporary directory to be used for the duration of
	// this test.
	tempDirName, err := ioutil.TempDir("", "bboltmemtest")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	path := filepath.Join(tempDirName, "testdb.db")

	//bdb, err := Open(path, 0600, nil)
	bdb, err := Open(path)
	if err != nil {
		t.Fatalf("error creating bbolt db: %v", err)
	}

	cleanup := func() {
		bdb.Close()
		os.RemoveAll(tempDirName)
	}

	return bdb, cleanup
}

func createAndPutKeys(t *testing.T) {
	t.Parallel()

	db, cleanup := createDb(t)
	defer cleanup()

	bucketName := "bucket"

	for i := 0; i < 100; i++ {
		//var key [16]byte
		var key [25*Megabyte]byte

		rand.Read(key[:])
		if _, err := db.Put(bucketName, strconv.Itoa(i), key[:]); err != nil {
			t.Fatal(err)
		}
	}
}

func TestManyDBs(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("%d", i), createAndPutKeys)
	}
}
