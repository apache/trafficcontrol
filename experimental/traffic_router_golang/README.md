<!--
    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.
-->

# Traffic Router 

This is a prototype of Traffic Router's HTTP side in Golang.

# How to build

To get this app running locally:

- Clone this repo
- Install Golang programming language ([instructions](https://golang.org/doc/install))


# Configuration

Sample configuration file(cfg.json) available in traffic_router_golang directory, please add coveragezone files to path specified in cfg.json
   

# Build

Compile and generate binary:

   - `cd traffic_router_goland`
   - `go mod vendor`
   - `go build` #This will generate binary file traffic_router_golang)

# Unit Test
    
   - Run `go test ./...` from traffic_router_golang directory
     
     ```$ go test ./...
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang     [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/availableservers    [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/cgsrch      [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/config      [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/coveragezone        [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crconfig    [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crconfigdsservers   [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crconfigpoller      [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crconfigregex       [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crstates    [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/crstatespoller      [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/fetch       [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/httpsrvr    [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/ipmap       [no test files]
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/nextcache   [no test files]
		ok      github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/quadtree    1.190s
		?       github.com/apache/trafficcontrol/v8/experimental/traffic_router_golang/toutil      [no test files]
     ```
