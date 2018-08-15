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
	"fmt"
	"strings"
	"sync"

	"github.com/op/go-logging"
)

var (
	// supported signer
	signers sync.Map

	// logger
	logger = logging.MustGetLogger("crypto/signer")
)

func init() {
	// register default signers
	RegisterSigner("RSA", &rsaSigner{})
	RegisterSigner("ECDSA", &ecdsaSigner{})
}

// Signer contains signing functions
type Signer interface {

	// Sign signs digest using PrivateKey k.
	Sign(k crypto.PrivateKey, digest []byte, opts crypto.SignerOpts) ([]byte, error)

	// Verify verifies signature against key k and digest
	//Verify(k crypto.PublicKey, signature, digest []byte) (bool, error)
}

// RegisterSigner stores signer into signers, if signerName is already registered, return error
func RegisterSigner(signerName string, s Signer) error {
	_, loaded := signers.LoadOrStore(strings.ToUpper(signerName), s)
	if loaded {
		// already registered
		logger.Warningf("signer %s already registered", signerName)
		return &ErrSignerAlreadyRegistered{signerName: signerName}
	}

	logger.Infof("signer %s registered", signerName)
	return nil
}

// GetHash return a hash function that already registered in hashes, if not return error
func GetSigner(signerName string) (Signer, error) {
	signer, ok := signers.Load(strings.ToUpper(signerName))
	if !ok {
		// not found
		logger.Warningf("signer %s not found", signerName)
		return nil, &ErrSignerNotFound{signerName: signerName}
	}

	return signer.(Signer), nil
}

// ErrSignerAlreadyRegistered indicated a signer already registered in to signers
type ErrSignerAlreadyRegistered struct {
	signerName string
}

// Error output error message
func (err *ErrSignerAlreadyRegistered) Error() string {
	return fmt.Sprintf("signer %s has already registered", err.signerName)
}

// ErrSignerNotFound indicated a signer can not found in signers
type ErrSignerNotFound struct {
	signerName string
}

// Error output error message
func (err *ErrSignerNotFound) Error() string {
	return fmt.Sprintf("signer %s not found", err.signerName)
}
