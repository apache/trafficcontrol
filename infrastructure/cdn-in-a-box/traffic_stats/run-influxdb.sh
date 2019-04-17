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

set -e
set -x
set -m

set-dns.sh
insert-self-into-dns.sh

source /to-access.sh

# Wait on SSL certificate generation
until [[ -f "$X509_CA_ENV_FILE" ]]
do
  echo "Waiting on Shared SSL certificate generation"
  sleep 3
done

# Source the CIAB-CA shared SSL environment
until [[ -n "$X509_GENERATION_COMPLETE" ]]
do
  echo "Waiting on X509 vars to be defined"
  sleep 1
  source "$X509_CA_ENV_FILE"
done

# Copy the CIAB-CA certificate to here so it can be added to the trust store
cp "$X509_CA_CERT_FULL_CHAIN_FILE" /usr/local/share/ca-certificates
update-ca-certificates

# Traffic Ops must be accepting connections before enroller can start
until nc -z $TO_FQDN $TO_PORT </dev/null >/dev/null && to-ping; do
  echo "Waiting for $TO_URL"
  sleep 5
done

to-enroll influxdb ALL "" "$INFLUXDB_PORT" "$INFLUXDB_PORT" || (while true; do echo "enroll failed."; sleep 3 ; done)

exec /entrypoint.sh influxd
