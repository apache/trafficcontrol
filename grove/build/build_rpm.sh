#!/bin/bash
BUILDDIR="$HOME/rpmbuild"

# prep build environment
rm -rf $BUILDDIR
mkdir -p $BUILDDIR/{BUILD,RPMS,SOURCES}
echo "$BUILDDIR" > ~/.rpmmacros

# build
go build -v

# tar
tar -cvzf $BUILDDIR/SOURCES/grove-0.1.tgz grove conf/grove.cfg build/grove.init

# build RPM
rpmbuild -ba build/grove.spec

# copy build RPM to .
cp $BUILDDIR/RPMS/x86_64/grove-0.1-1.x86_64.rpm .
