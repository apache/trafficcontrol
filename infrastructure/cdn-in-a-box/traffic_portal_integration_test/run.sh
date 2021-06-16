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

envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS)
for v in $envvars
do
	if [[ -z "${!v}" ]]; then echo "$v is unset"; exit 1; fi
done

source to-access.sh
check-skips

set-dns.sh
insert-self-into-dns.sh

TO_URL="https://$TO_FQDN:$TO_PORT"
until to-ping 2>/dev/null; do
   echo "waiting for Traffic Ops at '$TO_URL' fqdn '$TO_FQDN' host '$TO_HOST'"
   sleep 3
done

until [[ -e "${ENROLLER_DIR}/initial-load-done" ]]; do
	echo 'Waiting for Traffic Ops to finish seeding the Traffic Ops data so Traffic Portal will start...'
	sleep 3;
done

timeout_in_seconds=30
start_time="$(date +%s)"
date_in_seconds="$start_time"
TP_URL="https://$TP_FQDN:$TP_PORT"
until curl -sk "${TP_URL}/api/4.0/ping" || (( time - start_time >= timeout_in_seconds )); do
	echo "waiting for Traffic Portal at '$TP_URL' fqdn '$TP_FQDN' host '$TP_HOST'"
	time="$(date +%s)"
	sleep 3;
done

if (( time - start_time >= timeout_in_seconds )); then
	echo "Warning: Traffic Portal did not start after ${timeout_in_seconds} seconds.";
fi;

chrome_version="$(chromium-browser --version | grep -o '[0-9.]\+')"
nohup webdriver-manager start --versions.chrome "$chrome_version" &

selenium_port=4444
selenium_fqdn="http://localhost:${selenium_port}"
while ! curl -Lvsk "${selenium_fqdn}" 2>/dev/null >/dev/null; do
   echo "waiting for selenium server to start on '${selenium_fqdn}'"
   sleep 1
done

echo "$(jq "$(<<JQ_FILTERS cat
	.params.baseUrl = "https://$TP_FQDN" |
	.params.apiUrl = "https://$TP_FQDN/api/4.0" |
	.params.login.username = "$TO_ADMIN_USER" |
	.params.login.password = "$TO_ADMIN_PASSWORD"
JQ_FILTERS
)" config.json)" > config.json

cat config.json

./node_modules/.bin/tsc
./node_modules/.bin/protractor ./GeneratedCode/config.js
rc=$?

cp /portaltestresults/* /junit/

exit $rc
