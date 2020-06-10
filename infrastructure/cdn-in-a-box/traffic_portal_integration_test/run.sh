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
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

source to-access.sh

set-dns.sh
insert-self-into-dns.sh

TO_URL="https://$TO_FQDN:$TO_PORT"
while ! to-ping 2>/dev/null; do
   echo "waiting for trafficops at '$TO_URL' fqdn '$TO_FQDN' host '$TO_HOST'"
   sleep 3
done

nohup webdriver-manager start &

selenium_port=4444
selenium_fqdn="http://localhost:${selenium_port}"
while ! curl -Lvsk "${selenium_fqdn}" 2>/dev/null >/dev/null; do
   echo "waiting for selenium server to start on '${selenium_fqdn}'"
   sleep 1
done

cat conf.js

protractor conf.js --params.adminUser "$TO_ADMIN_USER" --params.adminPassword "$TO_ADMIN_PASSWORD"
rc=$?

cp /portaltestresults/* /junit/

exit $rc
