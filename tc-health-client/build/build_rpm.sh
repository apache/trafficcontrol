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
	HC_DIR="$(dirname "$scriptdir")";
	TC_DIR="$(dirname "$HC_DIR")";
	export HC_DIR TC_DIR;
	functions_sh="$TC_DIR/build/functions.sh";
	if [ ! -r "$functions_sh" ]; then
		echo "error: can't find $functions_sh" >&2;
		return 1;
	fi
	. "$functions_sh";
}

#----------------------------------------
initBuildArea() {
	echo "Initializing the build area for tc-health-client";
	(mkdir -p "$RPMBUILD"
	 cd "$RPMBUILD"
	 mkdir -p SPECS SOURCES RPMS SRPMS BUILD BUILDROOT) || { echo "Could not create $RPMBUILD: $?"; return 1; }

	local dest;
	dest=$(createSourceDir trafficcontrol-health-client);
	cd "$HC_DIR";

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
		go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.BuildTimestamp=$(date +'%Y-%m-%dT%H:%M:%S') -X main.Version=${TC_VERSION}-${BUILD_NUMBER}";
		buildManpage 'tc-health-client';
	)

	echo "build_rpm.sh lsing for logrotate";
	ls -lah .;
	ls -lah ./build;

  cp -Rp config "$dest";
  cp -Rp tmagent "$dest";
  cp -Rp util "$dest";
  cp -p tc-health-client "$dest";
  cp -p tc-health-client.go "$dest";
  cp -p tc-health-client.1 "$dest";
  cp -p tc-health-client.json "$dest"/tc-health-client.sample.json;
  cp -p tc-health-client.service "$dest";
  cp -p build/tc-health-client.logrotate "$dest";


	# include LICENSE in the tarball
	cp "${TC_DIR}/LICENSE" "$dest"

	tar -czvf "$dest".tgz -C "$RPMBUILD"/SOURCES "$(basename "$dest")";
	cp build/trafficcontrol-health-client.spec "$RPMBUILD"/SPECS/.;

	echo "The build area has been initialized.";
}

# buildManpage builds an app's manpage using pandoc.
# It takes 1 argument: the app name.
# The working directory must be of the app, and must contain a README.md formatted like a manpage.
buildManpage() {
	app="$1";
	desc="ATC tc-health-client Manual";
	# prepend the pandoc header to the readme
	printf "%s\n%s\n%s\n%s" "% ${app}(1) ${app} ${TC_VERSION} | ${desc}" "%" "% $(date '+%Y-%m-%d')" "$(cat ./README.md)" > README.md
	pandoc --standalone --to man README.md -o "${app}.1";
}

#----------------------------------------
importFunctions;
checkEnvironment go;
initBuildArea;
buildRpm trafficcontrol-health-client;
