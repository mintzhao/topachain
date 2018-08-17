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
package crypto

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	msg := bytes.NewBufferString("this is used for sha256 test").Bytes()
	digest := "069c725da1725e638c8f1d901d4c4245d8b68cf571bbe445cb7be8709e5b59a2"

	retHash, err := Hash(msg)
	assert.NoError(t, err)
	assert.EqualValues(t, digest, hex.EncodeToString(retHash))

	md5msg := bytes.NewBufferString("this is used for md5 test").Bytes()
	// MD5 ("this is used for md5 test") = 14c4063d61bd57528837784838ea5a79
	md5digest := "14c4063d61bd57528837784838ea5a79"

	md5retHash, err := Hash(md5msg, "md5")
	assert.NoError(t, err)
	assert.EqualValues(t, md5digest, hex.EncodeToString(md5retHash))
}

func TestSignVerify(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(t, err)
	assert.NotNil(t, privKey)

	msg := bytes.NewBufferString("this is a string used for test ecdsaSigner").Bytes()

	// sign
	sig, err := Sign(privKey, msg, nil)
	assert.Nil(t, err)
	assert.NotNil(t, sig)

	// verify
	ok, err := Verify(privKey.Public(), sig, msg, nil)
	assert.NoError(t, err)
	assert.Equal(t, true, ok)

	// rsa
	rsaPrivKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.NoError(t, err)
	assert.NotNil(t, rsaPrivKey)

	rsamsg := bytes.NewBufferString("this is a string used for test rsaSigner").Bytes()
	opts := &rsa.PSSOptions{Hash: crypto.SHA256, SaltLength: rsa.PSSSaltLengthEqualsHash}

	// sign
	rsasig, err := Sign(rsaPrivKey, rsamsg, opts, "rsa")
	assert.Nil(t, err)
	assert.NotNil(t, rsasig)

	// verify
	rsaok, err := Verify(rsaPrivKey.Public(), rsasig, rsamsg, opts, "rsa")
	assert.NoError(t, err)
	assert.Equal(t, true, rsaok)
}
