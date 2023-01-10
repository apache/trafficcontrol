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
	[ -n "$TC_DIR" ] || { echo "Cannot find repository root." >&2 ; exit 1; }
	export TC_DIR
	functions_sh="$TC_DIR/build/functions.sh"
	if [ ! -r "$functions_sh" ]; then
		echo "Error: Can't find $functions_sh"
		return 1
	fi
	. "$functions_sh"
}

#----------------------------------------
checkGroveEnvironment() {
	echo "Verifying the build configuration environment."
	local script scriptdir
	script=$(realpath "$0")
	scriptdir=$(dirname "$script")

	GROVETC_DIR='' GROVE_DIR='' GROVE_VERSION='' PACKAGE='' RPMBUILD='' DIST='' RPM=''
	GROVETC_DIR=$(dirname "$scriptdir")
	GROVE_DIR=$(dirname "$GROVETC_DIR")
	GROVE_VERSION="$(cat "${GROVE_DIR}/VERSION")"
	PACKAGE="grovetccfg"
	BUILD_NUMBER=${BUILD_NUMBER:-$(getBuildNumber)}
	RPMBUILD="${GROVE_DIR}/rpmbuild"
	DIST="${TC_DIR}/dist"
	RPM="${PACKAGE}-${GROVE_VERSION}-${BUILD_NUMBER}.${RHEL_VERSION}.$(rpm --eval %_arch).rpm"
	SRPM="${PACKAGE}-${GROVE_VERSION}-${BUILD_NUMBER}.${RHEL_VERSION}.src.rpm"
	GOOS="${GOOS:-linux}"
	RPM_TARGET_OS="${RPM_TARGET_OS:-$GOOS}"
	export GROVETC_DIR GROVE_DIR GROVE_VERSION PACKAGE BUILD_NUMBER RPMBUILD DIST RPM GOOS RPM_TARGET_OS

	echo "=================================================="
	echo "GO_VERSION: $GO_VERSION"
	echo "TC_DIR: $TC_DIR"
	echo "PACKAGE: $PACKAGE"
	echo "GROVE_DIR: $GROVE_DIR"
	echo "GROVETC_DIR: $GROVETC_DIR"
	echo "GROVE_VERSION: $GROVE_VERSION"
	echo "BUILD_NUMBER: $BUILD_NUMBER"
	echo "DIST: $DIST"
	echo "RPM: $RPM"
	echo "RPMBUILD: $RPMBUILD"
	echo "--------------------------------------------------"
}

# ---------------------------------------
initBuildArea() {
	cd "$GROVETC_DIR"

	# prep build environment
	[ -e "$RPMBUILD" ] && rm -rf "$RPMBUILD"
	[ ! -e "$RPMBUILD" ] || { echo "Failed to clean up rpm build directory '$RPMBUILD': $?" >&2; return 1; }
	(mkdir -p "$RPMBUILD"
	 cd "$RPMBUILD"
	 mkdir -p BUILD RPMS SOURCES) || { echo "Failed to create build directory '$RPMBUILD': $?" >&2; return 1; }
}

# ---------------------------------------
buildRpmGrove() {
	# build
	ldflags='-s -w'
	tags='osusergo netgo'
	go mod vendor -v || { echo "Failed to vendor go dependencies: $?" >&2; return 1; }
	go build -v -ldflags "${ldflags} -X main.Version=$GROVE_VERSION" -tags "$tags" || { echo "Failed to build $PACKAGE: $?" >&2; return 1; }

	# tar
	tar -cvzf "${RPMBUILD}/SOURCES/${PACKAGE}-${GROVE_VERSION}.tgz" ${PACKAGE}|| { echo "Failed to create archive for rpmbuild: $?" >&2; return 1; }

	# Work around bug in rpmbuild. Fixed in rpmbuild 4.13.
	# See: https://github.com/rpm-software-management/rpm/commit/916d528b0bfcb33747e81a57021e01586aa82139
	# Takes ownership of the spec file.
	spec=build/${PACKAGE}.spec
	spec_owner=$(stat -c%u $spec)
	spec_group=$(stat -c%g $spec)
	if ! id "$spec_owner" >/dev/null 2>&1; then
		chown "$(id -u):$(id -g)" build/${PACKAGE}.spec

		give_spec_back() {
		chown "${spec_owner}:${spec_group}" build/${PACKAGE}.spec
		}
		trap give_spec_back EXIT
	fi

	build_flags="-ba";
	if [[ "$NO_SOURCE" -eq 1 ]]; then
		build_flags="-bb";
	fi


	# build RPM with xz level 2 compression
	rpmbuild \
		--define "_topdir $RPMBUILD" \
		--define "version ${GROVE_VERSION}" \
		--define "build_number ${BUILD_NUMBER}.${RHEL_VERSION}" \
		--define "_target_os ${RPM_TARGET_OS}" \
		--define '%_source_payload w2.xzdio' \
		--define '%_binary_payload w2.xzdio' \
		$build_flags build/${PACKAGE}.spec ||
		{ echo "rpmbuild failed: $?" >&2; return 1; }


	rpmDest=".";
	srcRPMDest=".";
	if [[ "$SIMPLE" -eq 1 ]]; then
		rpmDest="grovetccfg.rpm";
		srcRPMDest="grovetccfg.src.rpm";
	fi

	# copy build RPM to .
	[ -d "$DIST" ] || mkdir -p "$DIST";

	cp -f "$RPMBUILD/RPMS/$(rpm --eval %_arch)/${RPM}" "$DIST/$rpmDest";
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

importFunctions
checkEnvironment -i go
checkGroveEnvironment
initBuildArea
buildRpmGrove
