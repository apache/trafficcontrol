#!/bin/bash

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

BUILDDIR="$HOME/rpmbuild"

VERSION=`cat ./../VERSION`.`git rev-list --all --count`

# prep build environment
rm -rf $BUILDDIR
mkdir -p $BUILDDIR/{BUILD,RPMS,SOURCES}
echo "$BUILDDIR" > ~/.rpmmacros

# get traffic_ops client
# godir=src/github.com/apache/trafficcontrol/traffic_ops/client
# ( mkdir -p "$godir" && \
#   cd "$godir" && \
#   cp -r ${GOPATH}/${godir}/* . && \
#   go get -v \
# ) || { echo "Could not build go program at $(pwd): $!"; exit 1; }

# build
go build -v

# tar
tar -cvzf $BUILDDIR/SOURCES/grovetccfg-${VERSION}.tgz grovetccfg

# build RPM
rpmbuild --define "version ${VERSION}" -ba build/grovetccfg.spec

# copy build RPM to .
cp $BUILDDIR/RPMS/x86_64/grovetccfg-${VERSION}-1.x86_64.rpm .
