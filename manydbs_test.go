package calmcache

import (
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

func createAndPutKeys(t *testing.T, numRounds int, size int64) {
	db, cleanup := createDb(t)
	defer cleanup()

	for i := 0; i < numRounds; i++ {
	    value := make([]byte, size)
	    rand.Read(value)
		if _, err := db.Put("bucket", strconv.Itoa(i), value); err != nil {
			t.Fatal(err)
		}
	}
}

func TestManyDBs(t *testing.T) {
	for i := 0; i < 12; i++ {
		createAndPutKeys(t, 12, 16)
	}
	for i := 0; i < 10; i++ {
		createAndPutKeys(t, 10, 25*Megabyte)
	}
}
