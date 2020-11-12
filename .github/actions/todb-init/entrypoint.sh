#!/bin/sh -l
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

set -e

download_go() {
  go_version="$(cat "${GITHUB_WORKSPACE}/GO_VERSION")"
  wget -O go.tar.gz "https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz"
  tar -C /usr/local -xzf go.tar.gz
  rm go.tar.gz
  export PATH="${PATH}:${GOROOT}/bin"
  go version
}
download_go
GOPATH="$(mktemp -d)"
export PATH="${GOPATH}/bin:${PATH}" GOPATH

apk add --no-cache git gcc gettext postgresql-client musl-dev
go get -v bitbucket.org/liamstask/goose/cmd/goose
mv $GOPATH/bin/goose /bin/ &&\
	apk del git gcc musl-dev

cd "traffic_ops/app/db"

mv /dbconf.yml ./

psql -d postgresql://traffic_ops:twelve@postgres:5432/traffic_ops < ./create_tables.sql >/dev/null
goose --env=test --path="$PWD" up
psql -d postgresql://traffic_ops:twelve@postgres:5432/traffic_ops < ./seeds.sql >/dev/null
psql -d postgresql://traffic_ops:twelve@postgres:5432/traffic_ops < ./patches.sql >/dev/null
