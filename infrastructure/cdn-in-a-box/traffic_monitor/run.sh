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

envvars=( TO_HOST TO_PORT TM_USER TM_PASSWORD TO_ADMIN_USER TO_ADMIN_PASSWORD)
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

while ! curl -sk $TO_URL/api/1.3/ping </dev/null; do
	echo "waiting for $TO_HOST:$TO_PORT"
	sleep 3
done

RESPONSE=$(curl -sk -d '{ "u":"'"$TM_USER"'", "p":"'"$TM_PASSWORD"'" }' $TO_URL/api/1.3/user/login)
while [[ "$RESPONSE" == '{"alerts":[{"text":"Invalid username or password.","level":"error"}]}' ]]; do
	echo "Waiting for availability of $TM_USER login"
	RESPONSE=$(curl -sk -d '{ "u":"'"$TM_USER"'", "p":"'"$TM_PASSWORD"'" }' $TO_URL/api/1.3/user/login)
	sleep 3
done
curl -ksc cookie.jar -d "{\"u\":\"$TO_ADMIN_USER\",\"p\":\"${TO_ADMIN_PASSWORD}\"}" $TO_URL/api/1.3/user/login
echo "Got Cookie: $(tail -n1 cookie.jar | tr '\t' ' ')"

# Gets our CDN ID
CDN=$(curl -ksb cookie.jar $TO_URL/api/1.3/cdns)
CDN=$(echo $CDN | tr '}' '\n' | grep CDN-in-a-Box | tr ',' '\n' | grep '"id"' | cut -d : -f2)
while [[ -z "$CDN" ]]; do
	echo "waiting for trafficops setup to complete..."
	sleep 3
	CDN=$(curl -ksb cookie.jar $TO_URL/api/1.3/cdns)
	CDN=$(echo $CDN | tr '}' '\n' | grep CDN-in-a-Box | tr ',' '\n' | grep '"id"' | cut -d : -f2)
done

# Now we upload a profile for later use
sed -ie "s;CDN_ID;$CDN;g" /profile.json
cat /profile.json
PROFILE=$(curl -ksb cookie.jar -d @/profile.json $TO_URL/api/1.3/profiles)
PROFILENAME=$(echo $PROFILE | tr ',' '\n' | grep '"name"' | cut -d : -f2 | tr -d '"')
PROFILEID=$(echo $PROFILE | tr ',{' '\n' | grep '"id"' | cut -d : -f2)
curl -ksb cookie.jar -d @/parameters.json $TO_URL/api/1.3/profiles/name/$PROFILENAME/parameters
echo

# Gets the location ID
location=$(curl -ksb cookie.jar $TO_URL/api/1.3/phys_locations)
while [[ "$location" == '{"response":[]}' ]]; do
	echo "Waiting for location setup"
	sleep 3
	location=$(curl -ksb cookie.jar $TO_URL/api/1.3/phys_locations)
done
location=$(echo $location | tr ']' '\n' | grep CDN_in_a_Box | tr ',' '\n' | grep '"id"' | cut -d ':' -f2)

# Gets the id of a RASCAL server type
TYPE=$(curl -ksb cookie.jar $TO_URL/api/1.3/types)
TYPE=$(echo $TYPE | tr '}' '\n' | grep '"RASCAL"' | tr ',' '\n' | grep '"id"' | cut -d : -f2)

# Gets the id of the 'ONLINE' status
ONLINE=$(curl -ksb cookie.jar $TO_URL/api/1.3/statuses)
ONLINE=$(echo $ONLINE | tr '}' '\n' | grep ONLINE | tr ',' '\n' | grep '"id"' | cut -d : -f2)

# Gets the cachegroup ID
CACHEGROUP=$(curl -ksb cookie.jar $TO_URL/api/1.3/cachegroups)
while [[ CACHEGROUP == '{"response":[]}' ]]; do
	echo "waiting for trafficops setup to complete..."
	sleep 3
	CACHEGROUP=$(curl -ksb cookie.jar $TO_URL/api/1.3/cachegroups)
done
CACHEGROUP=$(echo $CACHEGROUP | tr '{' '\n' | grep CDN_in_a_Box_Mid | tr ',' '\n' | grep '"id"' | cut -d : -f2)

# Now put it all together and send it up
sed -ie "s;MY_LOCATION;$location;g" /server.json
sed -ie "s;MY_TYPE;$TYPE;g" /server.json
sed -ie "s;MY_CDN_ID;$CDN;g" /server.json
sed -ie "s;MY_STATUS;$ONLINE;g" /server.json
sed -ie "s;CACHE_GROUP_ID;$CACHEGROUP;g" /server.json
sed -ie "s;MY_PROFILE_ID;$PROFILEID;g" /server.json
cat /server.json
curl -ksb cookie.jar -d @/server.json $TO_URL/api/1.3/servers
echo

touch /opt/traffic_monitor/var/log/traffic_monitor.log

cd /opt/traffic_monitor
/opt/traffic_monitor/bin/traffic_monitor -opsCfg /opt/traffic_monitor/conf/traffic_ops.cfg -config /opt/traffic_monitor/conf/traffic_monitor.cfg &
disown
exec tail -f /opt/traffic_monitor/var/log/traffic_monitor.log
