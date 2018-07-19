#!/usr/bin/bash

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
sed -ie "s;CDN_ID;$CDN;g" /profile.json
PROFILE=$(curl -ksb cookie.jar -d @/profile.json https://trafficportal/api/1.3/profiles)
PROFILENAME=$(echo $PROFILE | tr ',' '\n' | grep '"name"' | cut -d : -f2 | tr -d '"')
PROFILEID=$(echo $PROFILE | tr ',{' '\n' | grep '"id"' | cut -d : -f2)
curl -ksb cookie.jar -d @/parameters.json https://trafficportal/api/1.3/profiles/name/$PROFILENAME/parameters
echo

# Gets the location ID
location=$(curl -ksb cookie.jar https://trafficportal/api/1.3/phys_locations)
while [[ "$location" == '{"response":[]}' ]]; do
	echo "Waiting for location setup"
	sleep 3
	location=$(curl -ksb cookie.jar https://trafficportal/api/1.3/phys_locations)
done
location=$(echo $location | tr ']' '\n' | grep CDN_in_a_Box | tr ',' '\n' | grep '"id"' | cut -d ':' -f2)

# Gets the id of a MID server type
TYPE=$(curl -ksb cookie.jar https://trafficportal/api/1.3/types)
TYPE=$(echo $TYPE | tr '}' '\n' | grep '"EDGE"' | tr ',' '\n' | grep '"id"' | cut -d : -f2)

# Gets the id of the 'REPORTED' status
REPORTED=$(curl -ksb cookie.jar https://trafficportal/api/1.3/statuses)
REPORTED=$(echo $REPORTED | tr '}' '\n' | grep REPORTED | tr ',' '\n' | grep '"id"' | cut -d : -f2)

# Gets the cachegroup ID
CACHEGROUP=$(curl -ksb cookie.jar https://trafficportal/api/1.3/cachegroups)
while [[ "$CACHEGROUP" == '{"response":[]}' ]]; do
	echo "waiting for trafficops setup to complete..."
	sleep 3
	CACHEGROUP=$(curl -ksb cookie.jar https://trafficportal/api/1.3/cachegroups)
done
CACHEGROUP=$(echo $CACHEGROUP | tr '{' '\n' | grep CDN_in_a_Box_Edge | tr ',' '\n' | grep '"id"' | cut -d : -f2)

# Now put it all together and send it up
sed -ie "s;MY_LOCATION;$location;g" /server.json
sed -ie "s;MY_TYPE;$TYPE;g" /server.json
sed -ie "s;MY_CDN_ID;$CDN;g" /server.json
sed -ie "s;REPORTED_ID;$REPORTED;g" /server.json
sed -ie "s;CACHE_GROUP_ID;$CACHEGROUP;g" /server.json
sed -ie "s;MY_PROFILE_ID;$PROFILEID;g" /server.json
cat /server.json
SERVER=$(curl -ksb cookie.jar -d @/server.json https://trafficportal/api/1.3/servers)
SERVER=$(echo $SERVER | tr ',' '\n' | grep '"id"' | cut -d : -f2)


#finally, link this server to a delivery service
DSID=$(curl -ksb cookie.jar https://trafficportal/api/1.3/deliveryservices)
while [[ "$DSID" == '{"response":[]}' ]]; do
	echo "Waiting for delivery service creation..."
	sleep 3
	DSID=$(curl -ksb cookie.jar https://trafficportal/api/1.3/deliveryservices)
done
DSID=$(echo $DSID | tr '{' '\n' | grep '"CDN in a Box"' | tr ',' '\n' | grep '"id"' | cut -d : -f2)

curl -ksb cookie.jar -d "{\"dsId\":$DSID,\"servers\":[$SERVER],\"replace\":true}" https://trafficportal/api/1.2/deliveryserviceserver
echo

fg
