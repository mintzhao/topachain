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

	"golang.org/x/crypto/ed25519"
)

type ed25519Signer struct {
}

// Sign signs digest using PrivateKey k.
func (es *ed25519Signer) Sign(k crypto.PrivateKey, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return ed25519.Sign(k.(ed25519.PrivateKey), digest), nil
}

// Verify verifies signature against key k and digest
func (es *ed25519Signer) Verify(k crypto.PublicKey, signature, digest []byte, opts crypto.SignerOpts) (bool, error) {
	return ed25519.Verify(k.(ed25519.PublicKey), digest, signature), nil
}
