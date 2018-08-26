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
package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	ver, err := newVersion("")
	assert.EqualError(t, err, ErrEmptyVersion.Error())
	assert.Nil(t, ver)

	ver1, err := newVersion("1.0.0")
	assert.NoError(t, err)

	ver2, err := newVersion("1.0.1")
	assert.NoError(t, err)

	ver3, err := newVersion("1.0.1")
	assert.NoError(t, err)

	assert.Equal(t, -1, ver1.compare(ver2))
	assert.Equal(t, true, ver2.equal(ver3))
}

func TestAppVersion(t *testing.T) {
	version, err := NewAppVersion("", "")
	assert.EqualError(t, err, ErrEmptyVersion.Error())
	assert.Nil(t, version)

	version1, err := NewAppVersion("1.0.0", "0.9.1")
	assert.NoError(t, err)

	version2, err := NewAppVersion("1.2.0", "1.0.1")
	assert.NoError(t, err)

	version3, err := NewAppVersion("1.1.0", "0.9.5")
	assert.NoError(t, err)

	assert.Equal(t, false, version1.Compatible(version2))
	assert.Equal(t, true, version2.Compatible(version3))
}
