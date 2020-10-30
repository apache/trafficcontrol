#!/bin/bash
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
	. build/functions.sh
	if verify_and_set_go_version; then
		return
	fi
	go_version="$(cat "${GITHUB_WORKSPACE}/GO_VERSION")"
	wget -O go.tar.gz "https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz"
	echo "Extracting Go ${go_version}..."
	<<-'SUDO_COMMANDS' sudo sh
		set -o errexit
		go_dir="$(
			dirname "$(
				dirname "$(
					realpath "$(
						which go
						)")")")"
		mv "$go_dir" "${go_dir}.unused"
		tar -C /usr/local -xzf go.tar.gz
	SUDO_COMMANDS
	rm go.tar.gz
	export PATH="${PATH}:${GOROOT}/bin"
	go version
}

GOROOT=/usr/local/go
export GOROOT PATH="${PATH}:${GOROOT}/bin"
download_go

if [ -z "$INPUT_DIR" ]; then
	# There's a bug in "defaults" for inputs
	INPUT_DIR="./lib/..."
fi

GOROOT=/usr/local/go
export GOROOT PATH="${PATH}:${GOROOT}/bin"
download_go
export GOPATH="${HOME}/go"
org_dir="$GOPATH/src/github.com/apache"
repo_dir="${org_dir}/trafficcontrol"
if [[ ! -e "$repo_dir" ]]; then
	mkdir -p "$org_dir"
	cd
	mv "${GITHUB_WORKSPACE}" "${repo_dir}/"
	ln -s "$repo_dir" "${GITHUB_WORKSPACE}"
fi
cd "$repo_dir"

# Need to fetch golang.org/x/* dependencies
go mod vendor -v
go test -v $INPUT_DIR
exit $?
