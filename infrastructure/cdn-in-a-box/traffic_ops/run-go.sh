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
/config.sh

# Source the CIAB-CA shared SSL environment
until [[ -n "$X509_GENERATION_COMPLETE" ]]
do
  echo "Waiting on X509 vars to be defined"
  sleep 1
  source "$X509_CA_ENV_FILE"
done

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
demo1_sslkeys_verified=false
demo1_version=1
while [[ "$demo1_sslkeys_verified" = false ]]; do
   while true; do
     sslkeys_response=$(to-get "api/1.4/deliveryservices/xmlId/$ds_name/sslkeys?decode=true")
     echo "CDN SSLKeys=$sslkeys_response"
     [[ -n "$sslkeys_response" ]] && break
     sleep 2
   done
   demo1_crt="$(sed -n -e '/-----BEGIN CERTIFICATE-----/,$p' $X509_DEMO1_CERT_FILE | jq -s -R '.')"
   demo1_csr="$(sed -n -e '/-----BEGIN CERTIFICATE REQUEST-----/,$p' $X509_DEMO1_REQUEST_FILE | jq -s -R '.')"
   demo1_key="$(sed -n -e '/-----BEGIN PRIVATE KEY-----/,$p' $X509_DEMO1_KEY_FILE | jq -s -R '.')"
   demo1_json_request=$(jq -n \
                           --arg     cdn        "$CDN_NAME" \
                           --arg     hostname   "*.demo1.mycdn.ciab.test" \
                           --arg     dsname     "$ds_name" \
                           --argjson crt        "$demo1_crt" \
                           --argjson csr        "$demo1_csr" \
                           --argjson key        "$demo1_key" \
                           --argjson version    $demo1_version \
                          "{ cdn: \$cdn, 
                             certificate: { 
                               crt: \$crt, 
                               csr: \$csr,
                               key: \$key 
                             },
                             deliveryservice: \$dsname,
                             hostname: \$hostname,
                             key: \$dsname,
                             version: $demo1_version 
                          }")

   demo1_json_response=$(to-post 'api/1.4/deliveryservices/sslkeys/add' "$demo1_json_request")

   if [[ -n "$demo1_json_response" ]] ; then 
      sleep 2
      cdn_sslkeys_response=$(to-get "api/1.3/cdns/name/$CDN_NAME/sslkeys.json" | jq '.response[] | length')
      echo "cdn_sslkeys_response=$cdn_sslkeys_response"

      if [ -n "$cdn_sslkeys_response" ] ; then 
         if ((cdn_sslkeys_response==0)); then 
           sleep 2 # Submit it again because the first time doesn't work !
           demo1_json_response=$(to-post 'api/1.4/deliveryservices/sslkeys/add' "$demo1_json_request")
 
           if [[ -n "$demo1_json_response" ]] ; then 
              demo1_sslkeys_verified=true
           fi
        elif ((cdn_sslkeys_response>0)); then
           demo1_sslkeys_verified=true
        fi
      fi
   fi

   ((demo_version+=1)) 
done

### Automatic Queue/Snapshot ###
while [[ "$AUTO_SNAPQUEUE_ENABLED" = true ]] ; do
  # AUTO_SNAPQUEUE_SERVERS should be a comma delimited list of expected docker service names to be enrolled - see varibles.env
  expected_servers_json=$(echo "$AUTO_SNAPQUEUE_SERVERS" | tr ',' '\n' | jq -R . | jq -M -c -e -s '.|sort')
  expected_servers_list=$(jq -r -n --argjson expected "$expected_servers_json" '$expected|join(",")')
  expected_servers_total=$(jq -r -n --argjson expected "$expected_servers_json" '$expected|length')

  current_servers_json=$(to-get 'api/1.4/servers' 2>/dev/null | jq -c -e '[.response[] | .xmppId] | sort')
  [ -z "$current_servers_json" ] && current_servers_json='[]'
  current_servers_list=$(jq -r -n --argjson current "$current_servers_json" '$current|join(",")')
  current_servers_total=$(jq -r -n --argjson current "$current_servers_json" '$current|length')

  remain_servers_json=$(jq -n --argjson expected "$expected_servers_json" --argjson current "$current_servers_json" '$expected-$current')
  remain_servers_list=$(jq -r -n --argjson remain "$remain_servers_json" '$remain|join(",")')
  remain_servers_total=$(jq -r -n --argjson remain "$remain_servers_json" '$remain|length')

  echo "AUTO-SNAPQUEUE - Expected Servers ($expected_servers_total): $expected_servers_list"
  echo "AUTO-SNAPQUEUE - Current Servers ($current_servers_total): $current_servers_list"
  echo "AUTO-SNAPQUEUE - Remain Servers ($remain_servers_total): $remain_servers_list"

  if ((remain_servers_total == 0)) ; then
     echo "AUTO-SNAPQUEUE - All expected servers enrolled."
     sleep $AUTO_SNAPQUEUE_ACTION_WAIT
     echo "AUTO-SNAPQUEUE - Do automatic snapshot..."
     to-put 'api/1.3/cdns/2/snapshot'
     sleep $AUTO_SNAPQUEUE_ACTION_WAIT
     echo "AUTO-SNAPQUEUE - Do queue updates..."
     to-post 'api/1.3/cdns/2/queue_update' '{"action":"queue"}'
     break
  fi

  sleep $AUTO_SNAPQUEUE_POLL_INTERVAL
done

exec tail -f /dev/null
