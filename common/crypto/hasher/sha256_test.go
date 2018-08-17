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
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256Hasher_Hash(t *testing.T) {
	msg := bytes.NewBufferString("this is used for sha256 test").Bytes()
	digest := "069c725da1725e638c8f1d901d4c4245d8b68cf571bbe445cb7be8709e5b59a2"

	rethash, err := (&sha256Hasher{}).Hash(msg)
	assert.NoError(t, err)

	t.Logf("ret hash %x", rethash)

	assert.EqualValues(t, digest, hex.EncodeToString(rethash))
}

func BenchmarkSha256Hasher_Hash(b *testing.B) {
	msg := bytes.NewBufferString("this is used for sha256 test").Bytes()
	hasher := &sha256Hasher{}

	for i := 0; i < b.N; i++ {
		hasher.Hash(msg)
	}
}
