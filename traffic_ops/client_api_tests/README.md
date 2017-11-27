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

# Traffic Ops Client Integration Tests

The Traffic Ops Client Integration tests are used to validate the clients responses against those from the Traffic Ops API.  In order to run the tests you will need a Traffic Ops instance with at least one of each of the following:  CDN, Delivery Service, Type, Cachegroup, User and Server.

## Running the Integration Tests
The integration tests are run using `go test`, however, there are some flags that need to be provided in order for the tests to work.  The flags are:

* toURL - The URL to Traffic Ops.  Default is "http://localhost:3000".
* toUser - The Traffic Ops user to use.  Default is "admin".
* toPass - The password of the user provided.  Deafault is "password".

Example command to run the tests: `go test -v -toUrl=https://to.kabletown.net -toUser=myUser -toPass=myPass`

*It can take serveral minutes for the integration tests to complete, so using the `-v` flag is recommended to see progress.*
