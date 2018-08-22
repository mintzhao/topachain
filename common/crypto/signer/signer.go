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
	"strings"
	"sync"

	"github.com/op/go-logging"
	"github.com/pkg/errors"
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
	RegisterSigner("ED25519", &ed25519Signer{})
}

// Signer contains signing functions
type Signer interface {

	// Sign signs msg using PrivateKey k.
	Sign(k crypto.PrivateKey, msg []byte, opts crypto.SignerOpts) ([]byte, error)

	// Verify verifies signature against key k and digest
	Verify(k crypto.PublicKey, signature, msg []byte, opts crypto.SignerOpts) (bool, error)
}

// RegisterSigner stores signer into signers, if signerName is already registered, return error
func RegisterSigner(signerName string, s Signer) error {
	_, loaded := signers.LoadOrStore(signerNameFmt(signerName), s)
	if loaded {
		// already registered
		logger.Warningf("signer %s already registered", signerName)
		return ErrSignerAlreadyRegistered
	}

	logger.Infof("signer %s registered", signerName)
	return nil
}

// DeRegisterSigner delete signer from signers, SHOULD ONLY USED IN TEST
func DeRegisterSigner(signerName string) {
	signers.Delete(signerNameFmt(signerName))
}

// GetHash return a hash function that already registered in hashes, if not return error
func GetSigner(signerName string) (Signer, error) {
	signer, ok := signers.Load(signerNameFmt(signerName))
	if !ok {
		// not found
		logger.Warningf("signer %s not found", signerName)
		return nil, ErrSignerNotFound
	}

	return signer.(Signer), nil
}

func signerNameFmt(signerName string) string {
	return strings.ToUpper(strings.TrimSpace(signerName))
}

var (
	// ErrSignerAlreadyRegistered indicated a signer already registered in to signers
	ErrSignerAlreadyRegistered = errors.New("signer already registered")

	// ErrSignerNotFound indicated a signer can not found in signers
	ErrSignerNotFound = errors.New("signer not found")

	// ErrNilSignerOptions indicated a nil signer options
	ErrNilSignerOptions = errors.New("invalid options, it must not be nil")

	// ErrInvalidSignerOptions indicated invalid signer options
	ErrInvalidSignerOptions = errors.New("invalid options")
)
