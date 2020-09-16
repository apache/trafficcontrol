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

# Traffic Ops API Tests

The Traffic Ops Client API tests are used to validate the clients responses against those from the Traffic Ops API.  

The v1 tests are for regression purposes, and the v2 tests were forked from them when Traffic Ops API v2 was merged. All further feature development will only occur in v2.

In order to run the tests you will need the following:

1. Port access to both the Postgres port (usually 5432) that your Traffic Ops instance is using as well as the Traffic Ops configured port (usually 443 or 60443).

2. An instance of Postgres running with a `to_test` database that has empty tables.

    To get your to_test database setup do the following:
    
    `$ cd trafficcontrol/traffic_ops/app`
    
    `$ db/admin --env=test reset` 

    NOTE on passwords:
    Check that the passwords defined defined for your `to_test` database match 
    here: `trafficcontrol/traffic_ops/app/conf/test/database.conf`
    and here: `traffic-ops-test.conf` 

    The Traffic Ops users will be created by the tool for accessing the API once the database is accessible.
    
    Note that for the database to successfully set up the tables and run the migrations, you will need `goose`.
    On your local machines(Macs), DO NOT install `goose` using `homebrew`. It will install a different version
    `goose`. Instead, install it using the script located at
    `traffic_ops/install/bin/install_goose.sh`. 
    To test if `goose` migrated everything correctly, you can run the following command from the `traffic_ops/app`
    directory:
    `goose -env=test status "user=traffic_ops dbname=to_test sslmode=disable"`
    
    The result should be something similar to:
    ```
     goose: status for environment 'test'
     Applied At                  Migration
     =======================================
     Thu Apr  9 21:55:29 2020 -- 20181206000000_create_monitor_snapshots.sql
     Thu Apr  9 21:55:29 2020 -- 20190219000000_add_consistent_hash_regex.sql
     Thu Apr  9 21:55:29 2020 -- 20190319000000_add_max_origin_connections.sql
     Thu Apr  9 21:55:30 2020 -- 20190513000000_add-allowed_query_keys.sql
     Thu Apr  9 21:55:30 2020 -- 20191004000000_add_server_capabilities.sql
     Thu Apr  9 21:55:30 2020 -- 20191005000000_add_server_server_capability.sql
     Thu Apr  9 21:55:30 2020 -- 20191024000000_add_deliveryservices_required_capability.sql
     Thu Apr  9 21:55:30 2020 -- 20191215000000_add_ecs_enabled.sql
     Thu Apr  9 21:55:30 2020 -- 20200218000000_add_dns_challenges.sql
     Thu Apr  9 21:55:30 2020 -- 20200227000000_add_ds_slice_block_size.sql
     Thu Apr  9 21:55:30 2020 -- 20200313000000_add_server_ip_toggles.sql
     ```
    
    For more info see: http://trafficcontrol.apache.org/docs/latest/development/traffic_ops.html?highlight=reset

3. A running Traffic Ops instance running with the `secure` (https) and is pointing to the `to_test` 
   database by running in `MOJO_MODE=test` which will point to your `to_test` database.
    To get your to_test database setup do the following:
    
   	`$ export MOJO_MODE=test`  
   	
   	`$ cd trafficcontrol/traffic_ops/app`
   	
    `$ bin/start.pl --secure`

4. A running Traffic Ops Golang proxy pointing to the to_test database.
	`$ cd trafficcontrol/traffic_ops/traffic_ops_golang`
	`$ cp ../app/conf/cdn.conf $HOME/cdn.conf`
	change `traffic_ops_golang->port` to 8443

    `$ go build && ./traffic_ops_golang -cfg $HOME/cdn.conf -dbcfg ../app/conf/test/database.conf`
    
    In your local development environment, if the above command fails for an error similar to 
    `ERROR: traffic_ops_golang.go:193: 2020-04-10T10:55:53.190298-06:00: cannot open /etc/pki/tls/certs/localhost.crt for read: open /etc/pki/tls/certs/localhost.crt: no such file or directory`
    you might not have the right certificates at the right locations. Follow the procedure listed
    [here](https://traffic-control-cdn.readthedocs.io/en/latest/admin/traffic_ops.html#id12) to fix it. 
## Running the API Tests
The integration tests are run using `go test`, however, there are some flags that need to be provided in order for the tests to work.  

The flags are:

* usage - API Test tool usage
* cfg - the config file needed to run the tests
* env - Environment variables that can be used to override specific config options that are specified in the config file
* env_vars - Show environment variables that can be overridden
* test_data - traffic control
* run - Go runtime flag for executing a specific test case

Example command to run the tests: 
`TO_URL=https://localhost:8443 go test -v -cfg=traffic-ops-test.conf -run TestCDNs`



* It can take several minutes for the API tests to complete, so using the `-v` flag is recommended to see progress.*
