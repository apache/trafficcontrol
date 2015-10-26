#!/bin/bash

#
# Copyright 2015 Comcast Cable Communications Management, LLC
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
function getTrafficControlDir() {
	local script=$(readlink -f "$0")
	local scriptdir=$(dirname "$script")
	export TS_DIR=$(dirname "$scriptdir")
	export TC_DIR=$(dirname "$TS_DIR")

	functions_sh="$TC_DIR/build/functions.sh"
	if [[ ! -r $functions_sh ]]; then
		echo "Error: Can't find $functions_sh"
		exit 1
	fi
	. "$functions_sh"
}

#----------------------------------------
function buildRpm () {
	echo "Building the rpm."

	cd "$RPMBUILD" && \
		rpmbuild --define "_topdir $(pwd)" \
			 --define "traffic_control_version $TC_VERSION" \
			 --define "build_number $BUILD_NUMBER" -ba "SPECS/$PACKAGE.spec"

	if [[ $? -ne  0 ]]; then
		echo -e "\nRPM BUILD FAILED.\n\n"
		exit 1
	fi
	echo
	echo "========================================================================================"
	echo "RPM BUILD SUCCEEDED, See $DIST/$RPM for the newly built rpm."
	echo "========================================================================================"
	echo

	mkdir -p "$DIST" || { echo "Could not create $DIST: $!"; exit 1; }

	/bin/cp "$RPMBUILD"/RPMS/*/*.rpm "$DIST/." || { echo "Could not copy rpm to $DIST: $!"; exit 1; }
	/bin/cp "$RPMBUILD"/SRPMS/*/*.rpm "$DIST/." || { echo "Could not copy source rpm to $DIST: $!"; exit 1; }
}

#----------------------------------------
function checkEnvironment() {
	echo "Verifying the build configuration environment."

	# 
	# get traffic_control src path -- relative to build_rpm.sh script
	export PACKAGE="traffic_ops_ort"
	export TC_VERSION=$(getVersion "$TC_DIR")
	export BUILD_NUMBER=${BUILD_NUMBER:-$(getBuildNumber)}
	export WORKSPACE=${WORKSPACE:-$TC_DIR}
	export RPMBUILD="$WORKSPACE/rpmbuild"
	export DIST="$WORKSPACE/dist"
	export RPM="${PACKAGE}-${TC_VERSION}-${BUILD_NUMBER}.x86_64.rpm"
	export IN_GIT=$(isInGitTree)

	echo "Build environment has been verified."

	echo "=================================================="
	echo "WORKSPACE: $WORKSPACE"
	echo "BUILD_NUMBER: $BUILD_NUMBER"
	echo "TC_VERSION: $TC_VERSION"
	echo "RPM: $RPM"
	echo "--------------------------------------------------"
}

# ---------------------------------------
function initBuildArea() {
	echo "Initializing the build area."
	mkdir -p "$RPMBUILD"/{SPECS,SOURCES,RPMS,SRPMS,BUILD,BUILDROOT} || { echo "Could not create $RPMBUILD: $!"; exit 1; }

	/bin/cp -r "$TS_DIR"/build/*.spec "$RPMBUILD"/SPECS/. || { echo "Could not copy spec files: $!"; exit 1; }

	# build the go scripts for database initialization and tm testing.

	# tar/gzip the source
	local target="$PACKAGE-$TC_VERSION"
	local targetpath="$RPMBUILD/SOURCES/$target"
	mkdir -p "$targetpath"
	/bin/cp -p "$TS_DIR"/bin/*.pl "$targetpath"/. || { echo "Could not copy $target files: $!"; exit 1; }


	tar -czvf "$targetpath.tgz" -C "$RPMBUILD/SOURCES" "$target" || { echo "Could not create tar archive $targetpath.tgz: $!"; exit 1; }

	echo "The build area has been initialized."
}

# ---------------------------------------

checkEnvironment
initBuildArea
buildRpm

GIT_SHORT_REVISION=`git rev-parse --short HEAD`; export GIT_SHORT_REVISION
export WORKSPACE=/vol1/jenkins/jobs
GOPATH=$HOME/go; export GOPATH

#. ~/.bash_profile

##build traffic_ops client
/usr/bin/rsync -av --delete $WORKSPACE/traffic_stats/workspace/traffic_ops/client/ $GOPATH/src/github.com/comcast/traffic_control/traffic_ops/client/
cd $GOPATH/src/github.com/comcast/traffic_control/traffic_ops/client/
/usr/local/bin/go install

##build influxdb client
#cd $GOPATH/src/github.com/influxdb/influxdb/client/
#git pull
#/usr/local/bin/go install


#traffic_stats
/usr/bin/rsync -av --delete $WORKSPACE/traffic_stats/workspace/traffic_stats/ $GOPATH/src/github.com/comcast/traffic_control/traffic_stats/
cd $GOPATH/src/github.com/comcast/traffic_control/traffic_stats
/usr/local/bin/go get

sed -i -e "s/@VERSION@/$VERSION/g" traffic_stats.spec
sed -i -e "s/@RELEASE@/$GIT_SHORT_REVISION/g" traffic_stats.spec
#go get all
/usr/local/bin/go build traffic_stats.go


rm -rf $WORKSPACE/traffic_stats/SOURCES/traffic_stats*
rm -rf $WORKSPACE/traffic_stats/RPMS/x86_64/traffic_stats*

mkdir -p ${targetpath}/opt/traffic_stats
mkdir -p ${targetpath}/opt/traffic_stats/bin
mkdir -p ${targetpath}/opt/traffic_stats/conf
mkdir -p ${targetpath}/opt/traffic_stats/var/log
mkdir -p ${targetpath}/etc/init.d
mkdir -p ${targetpath}/etc/logrotate.d

cp $GOPATH/src/github.com/comcast/traffic_control/traffic_stats/traffic_stats ${targetpath}/opt/traffic_stats/bin
cp $WORKSPACE/traffic_stats/workspace/traffic_stats/traffic_stats.cfg ${targetpath}/opt/traffic_stats/conf/traffic_stats.cfg
cp $WORKSPACE/traffic_stats/workspace/traffic_stats/traffic_stats_seelog.xml ${targetpath}/opt/traffic_stats/conf
cp $WORKSPACE/traffic_stats/workspace/traffic_stats/traffic_stats.init ${targetpath}/etc/init.d/traffic_stats
cp $WORKSPACE/traffic_stats/workspace/traffic_stats/traffic_stats.logrotate ${targetpath}/etc/logrotate.d/traffic_stats

cd "$RPMBUILD/SOURCES"
tar -zcvf ${target}.tar.gz ${target}
cd "$RPMBUILD"

rpmbuild -b traffic_stats.spec

