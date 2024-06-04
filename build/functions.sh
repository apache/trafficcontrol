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

# macOS's version of realpath does not resolve symlinks, so we add a function
# for it.
get_realpath() {
	local bin
	local found=''
	first_realpath="$(type -P realpath)"
	for bin in $(type -aP grealpath realpath | uniq); do
		if "$bin" -e . >/dev/null 2>&1; then
			found=y
			break
		fi
	done
	if [[ -n "$found" ]]; then
		if [[ "$first_realpath" == "$bin" ]]; then
			# Default realpath works.
			return
		fi
		realpath_path="$bin"
		# by default, macOS does not have realpath
		eval "$(<<FUNCTION cat
		realpath() {
			"$realpath_path" "\$@"
		}
FUNCTION
		)"
		export -f realpath
	else
			cat <<'MESSAGE'
GNU realpath is required to build Apache Traffic Control if your
realpath binary does not support the -e flag, as is the case on BSD-like
operating systems like macOS. Install it by running the following
command:
    brew install coreutils
MESSAGE
			exit 1
	fi
}
get_realpath

if { ! stat -c%u . >/dev/null && stat -f%u .; } >/dev/null 2>&1; then
	#BSD stat uses -f as its formatting flag instead of -c
	stat() {
		local format=''
		while getopts c: option; do
			case "$option" in
				c) format="$OPTARG";;
				*) return 1
			esac
		done
		shift $(( OPTIND - 1 ))
		unset OPTIND
		$(which stat) "-f${format}" "$@"
	}
	export -f stat
fi

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
		# The number of commits since the last tag
		if ! commits="$(git describe --long --tags \
			--match='RELEASE-[0-9].[0-9].[0-9]' \
			--match='RELEASE-[0-9][0-9].[0-9][0-9].[0-9][0-9]' \
			--match='v[0-9].[0-9].[0-9]' \
			--match='v[0-9][0-9].[0-9][0-9].[0-9][0-9]' |
			awk -F- '{print $(NF-1)}')";
		then
			commits=0
		fi
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
getGoVersion() {
	local directory="$1"
	local go_version_file="$directory/GO_VERSION"
	[ -r "$go_version_file" ] || { echo "Could not read $go_version_file: $!"; return 1; }
	cat "$go_version_file"
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
	local releasever=${RHEL_VERSION:-}
	if [ -n "$releasever" ]; then
		echo "el${releasever}"
		return
	fi

	local redhat_release=/etc/redhat-release
	local default_version=7
	if [ -e $redhat_release ]; then
		releasever="$(rpm -q --qf '%{version}' -f $redhat_release)"
		releasever="${releasever%%.*}"
	else
		echo "${redhat_release} not found, defaulting to major release ${default_version}" >/dev/stderr
		releasever=${default_version}
	fi;

	echo "el${releasever}"
}

# ---------------------------------------
getCommit() {
	local buildNum
	buildNum="$(getBuildNumber)"
	echo "${buildNum%.*}"
}

# ---------------------------------------
checkEnvironment() {
	include_programs='' exclude_programs=''
	while getopts i:e: option; do
		case "$option" in
			i) include_programs="$OPTARG";;
			e) exclude_programs="$OPTARG";;
			*) return 1
		esac
	done
	unset OPTIND
	include_file="$(mktemp)"
	exclude_file="$(mktemp)"
	printf '%s\n' git rpmbuild  ${include_programs//,/ } | sort >"$include_file"
	printf '%s\n' ${exclude_programs//,/ } | sort >"$exclude_file"
	programs="$(comm -23 "$include_file" "$exclude_file")"
	rm "$include_file" "$exclude_file"

	# verify required tools available in path -- extra tools required by subsystem are passed in
	for program in $programs; do
		if ! type "$program"; then
			echo "$program not found in PATH"
			return 1
		fi
	done
	# verify git version
	requiredGitVersion=1.7.12
	if ! versionOk "$(git --version | tr -dc 0-9. )" "$requiredGitVersion"; then
		echo "$(git --version) must be at least $requiredGitVersion"
		return 1
	fi

	TC_VERSION='' BUILD_NUMBER='' RPMBUILD='' DIST=''
	TC_VERSION="$(getVersion "$TC_DIR")"
	BUILD_NUMBER="$(getBuildNumber)"
	GO_VERSION="$(getGoVersion "$TC_DIR")"
	RHEL_VERSION="$(getRhelVersion)"
	WORKSPACE="${WORKSPACE:-$TC_DIR}"
	RPMBUILD="$WORKSPACE/rpmbuild"
	GOOS="${GOOS:-linux}"
	RPM_TARGET_OS="${RPM_TARGET_OS:-$GOOS}"
	DIST="$WORKSPACE/dist"
	export TC_VERSION BUILD_NUMBER GO_VERSION RHEL_VERSION WORKSPACE RPMBUILD GOOS RPM_TARGET_OS DIST

	mkdir -p "$DIST" || { echo "Could not create ${DIST}: ${?}"; return 1; }

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
		local arch
		arch="$(rpm --eval %_arch)"
		rpm="${pre}.${arch}.rpm"
		local srpm="${pre}.src.rpm"
		echo "Building the rpm."
		{ set +o nounset
		set -- # Clear arguments for reuse
		if [ "$DEBUG_BUILD" = true ]; then
			echo 'RPM will not strip binaries before packaging.';
			set -- "$@" --define '%__os_install_post %{nil}' # Do not strip binaries before packaging
		fi;
		set -- "$@" --define '%_source_payload w2.xzdio' # xz level 2 compression for text files
		set -- "$@" --define '%_binary_payload w2.xzdio' # xz level 2 compression for binary files
		set -o nounset; }

		build_flags="-ba";
		if [[ "$NO_SOURCE" -eq 1 ]]; then
			build_flags="-bb";
		fi

		pushd "$RPMBUILD";

		rpmbuild --define "_topdir $(pwd)" \
			--define "traffic_control_version $TC_VERSION" \
			--define "go_version $GO_VERSION" \
			--define "commit $(getCommit)" \
			--define "build_number $BUILD_NUMBER.$RHEL_VERSION" \
			--define "rhel_vers $RHEL_VERSION" \
			--define "_target_os $RPM_TARGET_OS" \
			"$build_flags" SPECS/$package.spec \
			"$@";
		code=$?
		if [[ "$code" -ne 0 ]]; then
			echo "RPM BUILD FAILED: $code" >&2;
			return $code;
		fi

		echo
		echo "========================================================================================"
		echo "RPM BUILD FOR $package SUCCEEDED, See $DIST/$rpm for the newly built rpm."
		echo "========================================================================================"
		echo

		rpmDest=".";
		srcRPMDest=".";
		if [[ "$SIMPLE" -eq 1 ]]; then
			rpmDest="${package}.rpm";
			srcRPMDest="${package}.src.rpm";
		fi

		cp -f "$RPMBUILD/RPMS/${arch}/$rpm" "$DIST/$rpmDest";
		code="$?";
		if [[ "$code" -ne 0 ]]; then
			echo "Could not copy $rpm to $DIST: $code" >&2;
			return "$code";
		fi

		if [[ "$NO_SOURCE" -eq 1 ]]; then
			return 0;
		fi

		cp -f "$RPMBUILD/SRPMS/$srpm" "$DIST/$srcRPMDest";
		code="$?";
		if [[ "$code" -ne 0 ]]; then
			echo "Could not copy $srpm to $DIST: $code" >&2;
			return "$code";
		fi
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
# verify whether the minor-level version of the go compiler is greater or equal
# to the minor-level version in the GO_VERSION file: returns 0 if if not.
# returns 1 if it is.
#
verify_and_set_go_version() {
	if [ -v GO_VERSION ]; then
		GO_VERSION="$(getGoVersion .)"
	else
		GO_VERSION=''
	fi
	local major_version="$(echo "$GO_VERSION" | awk -F. '{print $1}')"
	local minor_version="$(echo "$GO_VERSION" | awk -F. '{print $2}')"
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
			group_1="$(echo "$go_version" | sed -E "s/${version_pattern}/\1/")"
			group_2="$(echo "$go_version" | sed -E "s/${version_pattern}/\2/")"
			if [ ! "$group_1" -ge "$major_version" ] || [ ! "$group_2" -ge "$minor_version" ]; then
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
		echo "ERROR: this build needs go ${major_version}.${minor_version} or greater and no usable go compiler was found, found GO_VERSION: $GO_VERSION"
		unset GO_VERSION
		return 1
	fi
}
