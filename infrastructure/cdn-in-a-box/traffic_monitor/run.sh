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

envvars=( TO_HOST TO_PORT TM_USER TM_PASSWORD)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

TO_URL="https://$TO_HOST:$TO_PORT"
cat > /opt/traffic_monitor/conf/traffic_ops.cfg <<- ENDOFMESSAGE
{
	"username": "$TM_USER",
	"password": "$TM_PASSWORD",
	"url": "$TO_URL",
	"insecure": true,
	"cdnName": "CDN-in-a-Box",
	"httpListener": ":80"
}
ENDOFMESSAGE

source /to-access.sh

while ! to-ping 2>/dev/null; do
	echo "waiting for traffic_ops..."
	sleep 3
done

export TO_USER=$TO_ADMIN_USER
export TO_PASSWORD=$TO_ADMIN_PASSWORD

to-enroll $(hostname -s)

# There's a race condition with setting the TM credentials and TO actually creating
# the TM user
until to-get "api/1.3/users?username=$TM_USER" 2>/dev/null | jq -c -e '.response[].username|length'; do
	echo "waiting for TM_USER creation..."
	sleep 3
done

# now that TM_USER is available,  use that for all further operations
export TO_USER="$TM_USER"
export TO_PASSWORD="$TM_PASSWORD"

export TO_USER=$TO_ADMIN_USER
export TO_PASSWORD=$TO_ADMIN_PASSWORD

touch /opt/traffic_monitor/var/log/traffic_monitor.log

cd /opt/traffic_monitor
/opt/traffic_monitor/bin/traffic_monitor -opsCfg /opt/traffic_monitor/conf/traffic_ops.cfg -config /opt/traffic_monitor/conf/traffic_monitor.cfg &
disown
exec tail -f /opt/traffic_monitor/var/log/traffic_monitor.log
