#!/usr/bin/bash

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
set -x
set -m

traffic_server start &

INFO=$(ifconfig eth0 | grep 'inet ' | tr -s ' ')
myIP=$(echo $INFO | cut -d' ' -f2)
device=eth0
gateway=$(route -n | grep $device | grep -E '^0\.0\.0\.0' | tr -s ' ' | cut -d ' ' -f2)
mtu=$(ip addr show | grep $device | head -n 1 | cut -d ' ' -f5)
netmask=$(echo $INFO | cut -d ' ' -f4)

sed -ie "s;MY_HOSTNAME;$(hostname -s);g" /server.json
sed -ie "s;MY_DOMAINNAME;$(dnsdomainname);g" /server.json
sed -ie "s;MY_IFACE_NAME;$device;g" /server.json
sed -ie "s;MY_MTU;$mtu;g" /server.json
sed -ie "s;MY_GATEWAY;$gateway;g" /server.json
sed -ie "s;MY_NETMASK;$netmask;g" /server.json
sed -ie "s;MY_IP;$myIP;g" /server.json

while ! nc trafficops 6443 </dev/null; do
	echo "waiting for traffic_ops_golang:6443"
	sleep 3
done

# Now network things. First need to authenticate
curl -ksc cookie.jar -d '{"u":"admin","p":"twelve"}' https://trafficportal/api/1.3/user/login
echo

# Gets our CDN ID
CDN=$(curl -ksb cookie.jar https://trafficportal/api/1.3/cdns)
CDN=$(echo $CDN | tr '}' '\n' | grep CDN-in-a-Box | tr ',' '\n' | grep '"id"' | cut -d : -f2)
while [[ -z "$CDN" ]]; do
	echo "waiting for trafficops setup to complete..."
	sleep 3
	CDN=$(curl -ksb cookie.jar https://trafficportal/api/1.3/cdns)
	CDN=$(echo $CDN | tr '}' '\n' | grep CDN-in-a-Box | tr ',' '\n' | grep '"id"' | cut -d : -f2)
done


# Now we upload a profile for later use
sed -ie "s;CDN_ID;$CDN;g" /mid_profile.json
PROFILE=$(curl -ksb cookie.jar -d @/mid_profile.json https://trafficportal/api/1.3/profiles)
PROFILENAME=$(echo $PROFILE | tr ',' '\n' | grep '"name"' | cut -d : -f2 | tr -d '"')
PROFILEID=$(echo $PROFILE | tr ',{' '\n' | grep '"id"' | cut -d : -f2)
curl -ksb cookie.jar -d @/mid_parameters.json https://trafficportal/api/1.3/profiles/name/$PROFILENAME/parameters
echo

# Gets the location ID
location=$(curl -ksb cookie.jar https://trafficportal/api/1.3/phys_locations)
location=$(echo $location | tr ']' '\n' | grep CDN_in_a_Box | tr ',' '\n' | grep '"id"' | cut -d ':' -f2)

# Gets the id of a MID server type
TYPE=$(curl -ksb cookie.jar https://trafficportal/api/1.3/types)
TYPE=$(echo $TYPE | tr '}' '\n' | grep '"MID"' | tr ',' '\n' | grep '"id"' | cut -d : -f2)

# Gets the id of the 'REPORTED' status
REPORTED=$(curl -ksb cookie.jar https://trafficportal/api/1.3/statuses)
REPORTED=$(echo $REPORTED | tr '}' '\n' | grep REPORTED | tr ',' '\n' | grep '"id"' | cut -d : -f2)

# Gets the cachegroup ID
CACHEGROUP=$(curl -ksb cookie.jar https://trafficportal/api/1.3/cachegroups)
while [[ CACHEGROUP == '{"response":[]}' ]]; do
	echo "waiting for trafficops setup to complete..."
	sleep 3
	CACHEGROUP=$(curl -ksb cookie.jar https://trafficportal/api/1.3/cachegroups)
done
CACHEGROUP=$(echo $CACHEGROUP | tr '{' '\n' | grep CDN_in_a_Box | tr ',' '\n' | grep '"id"' | cut -d : -f2)

# Now put it all together and send it up
sed -ie "s;MY_LOCATION;$location;g" /server.json
sed -ie "s;MY_TYPE;$TYPE;g" /server.json
sed -ie "s;MY_CDN_ID;$CDN;g" /server.json
sed -ie "s;REPORTED_ID;$REPORTED;g" /server.json
sed -ie "s;CACHE_GROUP_ID;$CACHEGROUP;g" /server.json
sed -ie "s;MY_PROFILE_ID;$PROFILEID;g" /server.json
cat /server.json
curl -ksb cookie.jar -d @/server.json https://trafficportal/api/1.3/servers
echo

fg
