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

trap 'echo "Error on line ${LINENO} of ${0}"; exit 1' ERR;
set -o xtrace
set -o errexit -o nounset -o pipefail
docker-compose up -d

# Constants
declare -r cookie_name=dev-ciab-cookie

# Get atc-ready function
source dev/atc.dev.sh

export -f atc-ready
echo 'Waiting until Traffic Ops is ready to accept requests...'
if ! timeout 10m bash -c 'atc-ready -w'; then
	echo 'Traffic Ops was not available within 10 minutes!'
	trap - ERR
	echo 'Exiting...'
	exit 1
fi

source infrastructure/cdn-in-a-box/traffic_ops/to-access.sh

# Log in
#login_body="$(<<<{} jq --arg TO_USER "$TO_USER" --arg TO_PASSWORD "$TO_PASSWORD" '.u = $TO_USER | .p = $TO_PASSWORD')"
#to-post user/login "$login_body"

declare -A service_by_hostname
service_by_hostname[trafficrouter]=trafficrouter
service_by_hostname[edge]=t3c

for hostname in trafficrouter edge; do
	container_id="$(docker-compose ps -q "${service_by_hostname[$hostname]}")"
	interface="$(<<'JSON' jq
	{
		"mtu": 1500,
		"monitor": true,
		"ipAddresses": [],
		"name": "eth0"
	}
JSON
	)"
	docker_network="$(docker network inspect dev.ciab.test)"
	for ip_address_field in IPv4Address IPv6Address; do
		ip_address="$(<<<"$docker_network" jq -r --arg CONTAINER_ID "$container_id" --arg IP_ADDRESS_FIELD "$ip_address_field" '.[0].Containers[$CONTAINER_ID][$IP_ADDRESS_FIELD]')"
		if [[ "$ip_address" == null ]]; then
			echo "Could not find ${ip_address_field} for ${hostname} service!"
			exit 1
		fi
		interface="$(<<<"$interface" jq --arg IP_ADDRESS "$ip_address" '.ipAddresses += [{} | .address = $IP_ADDRESS | .serviceAddress = true]')"
	done


	# Get Traffic Router server JSON
	server="$(to-get "api/$TO_API_VERSION/servers?hostName=${hostname}" | jq '.response[0]')"
	if [[ -z "$server" ]]; then
		echo "Could not get JSON for server ${hostname}"
		exit 1
	fi

	# Update Traffic Router's interface with its IP addresses
	server="$(<<<"$server" jq ".interfaces = [${interface}]")"
	server_id="$(<<<"$server" jq .id)"
	if ! to-put "api/$TO_API_VERSION/servers/${server_id}" "$server"; then
		echo "Could not update server ${hostname} with ${server}"
	fi
done

# Snapshot
cdn_id="$(<<<"$server" jq .cdnId)"
to-put "api/$TO_API_VERSION/snapshot?cdnID=${cdn_id}"

echo "Waiting for Traffic Monitor to serve a snapshot..."
if ! timeout 10m curl \
	--retry 99999 \
	--retry-delay 5 \
	--show-error \
	-fIso/dev/null \
	http://localhost/publish/CrConfig; then
	echo "CrConfig was not available from Traffic Monitor within 10 minutes!"
	trap - ERR
	echo 'Exiting...'
	exit 1
fi


deliveryservice=cdn.dev-ds.ciab.test
echo "Waiting for Delivery Service ${deliveryservice} to be available..."
if ! timeout 10m bash -c 'atc-ready -d'; then
	echo "Delivery Service ${deliveryservice} was not available within 10 minutes!"
	trap - ERR
	echo 'Exiting...'
	exit 1
fi

http_result=0 dns_result=0

http_requests_threshold=1200
dns_requests_threshold=2500
# Compile the tests
go test -c ./traffic_router/ultimate-test-harness
if ! ./ultimate-test-harness.test -test.v -test.run=^TestHTTPLoad$ -http_requests_threshold "$http_requests_threshold"; then
	http_result=1
fi

if ! ./ultimate-test-harness.test -test.v -test.run=^TestDNSLoad$ -dns_requests_threshold "$dns_requests_threshold"; then
	dns_result=1
fi
if [[ $http_result -eq 0 && $dns_result -eq 0 ]]; then echo
	echo Tests passed!
else
	exit_code=$?
	echo Tests failed!
	exit $exit_code
fi
