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
set -o errexit -o nounset -o pipefail -o xtrace;

#----------------------------------------
importFunctions() {
	local script scriptdir;
	script="$(realpath "$0")";
	scriptdir="$(dirname "$script")";
	T3C_DIR="$(dirname "$scriptdir")";
	TC_DIR="$(dirname "$T3C_DIR")";
	export T3C_DIR TC_DIR;
	functions_sh="$TC_DIR/build/functions.sh";
	if [ ! -r "$functions_sh" ]; then
		echo "error: can't find $functions_sh" >&2;
		return 1;
	fi
	. "$functions_sh";
}

#----------------------------------------
initBuildArea() {
	echo "Initializing the build area for t3c";
	(mkdir -p "$RPMBUILD"
	 cd "$RPMBUILD"
	 mkdir -p SPECS SOURCES RPMS SRPMS BUILD BUILDROOT) || { echo "Could not create $RPMBUILD: $?"; return 1; }

	local dest;
	dest=$(createSourceDir trafficcontrol-cache-config);
	cd "$T3C_DIR";

	echo "PATH: $PATH";
	echo "GOPATH: $GOPATH";
	go version;
	go env;

	go mod vendor -v;

	gcflags=''
	ldflags=''
	export CGO_ENABLED=0
	{ set +o nounset;
	if [ "$DEBUG_BUILD" = true ]; then
		echo 'DEBUG_BUILD is enabled, building without optimization or inlining...';
		gcflags="${gcflags} all=-N -l";
	else
		ldflags="${ldflags} -s -w"; # strip binary
	fi;
	set -o nounset; }

	(
		cd t3c;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c';
	)

	(
		cd t3c-apply;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c-apply';
	)

	(
		cd t3c-generate;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c-generate';
	)

	(
		cd t3c-request;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c-request';
	)

	(
		cd t3c-update;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c-update';
	)

	(
		cd t3c-check;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c-check';
	)

	(
		cd t3c-check-refs;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c-check-refs';
	)

	(
		cd t3c-check-reload;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c-check-reload';
	)

	(
		cd t3c-diff;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c-diff';
	)

	(
		cd t3c-tail;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c-tail';
	)

	(
		cd t3c-preprocess;
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.GitRevision=$(git rev-parse HEAD) -X main.BuildTimestamp=$(date +'%Y-%M-%dT%H:%M:%s') -X main.Version=${TC_VERSION}";
		buildManpage 't3c-preprocess';
	)

	mkdir -p "${dest}/build";

	echo "build_rpm.sh lsing for logrotate";
	ls -lah .;
	ls -lah ./build;

	cp -p build/atstccfg.logrotate "$dest"/build;

	# include LICENSE in the tarball
	cp "${TC_DIR}/LICENSE" "$dest"

	tar -czvf "$dest".tgz -C "$RPMBUILD"/SOURCES "$(basename "$dest")";
	cp build/trafficcontrol-cache-config.spec "$RPMBUILD"/SPECS/.;
	cp build/atstccfg.logrotate "$RPMBUILD"/.;

	echo "The build area has been initialized.";
}

# buildManpage builds an app's manpage using pandoc.
# It takes 1 argument: the app name.
# The working directory must be of the app, and must contain a README.md formatted like a manpage.
buildManpage() {
	app="$1";
	desc="ATC t3c Manual";
	# prepend the pandoc header to the readme
	printf "%s\n%s\n%s\n%s" "% ${app}(1) ${app} ${TC_VERSION} | ${desc}" "%" "% $(date '+%Y-%m-%d')" "$(cat ./README.md)" > README.md
	pandoc --standalone --to man README.md -o "${app}.1";
}

#----------------------------------------
importFunctions;
checkEnvironment go;
initBuildArea;
buildRpm trafficcontrol-cache-config;
