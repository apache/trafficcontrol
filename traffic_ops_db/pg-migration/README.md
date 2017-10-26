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

# Converts existing mysql `traffic_ops` database to postgres

* Requires a fairly recent ( 17.05.0-ce ) version of `docker-engine` and `docker-compose`.

* Modify the mysql-to-postgres.env file for the parameters in your Traffic Ops environment

* Ensure that your new Postgres service is running (local or remote)

* A sample Postgres Docker container has been provided for testing
  1. `cd ../docker`
  2. `$ sh todb.sh run` - to download/start your Postgres Test Docker container
  3. `$ sh todb.sh setup` - to create your new 'traffic_ops' role and 'traffic_ops' database

* Run the Mysql to Postgres Migration Docker flow
  1. `$ cd ../pg-migration`
  2. `$ sh migrate.sh`
