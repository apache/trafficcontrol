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

set -eu

start_time=$(date +%s)

source to-access.sh

set-dns.sh
insert-self-into-dns.sh

while ! to-ping; do
    echo waiting for trafficops
    sleep 3
done

xmlID="$(<<<"$DS_HOSTS" sed 's/ .*//g')" # Only get the first xmlID
while true; do
    sleep 3
    if ! exampleURLs=($(to-get "/api/${TO_API_VERSION}/deliveryservices?xmlId=${xmlID}" | jq -r '.response[].exampleURLs[]')) ||
      [[ "${#exampleURLs[@]}" -eq 0 ]]; then
        echo waiting for delivery service ${xmlID} example URLs
        continue
    fi
    echo "example URLs: '${exampleURLs[*]}'"

    success="true"
    for u in "${exampleURLs[@]}"; do
        if ! status=$(curl -Lkvs --connect-timeout 2 -m5 -o /dev/null -w "%{http_code}" "$u") ||
          [[ "$status" -ne 200 ]]; then
            echo "failed to curl delivery service example URL '$u' got status code '$status'"
            success="false"
            break
        fi
        echo "successfully curled delivery service example URL '$u'"
    done
    if [[ "$success" == "true" ]]; then
        echo "successfully curled all delivery service example URLs '${exampleURLs[*]}'"
        break
    fi
done
xmlID="$(<<<"$DS_HOSTS" sed -E 's/.* (\w+) .*/\1/g')" # Only get the second xmlID
while true; do
  sleep 3
  if ! exampleURLs=($(to-get "/api/${TO_API_VERSION}/deliveryservices?xmlId=${xmlID}" | jq -r '.response[].exampleURLs[]')) ||
    [[ "${#exampleURLs[@]}" -eq 0 ]]; then
      echo waiting for delivery service ${xmlID} example URLs
      continue
  fi
  echo "example URLs: '${exampleURLs[*]}'"
  dsDomain="${exampleURLs[0]/*\/}"
  if ! digResponse="$(dig +short "@${TR_FQDN}" -t CNAME "$dsDomain")" ||
    [[ "$digResponse" != "$FEDERATION_CNAME" ]]; then
    echo "Failed to query Delivery Service ${dsDomain} for Federation CNAME ${FEDERATION_CNAME}"
    success=false
  fi
    if [[ "$success" == true ]]; then
        echo "successfully queried Delivery Service ${dsDomain} for Federation CNAME '${FEDERATION_CNAME}'"
        break
    fi
done

end_time=$(date +%s)
delta=$((end_time - start_time))
echo "completed readiness check in $delta seconds"
