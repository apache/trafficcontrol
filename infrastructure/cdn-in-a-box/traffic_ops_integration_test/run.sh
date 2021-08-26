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

# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS)
set -ex
for v in $envvars
do
	if [[ -z "${!v}" ]]; then echo "$v is unset"; exit 1; fi
done

source to-access.sh
check-skips

set-dns.sh
insert-self-into-dns.sh

TO_URL="https://$TO_FQDN:$TO_PORT"
while ! to-ping 2>/dev/null; do
	echo waiting for trafficops
	sleep 3
done

# if [[ -r ]]; then; source config.sh; done    traffic_ops/run.sh
# if [[ -x ]]; then;./config.sh; done          traffic_ops/run-go.sh
source config.sh

exit_code=0
for api_version in v{2..5}; do
	./traffic_ops_${api_version}_integration_test -test.v -cfg=traffic-ops-test.conf -fixtures=tc-fixtures-${api_version}.json 2>&1 | ./go-junit-report --package-name=golang.test.toapi.${api_version} --set-exit-code > /junit/golang.test.toapi.${api_version}.xml && find /junit -type 'f' | xargs chmod 664
	declare ${api_version}_exit_code=$?
done

cat /junit/golang.test.toapi.v{2..5}.xml


if [[ $v2_exit_code -eq 0 && $v3_exit_code -eq 0 && $v4_exit_code -eq 0 && $v5_exit_code -eq 0 ]]; then
	echo "TO API tests success"
else
	echo "TO API tests failed"
	exit 1
fi
