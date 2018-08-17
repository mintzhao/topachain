#!/bin/bash
#
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

set -e

declare -a arr=(
)

for i in "${arr[@]}"
do
    echo "Checking $i"
    go vet $i/...
    OUTPUT="$(goimports -srcdir $GOPATH/src/github.com/mintzhao/topachain -l $i)"
    if [[ $OUTPUT ]]; then
	echo "The following files contain goimports errors"
	echo $OUTPUT
	echo "The goimports command must be run for these files"
	exit 1
    fi
done
