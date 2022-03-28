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

export DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 # build Docker images faster

trap 'echo "Error on line ${LINENO} of ${0}"; exit 1' ERR;
set -o xtrace
set -o errexit -o nounset -o pipefail
docker-compose up -d

# Constants
declare -r cookie_name=dev-ciab-cookie

# Set TO_USER, TO_PASSWORD, and TO_URL environment variables and get atc-ready function
source dev/atc.dev.sh

store_dev_ciab_logs() {
	echo 'Storing Dev CDN-in-a-Box logs...';
	mkdir -p dev/logs;
	for service in $(docker-compose ps --services); do
		docker-compose logs --no-color --timestamps "$service" >"dev/logs/${service}.log";
	done;
}

export -f atc-ready
if ! timeout 10m <<'BASH_COMMANDS' bash; then
set -o errexit -o nounset
until atc-ready; do
	echo 'Waiting until Traffic Ops is ready to accept requests...'
	sleep 3
done
echo 'Traffic Ops is ready to accept requests!'
BASH_COMMANDS
	echo 'Traffic Ops was not available within 10 minutes!'
	store_dev_ciab_logs
	trap - ERR
	echo 'Exiting...'
	exit 1
fi

to-req() {
	endpoint="$1"
	shift
	local curl_command=(curl --insecure --silent --cookie-jar "$cookie_name" --cookie "$cookie_name" "${TO_URL}/api/${API_VERSION}")
	"${curl_command[@]}${endpoint}" "$@" | jq
}

# Log in
login_body="$(<<<{} jq --arg TO_USER "$TO_USER" --arg TO_PASSWORD "$TO_PASSWORD" '.u = $TO_USER | .p = $TO_PASSWORD')"
to-req /user/login --data "$login_body"

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
	for ip_address_field in IPv4Address IPv6Address; do
		ip_address="$(docker network inspect dev.ciab.test |
			jq -r --arg TR_CONTAINER_ID "$container_id" --arg IP_ADDRESS_FIELD "$ip_address_field" '.[0].Containers[$TR_CONTAINER_ID][$IP_ADDRESS_FIELD]')"
		interface="$(<<<"$interface" jq --arg IP_ADDRESS "$ip_address" '.ipAddresses += [{} | .address = $IP_ADDRESS | .serviceAddress = true]')"
	done


	# Get Traffic Router server JSON
	server="$(to-req "/servers?hostName=${hostname}" | jq '.response[0]')"
	if [[ -z "$server" ]]; then
		echo "Could not get JSON for server ${hostname}"
		exit 1
	fi

	# Update Traffic Router's interface with its IP addresses
	server="$(<<<"$server" jq ".interfaces = [${interface}]")"
	server_id="$(<<<"$server" jq .id)"
	if ! to-req "/servers/${server_id}" --request PUT --data "$server"; then
		echo "Could not update server ${hostname} with ${server}"
	fi
done

# Snapshot
cdn_id="$(<<<"$server" jq .cdnId)"
to-req "/snapshot?cdnID=${cdn_id}" --request PUT

http_result=0 dns_result=0
# Compile the tests
go test -c ./traffic_router/ultimate-test-harness
if ! ./ultimate-test-harness.test -test.v -test.run=^TestHTTPLoad$ -http_requests_threshold=5000; then
	http_result=1
fi

if ! ./ultimate-test-harness.test -test.v -test.run=^TestDNSLoad$ -dns_requests_threshold=20500; then
	dns_result=1
fi
if [[ $http_result -eq 0 && $dns_result -eq 0 ]]; then echo
	echo Tests passed!
else
	exit_code=$?
	echo Tests failed!
	exit $exit_code
fi
