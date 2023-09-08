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

trap 'echo "Error on line ${LINENO} of ${0}"; exit 1' ERR
set -o errexit -o nounset -o pipefail -o xtrace -o monitor
env > /ciab.env

mkdir /tmp/ort

set-dns.sh
insert-self-into-dns.sh

source /to-access.sh

# Wait on SSL certificate generation
until [[ -f "$X509_CA_ENV_FILE" ]]
do
	echo "Waiting on Shared SSL certificate generation"
	sleep 3
done

# Source the CIAB-CA shared SSL environment
until [[ -v X509_GENERATION_COMPLETE && -n "$X509_GENERATION_COMPLETE" ]]
do
	echo "Waiting on X509 vars to be defined"
	sleep 1
	source "$X509_CA_ENV_FILE"
done

# Trust the CIAB-CA at the System level
cp "$X509_CA_CERT_FULL_CHAIN_FILE" /etc/pki/ca-trust/source/anchors
update-ca-trust extract

while ! to-ping 2>/dev/null; do
	echo "waiting for Traffic Ops"
	sleep 5
done

export TO_USER=$TO_ADMIN_USER
export TO_PASSWORD=$TO_ADMIN_PASSWORD

# wait until the CDN has been registered
found=
while [[ -z $found ]]; do
	echo 'waiting for enroller setup'
	sleep 3
	found=$(to-get api/4.0/cdns?name="$CDN_NAME" | jq -r '.response[].name')
done

for f in /opt/init.d/*; do
	echo "$f"
	source "$f"
done

# Wait for SSL keys to exist
until [[ $(to-get "api/4.0/cdns/name/$CDN_NAME/sslkeys" | jq '.response | length') -ge 2 ]]; do
	echo 'waiting for SSL keys to exist'
	sleep 3
done
mkdir -p /tmp/trafficcontrol-cache-config
mkdir -p /opt/cache/etc/varnish

# hostname is already defined in /etc/init.d/99-run.sh
hostname="${hostname//-/_}" # replace - with _
hostname="${hostname^^}" # uppercase
debug_variable_name="T3C_DEBUG_COMPONENT_${hostname}"
debug_binary="${!debug_variable_name}"
if ! type -p "$debug_binary"; then
	t3c apply --cache=varnish --trafficserver-home=/opt/cache --run-mode=badass --traffic-ops-url="$TO_URL" --traffic-ops-user="$TO_USER" --traffic-ops-password="$TO_PASSWORD" --git=yes -vv || { echo "Failed"; }
fi
vstats -port 2000 >> /var/log/vstats.log 2>&1 &

envsubst < "/etc/cron.d/traffic_ops_ort-cron-template" > "/etc/cron.d/traffic_ops_ort-cron" && rm -f "/etc/cron.d/traffic_ops_ort-cron-template"
chmod "0644" "/etc/cron.d/traffic_ops_ort-cron" && crontab "/etc/cron.d/traffic_ops_ort-cron"

crond -im off

varnishlog
