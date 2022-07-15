#!/bin/bash
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

# Script for running the Dockerfile for Traffic Vault.
# The Dockerfile sets up a Docker image which can be used for any new container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set (ordinarily by `docker run -e` arguments):
# ADMIN_PASS
# USER_PASS
# CERT_COUNTRY
# CERT_STATE
# CERT_CITY
# CERT_COMPANY
# TRAFFIC_OPS_URI
# TRAFFIC_OPS_USER
# TRAFFIC_OPS_PASS
# DOMAIN
# IP
# GATEWAY
# CREATE_TO_DB_ENTRY (If set to yes, create the TO db entry for this server if set to no, assume it it already there)

start() {
	/etc/init.d/riak restart
	exec tail -f /var/log/riak/console.log
}

init() {
	TMP_TO_COOKIE="$(curl -v -s -k -X POST --data '{ "u":"'"$TRAFFIC_OPS_USER"'", "p":"'"$TRAFFIC_OPS_PASS"'" }' $TRAFFIC_OPS_URI/api/4.0/user/login 2>&1 | grep 'Set-Cookie' | sed -e 's/.*mojolicious=\(.*\); expires.*/\1/')"
	echo "Got Cookie: $TMP_TO_COOKIE"
#	curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" -F "filename=profile.traffic_vault.traffic_ops" -F "profile_to_import=@/profile.traffic_vault.traffic_ops" $TRAFFIC_OPS_URI/profile/doImport	

	# \todo figure out a better way to get the IP, domain (Docker network name), gateway
	TMP_IP=$IP
	TMP_DOMAIN=$DOMAIN
	TMP_GATEWAY=$GATEWAY

	echo "Got IP: $TMP_IP"
	echo "Got Domain: $TMP_DOMAIN"
	echo "Got Gateway: $TMP_GATEWAY"

	if [ "$CREATE_TO_DB_ENTRY" = "YES" ] ; then
		TMP_CACHEGROUP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="mid-east"]; print match[0]')"
		echo "Got cachegroup ID: $TMP_CACHEGROUP_ID"
	
		TMP_SERVER_TYPE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="RIAK"]; print match[0]')"
		echo "Got server type ID: $TMP_SERVER_TYPE_ID"
	
		TMP_SERVER_PROFILE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/profiles.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="RIAK_ALL"]; print match[0]')"
		echo "Got server profile ID: $TMP_SERVER_PROFILE_ID"
	
		TMP_PHYS_LOCATION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/phys_locations.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="plocation-nyc-1"]; print match[0]')"
		echo "Got phys location ID: $TMP_PHYS_LOCATION_ID"
	
		TMP_CDN_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/cdns.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="cdn"]; print match[0]')"
		echo "Got cdn ID: $TMP_CDN_ID"
	
		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "host_name=$HOSTNAME" --data-urlencode "domain_name=$TMP_DOMAIN" --data-urlencode "interface_name=eth0" --data-urlencode "ip_address=$TMP_IP" --data-urlencode "ip_netmask=255.255.0.0" --data-urlencode "ip_gateway=$TMP_GATEWAY" --data-urlencode "interface_mtu=9000" --data-urlencode "cdn=$TMP_CDN_ID" --data-urlencode "cachegroup=$TMP_CACHEGROUP_ID" --data-urlencode "phys_location=$TMP_PHYS_LOCATION_ID" --data-urlencode "type=$TMP_SERVER_TYPE_ID" --data-urlencode "profile=$TMP_SERVER_PROFILE_ID" --data-urlencode "tcp_port=8088" $TRAFFIC_OPS_URI/server/create

	fi

	TMP_SERVER_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/servers.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["hostName"]=="'"$HOSTNAME"'"]; print match[0]')"
	echo "Got server ID: $TMP_SERVER_ID"

	curl -v -k -H "Content-Type: application/x-www-form-urlencoded" -H "Cookie: mojolicious=$TMP_TO_COOKIE" -X POST --data-urlencode "id=$TMP_SERVER_ID" --data-urlencode "status=ONLINE" $TRAFFIC_OPS_URI/server/updatestatus

	openssl req -newkey rsa:2048 -nodes -keyout /etc/riak/certs/server.key -x509 -days 365 -out /etc/riak/certs/server.crt -subj "/C=$CERT_COUNTRY/ST=$CERT_STATE/L=$CERT_CITY/O=$CERT_COMPANY"
	cp /etc/riak/certs/server.crt /etc/riak/certs/ca-bundle.crt

	/etc/init.d/riak restart
	riak-admin security enable
	riak-admin security add-group admins
	riak-admin security add-group keysusers
	riak-admin security add-user admin password=$ADMIN_PASS groups=admins
	riak-admin security add-user riakuser password=$USER_PASS groups=keysusers
	riak-admin security add-source riakuser 0.0.0.0/0 password
	riak-admin security add-source admin 0.0.0.0/0 password
	riak-admin security grant riak_kv.list_buckets,riak_kv.list_keys,riak_kv.get,riak_kv.put,riak_kv.delete on any to admins
	riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default ssl to keysusers
	riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default dnssec to keysusers
	riak-admin security grant riak_kv.get,riak_kv.put,riak_kv.delete on default url_sig_keys to keysusers

	echo "INITIALIZED=1" >> /etc/environment
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
