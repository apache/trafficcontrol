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

set -o errexit -o nounset
trap 'echo "Error on line ${LINENO} of ${0}"; exit 1' ERR

# download_go downloads and installs the GO version specified in GO_VERSION
download_go() {
	. build/functions.sh
	if verify_and_set_go_version; then
		return
	fi
	go_version="$(cat "${GITHUB_WORKSPACE}/GO_VERSION")"
	wget -O go.tar.gz "https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz" --no-verbose
	echo "Extracting Go ${go_version}..."
	<<-'SUDO_COMMANDS' sudo sh
		set -o errexit
    go_dir="$(command -v go | xargs realpath | xargs dirname | xargs dirname)"
		mv "$go_dir" "${go_dir}.unused"
		tar -C /usr/local -xzf go.tar.gz
	SUDO_COMMANDS
	rm go.tar.gz
	go version
}

GOROOT=/usr/local/go
export PATH="${PATH}:${GOROOT}/bin"
export GOPATH="${HOME}/go"

download_go

# update all golang.org/x dependencies in go.mod/go.sum
go get -u \
	golang.org/x/crypto \
	golang.org/x/net \
	golang.org/x/sys \
	golang.org/x/text \
	golang.org/x/xerrors

# update vendor/modules.txt
go mod vendor -v
