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

cd ../../../
./pkg traffic_monitor_build
rpm=`ls dist | grep monitor | grep -v log | grep -v src | grep "$(git rev-parse --short=8 HEAD)"`
if [ $? -ne 0 ]; then
  echo "Unable to build TM"
  exit 1;
fi

cp "dist/$rpm"  "traffic_monitor/tests/_integration/tm/traffic_monitor.rpm"
cd -
docker compose build
