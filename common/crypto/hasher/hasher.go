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
package hasher

import (
	"fmt"
	hash2 "hash"
	"strings"
	"sync"

	"github.com/op/go-logging"
)

var (
	// supported hash func
	hashes sync.Map

	// logger
	logger = logging.MustGetLogger("crypto/hasher")
)

func init() {
	// register default hash func
	RegisterHasher("MD5", &md5Hasher{})
	RegisterHasher("SHA256", &sha256Hasher{})
}

// hash hashes messages msg using hf
func hash(hf hash2.Hash, msg []byte) ([]byte, error) {
	hf.Reset()
	hf.Write(msg)
	return hf.Sum(nil), nil
}

// hasher contains hash related functions
type Hasher interface {

	// Hash hashes messages msg, can customize hash function, default is SHA256.
	Hash(msg []byte) ([]byte, error)
}

// RegisterHasher stores hash function into hashes, if hashName is already registered, return error
func RegisterHasher(hashName string, hasher Hasher) error {
	_, loaded := hashes.LoadOrStore(strings.ToUpper(hashName), f)
	if loaded {
		// already registered
		logger.Warningf("hasher %s already registered", hashName)
		return &ErrHashAlreadyRegistered{hashName: hashName}
	}

	logger.Infof("hasher %s registered", hashName)
	return nil
}

// GetHasher return a hash function that already registered in hashes, if not return error
func GetHasher(hashName string) (Hasher, error) {
	f, ok := hashes.Load(strings.ToUpper(hashName))
	if !ok {
		// not found
		logger.Warningf("hasher %s not found", hashName)
		return nil, &ErrHashNotFound{hashName: hashName}
	}

	return f.(Hasher), nil
}

// ErrHashAlreadyRegistered indicated a hash function already registered in to hashes
type ErrHashAlreadyRegistered struct {
	hashName string
}

// Error output error message
func (err *ErrHashAlreadyRegistered) Error() string {
	return fmt.Sprintf("hash %s has already registered", err.hashName)
}

// ErrHashNotFound indicated a hash function can not found in hashes
type ErrHashNotFound struct {
	hashName string
}

// Error output error message
func (err *ErrHashNotFound) Error() string {
	return fmt.Sprintf("hash %s not found", err.hashName)
}
