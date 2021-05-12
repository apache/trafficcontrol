#!/bin/bash

#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

export CGO_ENABLED=0
export GOOS=linux

if [ ! -f "traffic_monitor.rpm" ]; then
  echo "Unable to find traffic_monitor.rpm."
  echo "Run './pkg traffic_monitor_build' and copy TM rpm to this directory."
  exit 1
fi

cd tools/testto
rm testto
go build
cd - > /dev/null

cd tools/testcaches
rm testcaches
go build
cd - > /dev/null

cd tests/_integration
rm traffic_monitor_integration_test
go test -c -o traffic_monitor_integration_test
cd - > /dev/null
