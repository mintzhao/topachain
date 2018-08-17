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
package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterHasher(t *testing.T) {
	assert.EqualError(t, RegisterHasher("SHA256", &sha256Hasher{}), (&ErrHasherAlreadyRegistered{hasherName: "SHA256"}).Error())
	assert.NoError(t, RegisterHasher("test", &sha256Hasher{}))

	_, err := GetHasher("test")
	assert.NoError(t, err)

	DeRegisterHasher("test")
}

func TestGetHasher(t *testing.T) {
	_, err := GetHasher("test")
	assert.EqualError(t, err, (&ErrHasherNotFound{"test"}).Error())

	_, err = GetHasher("sha256")
	assert.NoError(t, err)
}
