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
package config

import "github.com/spf13/viper"

// GetHasherName returns customize hash function, default is SHA256
func GetHasherName() string {
	hasherName := viper.GetString("crypto.hash")

	if hasherName == "" {
		hasherName = "SHA256"
	}

	return hasherName
}

// GetSignerName returns customize signer, default is ECDSA
func GetSignerName() string {
	signerName := viper.GetString("crypto.sign")

	if signerName == "" {
		signerName = "ECDSA"
	}

	return signerName
}
