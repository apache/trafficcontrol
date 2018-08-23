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

IP=$(ip addr | grep 'global' | grep -v "inet6" | head -n1 | awk '{print $2}' | cut -f1 -d'/')
NETMASK=$(ifconfig eth0 | grep 'inet ' | tr -s ' ' | cut -d ' ' -f4)
GATEWAY="$(ip route | grep default | awk '{print $3}')"
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

sed -ie "s;MY_HOSTNAME;$(hostname -s);g" /server.json
sed -ie "s;MY_DOMAINNAME;$(dnsdomainname);g" /server.json
sed -ie "s;MY_GATEWAY;$GATEWAY;g" /server.json
sed -ie "s;MY_NETMASK;$NETMASK;g" /server.json
sed -ie "s;MY_IP;$IP;g" /server.json

source /to-access.sh

while ! to-ping 2>/dev/null; do
	echo "waiting for traffic_ops..."
	sleep 3
done

export TO_USER=$TO_ADMIN_USER
export TO_PASSWORD=$TO_ADMIN_PASSWORD

# There's a race condition with setting the TM credentials and TO actually creating
# the TM user
until to-get api/1.3/users?username="$TM_USER" | jq -c -e '.response[].username|length'; do
	echo "waiting for TM_USER creation..."
	sleep 3
done

# now that TM_USER is available,  use that for all further operations
export TO_USER="$TM_USER"
export TO_PASSWORD="$TM_PASSWORD"

# Gets our CDN ID
CDN=$(to-get api/1.3/cdns | jq '.response|.[]|select(.name=="CDN-in-a-Box")|.id')
while [[ -z "$CDN" ]]; do
	echo "waiting for traffic_ops setup to complete..."
	sleep 3
	CDN=$(to-get api/1.3/cdns | jq '.response|.[]|select(.name=="CDN-in-a-Box")|.id')
done

# Now we upload a profile for later use
sed -ie "s;CDN_ID;$CDN;g" /profile.json
cat /profile.json
PROFILE=$(to-post api/1.3/profiles /profile.json | jq '.response')
PROFILENAME=$(echo $PROFILE | jq '.name' | tr -d '"')
PROFILEID=$(echo $PROFILE | jq '.id')
to-post api/1.3/profiles/name/$PROFILENAME/parameters /parameters.json
echo

# Gets the location ID
location=$(to-get api/1.3/phys_locations | jq '.response|.[]|select(.name=="CDN_in_a_Box")|.id')
while [[ -z "$location" ]]; do
	echo "Waiting for location setup"
	sleep 3
	location=$(to-get api/1.3/phys_locations | jq '.response|.[]|select(.name=="CDN_in_a_Box")|.id')
done

# Gets the id of a RASCAL server type
TYPE=$(to-get api/1.3/types | jq '.response|.[]|select(.name=="RASCAL")|.id')

# Gets the id of the 'ONLINE' status
ONLINE=$(to-get api/1.3/statuses | jq '.response|.[]|select(.name=="ONLINE")|.id')

# Gets the cachegroup ID
CACHEGROUP=$(to-get api/1.3/cachegroups | jq '.response|.[]|select(.name=="CDN_in_a_Box_Mid")|.id')
while [[ -z "$CACHEGROUP" ]]; do
	echo "waiting for trafficops setup to complete..."
	sleep 3
	CACHEGROUP=$(to-get api/1.3/cachegroups | jq '.response|.[]|select(.name=="CDN_in_a_Box_Mid")|.id')
done

# Now put it all together and send it up
sed -ie "s;MY_LOCATION;$location;g" /server.json
sed -ie "s;MY_TYPE;$TYPE;g" /server.json
sed -ie "s;MY_CDN_ID;$CDN;g" /server.json
sed -ie "s;MY_STATUS;$ONLINE;g" /server.json
sed -ie "s;CACHE_GROUP_ID;$CACHEGROUP;g" /server.json
sed -ie "s;MY_PROFILE_ID;$PROFILEID;g" /server.json
cat /server.json
to-post api/1.3/servers /server.json
echo

export TO_USER=$TO_ADMIN_USER
export TO_PASSWORD=$TO_ADMIN_PASSWORD
. /to-access.sh
to-enroll $(hostname -s)

touch /opt/traffic_monitor/var/log/traffic_monitor.log

cd /opt/traffic_monitor
/opt/traffic_monitor/bin/traffic_monitor -opsCfg /opt/traffic_monitor/conf/traffic_ops.cfg -config /opt/traffic_monitor/conf/traffic_monitor.cfg &
disown
exec tail -f /opt/traffic_monitor/var/log/traffic_monitor.log
