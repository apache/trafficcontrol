#!/usr/bin/bash

<<<<<<< HEAD
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

=======
>>>>>>> Added edge cache
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

while ! nc $TO_HOST $TO_PORT </dev/null; do
	echo "waiting for Traffic Ops"
	sleep 3
done

# Now network things. First need to authenticate
curl -ksc cookie.jar -d "{\"u\":\"$TO_ADMIN_USER\",\"p\":\"$TO_ADMIN_PASSWORD\"}" https://$TP_HOST/api/1.3/user/login
echo

# Gets our CDN ID
CDN=$(curl -ksb cookie.jar https://$TP_HOST/api/1.3/cdns | jq '.response|.[]|select(.name=="CDN-in-a-Box")|.id')
while [[ -z "$CDN" ]]; do
	echo "waiting for trafficops setup to complete..."
	sleep 3
	CDN=$(curl -ksb cookie.jar https://$TP_HOST/api/1.3/cdns | jq '.response|.[]|select(.name=="CDN-in-a-Box")|.id')
done


# Now we upload a profile for later use
sed -ie "s;CDN_ID;$CDN;g" /profile.json
PROFILE=$(curl -ksb cookie.jar -d @/profile.json https://$TP_HOST/api/1.3/profiles | jq '.response')
PROFILENAME=$(echo $PROFILE | jq '.name' | tr -d '"')
PROFILEID=$(echo $PROFILE | jq '.id')
curl -ksb cookie.jar -d @/parameters.json https://$TP_HOST/api/1.3/profiles/name/$PROFILENAME/parameters
echo

# Gets the location ID
location=$(curl -ksb cookie.jar https://$TP_HOST/api/1.3/phys_locations | jq '.response|.[]|select(.name=="CDN_in_a_Box")|.id')
while [[ -z "$location" ]]; do
	echo "Waiting for location setup"
	sleep 3
	location=$(curl -ksb cookie.jar https://$TP_HOST/api/1.3/phys_locations | jq '.response|.[]|select(.name=="CDN_in_a_Box")|.id')
done

# Gets the id of a MID server type
TYPE=$(curl -ksb cookie.jar https://$TP_HOST/api/1.3/types | jq '.response|.[]|select(.name=="EDGE")|.id')

# Gets the id of the 'REPORTED' status
REPORTED=$(curl -ksb cookie.jar https://$TP_HOST/api/1.3/statuses | jq '.response|.[]|select(.name=="REPORTED")|.id')

# Gets the cachegroup ID
CACHEGROUP=$(curl -ksb cookie.jar https://$TP_HOST/api/1.3/cachegroups | jq '.response|.[]|select(.name=="CDN_in_a_Box_Edge")|.id')
while [[ -z "$CACHEGROUP" ]]; do
	echo "waiting for trafficops setup to complete..."
	sleep 3
	CACHEGROUP=$(curl -ksb cookie.jar https://$TP_HOST/api/1.3/cachegroups | jq '.response|.[]|select(.name=="CDN_in_a_Box_Edge")|.id')
done

# Now put it all together and send it up
sed -ie "s;MY_LOCATION;$location;g" /server.json
sed -ie "s;MY_TYPE;$TYPE;g" /server.json
sed -ie "s;MY_CDN_ID;$CDN;g" /server.json
sed -ie "s;REPORTED_ID;$REPORTED;g" /server.json
sed -ie "s;CACHE_GROUP_ID;$CACHEGROUP;g" /server.json
sed -ie "s;MY_PROFILE_ID;$PROFILEID;g" /server.json
cat /server.json
SERVER=$(curl -ksb cookie.jar -d @/server.json https://$TP_HOST/api/1.3/servers | jq '.response.id')


#finally, link this server to a delivery service
DSID=$(curl -ksb cookie.jar https://$TP_HOST/api/1.3/deliveryservices | jq '.response|.[]|select(.displayName=="CDN in a Box")|.id')
while [[ -z "$DSID" ]]; do
	echo "Waiting for delivery service creation..."
	sleep 3
	DSID=$(curl -ksb cookie.jar https://$TP_HOST/api/1.3/deliveryservices | jq '.response|.[]|select(.displayName=="CDN in a Box")|.id')
done

curl -ksb cookie.jar -d "{\"dsId\":$DSID,\"servers\":[$SERVER],\"replace\":true}" https://$TP_HOST/api/1.2/deliveryserviceserver
echo

fg
