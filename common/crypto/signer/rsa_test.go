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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRsaSigner_Sign(t *testing.T) {
	signer := &rsaSigner{}

	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.Nil(t, err)
	assert.NotNil(t, privKey)

	opts := &rsa.PSSOptions{
		Hash:       crypto.SHA256,
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	}

	msg := bytes.NewBufferString("this is a string for rsa signer test").Bytes()

	// sign error
	_, err = signer.Sign(privKey, msg, nil)
	assert.EqualError(t, err, ErrNilSignerOptions.Error())

	// sign
	sig, err := signer.Sign(privKey, msg, opts)
	assert.Nil(t, err)
	assert.NotEmpty(t, sig)

	// verify
	_, err = signer.Verify(privKey.Public(), sig, msg, nil)
	assert.EqualError(t, err, ErrNilSignerOptions.Error())

	// verify
	ok, err := signer.Verify(privKey.Public(), sig, msg, opts)
	assert.NoError(t, err)
	assert.Equal(t, true, ok)

	// verify
	notok, err := signer.Verify(privKey.Public(), sig, msg, &rsa.PSSOptions{Hash: crypto.MD5})
	assert.Error(t, err)
	assert.Equal(t, false, notok)
}

func TestRsaSigner_Verify(t *testing.T) {
	signer := &rsaSigner{}

	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	assert.Nil(t, err)
	assert.NotNil(t, privKey)

	opts := &rsa.PSSOptions{
		Hash:       crypto.SHA256,
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	}

	msg := bytes.NewBufferString("this is a string for rsa signer test").Bytes()

	// sign
	sig, err := signer.Sign(privKey, msg, opts)
	assert.Nil(t, err)
	assert.NotEmpty(t, sig)

	// verify ok
	ok, err := signer.Verify(privKey.Public(), sig, msg, opts)
	assert.NoError(t, err)
	assert.Equal(t, true, ok)

	// verify notok
	_, err = signer.Verify(privKey.Public(), sig, msg, nil)
	assert.Error(t, err)

	// verify notok
	notok, err := signer.Verify(privKey.Public(), sig[1:], msg, opts)
	assert.Error(t, err)
	assert.Equal(t, false, notok)
}

func BenchmarkRsaSigner_Sign(b *testing.B) {
	signer := &rsaSigner{}

	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		b.FailNow()
	}

	opts := &rsa.PSSOptions{
		Hash:       crypto.SHA256,
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	}

	msg := bytes.NewBufferString("this is a string for rsa signer test").Bytes()

	// sign
	for i := 0; i <= b.N; i++ {
		signer.Sign(privKey, msg, opts)
	}
}

func BenchmarkRsaSigner_Verify(b *testing.B) {
	signer := &rsaSigner{}

	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		b.FailNow()
	}

	opts := &rsa.PSSOptions{
		Hash:       crypto.SHA256,
		SaltLength: rsa.PSSSaltLengthEqualsHash,
	}

	msg := bytes.NewBufferString("this is a string for rsa signer test").Bytes()

	// sign
	sig, err := signer.Sign(privKey, msg, opts)
	if err != nil {
		b.FailNow()
	}

	for i := 0; i < b.N; i++ {
		signer.Verify(privKey.Public(), sig, msg, opts)
	}
}
