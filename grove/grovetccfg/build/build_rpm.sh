#!/bin/bash
BUILDDIR="$HOME/rpmbuild"

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
tar -cvzf $BUILDDIR/SOURCES/grovetccfg-0.1.tgz grovetccfg

# build RPM
rpmbuild -ba build/grovetccfg.spec

# copy build RPM to .
cp $BUILDDIR/RPMS/x86_64/grovetccfg-0.1-1.x86_64.rpm .
