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
# The following environment variables must be set, ordinarily by `docker run -e` arguments:
# DB_SERVER
# DB_PORT
# DB_ROOT_PASS
# DB_USER
# DB_USER_PASS
# DB_NAME
# ADMIN_USER
# ADMIN_PASS
# CERT_COUNTRY
# CERT_STATE
# CERT_CITY
# CERT_COMPANY
# DOMAIN

# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS CERT_COUNTRY CERT_STATE CERT_CITY CERT_COMPANY DOMAIN)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

cat <<-EOF >/opt/traffic_ops/app/conf/cdn.conf
{
    "hypnotoad" : {
        "listen" : [
            "https://[::]:60443?cert=$crt&key=$key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED"
        ],
        "user" : "trafops",
        "group" : "trafops",
        "heartbeat_timeout" : 20,
        "pid_file" : "/var/run/traffic_ops.pid",
        "workers" : 12
    },
    "traffic_ops_golang" : {
        "port" : "443",
        "proxy_timeout" : 60,
        "proxy_keep_alive" : 60,
        "proxy_tls_timeout" : 60,
        "proxy_read_header_timeout" : 60,
        "read_timeout" : 60,
        "read_header_timeout" : 60,
        "write_timeout" : 60,
        "idle_timeout" : 60,
        "log_location_error": "/var/log/traffic_ops/error.log",
        "log_location_warning": "/var/log/traffic_ops/warning.log",
        "log_location_info": "/var/log/traffic_ops/info.log",
        "log_location_debug": "/var/log/traffic_ops/debug.log",
        "log_location_event": "/var/log/traffic_ops/access.log",
        "max_db_connections": 20,
        "backend_max_connections": {
            "mojolicious": 4
        }
    },
    "cors" : {
        "access_control_allow_origin" : "*"
    },
    "to" : {
        "base_url" : "http://localhost:3000",
        "email_from" : "no-reply@traffic-ops-domain.com",
        "no_account_found_msg" : "A Traffic Ops user account is required for access. Please contact your Traffic Ops user administrator."
    },
    "portal" : {
        "base_url" : "http://localhost:8080/!#/",
        "email_from" : "no-reply@traffic-portal-domain.com",
        "pass_reset_path" : "user",
        "user_register_path" : "user"
    },
    "secrets" : [
        "mONKEYDOmONKEYSEE."
    ],
    "geniso" : {
        "iso_root_path" : "/opt/traffic_ops/app/public"
    },
    "inactivity_timeout" : 60
}
EOF

cat <<-EOF >/opt/traffic_ops/app/conf/production/database.conf
{
        "description": "Local PostgreSQL database on port 5432",
        "dbname": "$DB_NAME",
        "hostname": "$DB_SERVER",
        "user": "$DB_USER",
        "password": "$DB_USER_PASS",
        "port": "$DB_PORT",
        "ssl": false,
        "type": "Pg"
}
EOF

cat <<-EOF >/opt/traffic_ops/app/db/dbconf.yml
version: "1.0"
name: dbconf.yml

production:
  driver: postgres
  open: host=$DB_SERVER port=$DB_PORT user=$DB_USER password=$DB_USER_PASS dbname=$DB_NAME sslmode=disable
EOF

key=/server.key
crt=/server.crt

openssl req -newkey rsa:2048 -nodes -keyout $key -x509 -days 365 -out $crt -subj "/C=$CERT_COUNTRY/ST=$CERT_STATE/L=$CERT_CITY/O=$CERT_COMPANY"
chown trafops:trafops $key $crt
