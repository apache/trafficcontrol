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
set -o errexit -o nounset
cd "${GITHUB_WORKSPACE}"

mv -- dist/*/*.rpm dist/

ciab_dir=infrastructure/cdn-in-a-box
architecture="$(uname -m)"
mv "dist/trafficserver-"*".el${RHEL_VERSION}.${architecture}.rpm" "${ciab_dir}/cache/trafficserver.rpm"

cd "$ciab_dir"

# Make all targets except cache/trafficserver.rpm
make \
	cache/trafficcontrol-cache-config.rpm \
	traffic_monitor/traffic_monitor.rpm \
	traffic_ops/traffic_ops.rpm \
	traffic_portal/traffic_portal.rpm \
	traffic_ops/traffic_ops.rpm \
	traffic_router/tomcat.rpm \
	traffic_router/traffic_router.rpm \
	traffic_stats/traffic_stats.rpm

export DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 # Cuts ~11 mins build time down to ~6 mins
docker-compose pull &
docker-compose build --parallel
