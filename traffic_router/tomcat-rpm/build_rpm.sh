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
checkEnvironment() {
	echo "Verifying the build configuration environment."
	local script scriptdir
	script="$(realpath "$0")"
	scriptdir="$(dirname "$script")"
	TR_DIR='' TC_DIR=''
	TR_DIR="$(dirname "$scriptdir")"
	TC_DIR="$(dirname "$TR_DIR")"
	export TR_DIR TC_DIR
	functions_sh="$TC_DIR/build/functions.sh"
	if [ ! -r "$functions_sh" ]; then
		echo "Error: Can't find $functions_sh"
		exit 1
	fi
	. "$functions_sh"

	#
	# get traffic_control src path -- relative to build_rpm.sh script
	export PACKAGE="tomcat"
	export WORKSPACE=${WORKSPACE:-$TC_DIR}
	export RPMBUILD="$WORKSPACE/rpmbuild"
	export RPM_TARGET_OS="${RPM_TARGET_OS:-linux}"
	export DIST="$WORKSPACE/dist"
	# Forcing BUILD NUMBER to 1 since this is outside the tree and related to Tomcat Release
	export BUILD_NUMBER=1
	export RPM="${PACKAGE}-${TOMCAT_VERSION}-${BUILD_NUMBER}.${RHEL_VERSION}.noarch.rpm"
	export SRPM="${PACKAGE}-${TOMCAT_VERSION}-${BUILD_NUMBER}.${RHEL_VERSION}.src.rpm"


	echo "=================================================="
	echo "WORKSPACE: $WORKSPACE"
	echo "TOMCAT_VERSION: $TOMCAT_VERSION"  #defined in traffic_router
	echo "BUILD_NUMBER: $BUILD_NUMBER"      #defined in traffic_router
	echo "RPM: $RPM"
	echo "--------------------------------------------------"
}

# ---------------------------------------
initBuildArea() {
	echo "Initializing the build area."
	(mkdir -p "$RPMBUILD"
	 cd "$RPMBUILD"
	 mkdir -p SPECS SOURCES RPMS SRPMS BUILD BUILDROOT) || { echo "Could not create $RPMBUILD: $?"; return 1; }
	export VERSION=$TOMCAT_VERSION

	echo "Downloading Tomcat $VERSION..."
	curl -fo "${RPMBUILD}/SOURCES/apache-tomcat-${VERSION}.tar.gz" "https://archive.apache.org/dist/tomcat/tomcat-${VERSION%%.*}/v${VERSION}/bin/apache-tomcat-${VERSION}.tar.gz" || \
	{ echo "Could not download Tomcat $VERSION: $?"; exit 1; }

	cp "$TR_DIR/tomcat-rpm/tomcat.service" "$RPMBUILD/SOURCES/" || { echo "Could not copy source files: $?"; exit 1; }
	cp "$TR_DIR/tomcat-rpm/tomcat.spec" "$RPMBUILD/SPECS/" || { echo "Could not copy spec files: $?"; exit 1; }

	echo "The build area has been initialized."
}

#----------------------------------------
buildRpmTomcat () {
	export SPEC_FILE_NAME=tomcat.spec
	buildRpmForEl 7
}

buildRpmForEl () {
	echo "Building the rpm for ${RHEL_VERSION}."

	cd "$RPMBUILD"

	build_flags="-ba";
	if [[ "$NO_SOURCE" -eq 1 ]]; then
		build_flags="-bb";
	fi


	# build RPM with xz level 2 compression
	rpmbuild --define "_topdir $(pwd)" \
		--define "build_number $BUILD_NUMBER.$RHEL_VERSION" \
		--define "tomcat_version $TOMCAT_VERSION" \
		--define "_target_os ${RPM_TARGET_OS}" \
		--define '%_source_payload w2.xzdio' \
		--define '%_binary_payload w2.xzdio' \
		$build_flags SPECS/$SPEC_FILE_NAME ||
		{ echo "RPM BUILD FAILED: $?"; exit 1; }
	local rpm
	local srpm
	rpm="$(find ./RPMS -name "$RPM")"
	srpm="$(find ./SRPMS -name "$SRPM")";
	if [ -z "$rpm" ]; then
		echo "Could not find rpm file $RPM in $(pwd)"
		exit 1;
	fi
	echo "========================================================================================"
	echo "RPM BUILD SUCCEEDED, See $DIST/$RPM for the newly built rpm."
	echo "========================================================================================"
	echo
	mkdir -p "$DIST" || { echo "Could not create $DIST: $?"; exit 1; }

	rpmDest="."
	srcRPMDest="."
	if [[ "$SIMPLE" -eq 1 ]]; then
		rpmDest="tomcat.rpm";
		srcRPMDest="tomcat.src.rpm";
	fi

	cp -f "$RPMBUILD/RPMS/noarch/${RPM}" "$DIST/$rpmDest";
	code="$?";
	if [[ "$code" -ne 0 ]]; then
		echo "Could not copy $rpm to $DIST: $code" >&2;
		return "$code";
	fi

	if [[ "$NO_SOURCE" -eq 1 ]]; then
		return 0;
	fi

	cp -f "$RPMBUILD/SRPMS/${SRPM}" "$DIST/$srcRPMDest";
	code="$?";
	if [[ "$code" -ne 0 ]]; then
		echo "Could not copy $srpm to $DIST: $code" >&2;
		return "$code";
	fi
}

checkEnvironment -i curl
initBuildArea
buildRpmTomcat
