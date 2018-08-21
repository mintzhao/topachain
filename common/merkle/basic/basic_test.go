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
package basic

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/mintzhao/topachain/common/crypto/hasher"
	"github.com/stretchr/testify/assert"
)

type testvalue []byte

func (tv testvalue) Hash() ([]byte, error) {
	h := sha256.New()
	h.Write(tv)

	return h.Sum(nil), nil
}

func (tv testvalue) Equals(other Value) bool {
	return bytes.Equal(tv, other.(testvalue))
}

var (
	testValues = []Value{
		testvalue("1"),
		testvalue("2"),
		testvalue("3"),
	}
)

func TestBasicMerkletree(t *testing.T) {
	tree, err := New(testValues, &hasher.MD5Hasher{})
	assert.NoError(t, err)

	assert.Equal(t, true, tree.VerifyRoot())
	assert.Equal(t, true, tree.VerifyValue(tree.root.hash, testvalue("1")))
	assert.Equal(t, true, tree.VerifyValue(tree.root.hash, testvalue("2")))
	assert.Equal(t, true, tree.VerifyValue(tree.root.hash, testvalue("3")))
	assert.Equal(t, false, tree.VerifyValue(tree.root.hash, testvalue("4")))

	root1 := tree.Root()

	assert.NoError(t, tree.Reconstruct())
	assert.Equal(t, root1, tree.Root())

	assert.NoError(t, tree.Reconstruct([]interface{}{
		testvalue("2"),
		testvalue("3"),
		testvalue("4"),
		testvalue("5"),
	}...))

	assert.NotEqual(t, root1, tree.Root())
}

func BenchmarkBasicMerkletree_Reconstruct(b *testing.B) {
	tree, err := New(testValues, &hasher.MD5Hasher{})
	if err != nil {
		b.FailNow()
	}

	for i := 0; i < b.N; i++ {
		tree.Reconstruct()
	}
}

func BenchmarkBasicMerkletree_VerifyRoot(b *testing.B) {
	tree, err := New(testValues, &hasher.MD5Hasher{})
	if err != nil {
		b.FailNow()
	}

	for i := 0; i < b.N; i++ {
		tree.VerifyRoot()
	}
}

func BenchmarkBasicMerkletree_VerifyValue(b *testing.B) {
	tree, err := New(testValues, &hasher.MD5Hasher{})
	if err != nil {
		b.FailNow()
	}

	for i := 0; i < b.N; i++ {
		tree.VerifyValue(tree.Root(), testvalue("2"))
	}
}
