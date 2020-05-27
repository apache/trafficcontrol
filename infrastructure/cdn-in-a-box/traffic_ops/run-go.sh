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

# Setting the monitor shell option enables job control, which we need in order
# to bring traffic_ops_golang back to the foreground.
set -o xtrace -o monitor;

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

(
maxtries=10
for ((tries = 0; tries < maxtries; tries++)); do
        if nc -zvw2 "$SMTP_FQDN" "$SMTP_PORT"; then
          echo "${SMTP_FQDN}:${SMTP_PORT} was found."
          break;
        fi;
        echo "waiting for ${SMTP_FQDN}:${SMTP_PORT}"
        sleep 3
done
if (( tries == maxtries )); then
  echo "SMTP service was not found at ${SMTP_FQDN}:${SMTP_PORT} after ${maxtries} tries. Skipping..."
fi
)

cd /opt/traffic_ops/app

CDNCONF=/opt/traffic_ops/app/conf/cdn.conf
DBCONF=/opt/traffic_ops/app/conf/production/database.conf
RIAKCONF=/opt/traffic_ops/app/conf/production/riak.conf
mkdir -p /var/log/traffic_ops
touch /var/log/traffic_ops/traffic_ops.log

# enroll in the background so traffic_ops_golang can run in foreground
TO_USER=$TO_ADMIN_USER TO_PASSWORD=$TO_ADMIN_PASSWORD to-enroll $(hostname -s) &

traffic_ops_golang_command=(./bin/traffic_ops_golang -cfg "$CDNCONF" -dbcfg "$DBCONF" -riakcfg "$RIAKCONF");
if [[ "$TO_DEBUG_ENABLE" == true ]]; then
  traffic_ops_golang_command=(dlv '--accept-multiclient' '--continue' '--listen=:2345' '--headless=true' '--api-version=2' exec
    "${traffic_ops_golang_command[0]}" -- "${traffic_ops_golang_command[@]:1}");
fi;
"${traffic_ops_golang_command[@]}" &

to-enroll "to" ALL || (while true; do echo "enroll failed."; sleep 3 ; done)

while true; do
  echo "Verifying that edge was associated to delivery service..."

  edge_name="$(to-get 'api/2.0/servers?hostName=edge' 2>/dev/null | jq -r -c '.response[0]|.hostName')"
  ds_name=$(to-get 'api/2.0/deliveryservices' 2>/dev/null | jq -r -c '.response[] | select(.cdnName == "'"$CDN_NAME"'").xmlId')
  ds_id=$(to-get 'api/2.0/deliveryservices' 2>/dev/null | jq -r -c '.response[] | select(.cdnName == "'"$CDN_NAME"'").id')
  edge_verify=$(to-get "/api/2.0/deliveryservices/$ds_id/servers" | jq -r '.response[]|.hostName')

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

fg '"${traffic_ops_golang_command[@]}"'; # Bring traffic_ops_golang to foreground
fg; # Bring to-enroll to foreground if it is still running

if [[ "$TO_DEBUG_ENABLE" == true ]]; then
  tail -f /dev/null;
fi;
