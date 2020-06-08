package calmcache

import (
	"testing"
	"bytes"
	"io/ioutil"

	assert "github.com/stretchr/testify/require"
)

const (
	cachePath = "testdata/cache"
	bucket = "testbucket"
	key = "testkey"
)

func TestOpen(t *testing.T) {
	c, err := Open(cachePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		c.DeleteCache()
	}()
	err = c.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestSize(t *testing.T) {
	c, err := Open(cachePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		c.DeleteCache()
	}()

	size, err := c.DB.Size()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, size == 0)

	inBytes := []byte{'1', '2', '3'}
	_, err = c.Put(bucket, key, inBytes)
	if err != nil {
		t.Fatal(err)
	}
	size, err = c.DB.Size()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, size == 3)

	err = c.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMultiple(t *testing.T) {
	c, err := Open(cachePath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		c.DeleteCache()
	}()

	c.TestSelect(t)
	c.TestPut(t)
	c.TestPutAndDelete(t)
	c.TestReaderAndGetPath(t)

	err = c.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func (c *Cache) TestSelect(t *testing.T) {

	err := c.DeleteBucket(bucket)
	if err != nil {
		t.Fatal(err)
	}

	inBytes := []byte{'1', '2', '3'}
	_, err = c.Put(bucket, key, inBytes)
	if err != nil {
		t.Fatal(err)
	}
	outBytes, err := c.Get(bucket, key)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "123", string(outBytes))

	allKeys, err := c.AllKeys(bucket)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []string{"testkey"}, allKeys)
}

func (c *Cache) TestPut(t *testing.T) {

	// Test success cases.
	testCases := []struct {
		keys       []string
		bytes      []byte
		beforeKeys []string
		popKey     string
		afterKeys  []string
	}{
		{[]string{"testkey"}, []byte("123"), []string{"testkey"}, "testkey", []string{}},
		{[]string{"testkey", "testkey2"}, []byte("123"), []string{"testkey", "testkey2"}, "testkey", []string{"testkey2"}},
		{[]string{"testkey", "testkey2", "testkey3"}, []byte("123"), []string{"testkey", "testkey2", "testkey3"}, "testkey", []string{"testkey2", "testkey3"}},

		{[]string{"testkey"}, []byte("123456"), []string{"testkey"}, "testkey", []string{}},
		{[]string{"testkey"}, []byte(""), []string{"testkey"}, "testkey", []string{}},
	}

	for _, testCase := range testCases {
		err := c.DeleteBucket(bucket)
		if err != nil {
			t.Fatal(err)
		}

		for _, key := range testCase.keys {
			_, err = c.Put(bucket, key, testCase.bytes)
			if err != nil {
				t.Fatal(err)
			}
		}

		allKeys, err := c.AllKeys(bucket)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testCase.beforeKeys, allKeys)
	}
}

func (c *Cache) TestPutAndDelete(t *testing.T) {

	err := c.DeleteBucket(bucket)
	if err != nil {
		t.Fatal(err)
	}

	// Test success cases.
	testCases := []struct {
		key        string
		value      []byte
		beforeKeys []string
		afterKeys  []string
	}{
		{"testkey", []byte{'1', '2', '3'}, []string{"testkey"}, []string{}},
		{"testkey", []byte("456"), []string{"testkey"}, []string{}},
	}

	for _, testCase := range testCases {
		err := c.DeleteBucket(bucket)
		if err != nil {
			t.Fatal(err)
		}

		_, err = c.Put(bucket, testCase.key, testCase.value)
		if err != nil {
			t.Fatal(err)
		}

		fileKeys, err := c.AllKeys(bucket)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testCase.beforeKeys, fileKeys)

		_, err = c.Delete(bucket, testCase.key)
		if err != nil {
			t.Fatal(err)
		}

		fileKeys, err = c.AllKeys(bucket)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testCase.afterKeys, fileKeys)
		err = c.DeleteBucket(bucket)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test that 10 entries are not too slow
	for i := 0; i < 10; i++ {
		for _, testCase := range testCases {
			err := c.DeleteBucket(bucket)
			if err != nil {
				t.Fatal(err)
			}

			_, err = c.Put(bucket, testCase.key, testCase.value)
			if err != nil {
				t.Fatal(err)
			}

			fileKeys, err := c.AllKeys(bucket)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.beforeKeys, fileKeys)

			_, err = c.Delete(bucket, testCase.key)
			if err != nil {
				t.Fatal(err)
			}

			fileKeys, err = c.AllKeys(bucket)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, testCase.afterKeys, fileKeys)
		}
	}
}

func (c *Cache) TestReaderAndGetPath(t *testing.T) {

	// Test success cases.
	testCases := []struct {
		keys       []string
		bytes      []byte
		beforeKeys []string
		popKey     string
		afterKeys  []string
	}{
		{[]string{"testkey"}, []byte("123"), []string{"testkey"}, "testkey", []string{}},
		{[]string{"testkey", "testkey2"}, []byte("123"), []string{"testkey", "testkey2"}, "testkey", []string{"testkey2"}},
		{[]string{"testkey", "testkey2", "testkey3"}, []byte("123"), []string{"testkey", "testkey2", "testkey3"}, "testkey", []string{"testkey2", "testkey3"}},

		{[]string{"testkey"}, []byte("123456"), []string{"testkey"}, "testkey", []string{}},
		{[]string{"testkey"}, []byte(""), []string{"testkey"}, "testkey", []string{}},
	}

	for _, testCase := range testCases {
		err := c.DeleteBucket(bucket)
		if err != nil {
			t.Fatal(err)
		}
		for _, key := range testCase.keys {
			_, err = c.PutWithReader(bucket, key, bytes.NewReader(testCase.bytes), int64(len(testCase.bytes)))
			if err != nil {
				t.Fatal(err)
			}

			OK, fullPath, _, err := c.GetPathAndLock(bucket, key) 
			if err != nil {
				c.GetPathUnlock()
				t.Fatal(err)
			}
			assert.True(t, OK)
			b, err := ioutil.ReadFile(fullPath)
			if err != nil {
				c.GetPathUnlock()
				t.Fatal(err)
			}
			c.GetPathUnlock()
			assert.Equal(t, b, testCase.bytes)
		}
		allKeys, err := c.AllKeys(bucket)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, testCase.beforeKeys, allKeys)
	}
}


