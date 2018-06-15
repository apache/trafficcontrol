#!/bin/bash
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

set -x


function importFunctions() {
	local script=$(readlink -f "$0")
	local scriptdir=$(dirname "$script")
	export ASTATS_DIR=$(dirname "$scriptdir")
	export TC_DIR=$(dirname "$TM_DIR")
	functions_sh="$TC_DIR/build/functions.sh"
	if [[ ! -r $functions_sh ]]; then
		echo "error: can't find $functions_sh"
		exit 1
	fi
	. "$functions_sh"
}

function installTsDeps() {
  ts_git_branch=${1:-master}
  ts_git_repo=${2:-"https://github.com/apache/trafficserver.git"}
  ATS_PREFIX=/tmp/trafficserver
  git clone $ts_git_repo --depth 1 --branch $ts_git_branch /tmp/trafficserver
  cd /tmp/trafficserver
  autoreconf -i 
  ./configure --enable-experimental-plugins --with-user=ats --with-group=ats 
  mkdir -p /opt/trafficserver/bin
  mkdir -p /opt/trafficserver/include/ts
  cp ${ATS_PREFIX}/tools/tsxs /opt/trafficserver/bin/tsxs
  cp ${ATS_PREFIX}/lib/ts/apidefs.h /opt/trafficserver/include/ts/apidefs.h
  cp ${ATS_PREFIX}/proxy/api/experimental.h /opt/trafficserver/include/experimental.h
  cp ${ATS_PREFIX}/mgmt/api/include/mgmtapi.h /opt/trafficserver/include/mgmtapi.h
  cp ${ATS_PREFIX}/proxy/api/ts/ts.h /opt/trafficserver/include/ts/ts.h
  cp ${ATS_PREFIX}/proxy/api/ts/remap.h /opt/trafficserver/include/ts/remap.h
  chmod +x /opt/trafficserver/bin/tsxs
}

function initBuildArea() {
	echo "Initializing the build area."
	mkdir -p "$RPMBUILD"/{SPECS,SOURCES,RPMS,SRPMS,BUILD,BUILDROOT} || { echo "Could not create $RPMBUILD: $?"; exit 1; }

	# tar/gzip the source
	local astats_dest=$(createSourceDir astats_over_http)
	cd "$ASTATS_DIR" || \
		 { echo "Could not cd to $ASTATS_DIR: $?"; exit 1; }
	rsync -av ./ "$astats_dest"/ || \
		 { echo "Could not copy to $astats_dest: $?"; exit 1; }

	cp -r "$ASTATS_DIR"/ "$astats_dest" || { echo "Could not copy $ASTATS_DIR to $astats_dest: $?"; exit 1; }

	tar -czvf "$astats_dest".tgz -C "$RPMBUILD"/SOURCES $(basename $astats_dest) || { echo "Could not create tar archive $astats_dest.tgz: $?"; exit 1; }

        echo "Required TS Version: ${ts_version}"
        if [[ -z "$ts_version" ]]; then
		$ts_version=""
		echo "Traffic Server not defined, skipping version constraint in RPM!!"
	fi
       	sed -e 's/__TRAFFIC_SERVER_VERSION__/'"$ts_version"'/' $ASTATS_DIR/astats_over_http.spec  >  $ASTATS_DIR/astats_over_http.spec.build
         
	mv "$ASTATS_DIR"/astats_over_http.spec.build "$RPMBUILD"/SPECS/astats_over_http.spec || { echo "Could not move temp spec file: $?"; exit 1; }

	echo "The build area has been initialized."
}


importFunctions
# Consider checking for compilers here (i.e gcc/clang)
checkEnvironment 
installTsDeps
initBuildArea
buildRpm astats_over_http
