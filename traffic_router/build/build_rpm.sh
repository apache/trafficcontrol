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
function importFunctions() {
	echo "Verifying the build configuration environment."
	local script=$(readlink -f "$0")
	local scriptdir=$(dirname "$script")
	export TR_DIR=$(dirname "$scriptdir")
	export TC_DIR=$(dirname "$TR_DIR")
	functions_sh="$TC_DIR/build/functions.sh"
	if [[ ! -r $functions_sh ]]; then
		echo "Error: Can't find $functions_sh"
		return 1
	fi
	. "$functions_sh"
}

#----------------------------------------
function buildRpmTrafficRouter () {
	echo "Building the rpm."

	export STARTUP_SCRIPT_DIR="/lib/systemd/system"
	export STARTUP_SCRIPT_LOC="../core/src/main/lib/systemd/system"

	cd "$TR_DIR" || { echo "Could not cd to $TR_DIR: $?"; return 1; }
	mvn -P rpm-build -Dmaven.test.skip=true -DminimumTPS=1 clean package ||  \
		{ echo "RPM BUILD FAILED: $?"; return 1; }

	local rpm=$(find -name \*.rpm)
	if [[ -z $rpm ]]; then
		echo "Could not find rpm file $RPM in $(pwd)"
		return 1;
	fi
	echo "========================================================================================"
	echo "RPM BUILD SUCCEEDED, See $DIST/$RPM for the newly built rpm."
	echo "========================================================================================"
	echo
	mkdir -p "$DIST" || { echo "Could not create $DIST: $?"; return 1; }

	cp "$rpm" "$DIST/." || { echo "Could not copy $rpm to $DIST: $?"; return 1; }

}


#----------------------------------------
function adaptEnvironment() {
	echo "Verifying the build configuration environment."
	# get traffic_control src path -- relative to build_rpm.sh script
	export PACKAGE="traffic_router"
	export TC_VERSION=$(getVersion "$TC_DIR")
	export BUILD_NUMBER=${BUILD_NUMBER:-$(getBuildNumber)}
	export BUILD_LOCK=$(getBuildNumber).${RHEL_VERSION}
	export WORKSPACE=${WORKSPACE:-$TC_DIR}
	export RPMBUILD="$WORKSPACE/rpmbuild"
	export DIST="$WORKSPACE/dist"
	export RPM="${PACKAGE}-${TC_VERSION}-${BUILD_NUMBER}.x86_64.rpm"
	export TOMCAT_VERSION=8.5
	export TOMCAT_RELEASE=32

	echo "=================================================="
	echo "WORKSPACE: $WORKSPACE"
	echo "BUILD_NUMBER: $BUILD_NUMBER"
	echo "BUILD_LOCK: $BUILD_LOCK"
	echo "TOMCAT_VERSION=$TOMCAT_VERSION"
	echo "TOMCAT_RELEASE=$TOMCAT_RELEASE"
	echo "TC_VERSION: $TC_VERSION"
	echo "RPM: $RPM"
	echo "--------------------------------------------------"
}

# ---------------------------------------
function initBuildArea() {
	echo "Initializing the build area."
	mkdir -p "$RPMBUILD"/{SPECS,SOURCES,RPMS,SRPMS,BUILD,BUILDROOT} || { echo "Could not create $RPMBUILD: $?"; return 1; }

	tr_dest=$(createSourceDir traffic_router)

	export MVN_CMD="mvn versions:set -DnewVersion=$TC_VERSION"
	echo $MVN_CMD
	(cd $TR_DIR; $MVN_CMD)
	cp -r "$TR_DIR"/{build,connector,core} "$tr_dest"/. || { echo "Could not copy to $tr_dest: $?"; return 1; }
	cp  "$TR_DIR"/pom.xml "$tr_dest" || { echo "Could not copy to $tr_dest: $?"; return 1; }

	# tar/gzip the source
	tar -czf "$tr_dest".tgz -C "$RPMBUILD/SOURCES" $(basename $tr_dest) || { echo "Could not create tar archive $tr_dest: $?"; return 1; }

	echo "The build area has been initialized."
}

#----------------------------------------
function buildRpmTomcat () {
	echo "Building the rpm for Tomcat."

	cd "$TR_DIR"/tomcat-rpm || { echo "Could not cd to $TR_DIR/tomcat-rpm: $?"; return 1; }
        ./build_rpm.sh
}

# ---------------------------------------

importFunctions
checkEnvironment
adaptEnvironment
initBuildArea
buildRpmTrafficRouter
buildRpmTomcat
