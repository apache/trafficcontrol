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
#
# Defines bash functions to consistently interact with the Traffic Ops API


export TO_URL=${TO_URL:-https://$TO_HOST:$TO_PORT}
export TO_USER=${TO_USER:-$TO_ADMIN_USER}
export TO_PASSWORD=${TO_PASSWORD:-$TO_ADMIN_PASSWORD}

export CURLOPTS=${CURLOPTS:--LfsS}
export CURLAUTH=${CURLAUTH:--k}
export COOKIEJAR=$(mktemp)


login=$(mktemp)

cleanup() {
	rm -f "$COOKIEJAR" "$login"
}

trap cleanup EXIT

cookie_current() {
	local cookiefile=$1
	[[ -s $cookiefile ]] || return 1

	# get expiration from cookiejar -- compare to current time
	exp=$(awk '/mojolicious/ {print $5}' $cookiefile | tail -n 1)
	cur=$(date +%s)

	# return value is the comparison itself
	(( $exp > $cur ))
}

to-auth() {
	# These are required
	if [[ -z $TO_URL || -z $TO_USER || -z $TO_PASSWORD ]]; then
		echo TO_URL TO_USER TO_PASSWORD must all be set
		return 1
	fi

	# if cookiejar is current, nothing to do..
	cookie_current $COOKIEJAR && return

	local url=$TO_URL/api/1.3/user/login
	local datatype='Accept: application/json'
	cat >"$login" <<-CREDS
{ "u" : "$TO_USER", "p" : "$TO_PASSWORD" }
CREDS
	res=$(curl $CURLAUTH $CURLOPTS -H "$datatype" --cookie "$COOKIEJAR" --cookie-jar "$COOKIEJAR" -X POST --data @"$login" "$url")
	if [[ $res != *"Successfully logged in."* ]]; then
		echo "Login failed: $res"
		return 1
	fi
}

to-ping() {
	# ping endpoint does not require authentication
	curl $CURLAUTH $CURLOPTS -X GET "$TO_URL/api/1.3/ping"
}

to-get() {
	to-auth && \
		curl $CURLAUTH $CURLOPTS --cookie "$COOKIEJAR" -X GET "$TO_URL/$1"
}

to-post() {
	if [[ -z "$2" ]]; then
		data=""
	elif [[ -f "$2" ]]; then
		data="--data @$2"
	else
		data="--data $2"
	fi
	to-auth && \
	    curl $CURLAUTH $CURLOPTS --cookie "$COOKIEJAR" -X POST $data "$TO_URL/$1"
}

to-put() {
	if [[ -z "$2" ]]; then
		data=""
	elif [[ -f "$2" ]]; then
		data="--data @$2"
	else
		data="--data $2"
	fi
	to-auth && \
	    curl $CURLAUTH $CURLOPTS --cookie "$COOKIEJAR" -X PUT $data "$TO_URL/$1"
}

to-delete() {
	to-auth && \
		curl $CURLAUTH $CURLOPTS --cookie "$COOKIEJAR" -X DELETE "$TO_URL/$1"
}

to-enroll() {
    local service=$1
    until nc enroller 443 </dev/null >/dev/null 2>&1; do 
        echo "waiting for enroller"
        sleep 5
    done

    action=${service:+?name=$service}
    curl -k -X POST https://enroller${action}
}
