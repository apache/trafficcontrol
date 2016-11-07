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

#----------------------------------------
function importFunctions() {
	echo "Verifying the build configuration environment."
	local script=$(readlink -f "$0")
	local scriptdir=$(dirname "$script")
	export TM_DIR=$(dirname "$scriptdir")
	export TC_DIR=$(dirname "$TM_DIR")
	functions_sh="$TC_DIR/build/functions.sh"
	if [[ ! -r $functions_sh ]]; then
		echo "Error: Can't find $functions_sh"
		exit 1
	fi
	. "$functions_sh"
}

#----------------------------------------
function buildRpmTrafficMonitor () {
	echo "Building the rpm."

	cd "$TM_DIR" || { echo "Could not cd to $TM_DIR: $?"; exit 1; }
	export TRAFFIC_CONTROL_VERSION="$TC_VERSION"
	export GIT_REV_COUNT=$(getRevCount)
	mvn clean package || { echo "RPM BUILD FAILED: $?"; exit 1; }

	local rpm=$(find -name \*.rpm)
	if [[ -z $rpm ]]; then
		echo "Could not find rpm file $RPM in $(pwd)"
		exit 1;
	fi
	echo
	echo "========================================================================================"
	echo "RPM BUILD SUCCEEDED, See $DIST/$RPM for the newly built rpm."
	echo "========================================================================================"
	echo
	mkdir -p "$DIST" || { echo "Could not create $DIST: $?"; exit 1; }

	cp "$rpm" "$DIST/." || { echo "Could not copy $RPM to $DIST: $?"; exit 1; }
}

# ---------------------------------------
function initBuildArea() {
	echo "Initializing the build area."
	mkdir -p "$RPMBUILD"/{SPECS,SOURCES,RPMS,SRPMS,BUILD,BUILDROOT} || { echo "Could not create $RPMBUILD: $?"; exit 1; }

	tm_dest=$(createSourceDir traffic_monitor)

	export TRAFFIC_CONTROL_VERSION="$TC_VERSION"
    export MVN_CMD="mvn versions:set -DnewVersion=$TRAFFIC_CONTROL_VERSION"
    echo $MVN_CMD
    $MVN_CMD
	cp -r "$TM_DIR"/{build,etc,src} "$tm_dest"/. || { echo "Could not copy to $tm_dest: $?"; exit 1; }
	cp  "$TM_DIR"/pom.xml "$tm_dest" || { echo "Could not copy to $tm_dest: $?"; exit 1; }

	tar -czf "$tm_dest.tgz" -C "$RPMBUILD"/SOURCES $(basename "$tm_dest") || { echo "Could not create tar archive $tm_dest.tgz: $?"; exit 1; }

	echo "The build area has been initialized."
}

# ---------------------------------------

importFunctions
checkEnvironment
initBuildArea
buildRpmTrafficMonitor
