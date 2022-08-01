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

# Script for running the Dockerfile for Traffic Monitor.
# The Dockerfile sets up a Docker image which can be used for any new container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set (ordinarily by `docker run -e` arguments):
# TO_HOST
# TO_PORT
# TM_USER
# TM_PASSWORD

# Check that env vars are set

set -e
set -x
set -m

envvars=( TO_HOST TO_PORT TM_PORT TM_USER TM_PASSWORD)
for v in $envvars
do
	if [[ -z "${!v}" ]]; then echo "$v is unset"; exit 1; fi
done

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
until [[ -n "$X509_GENERATION_COMPLETE" ]]
do
  echo "Waiting on X509 vars to be defined"
  sleep 1
  source "$X509_CA_ENV_FILE"
done

# Trust the CIAB-CA at the System level
cp $X509_CA_CERT_FULL_CHAIN_FILE /etc/pki/ca-trust/source/anchors
update-ca-trust extract

# Enroll with traffic ops
TO_URL="https://$TO_FQDN:$TO_PORT"
to-enroll tm $CDN_NAME || (while true; do echo "enroll failed."; sleep 3 ; done)

# Configure Traffic Monitor
cat > /opt/traffic_monitor/conf/traffic_ops.cfg <<- ENDOFMESSAGE
{
	"username": "$TM_USER",
	"password": "$TM_PASSWORD",
	"url": "$TO_URL",
	"insecure": true,
	"cdnName": "$CDN_NAME",
	"httpListener": ":$TM_PORT"
}
ENDOFMESSAGE

while ! to-ping 2>/dev/null; do
	echo "waiting for trafficops ($TO_URL)..."
	sleep 3
done


export TO_USER=$TO_ADMIN_USER
export TO_PASSWORD=$TO_ADMIN_PASSWORD

# There's a race condition with setting the TM credentials and TO actually creating
# the TM user
until to-get "api/4.0/users?username=$TM_USER" 2>/dev/null | jq -c -e '.response[].username|length'; do
	echo "waiting for TM_USER creation..."
	sleep 3
done

# now that TM_USER is available,  use that for all further operations
export TO_USER="$TM_USER"
export TO_PASSWORD="$TM_PASSWORD"

export TO_USER=$TO_ADMIN_USER
export TO_PASSWORD=$TO_ADMIN_PASSWORD

# Do not start until there a valid Snapshot has been taken
until [ $(to-get "/api/4.0/cdns/$CDN_NAME/snapshot" 2>/dev/null | jq -c -e '.response.config|length') -gt 0 ] ; do
	echo "Waiting on valid Snapshot...";
  	sleep 3;
done

if [[ "$TM_DEBUG_ENABLE" == true ]]; then
	day_in_ms=$(( 1000 * 60 * 60 * 24 )); # Timing out debugging after 1 day seems fair
	set -o allexport;
	HTTP_TIMEOUT_MS=$day_in_ms
	SERVER_READ_TIMEOUT_MS=$day_in_ms
	SERVER_WRITE_TIMEOUT_MS=$day_in_ms
	set +o allexport;
else
	set -o allexport;
	HTTP_TIMEOUT_MS=2000
	SERVER_READ_TIMEOUT_MS=10000
	SERVER_WRITE_TIMEOUT_MS=10000
	set +o allexport;
fi;

envsubst < /opt/traffic_monitor/conf/traffic_monitor.cfg.template > /opt/traffic_monitor/conf/traffic_monitor.cfg

traffic_monitor_command=(/opt/traffic_monitor/bin/traffic_monitor -opsCfg /opt/traffic_monitor/conf/traffic_ops.cfg -config /opt/traffic_monitor/conf/traffic_monitor.cfg);
if [[ "$TM_DEBUG_ENABLE" == true ]]; then
  dlv '--continue' '--listen=:2344' '--accept-multiclient=true' '--headless=true' '--api-version=2' exec \
    "${traffic_monitor_command[0]}" -- "${traffic_monitor_command[@]:1}" &
else
  "${traffic_monitor_command[@]}" &
fi;
tail -f /dev/null; # Keeps the container running indefinitely. The container health check (see dockerfile) will report whether Traffic Monitor is running.
