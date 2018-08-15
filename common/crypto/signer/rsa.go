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
	"crypto/rand"
	"crypto/rsa"
)

type rsaSigner struct {
}

// Sign signs digest using PrivateKey k.
func (rs *rsaSigner) Sign(k crypto.PrivateKey, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return k.(*rsa.PrivateKey).Sign(rand.Reader, digest, opts)
}

// Verify verifies signature against key k and digest
func (rs *rsaSigner) Verify(k crypto.PublicKey, signature, digest []byte) (bool, error) {
	return false, nil
}
