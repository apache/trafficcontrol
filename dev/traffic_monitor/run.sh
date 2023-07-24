#!/bin/sh
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
set -o xtrace
trap '[ $? -eq 0 ] && exit 0 || echo "Error on line ${LINENO} of ${0}"; exit 1' EXIT

cd "$TC/traffic_monitor"
user=trafficmonitor
uid="$(stat -c%u "$TC")"
gid="$(stat -c%g "$TC")"
if [[ "$(id -u)" != "$uid" ]]; then
	for dir in "${GOPATH}/bin" "${GOPATH}/pkg"; do
		if [[ -e "$dir" ]] && [[ "$(stat -c%u "$dir")" -ne "$uid" || "$(stat -c%g "$dir")" -ne "$gid" ]] ; then
			chown -R "${uid}:${gid}" "$dir"
		fi
	done

	adduser -Du"$uid" "$user"
	sed -Ei "s/^(${user}:.*:)[0-9]+(:)$/\1${gid}\2/" /etc/group
	exec su "$user" -- "$0"
fi

dlv --accept-multiclient --continue --listen=:81 --headless --api-version=2 debug -- --opsCfg="$TC/dev/traffic_monitor/ops.config.json" --config="$TC/dev/traffic_monitor/tm.config.json" &

# "static" files need to be watched since TM caches their contents, so it needs
# to be restarted to apply changes to them.
while inotifywait --include '\.go$' -e modify -r . ; do
	kill "$(netstat -nlp | grep ':80' | grep __debug_bin | head -n1 | tr -s ' ' | cut -d ' ' -f7 | cut -d '/' -f1)"
	kill "$(netstat -nlp | grep ':81' | grep dlv | head -n1 | tr -s ' ' | cut -d ' ' -f7 | cut -d '/' -f1)"
	dlv --accept-multiclient --continue --listen=:81 --headless --api-version=2 debug -- --opsCfg="$TC/dev/traffic_monitor/ops.config.json" --config="$TC/dev/traffic_monitor/tm.config.json" &
	# for whatever reason, without this the repeated call to inotifywait will
	# sometimes lose track of th current directory. It spits out:
	# Couldn't watch .: No such file or directory
	# which is a bit odd.
	sleep 0.5
done;
