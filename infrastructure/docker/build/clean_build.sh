#!/usr/bin/env sh

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

set -x
cp -a /trafficcontrol /tmp/. && \
	cd /tmp/trafficcontrol && \
	rm -rf dist && \
	ln -fs /trafficcontrol/dist dist &&
	((((./build/build.sh $1 2>&1; echo $? >&3) | tee ./dist/build-$1.log >&4) 3>&1) | (read x; exit $x)) 4>&1
