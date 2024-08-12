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

# PostgreSQL Traffic Vault Tests

This test is designed to test all the DB migrations for the PostgreSQL Traffic
Vault backend, optionally including a starting DB dump, as well as test the
`reencrypt` tool.

NOTE: this test is similar to the TODB test found in traffic_ops_db/test/docker.

## Running the test

1. Build a Traffic Ops rpm (from the project root): `./pkg traffic_ops_build`
2. Copy the rpm to ./traffic_ops.rpm
3. Optional: place a DB dump file in ./initdb.d (file name must end in "dump")
4. `docker compose build`
5. `docker compose up --exit-code-from trafficvault-db-admin`

## Notes about data directory and AES keys

The encrypted data in the ./data directory was taken from cdn-in-a-box and was
encrypted using the ./aes.key file. If the data needs to change for any reason
for test updates (if for instance the schema was updated), it will have to be
encrypted, and whatever key was used for encryption should replace ./aes.key.
Note that ./new-aes.key is only used for testing the `reencrypt` tool, so this
test key will probably never need to be changed.
