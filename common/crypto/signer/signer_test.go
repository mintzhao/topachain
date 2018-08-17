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
	"crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

// used for test only
type testsigner struct {
}

func (ts *testsigner) Sign(k crypto.PrivateKey, msg []byte, opts crypto.SignerOpts) ([]byte, error) {
	return nil, nil
}

func (ts *testsigner) Verify(k crypto.PublicKey, signature, msg []byte, opts crypto.SignerOpts) (bool, error) {
	return false, nil
}

func TestRegisterSigner(t *testing.T) {
	assert.EqualError(t, RegisterSigner("RSA", &rsaSigner{}), (&ErrSignerAlreadyRegistered{}).Error())
	assert.EqualError(t, RegisterSigner("RSA ", &rsaSigner{}), (&ErrSignerAlreadyRegistered{}).Error())
	assert.EqualError(t, RegisterSigner(" rsa  ", &rsaSigner{}), (&ErrSignerAlreadyRegistered{}).Error())

	tSigner := &testsigner{}
	assert.NoError(t, RegisterSigner("test", tSigner))
	retSigner, err := GetSigner("Test")
	assert.NoError(t, err)
	assert.Equal(t, tSigner, retSigner)

	// delete test signer
	DeRegisterSigner("test")
}

func TestGetSigner(t *testing.T) {
	_, err := GetSigner("Test")
	assert.EqualError(t, err, (&ErrSignerNotFound{signerName: "Test"}).Error())

	tSigner := &testsigner{}
	assert.NoError(t, RegisterSigner("test", tSigner))
	retSigner, err := GetSigner(" test ")
	assert.NoError(t, err)
	assert.Equal(t, tSigner, retSigner)

	// delete test signer
	DeRegisterSigner("test")

	_, err = GetSigner("rsa")
	assert.NoError(t, err)
}
