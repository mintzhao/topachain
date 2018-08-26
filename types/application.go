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
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// AppMetadata
type AppMetadata struct {
	Name    string
	Version *AppVersion
}

// AppConfig
type AppConfig struct {
	MasterAddress string // blockchain master node address
}

// Software versioning: https://en.wikipedia.org/wiki/Software_versioning

var (
	ErrEmptyVersion = errors.New("empty version string")
)

// version
type version struct {
	major, minor, build string
}

func newVersion(ver string) (*version, error) {
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

	return &version{
		major: major,
		minor: minor,
		build: build,
	}, nil
}

// String format output
func (v *version) String() string {
	return fmt.Sprintf("%s.%s.%s", v.major, v.minor, v.build)
}

// Compare returns an integer comparing two versions lexicographically.
// The result will be 0 if v==other, -1 if v < other, and +1 if v > other.
func (v *version) compare(other *version) int {
	if v == nil {
		return -1
	}

	if other == nil {
		return 1
	}

	if v.major > other.major {
		return 1
	}

	if v.major < other.major {
		return -1
	}

	if v.minor > other.minor {
		return 1
	}

	if v.minor < other.minor {
		return -1
	}

	if v.build > other.build {
		return 1
	}

	if v.build < other.build {
		return -1
	}

	return 0
}

// Equal check whether two version at same stage.
func (v *version) equal(other *version) bool {
	return v.compare(other) == 0
}

// AppVersion
// Features:
// 1. upgrade version value must larger than current version
// 2. newest application can compatible a specific version of application,
//    older version than Backwards shouldn't be used in the blockchain network.
type AppVersion struct {
	Version   *version
	Backwards *version // Backwards Compatibility
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

// String format av string output
func (av *AppVersion) String() string {
	return fmt.Sprintf("^%s/%s", av.Backwards, av.Version)
}
