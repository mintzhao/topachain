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
	"bytes"
	"context"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/mintzhao/topachain/common/database"
	"github.com/pkg/errors"
)

var (
	compositeKeySep = []byte{0x00}
)

type BadgerDB struct {
	db     *badger.DB
	ctx    context.Context
	cancel context.CancelFunc
}

func New(dir string) (database.Database, error) {
	opts := badger.DefaultOptions

	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open database")
	}

	ctx, cancel := context.WithCancel(context.Background())
	bdb := &BadgerDB{
		db:     db,
		ctx:    ctx,
		cancel: cancel,
	}
	go bdb.gc()

	return bdb, nil
}

// Open open database handler
func (db *BadgerDB) Open() error {
	return nil
}

// Get get key/value pair from database
func (db *BadgerDB) Get(bucket, key string) ([]byte, error) {
	var getBytes []byte
	if err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(constructCompositeKey(bucket, key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return database.ErrKeyNotFound
			}

			return err
		}

		buf, err := item.Value()
		if err != nil {
			return err
		}
		getBytes = bytes.NewBuffer(buf).Bytes()

		return nil
	}); err != nil {
		return nil, err
	}

	return getBytes, nil
}

// Set set key/value pair
func (db *BadgerDB) Set(bucket, key string, value []byte) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Set(constructCompositeKey(bucket, key), value)
	})
}

// Delete delete key/value from database
func (db *BadgerDB) Delete(bucket, key string) error {
	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(constructCompositeKey(bucket, key))
	})
}

// NewBatch returns a Batch interface to handle multi key/value pairs writes.
func (db *BadgerDB) NewBatch() (database.Batch, error) {
	return &badgerBatch{
		txn: db.db.NewTransaction(true),
	}, nil
}

// NewIterator iterate over a key prefix
func (db *BadgerDB) NewIterator(bucket, prefix string) (database.Iterator, error) {
	opts := badger.DefaultIteratorOptions
	txn := db.db.NewTransaction(false)
	it := txn.NewIterator(opts)

	prefixBytes := constructCompositeKey(bucket, prefix)
	it.Seek(prefixBytes)

	return &badgerIterator{
		txn:    txn,
		it:     it,
		prefix: prefixBytes,
	}, nil
}

// Close close database
func (db *BadgerDB) Close() error {
	db.cancel()
	return db.db.Close()
}

// Badger values need to be garbage collected, because of two reasons:
// Badger keeps values separately from the LSM tree. This means that the compaction operations that clean up the LSM tree do not touch the values at all. Values need to be cleaned up separately.
// Concurrent read/write transactions could leave behind multiple values for a single key, because they are stored with different versions. These could accumulate, and take up unneeded space beyond the time these older versions are needed.
func (db *BadgerDB) gc() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		again:
			err := db.db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		case <-db.ctx.Done():
			return
		}
	}
}

func constructCompositeKey(bucket string, key string) []byte {
	return append(append([]byte(bucket), compositeKeySep...), []byte(key)...)
}

func splitCompositeKey(compositeKey []byte) (string, string) {
	split := bytes.SplitN(compositeKey, compositeKeySep, 2)
	return string(split[0]), string(split[1])
}
