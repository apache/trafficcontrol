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
trap 'echo "Error on line ${LINENO} of ${0}"; exit 1' ERR;
set -o errexit -o nounset

# verify required environment inputs.
if [[ -z ${INPUT_OWNER} || -z ${INPUT_REPO} || -z ${INPUT_BRANCH} ]]; then
  echo "Error: missing required environment variables"
  exit 1
fi

# fetch the branch info
if ! _brinfo="$(curl --silent "${GITHUB_API_URL}/repos/${INPUT_OWNER}/${INPUT_REPO}/branches/${INPUT_BRANCH}")"; then
  echo "Error: failed to fetch branch info ${INPUT_BRANCH}"
  exit 2
fi

# parse out the commit sha
_sha="$(<<<"$_brinfo" jq -r .commit.sha)"

# verify the sha
if [[ -z "${_sha}" || "${_sha}" == "null" ]]; then
  echo "Error: could not parse the commit from branch ${INPUT_BRANCH}"
  exit 3
fi

echo "::set-output name=sha::${_sha}"

branch_prefix_pattern="^$(<<<"$INPUT_BRANCH" grep -o '.*\.' | sed 's/\./\\./g')[0-9.]+$"

if ! tags_info="$(curl --silent "${GITHUB_API_URL}/repos/${INPUT_OWNER}/${INPUT_REPO}/tags")"; then
	echo "Error: failed to fetch tag info ${INPUT_BRANCH}"
	exit 2
fi

latest_tag="$(<<<"$tags_info" jq -r --arg BRANCH_PREFIX_PATTERN "$branch_prefix_pattern" '.[] | .name | select(test($BRANCH_PREFIX_PATTERN))' |
	head -n1)"

echo "::set-output name=latest-tag::${latest_tag}"

exit 0
