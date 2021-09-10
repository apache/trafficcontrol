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
CACHE_CONFIG_DIRS=$(filter-out cache-config/testing/ cache-config/build/,$(wildcard cache-config/*/))
TO_SOURCE=$(filter-out %_test.go,$(wildcard traffic_ops/traffic_ops_golang/**.go))

.PHONY: lint unit all check

all: traffic_ops/traffic_ops_golang/traffic_ops_golang traffic_ops/app/db/admin

traffic_ops/traffic_ops_golang/traffic_ops_golang: $(TO_SOURCE)
	cd traffic_ops/traffic_ops_golang && go build

traffic_ops/app/db/admin: traffic_ops/app/db/admin.go
	cd $(dir $@) && go build -o $(notdir $@) .

check: unit lint

lint:
	golangci-lint run ./...

cache-config/t3c-check-refs/t3c-check-refs: cache-config/t3c-check-refs/config/config.go cache-config/t3c-check-refs/t3c-check-refs.go
	go build "github.com/apache/trafficcontrol/cache-config/t3c-check-refs"
	mv -f t3c-check-refs $@

unit: cache-config/t3c-check-refs/t3c-check-refs
	go test $(addsuffix ...,$(addprefix ./,$(CACHE_CONFIG_DIRS))) ./grove/... ./lib/... ./traffic_monitor/... ./traffic_ops/traffic_ops_golang/... ./traffic_stats/...
