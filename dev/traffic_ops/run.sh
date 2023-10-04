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

db_init() {
	while ! pg_isready -h db -p 5432 -d postgres; do
		echo "waiting for db on postgresql://db:5432/postgres";
		sleep 3;
	done

	cd "$TC"
	make traffic_ops/app/db/admin
	(
		cd "$TC/dev/traffic_ops"
		"$ADMIN" -c ./dbconf.yml -s "$TC/traffic_ops/app/db/create_tables.sql" -S "$TC/traffic_ops/app/db/seeds.sql" -p "$TC/traffic_ops/app/db/patches.sql" -m "$TC/traffic_ops/app/db/migrations" reset
		"$ADMIN" -c ./dbconf.yml -s "$TC/traffic_ops/app/db/create_tables.sql" -S "$TC/traffic_ops/app/db/seeds.sql" -p "$TC/traffic_ops/app/db/patches.sql" -m "$TC/traffic_ops/app/db/migrations" upgrade
		"$ADMIN" -c ./dbconf.yml -s "$TC/traffic_ops/app/db/create_tables.sql" -S "$TC/traffic_ops/app/db/seeds.sql" -p "$TC/traffic_ops/app/db/patches.sql" -m "$TC/traffic_ops/app/db/migrations" seed
		"$ADMIN" -v -c ./traffic.vault.dbconf.yml -s "$TC/traffic_ops/app/db/trafficvault/create_tables.sql" -m "$TC/traffic_ops/app/db/trafficvault/migrations" reset
		"$ADMIN" -v -c ./traffic.vault.dbconf.yml -s "$TC/traffic_ops/app/db/trafficvault/create_tables.sql" -m "$TC/traffic_ops/app/db/trafficvault/migrations" upgrade

		psql -d 'postgres://traffic_ops:twelve12@db:5432/traffic_ops_development?sslmode=disable' -f ./seed.psql
	)
}

user=trafficops
uid="$(stat -c%u "$TC")"
gid="$(stat -c%g "$TC")"
if [[ "$(id -u)" != "$uid" ]]; then
	# db/admin must be run as root (see apache/trafficcontrol#7202)
	if [[ $uid -ne 0 ]]; then
		db_init
		chown "${uid}:${gid}" traffic_ops/app/db/admin
	fi

	for dir in "${GOPATH}/bin" "${GOPATH}/pkg"; do
		if [[ -e "$dir" ]] && [[ "$(stat -c%u "$dir")" -ne "$uid" || "$(stat -c%g "$dir")" -ne "$gid" ]] ; then
			chown -R "${uid}:${gid}" "$dir"
		fi
	done

	adduser -Du"$uid" "$user"
	sed -Ei "s/^(${user}:.*:)[0-9]+(:)$/\1${gid}\2/" /etc/group
	exec su "$user" -- "$0"
fi

# On Docker Desktop, bind mounts are owned by root
if [[ $uid -eq 0 ]]; then
	db_init
fi

cd "$TC/traffic_ops/traffic_ops_golang"

dlv --accept-multiclient --continue --listen=:6444 --headless --api-version=2 debug -- --cfg=../../dev/traffic_ops/cdn.json --dbcfg=../../dev/traffic_ops/db.config.json &

while inotifywait --include '\.go$' -e modify -r . ; do
	kill "$(netstat -nlp | grep ':443' | grep __debug_bin | head -n1 | tr -s ' ' | cut -d ' ' -f7 | cut -d '/' -f1)"
	kill "$(netstat -nlp | grep ':6444' | grep dlv | head -n1 | tr -s ' ' | cut -d ' ' -f7 | cut -d '/' -f1)"
	dlv --accept-multiclient --continue --listen=:6444 --headless --api-version=2 debug -- --cfg=../../dev/traffic_ops/cdn.json --dbcfg=../../dev/traffic_ops/db.config.json &
	# for whatever reason, without this the repeated call to inotifywait will
	# sometimes lose track of th current directory. It spits out:
	# Couldn't watch .: No such file or directory
	# which is a bit odd.
	sleep 0.5
done;
