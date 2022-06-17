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

# Docker Integration Tests

The Docker Integration Tests are designed to run in a Docker container with Traffic Ops. Traffic Ops has the locations of all other CDN components.

Thus, an integration test for Traffic Ops should simply query `localhost`; an integration test for any other component should get that component's IP from Traffic Ops, and subsequently query that IP.

# Running the Integration Tests

To run the integration tests, set up Traffic Ops in Docker (see https://github.com/apache/trafficcontrol/tree/master/infrastructure/docker/README.md), and run `run-docker-integration-test.sh`.

# Creating a Test

To create a new integration test, create an executable file -- it can be a shell script, Go binary, Perl script, etc -- and place it in `/infrastructure/test/integration/docker-integration-tests/`.

Your executable should:
* Query Traffic Ops at localhost
* Query any other components at their IPs, which can be queried from Traffic Ops at `https://localhost/api/4.0/servers`
  * If you're not using a Traffic Ops client, you'll need to send a login cookie; for an example of how to get the cookie, see https://github.com/apache/trafficcontrol/tree/master/infrastructure/docker/traffic_ops/run.sh
* Return 0 and print nothing for success; return a nonzero code and print the error on failure

That's it! The `run-docker-integration-test.sh` script will automatically run your executable in the Docker Traffic Ops container, and report its failures.
