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
	hash2 "hash"
	"strings"
	"sync"

	_ "github.com/mintzhao/topachain/common/logging"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
)

var (
	// supported hash func
	hashes sync.Map

	// logger
	logger = logging.MustGetLogger("crypto/hasher")
)

func init() {
	// register default hash func
	RegisterHasher("MD5", &MD5Hasher{})
	RegisterHasher("SHA256", &SHA256Hasher{})
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
func RegisterHasher(hasherName string, hasher Hasher) error {
	_, loaded := hashes.LoadOrStore(hasherNameFmt(hasherName), hasher)
	if loaded {
		// already registered
		logger.Warningf("hasher %s already registered", hasherName)
		return ErrHasherAlreadyRegistered
	}

	logger.Debugf("hasher %s registered", hasherName)
	return nil
}

// DeRegisterHasher delete hasher form hashes, SHOULD ONLY USED IN TEST
func DeRegisterHasher(hasherName string) {
	hashes.Delete(hasherNameFmt(hasherName))
}

// GetHasher return a hash function that already registered in hashes, if not return error
func GetHasher(hasherName string) (Hasher, error) {
	f, ok := hashes.Load(hasherNameFmt(hasherName))
	if !ok {
		// not found
		logger.Warningf("hasher %s not found", hasherName)
		return nil, ErrHasherNotFound
	}

	return f.(Hasher), nil
}

func hasherNameFmt(hasherName string) string {
	return strings.ToUpper(strings.TrimSpace(hasherName))
}

var (
	// ErrHasherNotFound is returned when hasher not found
	ErrHasherNotFound = errors.New("hasher not found")

	// ErrHasherAlreadyRegistered is returned when hasher already registered
	ErrHasherAlreadyRegistered = errors.New("hasher already registered")
)
