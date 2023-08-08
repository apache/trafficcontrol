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
set -o xtrace

#----------------------------------------
importFunctions() {
	local script scriptdir
	script="$(realpath "$0")"
	scriptdir="$(dirname "$script")"
	TO_DIR="$(dirname "$scriptdir")"
	TC_DIR="$(dirname "$TO_DIR")"
	export TO_DIR TC_DIR
	functions_sh="$TC_DIR/build/functions.sh"
	if [ ! -r "$functions_sh" ]; then
		echo "error: can't find $functions_sh"
		return 1
	fi
	. "$functions_sh"
}

# ---------------------------------------
initBuildArea() {
	echo "Initializing the build area."
	(mkdir -p "$RPMBUILD"
	 cd "$RPMBUILD"
	 mkdir -p SPECS SOURCES RPMS SRPMS BUILD BUILDROOT) || { echo "Could not create $RPMBUILD: $?"; return 1; }

	local dest
	dest="$(createSourceDir traffic_ops)"
	cd "$TO_DIR" || \
		 { echo "Could not cd to $TO_DIR: $?"; return 1; }

	echo "PATH: $PATH"
	echo "GOPATH: $GOPATH"
	go version
	go env

	# get x/* packages (everything else should be properly vendored)
	go mod vendor -v ||
		{ echo "Could not vendor go module dependencies"; return 1; }

	# compile traffic_ops_golang
	cd traffic_ops_golang
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
	go build -v -gcflags "$gcflags" -ldflags "${ldflags} -X main.version=traffic_ops-${TC_VERSION}-${BUILD_NUMBER}.${RHEL_VERSION} -B 0x$(git rev-parse HEAD)" || \
								{ echo "Could not build traffic_ops_golang binary"; return 1; }
	cd -

	# compile db/admin
	(cd app/db
	go build -v -o admin -gcflags "$gcflags" -ldflags "$ldflags" || \
								{ echo "Could not build db/admin binary"; return 1;})

	# compile ToDnssecRefresh.go
	(cd app/bin/checks/DnssecRefresh
	go build -v -o ToDnssecRefresh -gcflags "$gcflags" -ldflags "$ldflags" || \
								{ echo "Could not build ToDnssecRefresh binary"; return 1;})

	# compile db/reencrypt
		(cd app/db/reencrypt
	go build -v -o reencrypt || \
								{ echo "Could not build reencrypt binary"; return 1;})

	# compile db/traffic_vault_migrate
		(cd app/db/traffic_vault_migrate
	go build -v -o traffic_vault_migrate || \
								{ echo "Could not build traffic_vault_migrate binary"; return 1;})

	# compile TO profile converter
	(cd install/bin/convert_profile
	go build -v -gcflags "$gcflags" -ldflags "$ldflags" || \
								{ echo "Could not build convert_profile binary"; return 1; })

	rsync -av etc install "$dest"/ || \
		 { echo "Could not copy to $dest: $?"; return 1; }
	if ! (cd app; rsync -av bin conf db script templates "${dest}/app"); then
		echo "Could not copy to $dest/app"
		return 1
	fi

	# include LICENSE in the tarball
	cp "${TC_DIR}/LICENSE" "$dest"

	tar -czvf "$dest".tgz -C "$RPMBUILD"/SOURCES "$(basename "$dest")" || \
		 { echo "Could not create tar archive $dest.tgz: $?"; return 1; }

	cp "$TO_DIR"/build/traffic_ops.spec "$RPMBUILD"/SPECS/. || \
		 { echo "Could not copy spec files: $?"; return 1; }

	PLUGINS="$(grep -l 'AddPlugin(' "${TO_DIR}/traffic_ops_golang/plugin/"*.go | grep -v 'func AddPlugin(' | xargs -I '{}' basename {} '.go')"
	export PLUGINS

	echo "The build area has been initialized."
}

# ---------------------------------------
importFunctions
checkEnvironment -i go,rsync
initBuildArea
buildRpm traffic_ops
