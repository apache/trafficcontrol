#!/usr/bin/env bash

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

VERSION=`cat ./VERSION`.`git rev-list --all --count`

# prep build environment
rm -rf $BUILDDIR
mkdir -p $BUILDDIR/{BUILD,RPMS,SOURCES}
echo "$BUILDDIR" > ~/.rpmmacros

# build
go build -v -ldflags "-X main.Version=$VERSION"

# tar
tar -cvzf $BUILDDIR/SOURCES/grove-${VERSION}.tgz grove conf/grove.cfg build/grove.init

# build RPM
rpmbuild --define "version ${VERSION}" -ba build/grove.spec

# copy build RPM to .
cp $BUILDDIR/RPMS/x86_64/grove-${VERSION}-1.x86_64.rpm .
