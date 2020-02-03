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

function importFunctions() {
	local script=$(readlink -f "$0")
	local scriptdir=$(dirname "$script")
	export TS_DIR=$(dirname "$scriptdir")
	export TC_DIR=$(dirname "$TS_DIR")
	functions_sh="$TC_DIR/build/functions.sh"
	if [[ ! -r $functions_sh ]]; then
		echo "error: can't find $functions_sh"
		exit 1
	fi
	. "$functions_sh"
}

#----------------------------------------
function initBuildArea() {
	echo "Initializing the build area."
	mkdir -p "$RPMBUILD"/{SPECS,SOURCES,RPMS,SRPMS,BUILD,BUILDROOT} || { echo "Could not create $RPMBUILD: $?"; exit 1; }

	# tar/gzip the source
	local ts_dest=$(createSourceDir traffic_stats)
	cd "$TS_DIR" || \
		 { echo "Could not cd to $TS_DIR: $?"; exit 1; }

        echo "PATH: $PATH"
        echo "GOPATH: $GOPATH"
        go version
        go env

        # get x/* packages (everything else should be properly vendored)
        go get -v golang.org/x/net/publicsuffix || \
                { echo "Could not get go package dependencies"; exit 1; }

        # compile traffic_stats
        go build -v || \
                { echo "Could not build traffic_stats binary"; exit 1; }

	# compile influx_db_tools
	pushd influxdb_tools
	go build -v sync/sync_ts_databases.go || \
                { echo "Could not build sync_ts_databases binary"; exit 1; }
	go build -v create/create_ts_databases.go || \
                { echo "Could not build create_ts_databases binary"; exit 1; }
	popd

	rsync -aLv ./ "$ts_dest"/ || \
		 { echo "Could not copy to $ts_dest: $?"; exit 1; }
	cp "$TS_DIR"/build/*.spec "$RPMBUILD"/SPECS/. || \
		 { echo "Could not copy spec files: $?"; exit 1; }

	tar -czvf "$ts_dest".tgz -C "$RPMBUILD"/SOURCES $(basename $ts_dest) || { echo "Could not create tar archive $ts_dest.tgz: $?"; exit 1; }
	cp "$TS_DIR"/build/*.spec "$RPMBUILD"/SPECS/. || { echo "Could not copy spec files: $?"; exit 1; }

	echo "The build area has been initialized."
}

function preBuildChecks() {
    if [[ -e "$TS_DIR"/traffic_stats ]]; then
        echo "Found $TS_DIR/traffic_stats, please remove before retrying to build"
        exit 1
    fi
}

# ---------------------------------------

importFunctions
preBuildChecks
checkEnvironment go
initBuildArea
buildRpm traffic_stats
