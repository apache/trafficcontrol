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

# Script for running the Dockerfile for Traffic Ops.
# The Dockerfile sets up a Docker image which can be used for any new Traffic Ops container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
source /etc/environment

start() {
   export PATH=/usr/local/go/bin:/opt/traffic_ops/go/bin:$PATH
   export GOPATH=/go
   go get -u golang.org/x/net/publicsuffix golang.org/x/crypto/scrypt
   TODB_HOSTNAME=db TODB_NAME=traffic_ops TO_URL=https://trafficops go test -v -cfg=../conf/traffic-ops-test.conf
}


set -x

usage() {
        echo "Usage: $(basename $0) <test dir> <test env> <host> <port>"
        echo "  e.g. $(basename $0) ./t test db 5432"
}

finish() {
        local st=$?
        [[ $st -ne 0 ]] && echo "Exiting with status $st"
        [[ -n $msg ]] && echo $msg
}

start
