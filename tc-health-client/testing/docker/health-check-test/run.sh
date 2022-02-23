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

GOPATH=/root/go; export GOPATH
PATH=$PATH:/usr/local/go/bin:; export PATH
TERM=xterm; export TERM

# set up some convienient links
/bin/ln -s /root/go/src/github.com/apache/trafficcontrol /trafficcontrol

# install go
if [[ -f /trafficcontrol/GO_VERSION ]]; then
  go_version=$(cat /trafficcontrol/GO_VERSION) && \
  curl -Lo go.tar.gz https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz && \
  tar -C /usr/local -xvzf go.tar.gz && \
  ln -s /usr/local/go/bin/go /usr/bin/go && \
  rm go.tar.gz
else
  echo "no GO_VERSION file, unable to install go"
  exit 1
fi

# write the to-creds file
. /variables.env
echo "#!/bin/bash" > /etc/to-creds
echo "TO_USER=${TO_ADMIN_USER}" >> /etc/to-creds
echo "TO_PASS=${TO_ADMIN_PASS}" >> /etc/to-creds
echo "TO_URL=${TO_URI}" >> /etc/to-creds
chmod +x /etc/to-creds

/usr/bin/ln -s /trafficcontrol/tc-health-client/testing/tests /tests
mkdir /tests/conf
/usr/bin/cp /trafficcontrol/cache-config/testing/ort-tests/conf/docker-edge-cache.conf /tests/conf
/usr/bin/cp /trafficcontrol/cache-config/testing/ort-tests/tc-fixtures.json /tests

cd "$(realpath /tests)"
touch test.log && chmod a+rw test.log && nohup tail -f test.log&
go mod vendor -v

# setup trafficserver config files
/usr/bin/cp /strategies.yaml /opt/trafficserver/etc/trafficserver/strategies.yaml
/usr/bin/cp /parent.config /opt/trafficserver/etc/trafficserver/parent.config
# start trafficserver
systemctl start trafficserver

go test --cfg=conf/docker-edge-cache.conf 2>&1 >> test.log

rm /tests/tc-fixtures.json

exit 0
