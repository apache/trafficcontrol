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
# TO_ADMIN_USER
# TO_ADMIN_PASS
# CERT_COUNTRY
# CERT_STATE
# CERT_CITY
# CERT_COMPANY
# TO_DOMAIN
# TRAFFIC_VAULT_PASS

# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS DB_NAME TO_ADMIN_USER TO_ADMIN_PASS CERT_COUNTRY CERT_STATE CERT_CITY CERT_COMPANY TO_DOMAIN TV_DB_NAME TV_SERVER TV_PORT TV_USER TV_USER_PASS TV_AES_KEY_LOCATION)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

start() {
  traffic_ops_golang_command=(/opt/traffic_ops/app/bin/traffic_ops_golang -cfg "$CDNCONF" -dbcfg "$DATABASECONF" );
  "${traffic_ops_golang_command[@]}" &
	exec tail -f $TO_LOG
}

# generates and saves SSL certificates, database (TO and TV) config files.
init() {
  # install certificates for TO
  openssl req -newkey rsa:2048 -nodes -keyout /etc/pki/tls/private/localhost.key -x509 -days 365 \
    -out /etc/pki/tls/certs/localhost.crt -subj "/C=$CERT_COUNTRY/ST=$CERT_STATE/L=$CERT_CITY/O=$CERT_COMPANY"
  cp /etc/pki/tls/certs/localhost.crt /etc/pki/tls/certs/ca-bundle.crt
  openssl rand -base64 32 > $TV_AES_KEY_LOCATION
  chown trafops:trafops $TV_AES_KEY_LOCATION

  # update the base_url in cdn.conf
  sed -i -e "s#http://localhost\:3000#http://${TO_HOSTNAME}\:443#" $CDNCONF
	sed -i -e 's#https://\[::\]#https://127\.0\.0\.1#' $CDNCONF
  sed -i -e 's#"use_ims": false,#"use_ims": true,#' $CDNCONF
  sed -i -e 's#"traffic_vault_backend": "",#"traffic_vault_backend": "postgres",#' $CDNCONF
  sed -i -e "s#\"traffic_vault_config\": {},#\"traffic_vault_config\": {\n            \"dbname\": \"$TV_DB_NAME\",\n            \"hostname\": \"$TV_SERVER\",\n            \"user\": \"$TV_USER\",\n            \"password\": \"$TV_USER_PASS\",\n            \"port\": $TV_PORT,\n            \"ssl\": false,\n            \"conn_max_lifetime_seconds\": 60,\n            \"max_connections\": 20,\n            \"max_idle_connections\": 15,\n            \"query_timeout_seconds\": 10,\n            \"aes_key_location\": \"$TV_AES_KEY_LOCATION\"\n        },#" $CDNCONF
  #
  cat > $DATABASECONF << EOM
{
  "type" : "Pg",
  "description" : "Pg database on localhost:5432",
  "port" : "$DB_PORT",
  "dbname" : "$DB_NAME",
  "password" : "$DB_USER_PASS",
  "hostname" : "$DB_SERVER",
  "user" : "$DB_USER"
}
EOM

  cat > $DBCONF << EOM
version: "1.0"
name: dbconf.yml

development:
  driver: postgres
  open: host=$DB_SERVER port=$DB_PORT user=$DB_USER password=$DB_USER_PASS dbname=to_development sslmode=disable

test:
  driver: postgres
  open: host=$DB_SERVER port=$DB_PORT user=$DB_USER password=$DB_USER_PASS dbname=to_test sslmode=disable

integration:
  driver: postgres
  open: host=$DB_SERVER port=$DB_PORT user=$DB_USER password=$DB_USER_PASS dbname=to_integration sslmode=disable

production:
  driver: postgres
  open: host=$DB_SERVER port=$DB_PORT user=$DB_USER password=$DB_USER_PASS dbname=$DB_NAME sslmode=disable
EOM


  cat > $TVCONF << EOM
version: "1.0"
name: dbconf.yml

development:
  driver: postgres
  open: host=$TV_SERVER port=$TV_PORT user=$TV_USER password=$TV_USER_PASS dbname=tv_development sslmode=disable

test:
  driver: postgres
  open: host=$TV_SERVER port=$TV_PORT user=$TV_USER password=$TV_USER_PASS dbname=tv_test sslmode=disable

integration:
  driver: postgres
  open: host=$TV_SERVER port=$TV_PORT user=$TV_USER password=$TV_USER_PASS dbname=tv_integration sslmode=disable

production:
  driver: postgres
  open: host=$TV_SERVER port=$TV_PORT user=$TV_USER password=$TV_USER_PASS dbname=$TV_DB_NAME sslmode=disable
EOM

  touch $LOG_DEBUG $LOG_ERROR $LOG_EVENT $LOG_INFO $LOG_WARN $TO_LOG
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi

i=0
sleep_time=3
while ! nc $DB_SERVER $DB_PORT </dev/null; do
  echo "waiting for $DB_SERVER:$DB_PORT" >> /var/log/traffic_ops/to_admin.log
  sleep $sleep_time
  let i=i+1
  if [ $i -gt 10 ]; then
    let d=i*sleep_time
    echo "$DB_SERVER:$DB_PORT is unavailable after $d seconds, giving up" >> /var/log/traffic_ops/to_admin.log
    exit 1
  fi
done

(
set -o errexit
# create the 'traffic_ops' database, tables and runs migrations
mkdir -p /var/log/traffic_ops/
touch /var/log/traffic_ops/to_admin.log /var/log/traffic_ops/tv_admin.log
tail -f /var/log/traffic_ops/to_admin.log /var/log/traffic_ops/tv_admin.log &
cd /opt/traffic_ops/app
(
db/admin --env=production reset
db/admin --env=production upgrade
db/admin --env=production seed
db/admin --trafficvault --env=production reset
db/admin --trafficvault --env=production upgrade
) >> /var/log/traffic_ops/to_admin.log 2>&1) || {
	echo Failed to run migrations;
	exit 1
}

# start traffic_ops
start
