#!/bin/sh

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

set -e
go mod vendor -v
touch coverprofile
covertmp="$(mktemp)"
for pkg in $(go list ./lib/... ./traffic_monitor/... ./traffic_stats/... ./traffic_ops/traffic_ops_golang/... ./cache-config/... | grep -v "/vendor/\|testing/ort-tests"); do
	tmp="$(mktemp)"
	go test -v --coverprofile="$tmp" "$pkg" | tee -a result.txt
	if [ -f "$tmp" ]; then
		gocovmerge coverprofile "$tmp" > "$covertmp"
		cp -f "$covertmp" coverprofile
		rm "$tmp"
	fi
done

rm "$covertmp"
go tool cover --func=coverprofile > /junit/coverage
rm coverprofile
go-junit-report --package-name=golang.test --set-exit-code <result.txt >/junit/golang.test.xml
rm result.txt
