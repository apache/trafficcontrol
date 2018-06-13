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

# Script for running the Dockerfile for Traffic Ops.
# The Dockerfile sets up a Docker image which can be used for any new Traffic Ops container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
export GOROOT=/usr/local/go
export GOPATH=/opt/traffic_ops/go
PATH=$PATH:$GOPATH/bin:$GOROOT/bin

export PGOPTIONS='--client-min-messages=warning'
echo "PERL5LIB: $PERL5LIB"
export USER=root

cat <<-EOF >/opt/traffic_ops/app/conf/test/database.conf
{
        "description": "Local PostgreSQL database on port 5432",
        "dbname": "$DB_NAME",
        "hostname": "$DB_SERVER",
        "user": "$DB_USER",
        "password": "$DB_PASSWORD",
        "port": "$DB_PORT",
        "ssl": false,
        "type": "Pg"
}
EOF

cat <<-EOF >/opt/traffic_ops/app/db/dbconf.yml
version: "1.0"
name: dbconf.yml

test:
  driver: postgres
  open: host=$DB_SERVER port=$DB_PORT user=$DB_USER password=$DB_USER_PASS dbname=$DB_NAME sslmode=disable
EOF

cd /opt/traffic_ops/app
./db/admin.pl --env=test reset
