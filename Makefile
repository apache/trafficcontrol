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

T3C_TARGETS := cache-config/t3c/t3c cache-config/t3c-apply/t3c-apply cache-config/t3c-check/t3c-check t3c-check-refs/t3c-check-refs t3c-check-reload/t3c-check-reload cache-config/t3c-diff/t3c-diff cache-config/t3c-generate/t3c-generate cache-config/t3c-preprocess/t3c-preprocess cache-config/t3c-request/t3c-request cache-config/t3c-update/t3c-update

.PHONY: lint unit all check clean

all: traffic_ops/app/db/admin $(T3C_TARGETS)

$(T3C_TARGETS):
	$(MAKE) -C cache-config/ $@

traffic_ops/app/db/admin: traffic_ops/app/db/admin.go
	cd $(dir $@) && go build -o $(notdir $@) .

check: unit lint

lint:
	golangci-lint run ./...

cache-config/t3c-check-refs/t3c-check-refs: cache-config/t3c-check-refs/config/config.go cache-config/t3c-check-refs/t3c-check-refs.go
	go build "github.com/apache/trafficcontrol/v8/cache-config/t3c-check-refs"
	mv -f t3c-check-refs $@

unit: cache-config/t3c-check-refs/t3c-check-refs
	go test $(addsuffix ...,$(addprefix ./,$(CACHE_CONFIG_DIRS))) ./grove/... ./lib/... ./traffic_monitor/... ./traffic_ops/traffic_ops_golang/... ./traffic_stats/...

clean:
	$(RM) traffic_ops/app/db/admin $(T3C_TARGETS)
