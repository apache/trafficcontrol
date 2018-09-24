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

# Source to-access functions and FQDN vars
source /to-access.sh

until [ -f "$X509_CA_DONE_FILE" ] ; do
   echo "Waiting on SSL certificate generation."
   sleep 2
done

# Write config files
if [[ -x /config.sh ]]; then
	/config.sh
fi

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

### Workaround: Start DeliveryService and Edge association
while true; do
  edge_name="$(to-get 'api/1.3/servers/hostname/edge/details' 2>/dev/null | jq -r -c '.response|.hostName')"
  ds_name="$(to-get 'api/1.3/deliveryservices' 2>/dev/null | jq -r -c '.response[].xmlId')"

  if [ -n "$edge_name" ] && [ "$ds_name" ] ; then
    tmp_file=$(mktemp)
    echo "{ \"xmlId\" : \"$ds_name\", \"serverNames\": [ \"$edge_name\" ] }" > $tmp_file
    cp $tmp_file /shared/enroller/deliveryservice_servers/
    break
  else 
    echo "Waiting for delivery service and edge to exist..."
  fi

  sleep 2
done

while true; do
  echo "Verifying that edge was associated to delivery service..."

  edge_name="$(to-get 'api/1.3/servers/hostname/edge/details' 2>/dev/null | jq -r -c '.response|.hostName')"
  ds_name="$(to-get 'api/1.3/deliveryservices' 2>/dev/null | jq -r -c '.response[].xmlId')"
  ds_id="$(to-get 'api/1.3/deliveryservices' 2>/dev/null | jq -r -c '.response[].id')"
  edge_verify=$(to-get "/api/1.2/deliveryservices/$ds_id/servers" | jq -r '.response[]|.hostName')

  if [[ $edge_verify = $edge_name ]] ; then
    break
  fi

  sleep 2
done

# Snapshot the CDN 
until [ $(to-get '/CRConfig-Snapshots/CDN-in-a-Box/CRConfig.json' 2>/dev/null | jq -c -e '.config|length') -gt 0 ] ; do 
  sleep 2
  echo "Do Snapshot for CDN-in-a-Box..."; 
  to-put 'api/1.3/cdns/2/snapshot'
done

# Queue Updates 
sleep 1
to-post 'api/1.3/cdns/2/queue_update' '{"action":"queue"}'

### Workaround: End DeliveryService and Edge association

exec tail -f /dev/null
