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
	"github.com/dgraph-io/badger"
	"github.com/mintzhao/topachain/common/database"
)

type badgerIterator struct {
	txn    *badger.Txn
	it     *badger.Iterator
	prefix []byte
}

// HasNext return true if iterator isn't over, otherwise false
func (iter *badgerIterator) HasNext() bool {
	return iter.it.ValidForPrefix(iter.prefix)
}

// Value return matched key/value pair
func (iter *badgerIterator) Value() (*database.KeyValePair, error) {
	item := iter.it.Item()
	k := item.Key()
	v, err := item.Value()
	if err != nil {
		return nil, err
	}

	bucket, key := splitCompositeKey(k)
	return &database.KeyValePair{
		Bucket: bucket,
		Key:    key,
		Value:  v,
	}, nil
}

// Next move pointer to next matched value
func (iter *badgerIterator) Next() error {
	iter.it.Next()
	return nil
}

// Close close iterator
func (iter *badgerIterator) Close() error {
	iter.it.Close()
	iter.txn.Discard()
	return nil
}
