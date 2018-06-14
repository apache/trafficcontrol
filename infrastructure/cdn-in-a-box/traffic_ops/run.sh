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
# The following environment variables are used to configure the database and traffic ops.  They must be set
# in ../variables.env for docker-compose to pick up the values:
# 
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

# TODO:  Unused -- should be removed?  TRAFFIC_VAULT_PASS

# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS CERT_COUNTRY CERT_STATE CERT_CITY CERT_COMPANY DOMAIN)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

# Write config files
set -x
if [[ -r /config.sh ]]; then
	. /config.sh
fi

while ! nc $DB_SERVER $DB_PORT </dev/null; do # &>/dev/null; do
        echo "waiting for $DB_SERVER $DB_PORT"
        sleep 3
done

TO_DIR=/opt/traffic_ops/app
cat conf/production/database.conf

export PERL5LIB=$TO_DIR/lib:$TO_DIR/local/lib/perl5
export PATH=/usr/local/go/bin:/opt/traffic_ops/go/bin:$PATH
export GOPATH=/opt/traffic_ops/go

cd $TO_DIR && \
	./db/admin.pl --env=production reset && \
	./db/admin.pl --env=production seed || echo "db setup failed!"

/adduser.pl $TO_ADMIN_USER $TO_ADMIN_PASSWORD admin | psql -U$DB_USER -h$DB_SERVER $DB_NAME || echo "adding traffic_ops admin user failed!"

cd $TO_DIR && $TO_DIR/local/bin/hypnotoad script/cdn
exec tail -f /var/log/traffic_ops/traffic_ops.log
