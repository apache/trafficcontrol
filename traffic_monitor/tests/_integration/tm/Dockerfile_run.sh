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

# The following environment variables must be set (ordinarily by `docker run -e` arguments):
# TO_URI
# TO_USER
# TO_PASS
# CDN

# Check that env vars are set
envvars=( TO_URI TO_USER TO_PASS CDN PORT )
for v in ${envvars[@]}
do
	if [[ -z ${!v} ]]; then echo "$v is unset"; exit 1; fi
done

start() {
	service traffic_monitor start
	touch /var/log/traffic_monitor/traffic_monitor.log
	exec tail -f /var/log/traffic_monitor/traffic_monitor.log
}

init() {
	mkdir -p /opt/traffic_monitor/conf
	cat > /opt/traffic_monitor/conf/traffic_monitor.cfg <<- EOF
		{
				"monitor_config_polling_interval_ms": 15000,
				"http_timeout_ms": 2000,
				"max_events": 200,
				"health_flush_interval_ms": 20,
				"stat_flush_interval_ms": 20,
				"log_location_access": "/var/log/traffic_monitor/access.log",
				"log_location_event": "/var/log/traffic_monitor/event.log",
				"log_location_error": "/var/log/traffic_monitor/traffic_monitor.log",
				"log_location_warning": "/var/log/traffic_monitor/traffic_monitor.log",
				"log_location_info": "null",
				"log_location_debug": "null",
				"serve_read_timeout_ms": 10000,
				"serve_write_timeout_ms": 10000,
				"static_file_dir": "/opt/traffic_monitor/static/"
		}
EOF

  cat > /opt/traffic_monitor/conf/traffic_ops.cfg <<- EOF
		{
				"username": "$TO_USER",
				"password": "$TO_PASS",
				"url": "$TO_URI",
				"insecure": true,
				"cdnName": "$CDN",
				"httpListener": ":$PORT"
				}
	EOF

	echo "INITIALIZED=1" >> /etc/environment
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
