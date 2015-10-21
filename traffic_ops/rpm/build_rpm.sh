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

#----------------------------------------
function buildRpm () {
	echo "Building the rpm."

	cd "$RPMBUILD" && \
		rpmbuild --define "_topdir $(pwd)" \
			 --define "traffic_ops_version $TC_VERSION" \
			 --define "build_number $BUILD_NUMBER" -ba SPECS/traffic_ops.spec

	if [[ $? -ne  0 ]]; then
		echo -e "\nRPM BUILD FAILED.\n\n"
		exit 1
	fi
	echo
	echo "========================================================================================"
	echo "RPM BUILD SUCCEEDED, See $DIST/$RPM for the newly built rpm."
	echo "========================================================================================"
	echo

	mkdir -p "$DIST" || { echo "Could not create $DIST: $!"; exit 1; }

	/bin/cp "$RPMBUILD"/RPMS/*/*.rpm "$DIST/." || { echo "Could not copy rpm to $DIST: $!"; exit 1; }
}


#----------------------------------------
function checkEnvironment() {
	echo "Verifying the build configuration environment."
	local script=$(readlink -f "$0")
	local scriptdir=$(dirname "$script")
	export TO_DIR=$(dirname "$scriptdir")
	export TC_DIR=$(dirname "$TO_DIR")
	functions_sh="$TC_DIR/rpm/functions.sh"
	if [[ ! -r $functions_sh ]]; then
		echo "Error: Can't find $functions_sh"
		exit 1
	fi
	. "$functions_sh"

	# 
	# get traffic_control src path -- relative to build_rpm.sh script
	export PACKAGE="traffic_ops"
	export TC_VERSION=$(getVersion "$TC_DIR")
	export BUILD_NUMBER=${BUILD_NUMBER:-$(getBuildNumber)}
	export WORKSPACE=${WORKSPACE:-$TC_DIR}
	export RPMBUILD="$WORKSPACE/rpmbuild"
	export DIST="$WORKSPACE/dist"
	export RPM="${PACKAGE}-${TC_VERSION}-${BUILD_NUMBER}.x86_64.rpm"
	export IN_GIT=$(isInGitTree)

	# verify required tools available in path
	for pgm in go ; do
		type $pgm 2>/dev/null || { echo "$pgm not found in PATH"; exit 1; }
	done
	echo "Build environment has been verified."

	echo "=================================================="
	echo "WORKSPACE: $WORKSPACE"
	echo "BUILD_NUMBER: $BUILD_NUMBER"
	echo "TC_VERSION: $TC_VERSION"
	echo "RPM: $RPM"
	echo "--------------------------------------------------"
}

# ---------------------------------------
function initBuildArea() {
	echo "Initializing the build area."
	/bin/rm -rf "$RPMBUILD" && \
		mkdir -p "$RPMBUILD"/{SPECS,SOURCES,RPMS,SRPMS,BUILD,BUILDROOT} || { echo "Could not create $RPMBUILD: $!"; exit 1; }

	/bin/cp -r "$TO_DIR"/rpm/*.spec "$RPMBUILD"/SPECS/. || { echo "Could not copy spec files: $!"; exit 1; }

	# build the go scripts for database initialization and tm testing.

	# tar/gzip the source
	local target="traffic_ops-$TC_VERSION"
	local targetpath="$RPMBUILD/SOURCES/$target"
	/bin/mkdir -p "$targetpath/app" || { echo "Could not create $targetpath"; exit 1; }
	cd "$TO_DIR" || { echo "Could not cd to $TO_DIR: $!"; exit 1; }
	for d in app/{bin,conf,cpanfile,db,lib,public,script,templates} doc etc install; do
		/bin/cp -r "$d" "$targetpath/$d" || { echo "Could not copy $d files to $targetpath: $!"; exit 1; }
	done

	# compile go executables used during postinstall
	local bldinstall="$RPMBUILD/BUILD/install"
	mkdir -p "$bldinstall"
	cp -r "$TO_DIR/install/go" "$bldinstall" || { echo "Could not copy to $bldinstall: $!"; exit 1; }
	cd "$bldinstall/go" || { echo "Could not cd to $bldinstall/go: $!"; exit 1; }

	export GOPATH=$(pwd)
	export GOBIN="$targetpath/install/bin"

	echo "Compiling go executables"
	for d in src/comcast.com/*; do
		if [[ ! -d "$d" ]]; then
			echo "Could not find $d"
			exit 1
		fi
		(cd "$d" && go get || { echo "Could not compile $d"; exit 1; } )
	done

	tar -czvf "$targetpath.tgz" -C "$RPMBUILD/SOURCES" "$target" || { echo "Could not create tar archive $targetpath.tgz: $!"; exit 1; }

	echo "The build area has been initialized."
}

# ---------------------------------------

checkEnvironment
initBuildArea
buildRpm
