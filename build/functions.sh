#!/bin/bash

#
# Copyright 2015 Comcast Cable Communications Management, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# ---------------------------------------
function getVersion() {
	local d="$1"
	local vf="$d/VERSION"
	cat "$vf" || { echo "Could not read $vf: $!"; exit 1; }
}

function getRevCount() {
	git rev-list HEAD 2>/dev/null | wc -l
}

# ---------------------------------------
function isInGitTree() {
	git rev-parse --is-inside-work-tree 2>/dev/null
}

# ---------------------------------------
function getBuildNumber() {
	local in_git=$(isInGitTree)
	if [[ $in_git ]]; then
		local commits=$(git rev-list HEAD | wc -l)
		local sha=$(git rev-parse --short=8 HEAD)
		echo "$commits.$sha"
	else
		# TODO: is this a good method for generating a build number in absence of git?
		tar cf - . | sha1sum || { echo "Could not produce sha1sum of tar'd directory"; exit 1; }
	fi
}
