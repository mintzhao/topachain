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

import "github.com/dgraph-io/badger"

type badgerBatch struct {
	txn *badger.Txn
}

// Set set key/value pair
func (bt *badgerBatch) Set(bucket, key string, value []byte) error {
	err := bt.txn.Set(constructCompositeKey(bucket, key), value)
	if err != nil {
		bt.txn.Discard()
	}

	return err
}

// Delete delete key/value from database
func (bt *badgerBatch) Delete(bucket, key string) error {
	err := bt.txn.Delete(constructCompositeKey(bucket, key))
	if err != nil {
		bt.txn.Discard()
	}

	return err
}

// Commit commit all operations to database
func (bt *badgerBatch) Commit() error {
	return bt.txn.Commit(nil)
}

// Release release the resources hold by Batch
func (bt *badgerBatch) Release() error {
	bt.txn.Discard()
	return nil
}
