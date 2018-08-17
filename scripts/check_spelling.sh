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

CHECK=$(git diff --name-only HEAD * | grep -v .png$ | grep -v .git | grep -v ^CHANGELOG \
  | grep -v ^vendor/ | grep -v ^build/ | sort -u)

if [[ -z "$CHECK" ]]; then
  CHECK=$(git diff-tree --no-commit-id --name-only -r $(git log -2 \
    --pretty=format:"%h") | grep -v .png$ | grep -v .git | grep -v ^CHANGELOG \
    | grep -v ^vendor/ | grep -v ^build/ | sort -u)
fi

echo "Checking changed go files for spelling errors ..."
errs=`echo $CHECK | xargs misspell -source=text`
if [ -z "$errs" ]; then
   echo "spell checker passed"
   exit 0
fi
echo "The following files are have spelling errors:"
echo "$errs"
exit 0
