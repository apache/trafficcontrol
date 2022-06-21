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

# Script for running the Dockerfile for Traffic Server Edge caches.
# The Dockerfile sets up a Docker image which can be used for any new container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set (ordinarily by `docker run -e` arguments):
# TRAFFIC_OPS_URI
# TRAFFIC_OPS_USER
# TRAFFIC_OPS_PASS
# IP
# DOMAIN
# GATEWAY
# KAFKA_URI

start() {
	chown ats:ats /dev/ram0
	chown ats:ats /dev/ram1

	/opt/trafficserver/bin/trafficserver start
	service hekad start
	exec tail -f /opt/trafficserver/var/log/trafficserver/traffic_server.stderr
}

init() {
	TMP_TO_COOKIE="$(curl -v -s -k -X POST --data '{ "u":"'"$TRAFFIC_OPS_USER"'", "p":"'"$TRAFFIC_OPS_PASS"'" }' $TRAFFIC_OPS_URI/api/4.0/user/login 2>&1 | grep 'Set-Cookie' | sed -e 's/.*mojolicious=\(.*\); expires.*/\1/')"
	echo "Got Cookie: $TMP_TO_COOKIE"

	# \todo figure out a better way to get the IP, domain (Docker network name), gateway
  TMP_IP=$IP
	TMP_DOMAIN=$DOMAIN
	TMP_GATEWAY=$GATEWAY

	TMP_CACHEGROUP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="edge-east"]; print match[0]')"		
	echo "Got cachegroup ID: $TMP_CACHEGROUP_ID"

	TMP_SERVER_TYPE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="EDGE"]; print match[0]')"
	echo "Got server type ID: $TMP_SERVER_TYPE_ID"

	TMP_SERVER_PROFILE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/profiles.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="EDGE1_532"]; print match[0]')"
	echo "Got server profile ID: $TMP_SERVER_PROFILE_ID"

	TMP_PHYS_LOCATION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/phys_locations.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="plocation-nyc-1"]; print match[0]')"
	echo "Got phys location ID: $TMP_PHYS_LOCATION_ID"

	TMP_CDN_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/cdns.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="cdn"]; print match[0]')"
	echo "Got cdn ID: $TMP_CDN_ID"

	curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "host_name=$HOSTNAME" --data-urlencode "domain_name=$TMP_DOMAIN" --data-urlencode "interface_name=eth0" --data-urlencode "ip_address=$TMP_IP" --data-urlencode "ip_netmask=255.255.0.0" --data-urlencode "ip_gateway=$TMP_GATEWAY" --data-urlencode "interface_mtu=9000" --data-urlencode "cdn=$TMP_CDN_ID" --data-urlencode "cachegroup=$TMP_CACHEGROUP_ID" --data-urlencode "phys_location=$TMP_PHYS_LOCATION_ID" --data-urlencode "type=$TMP_SERVER_TYPE_ID" --data-urlencode "profile=$TMP_SERVER_PROFILE_ID" --data-urlencode "tcp_port=80" $TRAFFIC_OPS_URI/server/create

	TMP_SERVER_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/servers.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["hostName"]=="'"$HOSTNAME"'"]; print match[0]')"
	echo "Got server ID: $TMP_SERVER_ID"

	curl -v -k -H "Content-Type: application/x-www-form-urlencoded" -H "Cookie: mojolicious=$TMP_TO_COOKIE" -X POST --data-urlencode "id=$TMP_SERVER_ID" --data-urlencode "status=REPORTED" $TRAFFIC_OPS_URI/server/updatestatus

	export PERL5LIB=/usr/local/lib/neto_io/lib/perl5 && \
  		/opt/ort/traffic_ops_ort.pl badass WARN $TRAFFIC_OPS_URI $TRAFFIC_OPS_USER:$TRAFFIC_OPS_PASS

	# if the container wasn't given ramdisk devices, configure ATS for directories
	if [ ! -e "/dev/ram0" ] || [ ! -e "/dev/ram1" ]; then
			sed -i -- "s/volume=/1G volume=/g" /opt/trafficserver/etc/trafficserver/storage.config

			mkdir /atscache
			chmod 777 /atscache
			chown ats:ats /atscache

			mkdir /atscache/disk0
			chmod 777 /atscache/disk0
			chown ats:ats /atscache/disk0
			ln -s /atscache/disk0 /dev/ram0

			mkdir /atscache/disk1
			chmod 777 /atscache/disk1
			chown ats:ats /atscache/disk1
			ln -s /atscache/disk1 /dev/ram1
	fi

	# \todo remove when TO/ort is changed from 1%
	sed -i -- "s/size=1%/size=50%/g" /opt/trafficserver/etc/trafficserver/volume.config

	sed -i -- "s/{{.KafkaUri}}/$KAFKA_URI/g" /etc/hekad/heka.toml

	echo "INITIALIZED=1" >> /etc/environment
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
