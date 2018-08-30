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

TO_URL="https://$TO_HOST:$TO_PORT"
# wait until the ping endpoint succeeds
while ! to-ping 2>/dev/null; do
   echo waiting for trafficops
   sleep 3
done

# NOTE: order dependent on foreign key references, e.g. tenants must be defined before users
endpoints="cdns divisions regions phys_locations tenants users cachegroups deliveryservices"

load_data_from() {
    local dir="$1"
    if [[ ! -d $dir ]] ; then
        echo "Failed to load data from '$dir': directory does not exist"
    fi

    local status=0
    for ep in $endpoints; do
        d="$dir/$ep"
        [[ -d $d ]] || continue
        echo "Loading data from $d"
        for f in "$d"/*.json; do
            [[ -r $f ]] || continue
            t=$(mktemp --tmpdir $ep-XXX.json)
            envsubst <"$f" >"$t"
            if ! to-post api/1.3/"$ep" "$t"; then
                echo POST api/1.3/"$ep" "$t" failed
                status=$?
            fi
            rm "$t"
        done
    done
    if [[ $status -ne 0 ]]; then
        exit $status
    fi


}

# First,  load required data at the top level
load_data_from /traffic_ops_data

# If TO_DATA is defined, load from subdirs with that name (space-separated)
if [[ -n $TO_DATA ]]; then
    for subdir in $TO_DATA; do
        load_data_from /traffic_ops_data/$subdir
    done
fi


