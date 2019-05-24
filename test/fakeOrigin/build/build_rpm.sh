#!/usr/bin/env bash

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

set -ex
env

BUILDDIR="$HOME/rpmbuild"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
pwd
cd $DIR/..
pwd

if [ -z "${VER_MAJOR+set}" ]; then
  VER_MAJOR=$(sed '1q;d' $DIR/../version/VERSION)
fi
if [ -z "${VER_MINOR+set}" ]; then
  VER_MINOR=$(sed '2q;d' $DIR/../version/VERSION)
fi
if [ -z "${VER_PATCH+set}" ]; then
  VER_PATCH=$(sed '3q;d' $DIR/../version/VERSION)
fi
if [ -z "${VER_DESC+set}" ]; then
  VER_DESC=$(sed '4q;d' $DIR/../version/VERSION)
fi
if [ -z "${VER_COMMIT+set}" ]; then
#  VER_COMMIT=$(sed '5q;d' $DIR/../version/VERSION)
  VER_COMMIT=$(git -C ${DIR}/../../.. rev-list --all --count)
fi
if [ -z "${BUILD_NUMBER+set}" ]; then
  BUILD_NUMBER=1
fi


VERSION="${VER_MAJOR}.${VER_MINOR}.${VER_PATCH}_${VER_DESC}_${VER_COMMIT}"

# prep build environment
mkdir -p $DIR/../dist
rm -rf $BUILDDIR
mkdir -p $BUILDDIR/{BUILD,RPMS,SOURCES}
echo "$BUILDDIR" > ~/.rpmmacros

# build
go build -v -ldflags "-X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerFull=${VERSION} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerMajor=${VER_MAJOR} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerMinor=${VER_MINOR} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerPatch=${VER_PATCH} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerDesc=${VER_DESC} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerCommit=${VER_COMMIT}"

# tar
tar -cvzf $BUILDDIR/SOURCES/fakeOrigin-${VERSION}-${BUILD_NUMBER}.tgz fakeOrigin build/config.json build/fakeOrigin.init build/fakeOrigin.logrotate example

# build RPM
rpmbuild --define "_topdir ${BUILDDIR}" --define "_version ${VERSION}" --define "_release ${BUILD_NUMBER}" -ba build/fakeOrigin.spec

# copy build RPM to ../dist
cp $BUILDDIR/RPMS/x86_64/*.rpm ./dist/

# Cross compile because we can
GOBINEXT=""
for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    if [[ "$GOOS" == "windows" ]]
    then
      GOBINEXT=".exe"
    else
      GOBINEXT=""
    fi
    GOOS=$GOOS GOARCH=$GOARCH go build -v -ldflags "-X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerFull=${VERSION} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerMajor=${VER_MAJOR} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerMinor=${VER_MINOR} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerPatch=${VER_PATCH} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerDesc=${VER_DESC} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerCommit=${VER_COMMIT}" -v
    zip -r $DIR/../dist/fakeOrigin-$VERSION-$GOOS-$GOARCH.zip fakeOrigin$GOBINEXT example
  done
done

# ARM Cross compile because we can
GOOS=linux
GOARCH=arm
for GOARM in 5 6 7; do
  GOOS=$GOOS GOARCH=$GOARCH GOARM=$GOARM go build -v -ldflags "-X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerFull=${VERSION} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerMajor=${VER_MAJOR} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerMinor=${VER_MINOR} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerPatch=${VER_PATCH} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerDesc=${VER_DESC} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerCommit=${VER_COMMIT}" -v
  zip -r $DIR/../dist/fakeOrigin-$VERSION-$GOOS-$GOARCH-ARM$GOARM.zip fakeOrigin example
done
GOARCH=arm64
GOOS=$GOOS GOARCH=$GOARCH go build -v -ldflags "-X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerFull=${VERSION} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerMajor=${VER_MAJOR} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerMinor=${VER_MINOR} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerPatch=${VER_PATCH} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerDesc=${VER_DESC} -X github.com/apache/trafficcontrol/test/fakeOrigin/version.VerCommit=${VER_COMMIT}" -v
zip -r $DIR/../dist/fakeOrigin-$VERSION-$GOOS-$GOARCH-ARM8.zip fakeOrigin example
