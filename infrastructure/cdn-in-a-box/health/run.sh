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
############################################################

set -e
set -x
set -m

envvars=( TO_HOST TO_PORT TM_PORT TM_USER TM_PASSWORD )
for v in $envvars
do
	if [[ -z "${!v}" ]]; then echo "$v is unset"; exit 1; fi
done

set-dns.sh
insert-self-into-dns.sh

source /to-access.sh

# Wait on SSL certificate generation
until [[ -f "$X509_CA_ENV_FILE" ]]
do
     echo "Waiting on Shared SSL certificate generation"
     sleep 3
done

# If /tmp/trafficcontrol does not already exist when running t3c-apply, t3c-apply will create it and fail silently
mkdir -p /tmp/trafficcontrol

# Source the CIAB-CA shared SSL environment
until [[ -n "$X509_GENERATION_COMPLETE" ]]
do
  echo "Waiting on X509 vars to be defined"
  sleep 1
  source "$X509_CA_ENV_FILE"
done

# Trust the CIAB-CA at the System level
cp $X509_CA_CERT_FULL_CHAIN_FILE /etc/pki/ca-trust/source/anchors
update-ca-trust extract

# Enroll
TO_URL="https://$TO_FQDN:$TO_PORT"
to-enroll tchealthclient ALL || (while true; do echo "enroll failed... sad face :("; sleep 3 ; done)
tail -f /var/log/trafficcontrol/tc-health-client.log