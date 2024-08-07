#!/bin/sh -l
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

set -ex;

store_ciab_logs() {
	echo 'Storing CDN-in-a-Box logs...';
	mkdir logs;
	for service in $($docker_compose ps --services --all); do
		$docker_compose logs --no-color --timestamps "$service" >"logs/${service}.log";
	done;
}

cd infrastructure/cdn-in-a-box;
logged_services='trafficrouter readiness';
other_services='dns edge enroller mid-01 mid-02 origin static trafficmonitor trafficops trafficstats';
docker_compose='docker compose -f ./docker-compose.yml -f ./docker-compose.readiness.yml';
$docker_compose up -d $logged_services $other_services;
$docker_compose logs -f $logged_services &

echo 'Waiting for the readiness container to exit...';
if ! timeout 12m $docker_compose logs -f readiness >/dev/null; then
	echo "CDN-in-a-Box didn't become ready within 12 minutes - exiting" >&2;
	exit_code=1;
	store_ciab_logs;
elif exit_code="$(docker inspect --format='{{.State.ExitCode}}' "$($docker_compose ps -q --all readiness)")"; [ "$exit_code" -ne 0 ]; then
	echo 'Readiness container exited with an error' >&2;
	store_ciab_logs;
fi;

exit "$exit_code";
