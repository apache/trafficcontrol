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

ROOTDIR=$(git rev-parse --show-toplevel)
[ ! -z "$ROOTDIR" ] || { echo "Cannot find repository root." >&2 ; exit 1; }

cd "$ROOTDIR/grove"

BUILDDIR="$ROOTDIR/grove/rpmbuild"
VERSION=`cat ./VERSION`.`git rev-list --all --count`

# prep build environment
[ -e $BUILDDIR ] && rm -rf $BUILDDIR
[ ! -e $BUILDDIR ] || { echo "Failed to clean up rpm build directory '$BUILDDIR': $?" >&2; exit 1; }
mkdir -p $BUILDDIR/{BUILD,RPMS,SOURCES} || { echo "Failed to create build directory '$BUILDDIR': $?" >&2; exit 1; }

# build
go get -v -d . || { echo "Failed to go get dependencies: $?" >&2; exit 1; }
go build -v -ldflags "-X main.Version=$VERSION" || { echo "Failed to build grove: $?" >&2; exit 1; }

# tar
tar -cvzf $BUILDDIR/SOURCES/grove-${VERSION}.tgz grove conf/grove.cfg build/grove.init build/grove.logrotate || { echo "Failed to create archive for rpmbuild: $?" >&2; exit 1; }

# Work around bug in rpmbuild. Fixed in rpmbuild 4.13.
# See: https://github.com/rpm-software-management/rpm/commit/916d528b0bfcb33747e81a57021e01586aa82139
# Takes ownership of the spec file.
spec=build/grove.spec
spec_owner=$(stat -c%u $spec)
spec_group=$(stat -c%g $spec)
if ! id $spec_owner >/dev/null 2>&1; then
	chown $(id -u):$(id -g) build/grove.spec

	function give_spec_back {
		chown ${spec_owner}:${spec_group} build/grove.spec
	}
	trap give_spec_back EXIT
fi

# build RPM
rpmbuild --define "_topdir $BUILDDIR" --define "version ${VERSION}" -ba build/grove.spec || { echo "rpmbuild failed: $?" >&2; exit 1; }

# copy build RPM to .
[ -e ../dist ] || mkdir -p ../dist
cp $BUILDDIR/RPMS/x86_64/grove-${VERSION}-1.x86_64.rpm ../dist
