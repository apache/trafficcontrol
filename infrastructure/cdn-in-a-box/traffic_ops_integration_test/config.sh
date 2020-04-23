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
envvars=( DB_NAME DB_FQDN DB_USER DB_USER_PASS DB_PORT TO_HOST TO_PORT TO_ADMIN_PASSWORD )
unset_vars=""
for v in "${envvars[@]}"; do
    if [[ -z "${!v}" ]]; then
        unset_vars="$unset_vars $v"
    fi
done
if [[ ! -z "$unset_vars" ]]; then
    echo "required env vars are unset:$unset_vars"
    exit 1
fi

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
        },
        "includeSystemTests": true
    },
    "trafficOps": {
        "URL": "https://$TO_HOST:$TO_PORT",
        "password": "$TO_ADMIN_PASSWORD",
        "users": {
            "disallowed": "disallowed",
            "operations": "operations",
            "admin": "admin",
            "federation": "federation",
            "portal": "portal",
            "readOnly": "readOnly",
            "extension": "extension"
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

