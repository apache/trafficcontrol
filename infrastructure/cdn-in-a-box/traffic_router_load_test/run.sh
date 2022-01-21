#!/bin/bash
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

# Check that env vars are set
envvars=( DB_SERVER DB_PORT DB_ROOT_PASS DB_USER DB_USER_PASS ADMIN_USER ADMIN_PASS)
set -ex
for v in $envvars
do
	if [[ -z "${!v}" ]]; then echo "$v is unset"; exit 1; fi
done

set-dns.sh
insert-self-into-dns.sh
source /to-access.sh

# Source the CIAB-CA shared SSL environment
until [[ -v 'X509_GENERATION_COMPLETE' ]]; do
	echo 'Waiting on X509 vars to be defined'
	sleep 1
	if [[ ! -e "$X509_CA_ENV_FILE" ]]; then
		continue
	fi
	source "$X509_CA_ENV_FILE"
done

# Copy the CIAB-CA certificate to the ca-certificates directory so it can be added to the trust store
cp "$X509_CA_CERT_FULL_CHAIN_FILE" /usr/local/share/ca-certificates/
update-ca-certificates

./ultimate-test-harness.test
