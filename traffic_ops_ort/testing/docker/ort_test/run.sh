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

# wait period to insure all containers are up and running.
WAIT="60"

#
# this seems to wake up the to container.
#
function ping_to {
  /opt/ort/t3c \
		"--traffic-ops-insecure=true" \
		"--traffic-ops-timeout-milliseconds=3000" \
		"--traffic-ops-user=$TO_ADMIN_USER" \
		"--traffic-ops-password=$TO_ADMIN_PASS" \
		"--traffic-ops-url=$TO_URI" \
		"--cache-host-name=atlanta-edge-03" \
		"--log-location-error=stderr" \
		"--log-location-info=stderr" \
		"--log-location-debug=stderr" \
		"--run-mode=badass" 
}

set -x
GOPATH=/root/go; export GOPATH
PATH=$PATH:/usr/local/go/bin:; export PATH
TERM=xterm; export TERM

# setup some convienient links
/bin/ln -s /root/go/src/github.com/apache/trafficcontrol /trafficcontrol
/bin/ln -s /trafficcontrol/traffic_ops_ort/testing/ort-tests /ort-tests

if [ -f /trafficcontrol/GO_VERSION ]; then
  go_version=$(cat /trafficcontrol/GO_VERSION) && \
      curl -Lo go.tar.gz https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz && \
        tar -C /usr/local -xvzf go.tar.gz && \
        ln -s /usr/local/go/bin/go /usr/bin/go && \
        rm go.tar.gz
else
  echo "no GO_VERSION file, unable to install go"
  exit 0
fi

# fetch dependent packages for tests
go mod vendor -v

if [ -f /systemctl.sh ]; then
  mv /bin/systemctl /bin/systemctl.save
  cp /systemctl.sh /bin/systemctl
  chmod 0755 /bin/systemctl
fi

cd /ort-tests
echo "Sleeping for $WAIT seconds to ensure all containers have initialized" >> test.log
sleep $WAIT
echo "Running tests" >> test.log

# wake up the to_server
ping_to
sleep 2

(touch test.log && tail -f test.log)&
go test -v -failfast -cfg=conf/docker-edge-cache.conf 2>&1 >> test.log
