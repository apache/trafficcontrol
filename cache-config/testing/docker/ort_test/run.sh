#!/bin/bash
#
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
#

#
# this seems to wake up the to container.
#
function ping_to {
	t3c \
		"apply" \
		"--traffic-ops-insecure=true" \
		"--traffic-ops-timeout-milliseconds=3000" \
		"--traffic-ops-user=$TO_ADMIN_USER" \
		"--traffic-ops-password=$TO_ADMIN_PASS" \
		"--traffic-ops-url=$TO_URI" \
		"--cache-host-name=atlanta-edge-03" \
		"-vv" \
		"--run-mode=badass"
}

GOPATH=/root/go; export GOPATH
PATH=$PATH:/usr/local/go/bin:; export PATH
TERM=xterm; export TERM

# setup some convienient links
/bin/ln -s /root/go/src/github.com/apache/trafficcontrol /trafficcontrol
/bin/ln -s /trafficcontrol/cache-config/testing/ort-tests /ort-tests

if [ -f /trafficcontrol/GO_VERSION ]; then
  go_version=$(cat /trafficcontrol/GO_VERSION) && \
      curl -Lo go.tar.gz https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz && \
        tar -C /usr/local -xzf go.tar.gz && \
        ln -s /usr/local/go/bin/go /usr/bin/go && \
        rm go.tar.gz
else
  echo "no GO_VERSION file, unable to install go"
  exit 1
fi

if [[ -f /systemctl.sh ]]; then
  mv /bin/systemctl /bin/systemctl.save
  cp /systemctl.sh /bin/systemctl
  chmod 0755 /bin/systemctl
fi

cd "$(realpath /ort-tests)"

# fetch dependent packages for tests
go mod vendor

cp /ort-tests/tc-fixtures.json /tc-fixtures.json
ATS_RPM=`basename /yumserver/test-rpms/trafficserver-[0-9]*.rpm |
  gawk 'match($0, /trafficserver\-(.+)\.rpm$/, arr) {print arr[1]}'`

echo "ATS_RPM: $ATS_RPM"

if [[ -z $ATS_RPM ]]; then
  echo "ERROR: No ATS RPM was found"
  exit 2
else
  echo "$(</ort-tests/tc-fixtures.json jq --arg ATS_RPM "$ATS_RPM" '.profiles[] |= (
    select(.params != null).params[] |= (
      select(.configFile == "package" and .name == "trafficserver").value = $ATS_RPM
    ))')" >/ort-tests/tc-fixtures.json.tmp
  if ! </ort-tests/tc-fixtures.json.tmp jq -r --arg ATS_RPM "$ATS_RPM" '.profiles[] |
    select(.params != null).params[] |
    select(.configFile == "package" and .name == "trafficserver")
    .value' |
      grep -qF "$ATS_RPM";
  then
    echo "ATS RPM version ${ATS_RPM} was not set"
    exit 2
  fi
fi

# wake up the to_server
ping_to

echo "waiting for all the to_server container to initialize."
i=0
sleep_time=3
while ! nc $TO_HOSTNAME $TO_PORT </dev/null; do
  echo "$waiting for $TO_HOSTNAME:$TO_PORT"
  sleep $sleep_time
  let i=i+1
  if [ $i -gt 10 ]; then
    let d=i*sleep_time
    echo "$TO_HOSTNAME:$TO_PORT is unavailable after $d seconds, giving up"
    exit 1
  fi
done

mv /ort-tests/tc-fixtures.json.tmp /tc-fixtures.json
(touch test.log && chmod a+rw test.log && tail -f test.log)&

go test --cfg=conf/docker-edge-cache.conf 2>&1 >> test.log
if [[ $? != 0 ]]; then
  echo "ERROR: ORT tests failure"
  exit 3
fi

exit 0
