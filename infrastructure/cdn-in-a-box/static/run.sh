#!/usr/bin/env bash
#
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

set -eux

. /to-access.sh
set-dns.sh
insert-self-into-dns.sh

# Wait on SSL certificate generation
until [[ -f "$X509_CA_ENV_FILE" ]]
do
     echo "Waiting on Shared SSL certificate generation"
     sleep 3
done

# Source the CIAB-CA shared SSL environment
until [[ -n "${X509_GENERATION_COMPLETE:-}" ]]; do
  echo "Waiting on X509 vars to be defined"
  sleep 1
  source "$X509_CA_ENV_FILE"
done

source /to-access.sh
cat "$X509_INFRA_KEY_FILE" "$X509_INFRA_CERT_FILE" > "/etc/lighttpd/${INFRA_FQDN}.pem"

conf_file=/etc/lighttpd/lighttpd.conf
echo "$(<"$conf_file" envsubst)" > "$conf_file"
lighttpd -t -f "$conf_file"
exec lighttpd -D -f "$conf_file"
