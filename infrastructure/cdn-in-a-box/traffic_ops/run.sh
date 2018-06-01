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

# Script for running the Dockerfile for Traffic Ops.
# The Dockerfile sets up a Docker image which can be used for any new Traffic Ops container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set, ordinarily by `docker run -e` arguments:
# DB_SERVER
# DB_PORT
# DB_ROOT_PASS
# DB_USER
# DB_USER_PASS
# DB_NAME
# ADMIN_USER
# ADMIN_PASS
# CERT_COUNTRY
# CERT_STATE
# CERT_CITY
# CERT_COMPANY
# DOMAIN

# TODO:  Unused -- should be removed?  TRAFFIC_VAULT_PASS

# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS CERT_COUNTRY CERT_STATE CERT_CITY CERT_COMPANY DOMAIN)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

# Write config files
set -x
if [[ -r /config.sh ]]; then
	. /config.sh
fi

while ! nc $DB_SERVER $DB_PORT </dev/null; do # &>/dev/null; do
        echo "waiting for $DB_SERVER $DB_PORT"
        sleep 3
done

while true; do
	echo "Checking for existence of role $DB_USER"
	psql -U postgres -h $DB_SERVER -p $DB_PORT postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='$DB_USER'" | grep -q 1 && break
	sleep 3
done

TO_DIR=/opt/traffic_ops/app
cat conf/production/database.conf

export PERL5LIB=$TO_DIR/lib:$TO_DIR/local/lib/perl5
cd $TO_DIR && ./db/admin.pl -env production reset

cd $TO_DIR && $TO_DIR/local/bin/hypnotoad script/cdn
exec tail -f /var/log/traffic_ops/traffic_ops.log

init() {
	TRAFFIC_OPS_URI="https://localhost"

	COOKIE="$(curl -v -s -k -X POST --data '{ "u":"'"$ADMIN_USER"'", "p":"'"$ADMIN_PASS"'" }' $TRAFFIC_OPS_URI/api/1.2/user/login 2>&1 | grep 'Set-Cookie' | sed -e 's/.*mojolicious=\(.*\); expires.*/\1/')"
	echo "Got cookie: $COOKIE"

	TMP_DOMAIN=$DOMAIN
	sed -i -- "s/{{.Domain}}/$TMP_DOMAIN/g" /profile.origin.traffic_ops
	echo "Got domain: $TMP_DOMAIN"

	echo "Importing origin"
	curl -v -k -X POST -H "Cookie: mojolicious=$COOKIE" -F "filename=profile.origin.traffic_ops" -F "profile_to_import=@/profile.origin.traffic_ops" $TRAFFIC_OPS_URI/profile/doImport

	curl -v -k -X POST -H "Cookie: mojolicious=$COOKIE" --data-urlencode "division.name=East" $TRAFFIC_OPS_URI/division/create
	TMP_DIVISION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$COOKIE" $TRAFFIC_OPS_URI/region/add | grep --color=never -oE "<option value=\"[0-9]+\">East</option>" | grep --color=never -oE "[0-9]+")"
	echo "Got division ID: $TMP_DIVISION_ID"

	curl -v -k -X POST -H "Cookie: mojolicious=$COOKIE" --data-urlencode "region.name=Eastish" --data-urlencode "region.division_id=$TMP_DIVISION_ID" $TRAFFIC_OPS_URI/region/create
	TMP_REGION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$COOKIE" $TRAFFIC_OPS_URI/api/1.2/regions.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="Eastish"]; print match[0]')"
	echo "Got region ID: $TMP_REGION_ID"

	TMP_CACHEGROUP_TYPE="$(curl -s -k -X GET -H "Cookie: mojolicious=$COOKIE" $TRAFFIC_OPS_URI/api/1.2/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="MID_LOC"]; print match[0]')"
	echo "Got cachegroup type ID: $TMP_CACHEGROUP_TYPE"

	curl -v -k -X POST -H "Cookie: mojolicious=$COOKIE" --data-urlencode "cg_data.name=mid-east" --data-urlencode "cg_data.short_name=east" --data-urlencode "cg_data.latitude=0" --data-urlencode "cg_data.longitude=0" --data-urlencode "cg_data.parent_cachegroup_id=-1" --data-urlencode "cg_data.type=$TMP_CACHEGROUP_TYPE" $TRAFFIC_OPS_URI/cachegroup/create
	TMP_CACHEGROUP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$COOKIE" $TRAFFIC_OPS_URI/api/1.2/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="mid-east"]; print match[0]')"
	echo "Got cachegroup ID: $TMP_CACHEGROUP_ID"

	TMP_CACHEGROUP_EDGE_TYPE="$(curl -s -k -X GET -H "Cookie: mojolicious=$COOKIE" $TRAFFIC_OPS_URI/api/1.2/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="EDGE_LOC"]; print match[0]')"
	echo "Got cachegroup type ID: $TMP_CACHEGROUP_EDGE_TYPE"

	curl -v -k -X POST -H "Cookie: mojolicious=$COOKIE" --data-urlencode "cg_data.name=edge-east" --data-urlencode "cg_data.short_name=eeast" --data-urlencode "cg_data.latitude=0" --data-urlencode "cg_data.longitude=0" --data-urlencode "cg_data.parent_cachegroup_id=$TMP_CACHEGROUP_ID" --data-urlencode "cg_data.type=$TMP_CACHEGROUP_EDGE_TYPE" $TRAFFIC_OPS_URI/cachegroup/create
	TMP_CACHEGROUP_EDGE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$COOKIE" $TRAFFIC_OPS_URI/api/1.2/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="edge-east"]; print match[0]')"
	echo "Got cachegroup edge ID: $TMP_CACHEGROUP_EDGE_ID"

	curl -v -k -X POST -H "Cookie: mojolicious=$COOKIE" --data-urlencode "location.name=plocation-nyc-1" --data-urlencode "location.short_name=nyc" --data-urlencode "location.address=1 Main Street" --data-urlencode "location.city=nyc" --data-urlencode "location.state=NY" --data-urlencode "location.zip=12345" --data-urlencode "location.poc=" --data-urlencode "location.phone=" --data-urlencode "location.email=no@no.no" --data-urlencode "location.comments=" --data-urlencode "location.region=$TMP_REGION_ID" $TRAFFIC_OPS_URI/phys_location/create

	echo "INITIALIZED=1" >> /etc/environment
}
