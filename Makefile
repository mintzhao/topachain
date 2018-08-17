# Copyright © 2018 Zhao Ming <mint.zhao.chiu@gmail.com>.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# -------------------------------------------------------------
# This makefile defines the following targets
#
#   - all (default) - builds all targets and runs all tests/checks
#   - checks - runs all tests/checks
#   - release - builds release packages for the host platform
#   - release-all - builds release packages for all target platforms
#   - unit-test - runs the go-test based unit tests
#   - gotools - installs go tools like golint
#   - linter - runs all code checks
#   - license - checks go sourrce files for Apache license header
#   - native - ensures all native binaries are available
#   - docker[-clean] - ensures all docker images are available[/cleaned]
#   - clean - cleans the build area
#   - clean-all - superset of 'clean' that also removes persistent state
#   - dist-clean - clean release packages for all target platforms
#   - unit-test-clean - cleans unit test state (particularly from docker)

PROJECT_NAME   = mintzhao/topachain
BASE_VERSION = 0.0.1
PREV_VERSION = 0.0.0
IS_RELEASE = false

ifneq ($(IS_RELEASE),true)
EXTRA_VERSION ?= snapshot-$(shell git rev-parse --short HEAD)
PROJECT_VERSION=$(BASE_VERSION)-$(EXTRA_VERSION)
else
PROJECT_VERSION=$(BASE_VERSION)
endif

DOCKER_NS ?= mintzhao
DOCKER_TAG=$(ARCH)-$(PROJECT_VERSION)

PKGNAME = github.com/$(PROJECT_NAME)
CGO_FLAGS = CGO_CFLAGS=" "
ARCH=$(shell uname -m)
MARCH=$(shell go env GOOS)-$(shell go env GOARCH)

# defined in cmd/version.go
METADATA_VAR = Version=$(PROJECT_VERSION)

GOBIN=$(abspath $(GOPATH)/bin)
GO_LDFLAGS = $(patsubst %,-X $(PKGNAME)/cmd/version.%,$(METADATA_VAR))

GO_TAGS ?=

export GO_LDFLAGS

# No sense rebuilding when non production code is changed
PROJECT_FILES = $(shell git ls-files  | grep -v ^test | grep -v ^unit-test | \
	grep -v ^bddtests | grep -v ^docs | grep -v _test.go$ | grep -v .md$ | \
	grep -v ^.git | grep -v ^examples | grep -v ^devenv | grep -v .png$ | \
	grep -v ^LICENSE )

IMAGES = topa
RELEASE_PLATFORMS = windows-amd64 darwin-amd64 linux-amd64 linux-ppc64le linux-s390x
RELEASE_PKGS = topa

pkgmap.topa      := $(PKGNAME)

# dep tools
DEPTOOLS = golint goimports misspell
DEPTOOLS_BIN = $(patsubst %,$(GOBIN)/%, $(DEPTOOLS))

dep.tools.golint := github.com/golang/lint/golint
dep.tools.goimports := golang.org/x/tools/cmd/goimports
dep.tools.misspell := github.com/client9/misspell/cmd/misspell

all: native checks

checks: spelling linter test

.PHONY: topa
topa: build/bin/topa

native: topa

test:
	@go test $(@go list ./... | @grep -v /vendor/)

build/bin/topa: $(PROJECT_FILES)
	@mkdir -p $(@D)
	@echo "$@"
	$(CGO_FLAGS) go build -o $(@F) -tags "$(GO_TAGS)" -ldflags "$(GO_LDFLAGS)" $(pkgmap.$(@F))
	@echo "Binary available as $@"
	@touch $@

dep.tools.%:
	$(eval TOOL = ${subst dep.tools.,,${@}})
	@echo "Building ${dep.tools.${TOOL}} -> $(TOOL)"
	go get -v ${dep.tools.${TOOL}}

.PHONY: spelling
spelling: dep.tools.misspell
	@scripts/check_spelling.sh

.PHONY: linter
linter: dep.tools.golint
	@echo "LINT: Running code checks.."
	@scripts/golinter.sh

.PHONY: changelog
changelog:
	@scripts/changelog.sh v$(PREV_VERSION) v$(BASE_VERSION)