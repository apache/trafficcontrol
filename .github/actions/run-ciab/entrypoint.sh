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

set -ex

docker-compose --version;
STARTING_POINT="$PWD";
cd infrastructure/cdn-in-a-box;
make traffic_ops/traffic_ops.rpm traffic_stats/traffic_stats.rpm traffic_monitor/traffic_monitor.rpm traffic_router/tomcat.rpm traffic_router/traffic_router.rpm;

time docker-compose -f ./docker-compose.yml -f ./docker-compose.readiness.yml -f ./docker-compose.traffic-ops-test.yml build integration edge mid origin readiness trafficops trafficops-perl dns enroller trafficrouter trafficstats trafficvault trafficmonitor;
time docker-compose -f ./docker-compose.yml -f ./docker-compose.readiness.yml up -d edge mid origin readiness trafficops trafficops-perl dns enroller trafficrouter trafficstats trafficvault trafficmonitor;

ret=$(timeout 10m docker wait cdn-in-a-box_readiness_1)
if [[ "$ret" -ne 0 ]]; then
	echo "CDN in a Box didn't become ready within 10 minutes - exiting" >&2;
	docker-compose -f ./docker-compose.readiness.yml logs;
	docker-compose -f ./docker-compose.yml -f ./docker-compose.readiness.yml down -v --remove-orphans;
	exit "$ret";
fi

docker-compose -f ./docker-compose.traffic-ops-test.yml up;
ls junit || echo "JUnit dir not found";
docker-compose -f ./docker-compose.yml -f ./docker-compose.readiness.yml -f ./docker-compose.traffic-ops-test.yml down -v --remove-orphans;
