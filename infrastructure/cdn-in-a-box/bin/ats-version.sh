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

trap 'echo "Error on line ${LINENO} of ${0}" >/dev/stderr; exit 1' ERR;
set -o errexit -o nounset -o pipefail

project=trafficserver
script_dir="$(dirname "$0")"
ats_version_file="${script_dir}/../cache/ATS_VERSION"

remote_ats_version() {
	local gitbox_url=https://gitbox.apache.org/repos/asf
	local repo="${project}.git"
	local branch refs commit last_tag release
	branch="$(grep 'ATS_VERSION=' "${script_dir}/../../../cache-config/testing/docker/variables.env" | cut -d= -f2)"
	refs="$(curl -fs "${gitbox_url}/${repo}/info/refs?service=git-upload-pack" |
		sed -E 's/^00[0-9a-f]{2}//g' |
		tr -d '\0')"
	commit="$(<<<"$refs" grep "refs/heads/${branch}$" | awk '{print $1}')"

	# $last_tag is the latest tag before the commit at the head of $branch.
	last_tag="$(<<<"$refs" grep -oE 'refs/tags/[0-9.]+$' |
		cut -d/ -f3 |
		grep -F "${branch::$((${#branch} - 1))}" |
		tail -n1)"

	# $release is the number of commits between $release to $branch.
	page_output="$(curl -fs "${gitbox_url}?p=${repo};a=shortlog;h=${branch};hp=${last_tag}")"
	release="$(<<<"$page_output" grep -c 'class="link"' || true)"
	<<<"${last_tag}-${release}.${commit:0:9}" tee "$ats_version_file"
}

# Reuse the ATS version file if it was generated within the last day
if [[ -n "$(find "$ats_version_file" -mtime -1 2>/dev/null)" ]]; then
	cat "$ats_version_file"
else
	remote_ats_version
fi
