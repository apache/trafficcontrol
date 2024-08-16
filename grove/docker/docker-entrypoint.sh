#!/usr/bin/env bash

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# The following environment variables must be set, ordinarily by `docker run -e` arguments:
envvars=( REMAP_PATH )
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done

[ ! -z $PORT ]        || PORT=80
[ ! -z $HTTPS_PORT ]  || HTTPS_PORT=443
[ ! -z $CACHE_BYTES ] || CACHE_BYTES=100000000 # 100mb

start() {
	service grove start
	exec tail -f /var/log/grove/error.log
}

init() {
	cat > "/etc/grove/grove.cfg" <<- ENDOFMESSAGE
{
  "rfc_compliant":            false,
  "port":                     $PORT,
  "https_port":               $HTTPS_PORT,
  "cache_size_bytes":         $CACHE_BYTES,
  "remap_rules_file":         "/etc/grove/remap.json",
  "concurrent_rule_requests": 100,
  "connection_close":         false,
  "interface_name":           "bond0",
  "cert_file":                "/etc/grove/cert.pem",
  "key_file":                 "/etc/grove/key.pem",

  "log_location_error":   "/var/log/grove/error.log",
  "log_location_warning": "/var/log/grove/error.log",
  "log_location_info":    "null",
  "log_location_debug":   "null",
  "log_location_event":   "/var/log/trafficserver/custom_ats_2.log",

  "parent_request_timeout_ms":                 10000,
  "parent_request_keep_alive_ms":              10000,
  "parent_request_max_idle_connections":       10000,
  "parent_request_idle_connection_timeout_ms": 10000,

  "server_read_timeout_ms":  5000,
  "server_write_timeout_ms": 5000,
  "server_idle_timeout_ms":  5000
}
ENDOFMESSAGE

	# TODO add Traffic Ops uri+user+pass+hostname as an option, rather than remap file
	if [[ ! -z $REMAP_PATH ]]; then
    cp $REMAP_PATH /etc/grove/remap.json
  fi
	mkdir -p /var/log/trafficserver
	mkdir -p /var/log/grove/
	touch /var/log/grove/error.log

	openssl req -newkey rsa:2048 -nodes -keyout /etc/grove/key.pem -x509 -days 3650 -out /etc/grove/cert.pem -subj "/C=US/ST=Colorado/L=Denver/O=MyCompany/CN=cdn.example.net"

	# TODO add server to Traffic Ops, with env vars

	echo "INITIALIZED=1" >> /etc/environment
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
