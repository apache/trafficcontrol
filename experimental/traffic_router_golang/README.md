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
   - `go mod download`
   - `go build` #This will generate binary file traffic_router_golang)

# Known issues in build
   - if you find errors similar to below,
	``` go build
	..\..\vendor\github.com\lestrrat-go\jwx\x25519\x25519.go:9:2: cannot find package "." in:
	        C:\users\moham\go\src\github.com\apache\trafficcontrol\vendor\golang.org\x\crypto\curve25519
	..\..\traffic_ops\toclientlib\toclientlib.go:36:2: cannot find package "." in:
	        C:\users\moham\go\src\github.com\apache\trafficcontrol\vendor\golang.org\x\net\publicsuffix```
    plese run `go mod vendor` which will sync vendor directory, now run go build will build and generate binary

# Unit Test
    
   - Run `go test ./...` from traffic_router_golang directory
     
     ```$ go test ./...
		?       traffic_router_golang   [no test files]
		?       traffic_router_golang/availableservers  [no test files]
		?       traffic_router_golang/cgsrch    [no test files]
		?       traffic_router_golang/config    [no test files]
		?       traffic_router_golang/coveragezone      [no test files]
		?       traffic_router_golang/crconfig  [no test files]
		?       traffic_router_golang/crconfigdsservers [no test files]
		?       traffic_router_golang/crconfigpoller    [no test files]
		?       traffic_router_golang/crconfigregex     [no test files]
		?       traffic_router_golang/crstates  [no test files]
		?       traffic_router_golang/crstatespoller    [no test files]
		?       traffic_router_golang/fetch     [no test files]
		?       traffic_router_golang/httpsrvr  [no test files]
		?       traffic_router_golang/ipmap     [no test files]
		?       traffic_router_golang/nextcache [no test files]
		ok      traffic_router_golang/quadtree  0.526s
		?       traffic_router_golang/toutil    [no test files]
```
