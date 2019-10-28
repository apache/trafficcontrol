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
# TO_HOST
# TO_PORT
# TO_PERL_HOST
# TO_PERL_PORT
# TP_HOST
#
# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS DOMAIN TO_PERL_HOST TO_PERL_PORT TO_HOST TO_PORT TP_HOST)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

until [[ -f "$X509_CA_ENV_FILE" ]]
do
  echo "Waiting on SSL certificate generation."
  sleep 2
done

# these expected to be stored in $X509_CA_ENV_FILE, but a race condition could render the contents
# blank until it gets sync'd.  Ensure vars defined before writing cdn.conf.
until [[ -n "$X509_GENERATION_COMPLETE" ]]
do
  echo "Waiting on X509 vars to be defined"
  sleep 1
  source "$X509_CA_ENV_FILE"
done

# Add the CA certificate to sysem TLS trust store
cp $X509_CA_CERT_FULL_CHAIN_FILE /etc/pki/ca-trust/source/anchors
update-ca-trust extract

crt="$X509_INFRA_CERT_FILE"
key="$X509_INFRA_KEY_FILE"

echo "crt=$crt"
echo "key=$key"

cat <<-EOF >/opt/traffic_ops/app/conf/cdn.conf
{
    "hypnotoad" : {
        "listen" : [
            "https://$TO_PERL_FQDN:$TO_PERL_PORT?cert=$crt&key=$key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED"
        ],
        "user" : "trafops",
        "group" : "trafops",
        "heartbeat_timeout" : 20,
        "pid_file" : "/var/run/traffic_ops.pid",
        "workers" : 12
    },
    "traffic_ops_golang" : {
        "insecure": true,
        "port" : "$TO_PORT",
        "proxy_timeout" : 60,
        "proxy_keep_alive" : 60,
        "proxy_tls_timeout" : 60,
        "proxy_read_header_timeout" : 60,
        "read_timeout" : 60,
        "read_header_timeout" : 60,
        "write_timeout" : 60,
        "idle_timeout" : 60,
        "log_location_error": "$TO_LOG_ERROR",
        "log_location_warning": "$TO_LOG_WARNING",
        "log_location_info": "$TO_LOG_INFO",
        "log_location_debug": "$TO_LOG_DEBUG",
        "log_location_event": "$TO_LOG_EVENT",
        "max_db_connections": 20,
        "backend_max_connections": {
            "mojolicious": 4
        },
        "whitelisted_oauth_urls": [],
        "oauth_client_secret": ""
    },
    "cors" : {
        "access_control_allow_origin" : "*"
    },
    "to" : {
        "base_url" : "https://$TO_FQDN",
        "email_from" : "no-reply@$INFRA_SUBDOMAIN.$TLD_DOMAIN",
        "no_account_found_msg" : "A Traffic Ops user account is required for access. Please contact your Traffic Ops user administrator."
    },
    "portal" : {
        "base_url" : "https://$TP_HOST.$INFRA_SUBDOMAIN.$TLD_DOMAIN/#!/",
        "email_from" : "no-reply@$INFRA_SUBDOMAIN.$TLD_DOMAIN",
        "pass_reset_path" : "user",
        "user_register_path" : "user"
    },
    "secrets" : [
        "$TO_SECRET"
    ],
    "geniso" : {
        "iso_root_path" : "/opt/traffic_ops/app/public"
    },
    "inactivity_timeout" : 60,
    "smtp" : {
        "enabled" : false,
        "user" : "",
        "password" : "",
        "address" : ""
    },
    "InfluxEnabled": true,
    "influxdb_conf_path": "/opt/traffic_ops/app/conf/production/influx.conf"
}
EOF

cat <<-EOF >/opt/traffic_ops/app/conf/production/database.conf
{
        "description": "Local PostgreSQL database on port 5432",
        "dbname": "$DB_NAME",
        "hostname": "$DB_FQDN",
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
  open: host=$DB_FQDN port=$DB_PORT user=$DB_USER password=$DB_USER_PASS dbname=$DB_NAME sslmode=disable
test:
  driver: postgres
  open: host=$DB_FQDN port=$DB_PORT user=$DB_USER password=$DB_USER_PASS dbname=to_test sslmode=disable
EOF

cat <<-EOF >/opt/traffic_ops/app/conf/production/riak.conf
{     "user": "$TV_RIAK_USER",
  "password": "$TV_RIAK_PASSWORD"
}
EOF

cat <<-EOF >/opt/traffic_ops/app/conf/production/influx.conf
{
    "user": "$INFLUXDB_ADMIN_USER",
    "password": "$INFLUXDB_ADMIN_PASSWORD",
    "secure": false
}
EOF
