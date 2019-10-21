#!/bin/bash

#
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

function importFunctions() {
	local script=$(readlink -f "$0")
	local scriptdir=$(dirname "$script")
	export TM_DIR=$(dirname "$scriptdir")
	export TC_DIR=$(dirname "$TM_DIR")
	functions_sh="$TC_DIR/build/functions.sh"
	if [[ ! -r $functions_sh ]]; then
		echo "error: can't find $functions_sh"
		exit 1
	fi
	. "$functions_sh"
}

#----------------------------------------
function initBuildArea() {
	echo "Initializing the build area."
	mkdir -p "$RPMBUILD"/{SPECS,SOURCES,RPMS,SRPMS,BUILD,BUILDROOT} || { echo "Could not create $RPMBUILD: $?"; exit 1; }

	# tar/gzip the source
	local tm_dest=$(createSourceDir traffic_monitor)
	cd "$TM_DIR" || \
		 { echo "Could not cd to $TM_DIR: $?"; exit 1; }

	echo "PATH: $PATH"
	echo "GOPATH: $GOPATH"
	go version
	go env

	# get x/* packages (everything else should be properly vendored)
	go get -v golang.org/x/net/publicsuffix golang.org/x/sys/unix || \
		{ echo "Could not get go package dependencies"; exit 1; }

	# compile traffic_monitor
	go build -v -ldflags "-X main.GitRevision=`git rev-parse HEAD` -X main.BuildTimestamp=`date +'%Y-%M-%dT%H:%M:%s'` -X main.Version=${TC_VERSION}" || \
		{ echo "Could not build traffic_monitor binary"; exit 1; }

	rsync -av ./ "$tm_dest"/ || \
		 { echo "Could not copy to $tm_dest: $?"; exit 1; }
	cp "$TM_DIR"/build/*.spec "$RPMBUILD"/SPECS/. || \
		 { echo "Could not copy spec files: $?"; exit 1; }

	tar -czvf "$tm_dest".tgz -C "$RPMBUILD"/SOURCES $(basename $tm_dest) || { echo "Could not create tar archive $tm_dest.tgz: $?"; exit 1; }
	cp "$TM_DIR"/build/*.spec "$RPMBUILD"/SPECS/. || { echo "Could not copy spec files: $?"; exit 1; }

	echo "The build area has been initialized."
}

function preBuildChecks() {
	if [[ -e "$TM_DIR"/traffic_monitor ]]; then
		echo "Found $TM_DIR/traffic_monitor, please remove before retrying to build"
		exit 1
	fi
}

# ---------------------------------------

importFunctions
preBuildChecks
checkEnvironment go
initBuildArea
buildRpm traffic_monitor
