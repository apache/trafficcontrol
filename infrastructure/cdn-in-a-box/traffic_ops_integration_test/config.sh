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

# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS DOMAIN TO_PERL_HOST TO_PERL_PORT TO_HOST TO_PORT TP_HOST)
for v in $envvars
do
    if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

cat <<-EOF >/opt/integration/app/traffic-ops-test.conf
{
    "default": {
        "logLocations": {
            "debug": "null",
            "error": "null",
            "event": "null",
            "info": "null",
            "warning": "null"
        },
        "session": {
            "timeoutInSecs": 60
        }
    },
    "trafficOps": {
        "URL": "https://$TO_HOST:$TO_PORT",
        "password": "$ADMIN_PASS",
        "users": {
            "disallowed": "disallowed",
            "operations": "operations",
            "admin": "admin",
            "federation": "federation",
            "portal": "portal",
            "readOnly": "readOnly"
        }
    },
    "trafficOpsDB": {
        "description": "Local PostgreSQL database on port 5432",
        "dbname": "$DB_NAME",
        "hostname": "$DB_FQDN",
        "user": "$DB_USER",
        "password": "$DB_USER_PASS",
        "port": "$DB_PORT",
        "ssl": false,
        "type": "Pg"
    }
}
EOF

