#!/usr/bin/env sh
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
# shellcheck shell=ash
trap 'exit_code=$?; [ $exit_code -ne 0 ] && echo "Error on line ${LINENO} of ${0}" >/dev/stderr; exit $exit_code' EXIT;
set -o errexit -o nounset -o pipefail;

#----------------------------------------
importFunctions() {
	local script scriptdir
	TM_DIR='' TC_DIR=''
	script=$(realpath "$0")
	scriptdir=$(dirname "$script")
	TM_DIR="$(dirname "$scriptdir")"
	TC_DIR="$(dirname "$TM_DIR")"
	export TM_DIR TC_DIR
	functions_sh="$TC_DIR/build/functions.sh"
	if [ ! -r "$functions_sh" ]; then
		echo "error: can't find $functions_sh"
		return 1
	fi
	. "$functions_sh"
}

#----------------------------------------
initBuildArea() {
	echo "Initializing the build area."
	(mkdir -p "$RPMBUILD"
	 cd "$RPMBUILD"
	 mkdir -p SPECS SOURCES RPMS SRPMS BUILD BUILDROOT) || { echo "Could not create $RPMBUILD: $?"; return 1; }

	# tar/gzip the source
	local tm_dest
	tm_dest="$(createSourceDir traffic_monitor)"
	cd "$TM_DIR" || \
		 { echo "Could not cd to $TM_DIR: $?"; return 1; }

	echo "PATH: $PATH"
	echo "GOPATH: $GOPATH"
	go version
	go env

	# get x/* packages (everything else should be properly vendored)
	go mod vendor -v || \
		{ echo "Could not vendor go module dependencies"; return 1; }

	# compile traffic_monitor
	gcflags=''
	ldflags="-X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}"
	export CGO_ENABLED=0
	{ set +o nounset;
	if [ "$DEBUG_BUILD" = true ]; then
		echo 'DEBUG_BUILD is enabled, building without optimization or inlining...';
		gcflags="${gcflags} all=-N -l";
	else
		ldflags="${ldflags} -s -w"; #strip binary
	fi;
	set -o nounset; }
	go build -v -gcflags "$gcflags" -ldflags "$ldflags" || \
		{ echo "Could not build traffic_monitor binary"; return 1; }

	cp -av ./ "$tm_dest"/ || \
		 { echo "Could not copy to $tm_dest: $?"; return 1; }
	cp -av "$TM_DIR"/build/*.spec "$RPMBUILD"/SPECS/. || \
		 { echo "Could not copy spec files: $?"; return 1; }

	# include LICENSE in the source RPM
	cp "${TC_DIR}/LICENSE" "$tm_dest"

	tar -czvf "$tm_dest".tgz -C "$RPMBUILD"/SOURCES "$(basename "$tm_dest")" || { echo "Could not create tar archive $tm_dest.tgz: $?"; return 1; }
	cp "$TM_DIR"/build/*.spec "$RPMBUILD"/SPECS/. || { echo "Could not copy spec files: $?"; return 1; }

	echo "The build area has been initialized."
}

preBuildChecks() {
	if [ -e "$TM_DIR"/traffic_monitor ]; then
		echo "Found $TM_DIR/traffic_monitor, please remove before retrying to build"
		return 1
	fi
}

# ---------------------------------------

importFunctions
preBuildChecks
checkEnvironment -i go,rsync
initBuildArea
buildRpm traffic_monitor
