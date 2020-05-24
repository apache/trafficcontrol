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

removeFirstArg() {
	shift
	echo "$@"
}

# ---------------------------------------
# versionOk checks version number against required version.
#   ``versionOk 1.2.3 2.0.4.7'' returns false value indicating
#       version you have is not at least version you need
#   if versionOk $haveversion $needversion; then
#      echo "Need at least version $needversion"; return 1
#   fi
versionOk() {
	local h="$1" n="$2"
	# string compare -- no need to do more if the same
	[ "$h" = "$n" ] && return 0

	# split into fields
	local have="${h//\./ }"
	local need="${n//\./ }"
	# cmp first entry of each array.  Bail when unequal.
	while [ -n "$have" ] && [ "$have" = "$need" ]; do
		# pop 1st entry from each
		have="$(removeFirstArg ${have})"
		need="$(removeFirstArg ${need})"
	done
	have="${have%% *}";
	need="${need%% *}";
	if [ "${have:-0}" -lt "${need:-0}" ]; then
		return 1
	fi
	return 0
}

# ---------------------------------------
getRevCount() {
	local buildNum
	buildNum=$(getBuildNumber)
	echo "${buildNum%.*}"
}

# ---------------------------------------
isInGitTree() {
	# ignore output -- use exit status
	git rev-parse --is-inside-work-tree >/dev/null 2>&1
}

# ---------------------------------------
getBuildNumber() {
	local in_git=''
	if isInGitTree; then
		local commits sha
		commits="$(git rev-list HEAD 2>/dev/null | wc -l)"
		sha="$(git rev-parse --short=8 HEAD)"
		echo "$commits.$sha"
	else
		# Expect it's from the released tarball -- if BUILD_NUMBER file is not present,  abort
		if [ ! -f "${TC_DIR}/BUILD_NUMBER" ]; then
			echo "Not in git repository and no BUILD_NUMBER present -- aborting!"
			return 1
		fi
		grep -v '^#' "${TC_DIR}/BUILD_NUMBER"
	fi
}

# ---------------------------------------
getVersion() {
	local d="$1"
	local vf="$d/VERSION"
	[ -r "$vf" ] || { echo "Could not read $vf: $!"; return 1; }
	cat "$vf"
}

# ---------------------------------------
getRhelVersion() {
				echo el$(rpm -q --qf "%{VERSION}" $(rpm -q --whatprovides redhat-release))
}

# ---------------------------------------
getCommit() {
	local buildNum
	buildNum="$(getBuildNumber)"
	echo "${buildNum%.*}"
}

# ---------------------------------------
checkEnvironment() {
	TC_VERSION='' BUILD_NUMBER='' RHEL_VERSION='' RPMBUILD='' DIST=''
	TC_VERSION="$(getVersion "$TC_DIR")"
	BUILD_NUMBER="$(getBuildNumber)"
	RHEL_VERSION="$(getRhelVersion)"
	WORKSPACE="${WORKSPACE:-$TC_DIR}"
	RPMBUILD="$WORKSPACE/rpmbuild"
	DIST="$WORKSPACE/dist"
	export TC_VERSION BUILD_NUMBER RHEL_VERSION WORKSPACE RPMBUILD DIST

	mkdir -p "$DIST" || { echo "Could not create ${DIST}: ${?}"; return 1; }

	# verify required tools available in path -- extra tools required by subsystem are passed in
	for pgm in git rpmbuild "$@"; do
		type "$pgm" 2>/dev/null || { echo "$pgm not found in PATH"; }
	done
	# verify git version
	requiredGitVersion=1.7.12
	if ! versionOk "$(git --version | tr -dc 0-9. )" "$requiredGitVersion"; then
		echo "$(git --version) must be at least $requiredGitVersion"
		return 1
	fi
	echo "Build environment has been verified."

	echo "=================================================="
	echo "WORKSPACE: $WORKSPACE"
	echo "BUILD_NUMBER: $BUILD_NUMBER"
	echo "RHEL_VERSION: $RHEL_VERSION"
	echo "TC_VERSION: $TC_VERSION"
	echo "--------------------------------------------------"
}

# ---------------------------------------
createSourceDir() {
	local target="${1}-${TC_VERSION}"
	local srcpath="$RPMBUILD/SOURCES/$target"
	mkdir -p "$srcpath" || { echo "Could not create $srcpath: $?"; return 1; }
	echo "$srcpath"
}

# ---------------------------------------
buildRpm() {
	for package in "$@"; do
		local pre="${package}-${TC_VERSION}-${BUILD_NUMBER}.${RHEL_VERSION}"
		local rpm
		rpm="${pre}.$(uname -m).rpm"
		local srpm="${pre}.src.rpm"
		echo "Building the rpm."
		if [[ "$DEBUG_BUILD" == true ]]; then
			echo 'RPM will not strip binaries before packaging.';
			echo '%__os_install_post %{nil}' >> /etc/rpm/macros; # Do not strip binaries before packaging
		fi;

		cd "$RPMBUILD" && \
			rpmbuild --define "_topdir $(pwd)" \
				 --define "traffic_control_version $TC_VERSION" \
				 --define "commit $(getCommit)" \
				 --define "build_number $BUILD_NUMBER.$RHEL_VERSION" \
				 -ba SPECS/$package.spec || \
				 { echo "RPM BUILD FAILED: $?"; return 1; }

		echo
		echo "========================================================================================"
		echo "RPM BUILD FOR $package SUCCEEDED, See $DIST/$rpm for the newly built rpm."
		echo "========================================================================================"
		echo

		cp "$RPMBUILD/RPMS/$(uname -m)/$rpm" "$DIST/." || { echo "Could not copy $rpm to $DIST: $?"; return 1; }
		cp "$RPMBUILD/SRPMS/$srpm" "$DIST/." || { echo "Could not copy $srpm to $DIST: $?"; return 1; }
	done
}

# ---------------------------------------
createTarball() {
	local projDir
	projDir="$(cd "$1"; pwd)"
	local projName=trafficcontrol
	local version
	version="$(getVersion "$TC_DIR")"
	local tarball="dist/apache-${projName}-${version}.tar.gz"
	local tarDir
	tarDir="$(basename "$tarball" .tar.gz)"
	local tarLink
	( trap 'rm -f "${projDir}/BUILD_NUMBER" "$projLink"' EXIT
	projLink="$(dirname "$projDir")/${tarDir}"
	ln -sf "$projDir" "$projLink"

	# Create a BUILDNUMBER file and add to tarball
	getBuildNumber >"${projDir}/BUILD_NUMBER"

	# create the tarball from BUILD_NUMBER and files tracked in git
	git ls-files |
		sed "s|^|${tarDir}/|g" |
		tar -czf "$tarball" -C "${projLink}/.." -T - "${tarDir}/BUILD_NUMBER"
		# ^ read files from stdin to avoid tar argument limit
	)
	echo "$tarball"
}

# ---------------------------------------
createDocsTarball() {
	local projDir
	projDir="$(cd "$1"; pwd)"
	local projName=trafficcontrol
	local version
	version="$(getVersion "$TC_DIR")"
	local tarball="dist/apache-$projName-$version-docs.tar.gz"
	local tarDir="${projDir}/docs/build/"

	# Create a BUILDNUMBER file and add to tarball
	getBuildNumber >"${tarDir}/BUILD_NUMBER"

	# create the tarball only from files in repo and BUILD_NUMBER
	tar -czf "$tarball" -C "$tarDir" .
	rm "${tarDir}/BUILD_NUMBER"
	echo "$tarball"
}

# ----------------------------------------
# verify if the go compiler is version 1.14 or higher, returns 0 if if not. returns 1 if it is.
#
verify_and_set_go_version() {
	GO_VERSION="none"
	GO="none"
	go_in_path="$(type -p go)"
	local group_1='' group_2=''
	for g in "$go_in_path" /usr/bin/go /usr/local/go/bin/go; do
		if [ -z "$g" ] || [ ! -x "$g" ]; then
			continue
		fi

		go_version="$($g version | awk '{print $3}')"

		version_pattern='.*go([1-9]+)\.([1-9]+).*'
		if echo "$go_version" | grep -E "$version_pattern"; then
			group_1="$(echo "$go_version" | sed -E "s/.*${version_pattern}/\1/")"
			group_2="$(echo "$go_version" | sed -E "s/${version_pattern}/\2/")"
			if [ ! "$group_1" -ge 1 ] || [ ! "$group_2" -ge 14 ]; then
				GO_VERSION="${group_1}.${group_2}"; export GO_VERSION
				echo "go version for $g is ${group_1}.${group_2}"
				continue
			fi
			GO_VERSION="${group_1}.${group_2}"; export GO_VERSION
			GO=$g; export GO
			PATH="$(dirname "$g"):${PATH}"; export PATH
			echo "go version for $g is ${group_1}.${group_2}"
			echo "will use $g"
			return 0
		fi
	done

	if [ "$GO" = none ]; then
		echo "ERROR: this build needs go 1.14 or greater and no usable go compiler was found, found GO_VERSION: $GO_VERSION"
		return 1
	fi
}

