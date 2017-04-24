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

GO_BINARY=/usr/local/go/bin/go

echo "Now installing goose"
export GOPATH=/opt/traffic_ops/go
mkdir -p $GOPATH

echo "GO_BINARY: $GO_BINARY"
$GO_BINARY get bitbucket.org/liamstask/goose/cmd/goose
$GO_BINARY get github.com/lib/pq

echo "Successfully installed goose to $GOPATH/bin/goose"
