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

# Traffic Ops Client / API Integration Tests

The Traffic Ops Client API tests are used to validate the clients responses against those from the Traffic Ops API.

## Setup

In order to run the tests you will need a running instance of Traffic Ops and Traffic Ops DB:

1. **Traffic Ops Database** configured port access
    - _Usually 5432 - should match the value set in database.conf and the **trafficOpsDB port** in traffic-ops-test.conf_
2. **Traffic Ops** configured port access
    - _Usually 443 or 60443 - should match the value set in cdn.conf and the **URL** in traffic-ops-test.conf_
3. Running Postgres instance with a database that matches the **trafficOpsDB dbname** value set in traffic-ops-test.conf
    - For example to set up the `to_test` database do the following:

         ```console
         $ cd trafficcontrol/traffic_ops/app
         $ db/admin --env=test reset
         ```

      The Traffic Ops users will be created by the tool for accessing the API once the database is accessible.

      To test if `db/admin` ran all migrations successfully, you can run the following command from the `traffic_ops/app` directory:

        ```console
        db/admin -env=test dbversion
        ```
      The result should be something similar to:
        ```
        dbversion 2021070800000000
        ```
      If migrations did not run successfully, you may see:
        ```
        dbversion 20181206000000 (dirty)
        ```
      Make sure **trafficOpsDB dbname** in traffic-ops-test.conf is set to: `to_test`

      For more info see: http://trafficcontrol.apache.org/docs/latest/development/traffic_ops.html?highlight=reset

4. A running Traffic Ops Golang instance pointing to the `to_test` database.

    ```console
	$ cd trafficcontrol/traffic_ops/traffic_ops_golang
    $ cp ../app/conf/cdn.conf $HOME/cdn.conf 
    $ go build && ./traffic_ops_golang -cfg $HOME/cdn.conf -dbcfg ../app/conf/test/database.conf
    ```
   Verify that the passwords defined for your `to_test` database match:
    - `trafficcontrol/traffic_ops/app/conf/test/database.conf`
    - `traffic-ops-test.conf`

## Running the API Tests

The integration tests are run using `go test` from the **traffic_ops/testing/api/** directory, however, there are some flags that need to be provided in order for the tests to work.

The flags are:

* cfg - The config file path (default traffic-ops-test.conf)
* fixtures - The test fixtures for the API test tool (default "tc-fixtures.json)
* includeSystemTests - Whether to enable tests that have environment dependencies beyond a database
* run - Go runtime flag for executing a specific test case

Example commands to run the tests:

Test all v3 tests:
```console
go test ./v3/... -v --cfg=../conf/traffic-ops-test.conf
```

Test all v4 tests:
```console
go test ./v4/... -v --cfg=../conf/traffic-ops-test.conf
```

Test all v5 tests:
```console
go test ./v5/... -v --cfg=../conf/traffic-ops-test.conf
```

Only Test a specific endpoint:
```console
go test ./v4/... -run ^TestCDNs$ -v -cfg=../conf/traffic-ops-test.conf
```

Only Test a specific endpoint method :
```console
go test ./v5/... -run ^TestCacheGroups$/GET -v -cfg=../conf/traffic-ops-test.conf
```

Get Test Coverage:
```console
go test ./v4/... -v --cfg=../conf/traffic-ops-test.conf -coverpkg=../../v4-client/... -coverprofile=cover.out
```
View the cover out file in your browser:
```console
go tool cover -html=cover.out
```

* go test -run flag matches a regexp so if you want to ensure that only a test named TestCDNs is run use ^TestCDNs$
* It can take several minutes for the API tests to complete, so using the `-v` flag is recommended to see progress.*
