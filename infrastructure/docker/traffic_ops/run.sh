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

# Script for running the Dockerfile for Traffic Ops.
# The Dockerfile sets up a Docker image which can be used for any new Traffic Ops container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set, ordinarily by `docker run -e` arguments:
# MYSQL_IP
# MYSQL_PORT
# MYSQL_ROOT_PASS
# MYSQL_TRAFFIC_OPS_PASS
# ADMIN_USER
# ADMIN_PASS
# CERT_COUNTRY
# CERT_STATE
# CERT_CITY
# CERT_COMPANY
# DOMAIN

# TODO:  Unused -- should be removed?  TRAFFIC_VAULT_PASS

# Check that env vars are set
envvars=( MYSQL_IP MYSQL_PORT MYSQL_ROOT_PASS MYSQL_TRAFFIC_OPS_PASS ADMIN_USER ADMIN_PASS CERT_COUNTRY CERT_STATE CERT_CITY CERT_COMPANY DOMAIN)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

start() {
		service traffic_ops start
		exec tail -f /var/log/traffic_ops/traffic_ops.log
}

init() {
		mysql -h $MYSQL_IP -P $MYSQL_PORT -u root -p$MYSQL_ROOT_PASS -e "GRANT ALL ON * . * TO 'traffic_ops'@'localhost' IDENTIFIED BY '$MYSQL_TRAFFIC_OPS_PASS';"

		printf '\
#!/usr/bin/expect -f\n \
set force_conservative 0  ;# set to 1 to force conservative mode even if\n\
                          ;# script wasn'\''t run conservatively originally\n\
if {$force_conservative} {\n\
	 set send_slow {1 .1}\n\
	 proc send {ignore arg} {\n\
	 sleep .1\n\
	 exp_send -s -- $arg\n\
	 }\n\
}\n\
set timeout -1\n\
spawn /opt/traffic_ops/install/bin/postinstall\n\
match_max 100000\n\
expect -exact "Hit ENTER to continue:  "\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Database type \[mysql\]:  "\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Database name \[traffic_ops_db\]:  "\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Database server hostname IP or FQDN \[localhost\]:  "\n\
sleep 0.5\n\
send -- "$env(MYSQL_IP)\\r"\n\
expect -exact "Database port number \[3306\]:  "\n\
sleep 0.5\n\
send -- "$env(MYSQL_PORT)\\r"\n\
expect -exact "Traffic Ops database user \[traffic_ops\]:  "\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Password for traffic_ops:  "\n\
sleep 0.5\n\
send -- "$env(MYSQL_TRAFFIC_OPS_PASS)\\r"\n\
expect -exact "Re-Enter Password for traffic_ops:  "\n\
sleep 0.5\n\
send -- "$env(MYSQL_TRAFFIC_OPS_PASS)\\r"\n\
expect -exact "Database server root (admin) user name \[root\]:  "\n\
sleep 0.5\n\
send -- "root\\r"\n\
expect -exact "Database server root password:  "\n\
sleep 0.5\n\
send -- "$env(MYSQL_ROOT_PASS)\\r"\n\
expect -exact "Is the above information correct (y/n) \[n\]:  "\n\
sleep 0.5\n\
send -- "y\\r"\n\
expect -exact "Traffic Ops url \[https://localhost\]:  "\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Human-readable CDN Name.  (No whitespace, please) \[kabletown_cdn\]:  "\n\
sleep 0.5\n\
send -- "cdn\\r"\n\
expect -exact "DNS sub-domain for which your CDN is authoritative \[cdn1.kabletown.net\]:  "\n\
sleep 0.5\n\
send -- "$env(DOMAIN)\\r"\n\
expect -exact "Fully qualified name of your CentOS 6.5 ISO kickstart tar file, or '\''na'\'' to skip and add files later \[/var/cache/centos65.tgz\]:  "\n\
sleep 0.5\n\
send -- "na\\r"\n\
expect -exact "Fully qualified location to store your ISO kickstart files \[/var/www/files\]:  "\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Is the above information correct (y/n) \[n\]:  "\n\
sleep 0.5\n\
send -- "y\\r"\n\
expect -exact "Administration username for Traffic Ops \[admin\]:  "\n\
sleep 0.5\n\
send -- "$env(ADMIN_USER)\\r"\n\
expect "Password for the admin user *:  "\n\
sleep 0.5\n\
send -- "$env(ADMIN_PASS)\\r"\n\
expect "Re-Enter Password for the admin user *:  "\n\
sleep 0.5\n\
send -- "$env(ADMIN_PASS)\\r"\n\
expect -exact "Do you wish to create an ldap configuration for access to traffic ops \[y/n\] ? \[n\]:  "\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Do want to add a new one (only 2 will be kept) \[y/n\] ? \[y\]:  "\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Do you want one generated for you \[y/n\] ? \[y\]:  "\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Hit Enter when you are ready to continue:  "\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Enter pass phrase for server.key:"\n\
sleep 0.5\n\
send -- "pass\\r"\n\
expect -exact "Verifying - Enter pass phrase for server.key:"\n\
sleep 0.5\n\
send -- "pass\\r"\n\
expect -exact "Enter pass phrase for server.key:"\n\
sleep 0.5\n\
send -- "pass\\r"\n\
expect -exact "Country Name (2 letter code) \[XX\]:"\n\
sleep 0.5\n\
send -- "$env(CERT_COUNTRY)\\r"\n\
expect -exact "State or Province Name (full name) \[\]:"\n\
sleep 0.5\n\
send -- "$env(CERT_STATE)\\r"\n\
expect -exact "Locality Name (eg, city) \[Default City\]:"\n\
sleep 0.5\n\
send -- "$env(CERT_CITY)\\r"\n\
expect -exact "Organization Name (eg, company) \[Default Company Ltd\]:"\n\
sleep 0.5\n\
send -- "$env(CERT_COMPANY)\\r"\n\
expect -exact "Organizational Unit Name (eg, section) \[\]:"\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Common Name (eg, your name or your server'\''s hostname) \[\]:"\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Email Address \[\]:"\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "A challenge password \[\]:"\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "An optional company name \[\]:"\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Enter pass phrase for server.key.orig:"\n\
sleep 0.5\n\
send -- "pass\\r"\n\
expect -exact "Install Cron entry to clean install .iso files older than 7 days? \[y/n\] \[n\]:"
send -- "\\r"\n\
sleep 0.5\n\
expect -exact "Health Polling Interval (milliseconds) \[8000\]"\n\
sleep 0.5\n
send -- "\\r"\n\
expect -exact "TLD SOA admin \[traffic_ops\]:"\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "TrafficServer Drive Prefix \[/dev/sd\]:"\n\
sleep 0.5\n\
send -- "/dev/ram\\r"\n\
expect -exact "TrafficServer RAM Drive Prefix \[/dev/ram\]:"\n\
sleep 0.5\n\
send -- "/dev/ram\\r"\n\
expect -exact "TrafficServer RAM Drive Letters (comma separated) \[0,1,2,3,4,5,6,7\]:"\n\
sleep 0.5\n\
send -- "1\\r"\n\
expect -exact "Health Threshold Load Average \[25\]:"\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Health Threshold Available Bandwidth in Kbps \[1750000\]:"\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Traffic Server Health Connection Timeout (milliseconds) \[2000\]:"\n\
sleep 0.5\n\
send -- "\\r"\n\
expect -exact "Shutdown Traffic Ops \[y/n\] \[n\]:  "\n\
sleep 0.5\n\
send -- "n\\r"\
' > postinstall.exp

		export TERM=xterm && export USER=root && expect postinstall.exp

		TRAFFIC_OPS_URI="https://localhost"

		TMP_TO_COOKIE="$(curl -v -s -k -X POST --data '{ "u":"'"$ADMIN_USER"'", "p":"'"$ADMIN_PASS"'" }' $TRAFFIC_OPS_URI/api/1.2/user/login 2>&1 | grep 'Set-Cookie' | sed -e 's/.*mojolicious=\(.*\); expires.*/\1/')"
		echo "Got cookie: $TMP_TO_COOKIE"

		TMP_DOMAIN=$DOMAIN
		sed -i -- "s/{{.Domain}}/$TMP_DOMAIN/g" /profile.origin.traffic_ops
		echo "Got domain: $TMP_DOMAIN"

		echo "Importing origin"
		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" -F "filename=profile.origin.traffic_ops" -F "profile_to_import=@/profile.origin.traffic_ops" $TRAFFIC_OPS_URI/profile/doImport

		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "division.name=East" $TRAFFIC_OPS_URI/division/create
		TMP_DIVISION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/region/add | grep --color=never -oE "<option value=\"[0-9]+\">East</option>" | grep --color=never -oE "[0-9]+")"
		echo "Got division ID: $TMP_DIVISION_ID"

		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "region.name=Eastish" --data-urlencode "region.division_id=$TMP_DIVISION_ID" $TRAFFIC_OPS_URI/region/create
		TMP_REGION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/regions.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="Eastish"]; print match[0]')"
		echo "Got region ID: $TMP_REGION_ID"

		TMP_CACHEGROUP_TYPE="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="MID_LOC"]; print match[0]')"
		echo "Got cachegroup type ID: $TMP_CACHEGROUP_TYPE"

		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "cg_data.name=mid-east" --data-urlencode "cg_data.short_name=east" --data-urlencode "cg_data.latitude=0" --data-urlencode "cg_data.longitude=0" --data-urlencode "cg_data.parent_cachegroup_id=-1" --data-urlencode "cg_data.type=$TMP_CACHEGROUP_TYPE" $TRAFFIC_OPS_URI/cachegroup/create
		TMP_CACHEGROUP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="mid-east"]; print match[0]')"
		echo "Got cachegroup ID: $TMP_CACHEGROUP_ID"

		TMP_CACHEGROUP_EDGE_TYPE="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="EDGE_LOC"]; print match[0]')"
		echo "Got cachegroup type ID: $TMP_CACHEGROUP_EDGE_TYPE"

		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "cg_data.name=edge-east" --data-urlencode "cg_data.short_name=eeast" --data-urlencode "cg_data.latitude=0" --data-urlencode "cg_data.longitude=0" --data-urlencode "cg_data.parent_cachegroup_id=$TMP_CACHEGROUP_ID" --data-urlencode "cg_data.type=$TMP_CACHEGROUP_EDGE_TYPE" $TRAFFIC_OPS_URI/cachegroup/create
		TMP_CACHEGROUP_EDGE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/1.2/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="edge-east"]; print match[0]')"
		echo "Got cachegroup edge ID: $TMP_CACHEGROUP_EDGE_ID"

		curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "location.name=plocation-nyc-1" --data-urlencode "location.short_name=nyc" --data-urlencode "location.address=1 Main Street" --data-urlencode "location.city=nyc" --data-urlencode "location.state=NY" --data-urlencode "location.zip=12345" --data-urlencode "location.poc=" --data-urlencode "location.phone=" --data-urlencode "location.email=no@no.no" --data-urlencode "location.comments=" --data-urlencode "location.region=$TMP_REGION_ID" $TRAFFIC_OPS_URI/phys_location/create

		echo "INITIALIZED=1" >> /etc/environment
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
