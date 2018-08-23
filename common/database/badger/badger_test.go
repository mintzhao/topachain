// Copyright © 2018 Zhao Ming <mint.zhao.chiu@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package badger

import (
	"os"
	"strings"
	"testing"

	"github.com/mintzhao/topachain/common/database"
	"github.com/stretchr/testify/assert"
)

func TestBadgerDB_Open(t *testing.T) {
	assert.NoError(t, os.MkdirAll("./testdata", os.ModeDir))
	defer assert.NoError(t, os.RemoveAll("./testdata"))

	db, err := New("./testdata")
	assert.NoError(t, err)
	assert.NoError(t, db.Open())
	db.Close()
}

func TestBadgerDB_Get(t *testing.T) {
	assert.NoError(t, os.MkdirAll("./testdata", os.ModeDir))
	defer assert.NoError(t, os.RemoveAll("./testdata"))

	db, err := New("./testdata")
	assert.NoError(t, err)
	assert.NoError(t, db.Open())

	_, errnotfound := db.Get("test", "key1")
	assert.EqualError(t, errnotfound, database.ErrKeyNotFound.Error())

	assert.NoError(t, db.Set("test", "key1", []byte("value1")))
	retval, err := db.Get("test", "key1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), retval)

	db.Close()
}

func TestBadgerDB_Set(t *testing.T) {
	assert.NoError(t, os.MkdirAll("./testdata", os.ModeDir))
	defer assert.NoError(t, os.RemoveAll("./testdata"))

	db, err := New("./testdata")
	assert.NoError(t, err)
	assert.NoError(t, db.Open())

	assert.NoError(t, db.Set("test", "key1", []byte("value1")))
	retval, err := db.Get("test", "key1")
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), retval)

	db.Close()
}

func TestBadgerDB_Delete(t *testing.T) {
	assert.NoError(t, os.MkdirAll("./testdata", os.ModeDir))
	defer assert.NoError(t, os.RemoveAll("./testdata"))

	db, err := New("./testdata")
	assert.NoError(t, err)
	assert.NoError(t, db.Open())

	assert.NoError(t, db.Set("test", "key1", []byte("value1")))
	_, founderr := db.Get("test", "key1")
	assert.NoError(t, founderr)

	assert.NoError(t, db.Delete("test", "key1"))

	_, errnotfound := db.Get("test", "key1")
	assert.EqualError(t, errnotfound, database.ErrKeyNotFound.Error())

	db.Close()
}

func TestBadgerDB_Close(t *testing.T) {
	assert.NoError(t, os.MkdirAll("./testdata", os.ModeDir))
	defer assert.NoError(t, os.RemoveAll("./testdata"))

	db, err := New("./testdata")
	assert.NoError(t, err)
	assert.NoError(t, db.Open())

	db.Close()
}

func TestBadgerDB_Batch(t *testing.T) {
	assert.NoError(t, os.MkdirAll("./testdata", os.ModeDir))
	defer assert.NoError(t, os.RemoveAll("./testdata"))

	db, err := New("./testdata")
	assert.NoError(t, err)
	assert.NoError(t, db.Open())

	batch, err := db.NewBatch()
	assert.NoError(t, err)

	assert.NoError(t, batch.Set("test", "key1", []byte("value1")))
	assert.NoError(t, batch.Set("test", "key2", []byte("value2")))
	assert.NoError(t, batch.Set("test", "key3", []byte("value3")))

	assert.NoError(t, batch.Commit())

	// query
	_, found1 := db.Get("test", "key1")
	assert.NoError(t, found1)
	_, found2 := db.Get("test", "key2")
	assert.NoError(t, found2)
	_, found3 := db.Get("test", "key3")
	assert.NoError(t, found3)
	_, found4 := db.Get("test", "key4")
	assert.EqualError(t, found4, database.ErrKeyNotFound.Error())

	db.Close()
}

func TestBadgerDB_Iterator(t *testing.T) {
	assert.NoError(t, os.MkdirAll("./testdata", os.ModeDir))
	defer assert.NoError(t, os.RemoveAll("./testdata"))

	db, err := New("./testdata")
	assert.NoError(t, err)
	assert.NoError(t, db.Open())

	batch, err := db.NewBatch()
	assert.NoError(t, err)

	assert.NoError(t, batch.Set("test", "key1", []byte("value1")))
	assert.NoError(t, batch.Set("test", "key2", []byte("value2")))
	assert.NoError(t, batch.Set("test", "key3", []byte("value3")))

	assert.NoError(t, batch.Commit())

	// iterator
	it, err := db.NewIterator("test", "key")
	assert.NoError(t, err)

	cnt := 0
	for it.HasNext() {
		kv, err := it.Value()
		assert.NoError(t, err)

		cnt++
		assert.Equal(t, "test", kv.Bucket)
		assert.True(t, strings.HasPrefix(kv.Key, "key"))

		assert.NoError(t, it.Next())
	}
	assert.Equal(t, cnt, 3)

	db.Close()
}
