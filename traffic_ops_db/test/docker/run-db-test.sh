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

# Script for running the Traffic Ops DB migration tests.
#
DB_SERVER=db
DB_PORT=5432
DB_USER=postgres
DB_USER_PASS=twelve
DB_NAME=traffic_ops
export PGHOST="$DB_SERVER" PGPORT="$DB_PORT" PGUSER="$DB_USER" PGDATABASE="$DB_NAME"

# Write config files
set -x
if [[ ! -r /db-config.sh ]]; then
	echo "/db-config.sh not found/readable"
	exit 1
fi
. /db-config.sh

postgresql_package="$(<<<"postgresql${POSTGRES_VERSION}" sed 's/\.//g' |
	sed -E 's/([0-9]{2})[0-9]+/\1/g'
)"
pg_isready=$(rpm -ql "$postgresql_package" | grep bin/pg_isready)
if [[ ! -x "$pg_isready" ]] ; then
	echo "Can't find pg_ready in ${postgresql_package}"
	exit 1
fi

while ! $pg_isready -h"$DB_SERVER" -p"$DB_PORT" -d "$DB_NAME"; do
	echo "waiting for db on $DB_SERVER $DB_PORT"
	sleep 3
done

echo "*:*:*:postgres:$DB_USER_PASS" > "${HOME}/.pgpass"
echo "*:*:*:traffic_ops:$DB_USER_PASS" >> "${HOME}/.pgpass"
chmod 0600 "${HOME}/.pgpass"

export TO_DIR=/opt/traffic_ops/app

export PATH=/usr/local/go/bin:/opt/traffic_ops/go/bin:$PATH
export GOPATH=/opt/traffic_ops/go

# gets the current DB version. On success, output the version number. On failure, output a failure message starting with 'failed'.
get_current_db_version() {
	local dbversion_output
	if ! dbversion_output="$(./db/admin --env=production dbversion 2>&1)"; then
		echo "failed to get dbversion: $dbversion_output"
		return
	fi
	local version=$(echo "$dbversion_output" | egrep '^dbversion [[:digit:]]+$' | awk '{print $2}')
	if [[ -z "$version" ]]; then
		echo "failed to get dbversion from output: $db_version_output"
		return
	fi
	echo "$version"
}

get_db_dumps() {
	find /db_dumps -name '*.dump'
}

db_is_empty=true

for d in $(get_db_dumps); do
	db_is_empty=false
	echo "checking integrity of DB dump: $d"
	pg_restore -l "$d" > /dev/null || { echo "invalid DB dump: $d. Unable to list contents"; exit 1; }
done

cd "$TO_DIR"
# This NEEDS to be updated if migrations are squashed. It should be the
# timestamp of the oldest extant migration.
# TODO: this can be determined automatically from an inspection of the migrations dir
FIRST_MIGRATION=2022100210472946

old_db_version=$FIRST_MIGRATION

if [[ "$db_is_empty" = false ]]; then
  old_db_version=$(get_current_db_version)
  [[ "$old_db_version" =~ ^failed ]] && { echo "get_current_db_version failed: $old_db_version"; exit 1; }
fi

# reset the DB if it is empty (i.e. no db.dump was provided)
if [[ "$old_db_version" -eq $FIRST_MIGRATION ]]; then
	db_is_empty=true
	./db/admin --env=production reset || { echo "DB reset failed!"; exit 1; }
fi

# applies migrations then performs seeding and patching
./db/admin --env=production upgrade || { echo "DB upgrade failed!"; exit 1; }

new_db_version=$(get_current_db_version)
[[ "$new_db_version" =~ ^failed ]] && { echo "get_current_db_version failed: $new_db_version"; exit 1; }

# downgrade the DB until the initial DB version
while [[ "$new_db_version" != "$old_db_version" ]]; do
	./db/admin --env=production down || { echo "DB downgrade failed!"; exit 1; }
	new_db_version=$(get_current_db_version)
	[[ "$new_db_version" =~ ^failed ]] && { echo "get_current_db_version failed: $new_db_version"; exit 1; }
done

# test full restoration of the initial DB dump
for d in $(get_db_dumps); do
	echo "testing restoration of DB dump: $d"
	dropdb --echo --if-exists "$DB_NAME" < "$d" > /dev/null || echo "Dropping DB ${DB_NAME} failed: $d"
	createdb --echo < "$d" > /dev/null || echo "Creating DB ${DB_NAME} failed: $d"
	pg_restore --verbose --clean --if-exists --exit-on-error -d "$DB_NAME" < "$d" > /dev/null || { echo "DB restoration failed: $d"; exit 1; }
done
