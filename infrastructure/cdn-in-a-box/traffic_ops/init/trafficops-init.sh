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

# Required env vars
# Check that env vars are set
set -x
for v in TO_HOST TO_PORT TO_ADMIN_USER TO_ADMIN_PASSWORD; do
    [[ -z $(eval echo \$$v) ]] || continue
    echo "$v is unset"
    exit 1
done

TO_URL="https://$TO_HOST:$TO_PORT"
# wait until the ping endpoint succeeds
while ! curl -k $TO_URL/api/1.3/ping; do
   echo waiting for trafficops
   sleep 3
done

export COOKIEJAR=/tmp/cookiejar.$(echo $TO_URL $TO_ADMIN_USER | md5sum | awk '{print $1}')

login() {
    local datadir=$(mktemp -d)
    local login="$datadir/login.json"
    local url=$TO_URL/api/1.3/user/login
    local datatype='Accept: application/json'
    cat > "$login"  <<-CREDS
    { "u" : "$TO_ADMIN_USER", "p" : "$TO_ADMIN_PASSWORD" }
CREDS

    res=$(curl -k -H "$datatype" --cookie "$COOKIEJAR" --cookie-jar "$COOKIEJAR" -X POST --data @"$login" "$url")
    rm -rf "$datadir"
    if [[ $res != *"Successfully logged in."* ]]; then
        echo $res
        return -1
    fi
}

login

# load json files using the API
for f in */*.json; do
    ep=$(dirname $f)
    curl -k -s --cookie "$COOKIEJAR" -X POST --data @"$f" "$TO_URL/api/1.3/$ep"
done
