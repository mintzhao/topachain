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
	"strings"

	"github.com/pkg/errors"
)

// Software versioning: https://en.wikipedia.org/wiki/Software_versioning
var (
	ErrEmptyVersion = errors.New("empty version string")
)

func newVersion(ver string) (*Version, error) {
	ver = strings.TrimSpace(ver)
	if ver == "" {
		return nil, ErrEmptyVersion
	}

	var major, minor, build string
	vslice := strings.SplitN(ver, ".", 3)
	if len(vslice) >= 1 {
		major = vslice[0]
	}
	if len(vslice) >= 2 {
		minor = vslice[1]
	}
	if len(vslice) >= 3 {
		build = vslice[2]
	}

	return &Version{
		Major: major,
		Minor: minor,
		Build: build,
	}, nil
}

// Compare returns an integer comparing two versions lexicographically.
// The result will be 0 if v==other, -1 if v < other, and +1 if v > other.
func (v *Version) compare(other *Version) int {
	if v == nil {
		return -1
	}

	if other == nil {
		return 1
	}

	if v.GetMajor() > other.GetMajor() {
		return 1
	}

	if v.GetMajor() < other.GetMajor() {
		return -1
	}

	if v.GetMinor() > other.GetMinor() {
		return 1
	}

	if v.GetMinor() < other.GetMinor() {
		return -1
	}

	if v.GetBuild() > other.GetBuild() {
		return 1
	}

	if v.GetBuild() < other.GetBuild() {
		return -1
	}

	return 0
}

// Equal check whether two version at same stage.
func (v *Version) equal(other *Version) bool {
	return v.compare(other) == 0
}

// NewAppVersion construct app version by version format output string
func NewAppVersion(version, backwards string) (*AppVersion, error) {
	ver, err := newVersion(version)
	if err != nil {
		return nil, err
	}

	back, err := newVersion(backwards)
	if err != nil {
		return nil, err
	}

	return &AppVersion{
		Version:   ver,
		Backwards: back,
	}, nil
}

// Compatible check whether av & other can compatible with each other
func (av *AppVersion) Compatible(other *AppVersion) bool {
	if av.Version.equal(other.Version) {
		return true
	}

	if av.Version.compare(other.Version) == 1 {
		return other.Version.compare(av.Backwards) >= 0
	}

	return av.Version.compare(other.Backwards) >= 0
}
