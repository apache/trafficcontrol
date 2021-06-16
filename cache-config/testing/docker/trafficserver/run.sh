#!/bin/bash
#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

function initBuildArea() {
  cd /root

  # prep build environment
  [ -e rpmbuild ] && rm -rf rpmbuild
  [ ! -e rpmbuild ] || { echo "Failed to clean up rpm build directory 'rpmbuild': $?" >&2; exit 1; }
  mkdir -p rpmbuild/{BUILD,BUILDROOT,RPMS,SPECS,SOURCES,SRPMS} || { echo "Failed to create build directory '$RPMBUILD': $?" >&2;
  exit 1; }
}

setowner() {
	own="$(stat -c%u:%g "$1")"
	shift
	chown -R "${own}" "$@"
}
trap 'exit_code=$?; setowner /trafficcontrol /trafficcontrol/dist; exit $exit_code' EXIT;

case ${ATS_VERSION:0:1} in
  8) cp /trafficserver-8.spec /trafficserver.spec
     ;;
  9) cp /trafficserver-9.spec /trafficserver.spec
     ;;
  *) echo "Unknown trafficserver version was specified"
     exit 1
     ;;
esac

echo "Building a RPM for ATS version: $ATS_VERSION"

# add the 'ats' user
id ats &>/dev/null || /usr/sbin/useradd -u 176 -r ats -s /sbin/nologin -d /

# setup the environment to use the devtoolset-9 tools.
if [[ "${RHEL_VERSION%%.*}" -le 7 ]]; then \
  source scl_source enable devtoolset-9
else
  source scl_source enable gcc-toolset-9
fi

initBuildArea

cd /root/rpmbuild/SOURCES
# clone the trafficserver repo
git clone https://github.com/apache/trafficserver.git

# build trafficserver version 9
rm -f /root/rpmbuild/RPMS/x86_64/trafficserver-*.rpm
cd trafficserver
git fetch --all
git checkout $ATS_VERSION
rpmbuild -bb /trafficserver.spec

echo "Build completed"

if [[ ! -d /trafficcontrol/dist ]]; then
  mkdir /trafficcontrol/dist
fi

case ${ATS_VERSION:0:1} in
  8) cp /root/rpmbuild/RPMS/x86_64/trafficserver-8*.rpm /trafficcontrol/dist
     ;;
  9) cp /root/rpmbuild/RPMS/x86_64/trafficserver-8*.rpm /trafficcontrol/dist
     ;;
  *) echo "Unknown trafficserver version was specified"
     exit 1
     ;;
esac 
