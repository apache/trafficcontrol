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

#----------------------------------------
function importFunctions() {
	echo "Verifying the build configuration environment."
	local script=$(readlink -f "$0")
	local scriptdir=$(dirname "$script")
	export TO_DIR=$(dirname "$scriptdir")
	export TC_DIR=$(dirname "$TO_DIR")
	functions_sh="$TC_DIR/build/functions.sh"
	if [[ ! -r $functions_sh ]]; then
		echo "Error: Can't find $functions_sh"
		exit 1
	fi
	. "$functions_sh"
}


function createSourceDir() {
	local target="$1-$TC_VERSION"
	local srcpath="$RPMBUILD/SOURCES/$target"
	mkdir -p "$srcpath" || { echo "Could not create $srcpath: $?"; exit 1; }
	echo "$srcpath"
}

# ---------------------------------------
function initBuildArea() {
	echo "Initializing the build area."
	mkdir -p "$RPMBUILD"/{SPECS,SOURCES,RPMS,SRPMS,BUILD,BUILDROOT} || { echo "Could not create $RPMBUILD: $?"; exit 1; }

	to_src=$(createSourceDir traffic_ops)
	cd "$TO_DIR" || { echo "Could not cd to $TO_DIR: $?"; exit 1; }
	for d in app/{bin,conf,cpanfile,db,lib,public,script,templates} doc etc install; do
		if [[ -d "$d" ]]; then
			mkdir -p "$to_src/$d" || { echo "Could not create $to_src/$d: $?"; exit 1; }
			cp -r "$d"/* "$to_src/$d" || { echo "Could not copy $d files to $to_src: $?"; exit 1; }
		else
			cp "$d" "$to_src/$d" || { echo "Could not copy $d to $to_src: $?"; exit 1; }
		fi

	done

	tar -czvf "$to_src.tgz" -C "$RPMBUILD"/SOURCES $(basename "$to_src") || { echo "Could not create tar archive $to_src.tgz: $?"; exit 1; }
	cp "$TO_DIR"/build/*.spec "$RPMBUILD"/SPECS/. || { echo "Could not copy spec files: $?"; exit 1; }

	echo "The build area has been initialized."
}

# ---------------------------------------
importFunctions
checkEnvironment
initBuildArea
buildRpm traffic_ops traffic_ops_ort
