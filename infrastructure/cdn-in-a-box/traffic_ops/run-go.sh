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
# DOMAIN

# TODO:  Unused -- should be removed?  TRAFFIC_VAULT_PASS

# Check that env vars are set

set -x
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS)
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

set-dns.sh
insert-self-into-dns.sh

# Source to-access functions and FQDN vars
source /to-access.sh

# Write config files
. /config.sh

while ! nc "$TO_PERL_FQDN" $TO_PERL_PORT </dev/null 2>/dev/null; do
        echo "waiting for $TO_PERL_FQDN:$TO_PERL_PORT"
        sleep 3
done

cd /opt/traffic_ops/app

CDNCONF=/opt/traffic_ops/app/conf/cdn.conf
DBCONF=/opt/traffic_ops/app/conf/production/database.conf
RIAKCONF=/opt/traffic_ops/app/conf/production/riak.conf
mkdir -p /var/log/traffic_ops
touch /var/log/traffic_ops/traffic_ops.log

# enroll in the background so traffic_ops_golang can run in foreground
TO_USER=$TO_ADMIN_USER TO_PASSWORD=$TO_ADMIN_PASSWORD to-enroll $(hostname -s) &

./bin/traffic_ops_golang -cfg $CDNCONF -dbcfg $DBCONF -riakcfg $RIAKCONF &

to-enroll "to" ALL || (while true; do echo "enroll failed."; sleep 3 ; done)

while true; do
  echo "Verifying that edge was associated to delivery service..."

  edge_name="$(to-get 'api/1.3/servers/hostname/edge/details' 2>/dev/null | jq -r -c '.response|.hostName')"
  ds_name=$(to-get 'api/1.3/deliveryservices' 2>/dev/null | jq -r -c '.response[] | select(.cdnName == "'"$CDN_NAME"'").xmlId')
  ds_id=$(to-get 'api/1.3/deliveryservices' 2>/dev/null | jq -r -c '.response[] | select(.cdnName == "'"$CDN_NAME"'").id')
  edge_verify=$(to-get "/api/1.2/deliveryservices/$ds_id/servers" | jq -r '.response[]|.hostName')

  if [[ $edge_verify = $edge_name ]] ; then
    break
  fi

  sleep 2
done

### Add SSL keys for demo1 delivery service
until [[ -s "$X509_DEMO1_CERT_FILE" && -s "$X509_DEMO1_REQUEST_FILE" && -s "$X509_DEMO1_KEY_FILE" ]]
do
  echo "Waiting on X509_DEMO1 files to exist"
  sleep 3
  source "$X509_CA_ENV_FILE"
done
to-add-sslkeys $CDN_NAME $ds_name "*.demo1.mycdn.ciab.test" $X509_DEMO1_CERT_FILE $X509_DEMO1_REQUEST_FILE $X509_DEMO1_KEY_FILE

### Automatic Queue/Snapshot ###
if [[ "$AUTO_SNAPQUEUE_ENABLED" = true ]]; then
  # AUTO_SNAPQUEUE_SERVERS should be a comma delimited list of expected docker service names to be enrolled - see varibles.env
  to-auto-snapqueue $AUTO_SNAPQUEUE_SERVERS $CDN_NAME
fi

exec tail -f /dev/null
