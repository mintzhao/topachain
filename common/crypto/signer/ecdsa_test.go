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
package signer

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcdsaSigner_Sign(t *testing.T) {
	signer := &ecdsaSigner{}

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(t, err)
	assert.NotNil(t, privKey)

	msg := bytes.NewBufferString("this is a string used for test ecdsaSigner").Bytes()

	// sign
	sig, err := signer.Sign(privKey, msg, nil)
	assert.Nil(t, err)
	assert.NotNil(t, sig)

	// verify ok
	ok, err := signer.Verify(privKey.Public(), sig, msg, nil)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)

	// verify not ok
	notok, err := signer.Verify(privKey.Public(), sig[1:], msg, nil)
	assert.NotNil(t, err)
	assert.Equal(t, false, notok)
}

func TestEcdsaSigner_Verify(t *testing.T) {
	signer := &ecdsaSigner{}

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(t, err)
	assert.NotNil(t, privKey)

	msg := bytes.NewBufferString("this is a string used for test ecdsaSigner").Bytes()

	// sign
	sig, err := signer.Sign(privKey, msg, nil)
	assert.Nil(t, err)
	assert.NotNil(t, sig)

	// verify ok
	ok, err := signer.Verify(privKey.Public(), sig, msg, nil)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)

	// verify not ok
	notok, err := signer.Verify(privKey.Public(), sig, bytes.NewBufferString("this is another string").Bytes(), nil)
	assert.Nil(t, err)
	assert.Equal(t, false, notok)
}

func BenchmarkEcdsaSigner_Sign(b *testing.B) {
	signer := &ecdsaSigner{}

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		b.FailNow()
	}

	msg := bytes.NewBufferString("this is a string used for ecdsa signer benchmark").Bytes()
	for i := 0; i <= b.N; i++ {
		signer.Sign(privKey, msg, nil)
	}
}

func BenchmarkEcdsaSigner_Verify(b *testing.B) {
	signer := &ecdsaSigner{}

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		b.FailNow()
	}

	msg := bytes.NewBufferString("this is a string used for ecdsa signer benchmark").Bytes()
	sig, err := signer.Sign(privKey, msg, nil)
	if err != nil {
		b.FailNow()
	}

	for i := 0; i <= b.N; i++ {
		signer.Verify(privKey.Public(), sig, msg, nil)
	}
}
