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

source /to-access.sh

while ! to-ping 2>/dev/null; do
	echo "waiting for Traffic Ops"
	sleep 3
done

# Gets our CDN ID
CDN=$(to-get api/1.3/cdns | jq '.response|.[]|select(.name=="CDN-in-a-Box")|.id')
while [[ -z "$CDN" ]]; do
	echo "waiting for trafficops setup to complete..."
	sleep 3
	CDN=$(to-get api/1.3/cdns | jq '.response|.[]|select(.name=="CDN-in-a-Box")|.id')
done


# Now we upload a profile for later use
sed -ie "s;CDN_ID;$CDN;g" /profile.json
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

# Gets the id of a MID server type
TYPE=$(to-get api/1.3/types | jq '.response|.[]|select(.name=="MID")|.id')

# Gets the id of the 'REPORTED' status
REPORTED=$(to-get api/1.3/statuses | jq '.response|.[]|select(.name=="REPORTED")|.id')

# Gets the cachegroup ID
CACHEGROUP=$(to-get api/1.3/cachegroups | jq '.response|.[]|select(.name=="CDN_in_a_Box_Mid")|.id')
while [[ -z CACHEGROUP ]]; do
	echo "waiting for trafficops setup to complete..."
	sleep 3
	CACHEGROUP=$(to-get api/1.3/cachegroups | jq '.response|.[]|select(.name=="CDN_in_a_Box_Mid")|.id')
done

# Now put it all together and send it up
sed -ie "s;MY_LOCATION;$location;g" /server.json
sed -ie "s;MY_TYPE;$TYPE;g" /server.json
sed -ie "s;MY_CDN_ID;$CDN;g" /server.json
sed -ie "s;REPORTED_ID;$REPORTED;g" /server.json
sed -ie "s;CACHE_GROUP_ID;$CACHEGROUP;g" /server.json
sed -ie "s;MY_PROFILE_ID;$PROFILEID;g" /server.json
cat /server.json
to-post api/1.3/servers /server.json
echo

fg
