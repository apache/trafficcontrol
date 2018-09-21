#!/usr/bin/env bash
#
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

# Required env vars
# Check that env vars are set
set -x
for v in TO_HOST TO_PORT TO_ADMIN_USER TO_ADMIN_PASSWORD; do
    [[ -z $(eval echo \$$v) ]] || continue
    echo "$v is unset"
    exit 1
done

. /to-access.sh

TO_URL="https://$TO_FQDN:$TO_PORT"
# wait until the ping endpoint succeeds
while ! to-ping 2>/dev/null; do
   echo waiting for trafficops
   sleep 3
done

# NOTE: order dependent on foreign key references, e.g. profiles must be loaded before parameters
endpoints="cdns types divisions regions phys_locations tenants users cachegroups deliveryservices profiles parameters servers deliveryservice_servers"
vars=$(awk -F = '/^\w/ {printf "$%s ",$1}' /variables.env)

load_data_from() {
    local dir="$1"
    if [[ ! -d $dir ]] ; then
        echo "Failed to load data from '$dir': directory does not exist"
    fi
    cd "$dir"
    local status=0
    for d in $endpoints; do
        [[ -d $d ]] || continue
        for f in "$d"/*.json; do 
            echo "Loading $f"
            envsubst "$vars" <$f  > "$ENROLLER_DIR"/$f
        done
    done
    if [[ $status -ne 0 ]]; then
        exit $status
    fi
    # After done loading all data
    touch "$ENROLLER_DIR/initial-load-done"
    cd -
}

# First,  load required data at the top level
load_data_from /traffic_ops_data

# Copy the free MaxMind GeoLite DB to TrafficOps public directory
tar -C /var/tmp -zxpvf /GeoLite2-City.tar.gz
geo_dir=$(find /var/tmp -maxdepth 1 -type d -name GeoLite2-City\*)
gzip -c "$geo_dir/GeoLite2-City.mmdb" > "$TO_DIR/public/GeoLite2-City.mmdb.gz"
chown trafops:trafops "$TO_DIR/public/GeoLite2-City.mmdb.gz"
