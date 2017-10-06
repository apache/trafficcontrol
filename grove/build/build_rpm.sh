#!/usr/bin/env bash
BUILDDIR="$HOME/rpmbuild"

VERSION=`cat ./VERSION`.`git rev-list --all --count`

# prep build environment
rm -rf $BUILDDIR
mkdir -p $BUILDDIR/{BUILD,RPMS,SOURCES}
echo "$BUILDDIR" > ~/.rpmmacros

# build
go build -v

# tar
tar -cvzf $BUILDDIR/SOURCES/grove-${VERSION}.tgz grove conf/grove.cfg build/grove.init

# build RPM
rpmbuild --define "version ${VERSION}" -ba build/grove.spec

# copy build RPM to .
cp $BUILDDIR/RPMS/x86_64/grove-${VERSION}-1.x86_64.rpm .
