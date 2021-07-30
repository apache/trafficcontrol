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

# script for generating the reencrypt config (given environment variables)

# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_USER DB_USER_PASS DB_NAME DB_PORT )
for v in $envvars
do
	if [[ -z "${!v}" ]]; then echo "$v is unset"; exit 1; fi
done

cat <<-EOF >/opt/traffic_ops/app/db/reencrypt/reencrypt.conf
{
    "dbname": "$DB_NAME",
    "hostname": "$DB_SERVER",
    "user": "$DB_USER",
    "password": "$DB_USER_PASS",
    "port": $DB_PORT,
    "ssl": false
}
EOF
