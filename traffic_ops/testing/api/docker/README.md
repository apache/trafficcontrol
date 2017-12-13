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

# Traffic Ops API Tests in Docker

The Traffic Ops Client API tests can also be run from Docker if necessary.  The only major difference is the TO API Test tool for Docker needs to be driven with 
environment variables. By default the docker-compose.yml will look for it's default environment variables in the `docker/traffic-ops-test.env`

When necessary the environment variables can be overridden via the command line docker-compose option as well like this:

  $ docker-compose run -e TODB_HOSTNAME=your_traffic_ops_db_hostname -e TO_URL=https://your_traffic_ops_ip:8443 api_tests

