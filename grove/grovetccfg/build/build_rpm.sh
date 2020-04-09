#!/usr/bin/env bash

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

#----------------------------------------
function importFunctions() {
	[ ! -z "$TC_DIR" ] || { echo "Cannot find repository root." >&2 ; exit 1; }
	export TC_DIR
	functions_sh="$TC_DIR/build/functions.sh"
	if [[ ! -r $functions_sh ]]; then
		echo "Error: Can't find $functions_sh"
		exit 1
	fi
	. "$functions_sh"
}

#----------------------------------------
function checkGroveEnvironment() {
	echo "Verifying the build configuration environment."
	local script=$(readlink -f "$0")
	local scriptdir=$(dirname "$script")

	export GROVETC_DIR=$(dirname "$scriptdir")
	export GROVE_DIR=$(dirname "$GROVETC_DIR")
	export GROVE_VERSION=`cat ${GROVE_DIR}/VERSION`
	export PACKAGE="grovetccfg"
	export BUILD_NUMBER=${BUILD_NUMBER:-$(getBuildNumber)}
	export RPMBUILD="${GROVE_DIR}/rpmbuild"
	export DIST="${TC_DIR}/dist"
	export RPM="${PACKAGE}-${GROVE_VERSION}-${BUILD_NUMBER}.x86_64.rpm"

	# grovetccfg needs to be built with go 1.14 or greater
	verify_and_set_go_version
	if [[ $? -ne 1 ]]; then
		exit 0
	fi

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
function initBuildArea() {
	cd "$GROVETC_DIR"

	# prep build environment
	[ -e $RPMBUILD ] && rm -rf $RPMBUILD
	[ ! -e $RPMBUILD ] || { echo "Failed to clean up rpm build directory '$RPMBUILD': $?" >&2; exit 1; }
	mkdir -p $RPMBUILD/{BUILD,RPMS,SOURCES} || { echo "Failed to create build directory '$RPMBUILD': $?" >&2; exit 1; }
}

# ---------------------------------------
function buildRpmGrove() {
	# build
	$GO get -v -d . || { echo "Failed to go get dependencies: $?" >&2; exit 1; }
	$GO build -v -ldflags "-X main.Version=$GROVE_VERSION" || { echo "Failed to build $PACKAGE: $?" >&2; exit 1; }

	# tar
	tar -cvzf $RPMBUILD/SOURCES/${PACKAGE}-${GROVE_VERSION}.tgz ${PACKAGE}|| { echo "Failed to create archive for rpmbuild: $?" >&2; exit 1; }

	# Work around bug in rpmbuild. Fixed in rpmbuild 4.13.
	# See: https://github.com/rpm-software-management/rpm/commit/916d528b0bfcb33747e81a57021e01586aa82139
	# Takes ownership of the spec file.
	spec=build/${PACKAGE}.spec
	spec_owner=$(stat -c%u $spec)
	spec_group=$(stat -c%g $spec)
	if ! id $spec_owner >/dev/null 2>&1; then
	  chown $(id -u):$(id -g) build/${PACKAGE}.spec

	  function give_spec_back {
		chown ${spec_owner}:${spec_group} build/${PACKAGE}.spec
	  }
	  trap give_spec_back EXIT
	fi

	# build RPM
	rpmbuild --define "_topdir $RPMBUILD" --define "version ${GROVE_VERSION}" --define "build_number ${BUILD_NUMBER}" -ba build/${PACKAGE}.spec || { echo "rpmbuild failed: $?" >&2; exit 1; }

	# copy build RPM to .
	[ -e $DIST ] || mkdir -p $DIST
	cp $RPMBUILD/RPMS/x86_64/${RPM} $DIST
}

importFunctions
checkGroveEnvironment
initBuildArea
buildRpmGrove

