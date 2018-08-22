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
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"math/big"

	"github.com/pkg/errors"
)

type ecdsaSigner struct {
}

type ecdsaSignature struct {
	R, S *big.Int
}

// Sign signs digest using PrivateKey k.
func (es *ecdsaSigner) Sign(k crypto.PrivateKey, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.(*ecdsa.PrivateKey), digest)
	if err != nil {
		return nil, err
	}

	return asn1.Marshal(ecdsaSignature{r, s})
}

// Verify verifies signature against key k and digest
func (es *ecdsaSigner) Verify(k crypto.PublicKey, signature, digest []byte, opts crypto.SignerOpts) (bool, error) {
	sig := new(ecdsaSignature)
	_, err := asn1.Unmarshal(signature, sig)
	if err != nil {
		return false, errors.Wrap(err, "unmarshal ecdsa signature error")
	}

	return ecdsa.Verify(k.(*ecdsa.PublicKey), digest, sig.R, sig.S), nil
}
