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

waitfor() {
    local endpoint="$1"
    local field="$2"
    local value="$3"

    while true; do
        v=$(to-get "api/1.4/$endpoint?$field=$value" | jq -r --arg field "$field" '.response[][$field]')
        if [[ $v == $value ]]; then
            break
        fi
        echo "waiting for $endpoint $field=$value"
        sleep 3
    done
}

# special cases -- any data type requiring specific data to already be available in TO should have an entry here.
# e,g. deliveryservice_servers requires both deliveryservice and all servers to be available
delayfor() {
    local f="$1"
    local d="${f%/*}"

    case $d in
        deliveryservice_servers)
            ds=$( jq -r .xmlId <"$f" )
            waitfor deliveryservices xmlId "$ds"
            for s in $( jq -r .serverNames[] <"$f" ); do
                waitfor servers hostName "$s"
            done
            ;;
    esac
}

load_data_from() {
    local dir="$1"
    if [[ ! -d $dir ]] ; then
        echo "Failed to load data from '$dir': directory does not exist"
    fi
    cd "$dir"
    local status=0
    for d in $endpoints; do
        [[ -d $d ]] || continue
        # Let containers know to write out server.json
        if [[ "$d" = "deliveryservice_servers" ]] ; then
           touch "$ENROLLER_DIR/initial-load-done"
           sync
        fi
        for f in "$d"/*.json; do
            echo "Loading $f"
            delayfor "$f"
            envsubst "$vars" <$f  > "$ENROLLER_DIR"/$f
            sync
        done
    done
    if [[ $status -ne 0 ]]; then
        exit $status
    fi
    cd -
}

# First,  load required data at the top level
load_data_from /traffic_ops_data
