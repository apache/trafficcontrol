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
		local commits=$(getRevCount)
		local sha=$(git rev-parse --short=8 HEAD)
		echo "$commits.$sha"
	else
		# TODO: is this a good method for generating a build number in absence of git?
		tar cf - . | sha1sum || { echo "Could not produce sha1sum of tar'd directory"; exit 1; }
	fi
}

function getCommit() {
	git rev-parse HEAD
}

# ---------------------------------------
function checkEnvironment {
	export TC_VERSION=$(getVersion "$TC_DIR")
	export BUILD_NUMBER=${BUILD_NUMBER:-$(getBuildNumber)}
	export WORKSPACE=${WORKSPACE:-$TC_DIR}
	export RPMBUILD="$WORKSPACE/rpmbuild"
	export DIST="$WORKSPACE/dist"

	mkdir -p "$DIST" || { echo "Could not create $DIST: $?"; exit 1; }

	# verify required tools available in path
	for pgm in go ; do
		type $pgm 2>/dev/null || { echo "$pgm not found in PATH"; exit 1; }
	done
	echo "Build environment has been verified."

	echo "=================================================="
	echo "WORKSPACE: $WORKSPACE"
	echo "BUILD_NUMBER: $BUILD_NUMBER"
	echo "TC_VERSION: $TC_VERSION"
	echo "--------------------------------------------------"
}

# ---------------------------------------
function createSourceDir() {
	local target="$1-$TC_VERSION"
	local srcpath="$RPMBUILD/SOURCES/$target"
	mkdir -p "$srcpath" || { echo "Could not create $srcpath: $?"; exit 1; }
	echo "$srcpath"
}

# ---------------------------------------
function buildRpm () {
	for package in "$@"; do
		local rpm="${package}-${TC_VERSION}-${BUILD_NUMBER}.$(uname -m).rpm"
		local srpm="${package}-${TC_VERSION}-${BUILD_NUMBER}.src.rpm"
		echo "Building the rpm."

		cd "$RPMBUILD" && \
			rpmbuild --define "_topdir $(pwd)" \
				 --define "traffic_control_version $TC_VERSION" \
				 --define "commit $(getCommit)" \
				 --define "build_number $BUILD_NUMBER" \
				 -ba SPECS/$package.spec || \
				 { echo "RPM BUILD FAILED: $?"; exit 1; }

		echo
		echo "========================================================================================"
		echo "RPM BUILD FOR $package SUCCEEDED, See $DIST/$rpm for the newly built rpm."
		echo "========================================================================================"
		echo

		cp "$RPMBUILD/RPMS/$(uname -m)/$rpm" "$DIST/." || { echo "Could not copy $rpm to $DIST: $?"; exit 1; }
		cp "$RPMBUILD/SRPMS/$srpm" "$DIST/." || { echo "Could not copy $srpm to $DIST: $?"; exit 1; }
	done
}

