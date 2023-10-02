#!/usr/bin/env bash
set -x
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
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS TO_ADMIN_USER TO_ADMIN_PASS CERT_COUNTRY CERT_STATE CERT_CITY CERT_COMPANY TO_DOMAIN)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

start() {
  traffic_ops_golang_command=(./bin/traffic_ops_golang -cfg "$CDNCONF" -dbcfg "$DATABASECONF");
  "${traffic_ops_golang_command[@]}" &
	exec tail -f $TO_LOG
}

# generates and saves SSL certificates and database config files.
init() {
  # install certificates for TO
  openssl req -newkey rsa:2048 -nodes -keyout /etc/pki/tls/private/localhost.key -x509 -days 365 \
    -out /etc/pki/tls/certs/localhost.crt -subj "/C=$CERT_COUNTRY/ST=$CERT_STATE/L=$CERT_CITY/O=$CERT_COMPANY"
  cp /etc/pki/tls/certs/localhost.crt /etc/pki/tls/certs/ca-bundle.crt

  # update the base_url in cdn.conf
  sed -i -e "s#http://localhost\:3000#http://${TO_HOSTNAME}\:443#" $CDNCONF
	sed -i -e 's#https://\[::\]#https://127\.0\.0\.1#' $CDNCONF
  sed -i -e 's#"use_ims": false,#"use_ims": true,#' $CDNCONF
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
  open: host=$DB_SERVER port=5432 user=traffic_ops password=$DB_USER_PASS dbname=to_development sslmode=disable

test:
  driver: postgres
  open: host=$DB_SERVER port=5432 user=traffic_ops password=$DB_USER_PASS dbname=to_test sslmode=disable

integration:
  driver: postgres
  open: host=$DB_SERVER port=5432 user=traffic_ops password=$DB_USER_PASS dbname=to_integration sslmode=disable

production:
  driver: postgres
  open: host=$DB_SERVER port=5432 user=traffic_ops password=$DB_USER_PASS dbname=traffic_ops sslmode=disable
EOM

  touch $LOG_DEBUG $LOG_ERROR $LOG_EVENT $LOG_INFO $LOG_WARN $TO_LOG
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi

# create the 'traffic_ops' database, tables and runs migrations
pushd /opt/traffic_ops/app;
db/admin --env=production reset > /admin.log 2>&1
db/admin --env=production upgrade >> /admin.log 2>&1
db/admin --env=production seed >> /admin.log 2>&1
popd;

# start traffic_ops
start
