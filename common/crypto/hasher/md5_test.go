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

func TestMd5Hasher_Hash(t *testing.T) {
	msg := bytes.NewBufferString("this is used for md5 test").Bytes()
	// MD5 ("this is used for md5 test") = 14c4063d61bd57528837784838ea5a79
	digest := "14c4063d61bd57528837784838ea5a79"

	rethash, err := (&md5Hasher{}).Hash(msg)
	assert.NoError(t, err)

	t.Logf("ret hash %x", rethash)

	assert.EqualValues(t, digest, hex.EncodeToString(rethash))
}

func BenchmarkMd5Hasher_Hash(b *testing.B) {
	msg := bytes.NewBufferString("this is used for md5 test").Bytes()
	hasher := &md5Hasher{}

	for i := 0; i < b.N; i++ {
		hasher.Hash(msg)
	}
}
