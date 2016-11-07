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

# Script for running the Dockerfile for Traffic Monitor.
# The Dockerfile sets up a Docker image which can be used for any new container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set (ordinarily by `docker run -e` arguments):
# TRAFFIC_OPS_URI
# TRAFFIC_OPS_USER
# TRAFFIC_OPS_PASS
#
#
# ! TAKE NOTE !
#
# If you are using Docker on Mac or Windows with Docker Machine, version 1.9.1 (the current version, as of this writing)
# DOES NOT WORK. A kernel patch in the VM causes Java normal main returns to hang. Unfortunately, we use Java.
#
# To destroy all your containers and downgrade to 1.9.0:
# $ docker-machine rm default
# $ docker-machine create -d virtualbox --virtualbox-boot2docker-url=https://github.com/boot2docker/boot2docker/releases/download/v1.9.0/boot2docker.iso default
#
# See:
# https://stackoverflow.com/questions/34266314/java-process-in-docker-container-doesnt-exit-on-end-of-main/34351251#34351251
# https://github.com/docker/docker/issues/18180

start() {
	service tomcat start
	touch /opt/traffic_monitor/var/log/traffic_monitor.log
	exec tail -f /opt/traffic_monitor/var/log/traffic_monitor.log
}

init() {
	TMP_TO_COOKIE="$(curl -v -s -k -X POST --data '{ "u":"'"$TRAFFIC_OPS_USER"'", "p":"'"$TRAFFIC_OPS_PASS"'" }' $TRAFFIC_OPS_URI/api/1.2/user/login 2>&1 | grep 'Set-Cookie' | sed -e 's/.*mojolicious=\(.*\); expires.*/\1/')"
	echo "Got Cookie: $TMP_TO_COOKIE"

  TMP_IP=$IP
	TMP_DOMAIN=$DOMAIN
	TMP_GATEWAY=$GATEWAY

	TMP_DOCKER_GATEWAY="$(route -n | grep -E "^0\.0\.0\.0[[:space:]]" | cut -f1 -d" " --complement | sed -e 's/^[ \t]*//' | cut -f1 -d" ")"
	echo "Got Docker gateway: $TMP_DOCKER_GATEWAY"

	if [ "$CREATE_TO_DB_ENTRY" = "YES" ] ; then
		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" -F "filename=Traffic_Monitor_Dockerfile_profile.traffic_ops" -F "profile_to_import=@/Traffic_Monitor_Dockerfile_profile.traffic_ops" $TRAFFIC_OPS_URI/profile/doImport

		TMP_CACHEGROUP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="mid-east"]; print match[0]')"
		echo "Got cachegroup ID: $TMP_CACHEGROUP_ID"
	
		TMP_SERVER_TYPE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="RASCAL"]; print match[0]')"
		echo "Got server type ID: $TMP_SERVER_TYPE_ID"
	
		TMP_SERVER_PROFILE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/profiles.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="RASCAL_CDN1"]; print match[0]')"
		echo "Got server profile ID: $TMP_SERVER_PROFILE_ID"
	
		TMP_PHYS_LOCATION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/phys_locations.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="plocation-nyc-1"]; print match[0]')"
		echo "Got phys location ID: $TMP_PHYS_LOCATION_ID"
	
		TMP_CDN_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/cdns.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="cdn"]; print match[0]')"
		echo "Got cdn ID: $TMP_CDN_ID"
	
		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "host_name=$HOSTNAME" --data-urlencode "domain_name=$TMP_DOMAIN" --data-urlencode "interface_name=eth0" --data-urlencode "ip_address=$TMP_IP" --data-urlencode "ip_netmask=255.255.0.0" --data-urlencode "ip_gateway=$TMP_GATEWAY" --data-urlencode "interface_mtu=9000" --data-urlencode "cdn=$TMP_CDN_ID" --data-urlencode "cachegroup=$TMP_CACHEGROUP_ID" --data-urlencode "phys_location=$TMP_PHYS_LOCATION_ID" --data-urlencode "type=$TMP_SERVER_TYPE_ID" --data-urlencode "profile=$TMP_SERVER_PROFILE_ID" --data-urlencode "tcp_port=80" $TRAFFIC_OPS_URI/server/create
	
		# \todo dynamically get edge/mid profiles and cachegroups
		TMP_EDGE_NAME="EDGE1_532"
		TMP_EDGE_CACHEGROUP="edge-east"
	
		TMP_ALLOW_IP_VAL="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/parameters/profile/${TMP_EDGE_NAME}.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["value"] for x in obj["response"] if x["name"]=="allow_ip" and x["configFile"] == "astats.config"]; print match[0] if len(match) > 0 else ""')"
		echo "Got existing allow_ip: $TMP_ALLOW_IP_VAL"
	
		TMP_ALLOW_IP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/parameters/profile/${TMP_EDGE_NAME}.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="allow_ip" and x["configFile"] == "astats.config"]; print match[0] if len(match) > 0 else ""')"
		echo "Got existing allow_ip id: $TMP_ALLOW_IP_ID"
	
		# Note we need to add the Docker gateway to the astats allow whitelist, 
		# for a CDN with the Traffic Monitor and Traffic Servers both within Docker,
		# because that gateway will appear as the request IP that Traffic Server sees.
	
		if [ ! -z "$TMP_ALLOW_IP_VAL" ] && [ ! -z "$TMP_ALLOW_IP_ID" ]; then 
				echo "Got allow_ip exists"
	
				# \todo dynamically check if IP and Docker Gateway are in the list, and don't add duplicates
				TMP_ALLOW_IP_VAL="${TMP_ALLOW_IP_VAL},$IP,$TMP_DOCKER_GATEWAY"
				echo "Got new allow_ip val: $TMP_ALLOW_IP_VAL"
				curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "parameter.name=allow_ip" --data-urlencode "parameter.config_file=astats.config" --data-urlencode "parameter.value=$TMP_ALLOW_IP_VAL" $TRAFFIC_OPS_URI/parameter/$TMP_ALLOW_IP_ID/update
		else
				echo "Got allow_ip doesn't exist"
	
				TMP_ALLOW_IP_VAL="$IP,$TMP_DOCKER_GATEWAY"
	
				TMP_RESPONSE="$(curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "parameter.name=allow_ip" --data-urlencode "parameter.config_file=astats.config" --data-urlencode "parameter.value=$TMP_ALLOW_IP_VAL" -D - $TRAFFIC_OPS_URI/parameter/create)"
				echo "Got parameter create response: $TMP_RESPONSE"
				TMP_LOCATION="$(printf "$TMP_RESPONSE" | grep "Location: /parameter/")"
				echo "DEBUG0 Got parameter location: $TMP_LOCATION"
				TMP_PARAMETER_ID="$(printf "$TMP_LOCATION" | cut -c 22-)"
				echo "Got parameter ID: $TMP_PARAMETER_ID"
	
				# \todo fix to dynamically get edge profile name
				TMP_PARAMETER_PROFILE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/profiles.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="'"$TMP_EDGE_NAME"'"]; print match[0]')"
				echo "Got parameter profile ID: $TMP_PARAMETER_PROFILE_ID"
	
				TMP_EDGE_CACHEGROUP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="'"$TMP_EDGE_CACHEGROUP"'"]; print match[0]')"
				echo "Got cachegroup ID: $TMP_CACHEGROUP_ID"
	
				curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "parameter=$TMP_PARAMETER_ID"  --data-urlencode "profile=$TMP_PARAMETER_PROFILE_ID" $TRAFFIC_OPS_URI/profileparameter/create
				curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "parameter=$TMP_PARAMETER_ID"  --data-urlencode "cachegroup=$TMP_CACHEGROUP_ID" $TRAFFIC_OPS_URI/cachegroupparameter/create
				curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "parameter=$TMP_PARAMETER_ID"  --data-urlencode "cachegroup=$TMP_EDGE_CACHEGROUP_ID" $TRAFFIC_OPS_URI/cachegroupparameter/create
		fi

	fi

	TMP_SERVER_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/servers.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["hostName"]=="'"$HOSTNAME"'"]; print match[0]')"
	echo "Got server ID: $TMP_SERVER_ID"

	curl -v -k -H "Content-Type: application/x-www-form-urlencoded" -H "Cookie: mojolicious=$TMP_TO_COOKIE" -X POST --data-urlencode "id=$TMP_SERVER_ID" --data-urlencode "status=ONLINE" $TRAFFIC_OPS_URI/server/updatestatus

	ls -ltr /opt 

	/opt/traffic_monitor/bin/traffic_monitor_config.pl $TRAFFIC_OPS_URI $TRAFFIC_OPS_USER:$TRAFFIC_OPS_PASS auto

	# uncomment if running Traffic Ops in debug, i.e. http and port 3000
 	# sed -i 's#https://#http://#g' /opt/traffic_monitor/conf/traffic_monitor_config.js

	echo "INITIALIZED=1" >> /etc/environment
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
