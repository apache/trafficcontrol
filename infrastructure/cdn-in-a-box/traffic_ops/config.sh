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
# TP_HOST
# TV_DB_NAME
# TV_DB_PORT
# TV_DB_SERVER
# TV_DB_USER
# TV_DB_USER_PASS
#
# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS DOMAIN TO_HOST TO_PORT TP_HOST TV_DB_NAME TV_DB_PORT TV_DB_SERVER TV_DB_USER TV_DB_USER_PASS)
for v in $envvars; do
  if [[ -z "${!v}" ]]; then echo "$v is unset"; exit 1; fi
done

until [[ -f "$X509_CA_ENV_FILE" ]]; do
  echo "Waiting on SSL certificate generation."
  sleep 2
done

# these expected to be stored in $X509_CA_ENV_FILE, but a race condition could render the contents
# blank until it gets sync'd.  Ensure vars defined before writing cdn.conf.
until [[ -v X509_GENERATION_COMPLETE && -n "$X509_GENERATION_COMPLETE" ]]; do
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

if [[ "$TO_DEBUG_ENABLE" == true ]]; then
  DEBUGGING_TIMEOUT=$(( 60 * 60 * 24 )); # Timing out debugging after 1 day seems fair
fi;

cdn_conf=/opt/traffic_ops/app/conf/cdn.conf
>"$cdn_conf" echo "$(jq -s '.[0] * .[1]' "$cdn_conf" <(cat <<-EOF
{
    "disable_auto_cert_deletion": false,
    "use_ims": true,
    "server_update_status_cache_refresh_interval_sec": 0,
    "user_cache_refresh_interval_sec": 0,
    "role_based_permissions": true,
    "traffic_ops_golang" : {
          "traffic_vault_backend": "$TV_BACKEND",
          "traffic_vault_config": {
            "dbname": "$TV_DB_NAME",
            "hostname": "$TV_DB_SERVER.$INFRA_SUBDOMAIN.$TLD_DOMAIN",
            "user": "$TV_DB_USER",
            "password": "$TV_DB_USER_PASS",
            "port": ${TV_DB_PORT:-5432},
            "conn_max_lifetime_seconds": ${DEBUGGING_TIMEOUT:-60},
            "max_connections": 500,
            "max_idle_connections": 30,
            "query_timeout_seconds": ${DEBUGGING_TIMEOUT:-60},
            "aes_key_location": "$TV_AES_KEY_LOCATION"
        },
        "cert" : "$crt",
        "key" : "$key",
        "proxy_timeout" : ${DEBUGGING_TIMEOUT:-60},
        "proxy_tls_timeout" : ${DEBUGGING_TIMEOUT:-60},
        "proxy_read_header_timeout" : ${DEBUGGING_TIMEOUT:-60},
        "read_timeout" : ${DEBUGGING_TIMEOUT:-60},
        "read_header_timeout" : ${DEBUGGING_TIMEOUT:-60},
        "request_timeout" : ${DEBUGGING_TIMEOUT:-60},
        "write_timeout" : ${DEBUGGING_TIMEOUT:-60},
        "idle_timeout" : ${DEBUGGING_TIMEOUT:-60},
        "log_location_error": "$TO_LOG_ERROR",
        "log_location_warning": "$TO_LOG_WARNING",
        "log_location_info": "$TO_LOG_INFO",
        "log_location_debug": "$TO_LOG_DEBUG",
        "log_location_event": "$TO_LOG_EVENT",
        "db_conn_max_lifetime_seconds": ${DEBUGGING_TIMEOUT:-60},
        "db_query_timeout_seconds": ${DEBUGGING_TIMEOUT:-20}
    },
    "to" : {
        "email_from" : "no-reply@$INFRA_SUBDOMAIN.$TLD_DOMAIN"
    },
    "portal" : {
        "base_url" : "https://$TP_HOST.$INFRA_SUBDOMAIN.$TLD_DOMAIN/#!/",
        "email_from" : "no-reply@$INFRA_SUBDOMAIN.$TLD_DOMAIN"
    },
    "smtp" : {
        "enabled" : true,
        "address" : "${SMTP_FQDN}:${SMTP_PORT}"
    },
    "InfluxEnabled": true,
    "influxdb_conf_path": "/opt/traffic_ops/app/conf/production/influx.conf",
    "lets_encrypt" : {
        "environment": "staging"
    }
}
EOF
))"

<<RIAK_CONF cat >/opt/traffic_ops/app/conf/production/riak.conf
{
  "MaxTLSVersion": "1.1",
  "password": "$TV_RIAK_PASSWORD",
  "user": "$TV_RIAK_USER"
}
RIAK_CONF

<<INFLUX_CONF cat >/opt/traffic_ops/app/conf/production/influx.conf
{
  "password": "$INFLUXDB_ADMIN_PASSWORD",
  "secure": false,
  "user": "$INFLUXDB_ADMIN_USER"
}
INFLUX_CONF

install_bin=/opt/traffic_ops/install/bin
input_json="${install_bin}/input.json"
echo "$(jq "$(<<'JQ_FILTER' envsubst
  ."/opt/traffic_ops/app/conf/cdn.conf"[] |= (
    (select(.config_var == "base_url") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${TO_URL}"
    else . end))
  ) |
  ."/opt/traffic_ops/app/conf/production/database.conf"[] |= (
    (select(.config_var == "dbname") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${DB_NAME}"
    else . end)) |
    (select(.config_var == "hostname") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${DB_FQDN}"
    else . end)) |
    (select(.config_var == "user") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${DB_USER}"
    else . end)) |
    (select(.config_var == "password") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${DB_USER_PASS}"
    else . end))
  ) |
  ."/opt/traffic_ops/app/conf/production/tv.conf"[] |= (
    (select(.config_var == "dbname") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${TV_DB_NAME}"
    else . end)) |
    (select(.config_var == "hostname") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${DB_FQDN}"
    else . end)) |
    (select(.config_var == "user") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${TV_DB_USER}"
    else . end)) |
    (select(.config_var == "password") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${TV_DB_USER_PASS}"
    else . end))
  ) |
  ."/opt/traffic_ops/install/data/json/openssl_configuration.json"[] |= (
    (select(.config_var == "genCert") |= with_entries(if .key | test("^[A-Z]") then .value =
      "no"
    else . end)) |
    (select(.config_var == "pgPassword") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${DB_USER_PASS}"
    else . end))
  ) |
  ."/opt/traffic_ops/install/data/json/profiles.json"[] |= (
    (select(.config_var == "tm.url") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${TO_URL}"
    else . end)) |
    (select(.config_var == "cdn_name") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${CDN_NAME}"
    else . end)) |
    (select(.config_var == "dns_subdomain") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${CDN_SUBDOMAIN}.${TLD_DOMAIN}"
    else . end))
  ) |
  ."/opt/traffic_ops/install/data/json/users.json"[] |= (
    (select(.config_var == "tmAdminUser") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${TO_ADMIN_USER}"
    else . end)) |
    (select(.config_var == "tmAdminPw") |= with_entries(if .key | test("^[A-Z]") then .value =
      "${TO_ADMIN_PASSWORD}"
    else . end))
  )
JQ_FILTER
)" "$input_json")" >"$input_json"

"${install_bin}/postinstall" --debug -a --cfile "$input_json" -n --no-restart-to
