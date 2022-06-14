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

set -o errexit

NAME="Traffic Portal Application"
NODE_BIN_DIR="/usr/bin"
PM2_BIN_DIR="/opt/traffic_portal/node_modules/pm2/bin"
APPLICATION_PATH="/opt/traffic_portal/server.js"
PIDFILE="/var/run/traffic_portal.pid"
LOGFILE="/var/log/traffic_portal/traffic_portal.log"
RESTART_TIME="2000"

envvars=(TO_SERVER TO_PORT DOMAIN)
for v in "${envvars}"; do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

CONF_DIR="/opt/traffic_portal/conf"

mkdir -p "${CONF_DIR}"

KEY_FILE="${CONF_DIR}/key.pem"
CERT_FILE="${CONF_DIR}/cert.pem"
openssl req -nodes -x509 -newkey rsa:4096 -keyout "${KEY_FILE}" -out "${CERT_FILE}" -days 365 -subj "/CN=${DOMAIN}"

CONF_FILE="/etc/traffic_portal/conf/config.js"

sed -i -e "/^\s*base_url:/ s@'.*'@'https://$TO_SERVER:$TO_PORT/api/'@" "${CONF_FILE}"
sed -i -e "/^\s*cert:/ s@'.*'@'${CERT_FILE}'@" "${CONF_FILE}"
sed -i -e "/^\s*key:/ s@'.*'@'${KEY_FILE}'@" "${CONF_FILE}"

props=/opt/traffic_portal/public/traffic_portal_properties.json
tmp=$(mktemp)

jq --arg TO_SERVER "$TO_SERVER:$TO_PORT" '.properties.api.baseUrl = "https://"+$TO_SERVER' <$props >$tmp
mv $tmp $props

# Add node to the path for situations in which the environment is passed.
PATH=$PM2_BIN_DIR:$NODE_BIN_DIR:$PATH
pm2 \
    start "$APPLICATION_PATH" \
    -p $PIDFILE \
    -l $LOGFILE \
    --restart-delay $RESTART_TIME \
    --time \
    --name "$NAME"

tail -f /dev/null
