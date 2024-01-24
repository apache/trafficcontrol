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

set -ex
export BROWSER_FOLDER="/experimental/traffic-portal/dist/traffic-portal/browser"

cd "${GITHUB_WORKSPACE}/traffic_ops/traffic_ops_golang"

truncate -s0 out.log
envsubst <../../.github/actions/tpv2-integration-tests/cdn.json >./cdn.conf

./traffic_ops_golang --cfg ./cdn.conf --dbcfg ../../.github/actions/tpv2-integration-tests/database.json > out.log 2>&1 &

timeout 3m bash <<TMOUT
	while ! curl -k "https://localhost:6443/api/4.0/ping" >/dev/null 2>&1; do
		echo "waiting for TO API"
		sleep 5
	done
TMOUT

cd "${GITHUB_WORKSPACE}/experimental/traffic-portal"
timeout 15m npm run e2e:ci
kill %%
