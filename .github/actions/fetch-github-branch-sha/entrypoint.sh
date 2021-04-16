#!/bin/bash
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

# verify required environment inputs.
if [[ -z ${INPUT_OWNER} || -z ${INPUT_REPO} || -z ${INPUT_BRANCH} ]]; then
  echo "Error: missing required environment variables"
  exit 1
fi

# fetch the branch info
_brinfo=`curl --silent https://api.github.com/repos/${INPUT_OWNER}/${INPUT_REPO}/branches/${INPUT_BRANCH}`
_rc=$?

if [[ ${_rc} -ne 0 ]]; then
  echo "Error: failed to fetch branch info ${INPUT_BRANCH}"
  exit 2
fi

# parse out the commit sha
_sha=`echo -E ${_brinfo} | awk '{
  if ($0 ~ /Branch not found/) {
    print "BADBRANCH"
  } 
  else if ($0 ~ /name/) {
    print $7
  }
}' | sed -e 's/[",]//g'`

# verify the sha
if [[ -z ${_sha} || ${_sha} == "BADBRANCH" ]]; then
  echo "Error: could not parse the commit from branch ${INPUT_BRANCH}"
  exit 3
fi

echo "::set-output name=sha::${_sha}"

exit 0

