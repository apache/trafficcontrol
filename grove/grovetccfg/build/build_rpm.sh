#!/bin/bash
BUILDDIR="$HOME/rpmbuild"

VERSION=`cat ./../VERSION`.`git rev-list --all --count`

# prep build environment
rm -rf $BUILDDIR
mkdir -p $BUILDDIR/{BUILD,RPMS,SOURCES}
echo "$BUILDDIR" > ~/.rpmmacros

# get traffic_ops client
# godir=src/github.com/apache/incubator-trafficcontrol/traffic_ops/client
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
