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
package database

import "github.com/pkg/errors"

// Database contains all the funcs to handle database
type Database interface {
	// Open open database handler
	Open() error

	// Get get key/value pair from database
	Get(bucket, key string) ([]byte, error)

	// Set set key/value pair
	Set(bucket, key string, value []byte) error

	// Delete delete key/value from database
	Delete(bucket, key string) error

	// NewBatch returns a Batch interface to handle multi key/value pairs writes.
	NewBatch() (Batch, error)

	// NewIterator iterate over a key prefix
	NewIterator(bucket, prefix string) (Iterator, error)

	// Close close database
	Close() error
}

// Batch collect operations together send to database atomically.
type Batch interface {
	// Set set key/value pair
	Set(bucket, key string, value []byte) error

	// Delete delete key/value from database
	Delete(bucket, key string) error

	// Commit commit all operations to database
	Commit() error

	// Release release the resources hold by Batch
	Release() error
}

// Iterator seek key by prefix
type Iterator interface {
	// HasNext return true if iterator isn't over, otherwise false
	HasNext() bool

	// Value return matched key/value pair
	Value() (*KeyValePair, error)

	// Next move pointer to next matched value
	Next() error

	// Close close iterator
	Close() error
}

// KeyValuePair
type KeyValePair struct {
	Bucket, Key string
	Value       []byte
}

var (
	// ErrKeyNotFound returned if no key stored in database
	ErrKeyNotFound = errors.New("key not found")
)
