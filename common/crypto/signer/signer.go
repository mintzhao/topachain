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
		return &ErrSignerAlreadyRegistered{signerName: signerName}
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
		return nil, &ErrSignerNotFound{signerName: signerName}
	}

	return signer.(Signer), nil
}

func signerNameFmt(signerName string) string {
	return strings.ToUpper(strings.TrimSpace(signerName))
}

// ErrSignerAlreadyRegistered indicated a signer already registered in to signers
type ErrSignerAlreadyRegistered struct {
	signerName string
}

// Error output error message
func (err *ErrSignerAlreadyRegistered) Error() string {
	return fmt.Sprintf("Signer %s has already registered.", err.signerName)
}

// ErrSignerNotFound indicated a signer can not found in signers
type ErrSignerNotFound struct {
	signerName string
}

// Error output error message
func (err *ErrSignerNotFound) Error() string {
	return fmt.Sprintf("Signer %s not found.", err.signerName)
}

// ErrNilSignerOptions indicated a nil signer options
type ErrNilSignerOptions struct {
}

// Error output error message
func (err *ErrNilSignerOptions) Error() string {
	return fmt.Sprint("Invalid options. It must not be nil.")
}

// ErrInvalidSignerOptions indicated invalid signer options
type ErrInvalidSignerOptions struct {
	opts crypto.SignerOpts
}

// Error output error message
func (err *ErrInvalidSignerOptions) Error() string {
	return fmt.Sprintf("Invalid options: %s.", err.opts)
}

// ErrUnmarshalSignature indicated unmarshal signature occur error
type ErrUnmarshalSignature struct {
	e error
}

// Error output error message
func (err *ErrUnmarshalSignature) Error() string {
	return fmt.Sprintf("Failed unmashalling signature: %s", err.e)
}
