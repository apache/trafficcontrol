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

#----------------------------------------
function importFunctions() {
	local script=$(readlink -f "$0")
	local scriptdir=$(dirname "$script")
	export TO_DIR=$(dirname "$scriptdir")
	export TC_DIR=$(dirname "$TO_DIR")
	functions_sh="$TC_DIR/build/functions.sh"
	if [[ ! -r $functions_sh ]]; then
		echo "error: can't find $functions_sh"
		exit 1
	fi
	. "$functions_sh"
}

#----------------------------------------
function initBuildArea() {
	echo "Initializing the build area for Traffic Ops ORT"
	mkdir -p "$RPMBUILD"/{SPECS,SOURCES,RPMS,SRPMS,BUILD,BUILDROOT} || { echo "Could not create $RPMBUILD: $?"; exit 1; }

	local dest=$(createSourceDir traffic_ops_ort)
	cp -p traffic_ops_ort.pl "$dest"
	cp -p supermicro_udev_mapper.pl "$dest"
	tar -czvf "$dest".tgz -C "$RPMBUILD"/SOURCES $(basename "$dest") || \
	    { echo "Could not create tape archive $dest.tgz: $?"; exit 1; }

	echo "The build area has been initialized."
}

#----------------------------------------
importFunctions
initBuildArea
buildRpm traffic_ops_ort
