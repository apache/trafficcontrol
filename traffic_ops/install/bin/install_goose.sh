#!/usr/bin/env bash

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

GO_DOWNLOADS_URL=https://storage.googleapis.com/golang

GO_TARBALL_VERSION=go1.8.1.linux-amd64.tar.gz
GO_TARBALL_URL=$GO_DOWNLOADS_URL/$GO_TARBALL_VERSION

GO_TARBALL_VERSION_SHA=a579ab19d5237e263254f1eac5352efcf1d70b9dacadb6d6bb12b0911ede8994
GO_TARBALL_VERSION_SHA_FILE=$GO_TARBALL_VERSION.sha256
GO_TARBALL_VERSION_SHA_URL=$GO_DOWNLOADS_URL/$GO_TARBALL_VERSION_SHA_FILE
INSTALL_DIR=/usr/local
GOROOT=$INSTALL_DIR/go
GO_BINARY=$GOROOT/bin/go

cd /tmp
rm $GO_TARBALL_VERSION
rm $GO_TARBALL_VERSION_SHA_FILE
curl -O $GO_TARBALL_URL
curl -O $GO_TARBALL_VERSION_SHA_URL

echo $GO_TARBALL_VERSION_SHA_FILE
sha256sum -c <(cat $GO_TARBALL_VERSION_SHA_FILE; echo " ./$GO_TARBALL_VERSION")

if [[ $? ]]; then
    cd /usr/local
    echo "Extracting go tarball to $INSTALL_DIR/go"
    tar -zxf /tmp/$GO_TARBALL_VERSION
else
    echo "Checksum failed please verify $GO_TARBALL_VERSION against $GO_TARBALL_VERSION_SHA_FILE"
fi

echo "Now installing goose"
export GOPATH=/opt/traffic_ops/go
mkdir -p $GOPATH

echo "GO_BINARY: $GO_BINARY"
$GO_BINARY get bitbucket.org/liamstask/goose/cmd/goose
$GO_BINARY get github.com/lib/pq

echo "Successfully installed goose to $GOPATH/bin/goose"
