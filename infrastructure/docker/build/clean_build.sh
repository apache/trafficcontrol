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


# Fix ownership of output files
#  $1 is file or dir with correct ownership
#  remaining args are files/dirs to be fixed, recursively
setowner() {
    own=$(stat -c '%u:%g' $1)
    shift
    [[ -n $@ ]] && chown -R ${own} "$@"
}

cleanup() {
    setowner /trafficcontrol /trafficcontrol/dist
}

trap cleanup EXIT

set -x

# set owner of dist dir -- cleans up existing dist permissions...
mkdir -p /tmp/go/{src,pkg,bin}
mkdir -p /tmp/go/src/github.com/apache/
export GOPATH=/tmp/go
cp -a /trafficcontrol /tmp/go/src/github.com/apache/. && \
	cd /tmp/go/src/github.com/apache/trafficcontrol && \
	rm -rf dist && \
	mkdir -p /trafficcontrol/dist && \
	ln -s /trafficcontrol/dist dist && \
	((((./build/build.sh $1 2>&1; echo $? >&3) | tee ./dist/build-$1.log >&4) 3>&1) | (read x; exit $x)) 4>&1
