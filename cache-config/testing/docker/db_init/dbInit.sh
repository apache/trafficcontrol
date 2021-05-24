#!/usr/bin/env bash
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

############################################################
# Script for creating the database user account for traffic
# ops. 
# Used while the Docker Image is initializing itself
############################################################

while ! nc $DB_SERVER $DB_PORT </dev/null; do # &>/dev/null; do
        echo "waiting for $DB_SERVER:$DB_PORT"
        sleep 3
done
psql -h $DB_SERVER -U postgres -c "CREATE USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_USER_PASS'"
createdb $DB_NAME -h $DB_SERVER -U postgres --owner $DB_USER
