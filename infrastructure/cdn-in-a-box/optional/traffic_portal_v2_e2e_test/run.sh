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
set -e

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

TP_URL="https://$TP2_FQDN:$TP2_PORT"
timeout 3m bash <<TMOUT
	while ! curl -k "${TP_URL}/api/4.0/ping" >/dev/null 2>&1; do
		echo "waiting for Traffic Portal at '$TP_URL' fqdn '$TP_FQDN' host '$TP_HOST'"
		sleep 3;
	done
TMOUT

cd /lang/traffic-portal

config=cypress/fixtures/to.config.json
echo "$(<"$config" jq --arg TPURL $TP_URL --arg TOURL https://$TO_FQDN:$TO_PORT '.toURL = $TPURL')" > "$config"
sync

npm run e2e:ci
rc=$?

cp -r ./nightwatch/junit/* /junit/

exit $rc
