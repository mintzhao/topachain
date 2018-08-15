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
	"crypto"
	"sync"

	"github.com/mintzhao/topachain/common/config"
	"github.com/mintzhao/topachain/common/crypto/hasher"
	"github.com/mintzhao/topachain/common/crypto/signer"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("crypto")

// CryptoInterface contains all the functions related to crypto, include hash, encrypt, sign etc.
type CryptoInterface interface {
	hasher.Hasher
	signer.Signer
}

// cryptoImpl implement CryptoInterface
type cryptoImpl struct {
	her  hasher.Hasher
	sger signer.Signer
}

// NewCrypto return a CryptoInterface instance based on configuration
func NewCrypto() (CryptoInterface, error) {
	hashName := config.GetHasherName()
	her, err := hasher.GetHasher(hashName)
	if err != nil {
		return nil, err
	}

	signerName := config.GetSignerName()
	sger, err := signer.GetSigner(signerName)
	if err != nil {
		return nil, err
	}

	return &cryptoImpl{
		her:  her,
		sger: sger,
	}, nil
}

// Hash is a global function to hashes message msg.
func Hash(msg []byte, hashName ...string) ([]byte, error) {
	if len(hashName) != 0 {
		her, err := hasher.GetHasher(hashName[0])
		if err != nil {
			return nil, err
		}

		return her.Hash(msg)
	}

	return getInstance().Hash(msg)
}

// Hash hashes messages msg
func (impl *cryptoImpl) Hash(msg []byte) ([]byte, error) {
	return impl.her.Hash(msg)
}

// Sign is a global function to signs digest using PrivateKey k.
func Sign(k crypto.PrivateKey, digest []byte, opts crypto.SignerOpts, signerName ...string) ([]byte, error) {
	if len(signerName) != 0 {
		sger, err := signer.GetSigner(signerName[0])
		if err != nil {
			return nil, err
		}

		return sger.Sign(k, digest, opts)
	}

	return getInstance().Sign(k, digest, opts)
}

// Sign signs digest using PrivateKey k.
func (impl *cryptoImpl) Sign(k crypto.PrivateKey, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return impl.sger.Sign(k, digest, opts)
}

var (
	instance CryptoInterface
	once     sync.Once
)

func getInstance() CryptoInterface {
	once.Do(func() {
		var err error
		instance, err = NewCrypto()
		if err != nil {
			logger.Panic(err)
		}
	})

	return instance
}
