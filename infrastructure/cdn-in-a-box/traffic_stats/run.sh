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

# Script for running the Dockerfile for Traffic Stats.
# The Dockerfile sets up a Docker image which can be used for any new container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set (ordinarily by `docker run -e` arguments):
# TO_HOST
# TO_PORT
# INFLUXDB_HOST

# Check that env vars are set

set -e
set -x
set -m

envvars=( TO_HOST TO_PORT INFLUXDB_HOST)
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
TSCONF=/opt/traffic_stats/conf/traffic_stats.cfg
to-enroll ts ALL || (while true; do echo "enroll failed."; sleep 3 ; done)

while ! to-ping 2>/dev/null; do
  echo "waiting for trafficops ($TO_URL)..."
  sleep 3
done

cat <<-EOF >$TSCONF
{
	"toUser": "$TO_ADMIN_USER",
	"toPasswd": "$TO_ADMIN_PASSWORD",
	"toUrl": "$TO_URL",
	"toRequestTimeout": 10,
	"influxUser": "$INFLUXDB_ADMIN_USER",
	"influxPassword": "$INFLUXDB_ADMIN_PASSWORD",
	"pollingInterval": 10,
	"publishingInterval": 30,
	"maxPublishSize": 10000,
	"statusToMon": "ONLINE",
	"seelogConfig": "/opt/traffic_stats/conf/traffic_stats_seelog.xml",
	"dailySummaryPollingInterval": 300,
	"cacheRetentionPolicy": "daily",
	"dsRetentionPolicy": "daily",
	"dailySummaryRetentionPolicy": "indefinite",
    "influxUrls": ["http://$INFLUXDB_HOST:$INFLUXDB_PORT"]
}
EOF

cat <<-EOF >/opt/traffic_stats/conf/traffic_stats_seelog.xml
<?xml version='1.0'?>
<seelog minlevel="debug">
    <outputs formatid="std:debug-short">
        <file path="/var/log/traffic_stats/traffic_stats.log" />
    </outputs>
</seelog>
EOF

touch /var/log/traffic_stats/traffic_stats.log

# Wait for influxdb
until nc $INFLUXDB_HOST $INFLUXDB_PORT </dev/null >/dev/null 2>&1; do
  echo "Waiting for influxdb to start..."
  sleep 3
done

/opt/traffic_stats/influxdb_tools/create_ts_databases -user $INFLUXDB_ADMIN_USER -password $INFLUXDB_ADMIN_PASSWORD -url http://$INFLUXDB_HOST:$INFLUXDB_PORT -replication 1

# Wait for traffic monitor
until nc $TM_FQDN $TM_PORT </dev/null >/dev/null 2>&1; do
  echo "Waiting for Traffic Monitor to start..."
  sleep 3
done

traffic_stats_command=(/opt/traffic_stats/bin/traffic_stats -cfg $TSCONF);
if [[ "$TS_DEBUG_ENABLE" == true ]]; then
  dlv '--continue' '--listen=:2346' '--accept-multiclient=true' '--headless=true' '--api-version=2' exec \
    "${traffic_stats_command[0]}" -- "${traffic_stats_command[@]:1}" &
else
  "${traffic_stats_command[@]}" &
fi;

exec tail -f /var/log/traffic_stats/traffic_stats.log

