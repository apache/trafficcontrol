#!/bin/bash
#
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
#

# ---------------------------------------
function initBuildArea() {
    echo "Initializing Traffic Portal build area"
    removeBuildArea
    mkdir -p $WORKSPACE/build/rpmbuild/traffic_portal
    rsync -av --exclude='build/rpmbuild/traffic_portal' $WORKSPACE/traffic_portal/ $WORKSPACE/build/rpmbuild/traffic_portal
}

# ---------------------------------------
function initTmpDir() {
    rm -rf /tmp/traffic_portal-$BRANCH
    mkdir -p /tmp/traffic_portal-$BRANCH
    mkdir -p /tmp/traffic_portal-$BRANCH/etc/init.d
    mkdir -p /tmp/traffic_portal-$BRANCH/etc/logrotate.d
    mkdir -p /tmp/traffic_portal-$BRANCH/etc/traffic_portal
    mkdir -p /tmp/traffic_portal-$BRANCH/opt/traffic_portal
    mkdir -p /tmp/traffic_portal-$BRANCH/opt/traffic_portal/public
    mkdir -p /tmp/traffic_portal-$BRANCH/opt/traffic_portal/server
}

# ---------------------------------------
function setupRelease() {
    echo "Setting up Traffic Portal release"
    # for the ant build
    touch /tmp/traffic_portal_release.properties
    echo -e "\narch=$ARCH\ntraffic_portal_version=$BRANCH\ntraffic_portal_build_number=$BUILD_NUMBER" > /tmp/traffic_portal_release.properties
}

# ---------------------------------------
function buildRpm() {
    echo "Building Traffic Portal RPM"
    echo "Installing Traffic Portal dependencies and build artifacts"
    cd $WORKSPACE/build/rpmbuild/traffic_portal
    /usr/bin/npm install
    /usr/bin/bower install
    /usr/bin/grunt dist

    # copies server.js file and config.js files to tmp
    cp $WORKSPACE/build/rpmbuild/traffic_portal/server/server.js /tmp/traffic_portal-$BRANCH/opt/traffic_portal/server
    cp -r $WORKSPACE/build/rpmbuild/traffic_portal/conf /tmp/traffic_portal-$BRANCH/etc/traffic_portal

    # copy init.d to tmp (creates the traffic portal service)
    cp -r $WORKSPACE/build/rpmbuild/traffic_portal/build/etc/init.d/traffic_portal /tmp/traffic_portal-$BRANCH/etc/init.d

    # logrotate for logs
    cp -r $WORKSPACE/build/rpmbuild/traffic_portal/build/etc/logrotate.d/traffic_portal /tmp/traffic_portal-$BRANCH/etc/logrotate.d
    cp -r $WORKSPACE/build/rpmbuild/traffic_portal/build/etc/logrotate.d/traffic_portal-access /tmp/traffic_portal-$BRANCH/etc/logrotate.d
    cp -r $WORKSPACE/build/rpmbuild/traffic_portal/app/dist/* /tmp/traffic_portal-$BRANCH/opt/traffic_portal

    # creates dynamic json file needed at runtime for traffic portal to display release info
    BUILD_DATE=$(date +'%Y-%m-%d %H:%M:%S')
    VERSION="\"Version\":\"$BRANCH\""
    BUILD_NUMBER="\"Build Number\":\"$BUILD_NUMBER\""
    BUILD_DATE="\"Build Date\":\"$BUILD_DATE\""
    JSON_VERSION="{\n$VERSION,\n$BUILD_NUMBER,\n$BUILD_DATE\n}"
    echo -e $JSON_VERSION > /tmp/traffic_portal-$BRANCH/opt/traffic_portal/public/traffic_portal_release.json

    # ant builds the rpm
    cd $WORKSPACE/build/rpmbuild/traffic_portal/build
    /usr/bin/ant -v -DTMP_DIR=/tmp/traffic_portal-$BRANCH

    # copy the rpm to the dist dir
    mkdir -p $WORKSPACE/dist
    cp $WORKSPACE/build/rpmbuild/traffic_portal/build/dist/*.rpm $WORKSPACE/dist
}

# ---------------------------------------
function removeBuildArea() {
    rm -rf $WORKSPACE/build/rpmbuild/traffic_portal
}

# ---------------------------------------
function cleanupTmpDir() {
    rm -rf /tmp/traffic_portal-$BRANCH
    rm /tmp/traffic_portal_release.properties
}

# ---------------------------------------
function getRevCount() {
    git rev-list HEAD 2>/dev/null | wc -l
}

# ---------------------------------------
# MAIN
# ---------------------------------------
if [ -z "$WORKSPACE" ]; then
	WORKSPACE=$(dirname $(dirname $(pwd)))
fi

if [ -z "$BRANCH" ]; then
    echo "'BRANCH' defaults to master if not defined"
    BRANCH=master
fi

if [ -z "$BUILD_NUMBER" ]; then
    echo "'BUILD_NUMBER' defaults to 000 if not defined"
    BUILD_NUMBER=$(getRevCount)
fi

ARCH="x86_64"

echo "=================================================="
echo "BRANCH: $BRANCH"
echo "BUILD_NUMBER: $BUILD_NUMBER"
echo "--------------------------------------------------"

initBuildArea
initTmpDir
setupRelease
buildRpm
removeBuildArea
cleanupTmpDir
