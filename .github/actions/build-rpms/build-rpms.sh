#!/usr/bin/env bash
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

set -o errexit -o nounset -o pipefail -o xtrace
trap 'echo "Error on line ${LINENO} of ${0}"; exit 1' ERR

export DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1

pkg_command=(./pkg -v)
if [[ "$ATC_COMPONENT" == "traffic_portal_v2" ]]; then
	pkg_command=(./pkg -v -o)
fi

# If the Action is being run on a Pull Request
if [[ "$GITHUB_REF" == refs/pull/*/merge ]]; then
	sudo apt-get install jq
	pr_number="$(<<<"$GITHUB_REF" grep -o '[0-9]\+')"
	for ((i = 0; i < 10; i++)); do
		response="$(curl "${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/pulls/${pr_number}/files")"
		if files_changed="$(<<<"$response" jq -r .[].filename)"; then
			break
		else
			echo "unexpected response body: ${response}"
		fi
		sleep 2
	done
else
	files_changed="$(git diff --name-only HEAD~4 HEAD --)" # Files changed by the last 4 commits
fi
if <<<"$files_changed" grep '^GO_VERSION$' -e '' || <<<"$files_changed" grep '^infrastructure/docker/build/'; then
	pkg_command+=(-b)
fi

if [[ -z "${ATC_COMPONENT:-}" ]]; then
	echo 'Missing environment variable ATC_COMPONENT' >/dev/stderr
	exit 1
fi
ATC_COMPONENT+='_build'

"${pkg_command[@]}" "$ATC_COMPONENT"
