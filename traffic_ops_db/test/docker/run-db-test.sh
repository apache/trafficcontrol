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
DB_USER=traffic_ops
DB_USER_PASS=twelve
DB_NAME=traffic_ops

# Write config files
set -x
if [[ ! -r /goose-config.sh ]]; then
	echo "/goose-config.sh not found/readable"
	exit 1
fi
. /goose-config.sh

pg_isready=$(rpm -ql postgresql96 | grep bin/pg_isready)
if [[ ! -x $pg_isready ]] ; then
    echo "Can't find pg_ready in postgresql96"
    exit 1
fi

while ! $pg_isready -h$DB_SERVER -p$DB_PORT -d $DB_NAME; do
        echo "waiting for db on $DB_SERVER $DB_PORT"
        sleep 3
done

export TO_DIR=/opt/traffic_ops/app

export PATH=/usr/local/go/bin:/opt/traffic_ops/go/bin:$PATH
export GOPATH=/opt/traffic_ops/go

# gets the current DB version. On success, output the version number. On failure, output a failure message starting with 'failed'.
get_current_db_version() {
    local dbversion_output=$(./db/admin --env=production dbversion 2>&1)
    if [[ $? -ne 0 ]]; then
        echo "failed to get dbversion: $dbversion_output"
        return
    fi
    local version=$(echo "$dbversion_output" | egrep '^goose: dbversion [[:digit:]]+$' | awk '{print $3}')
    if [[ -z "$version" ]]; then
        echo "failed to get dbversion from output: $db_version_output"
        return
    fi
    echo "$version"
}

get_db_dumps() {
    find /db_dumps -name '*.dump'
}

for d in $(get_db_dumps); do
    echo "checking integrity of DB dump: $d"
    pg_restore -l "$d" > /dev/null || { echo "invalid DB dump: $d. Unable to list contents"; exit 1; }
done

cd $TO_DIR
db_is_empty=false
old_db_version=$(get_current_db_version)
[[ "$old_db_version" =~ ^failed ]] && { echo "get_current_db_version failed: $old_db_version"; exit 1; }

# reset the DB if it is empty (i.e. no db.dump was provided)
if [[ "$old_db_version" -eq 0 ]]; then
    db_is_empty=true
    ./db/admin --env=production reset || { echo "DB reset failed!"; exit 1; }
fi

./db/admin --env=production upgrade || { echo "DB upgrade failed!"; exit 1; }

if ! ./db/admin -env=production load_schema ||
  ! ./db/admin -env=production load_schema; then
  echo 'Could not re-run create_tables.sql!'
  exit 1
fi;

new_db_version=$(get_current_db_version)
[[ "$new_db_version" =~ ^failed ]] && { echo "get_current_db_version failed: $new_db_version"; exit 1; }

run_db_downgrades=true

if [[ "$old_db_version" = "$new_db_version" ]]; then
    echo "new DB version matches old DB version, no downgrade migrations to test"
    run_db_downgrades=false
fi

if [[ "$db_is_empty" = true ]]; then
    echo "starting DB was empty, skipping DB downgrades"
    run_db_downgrades=false
fi

if [[ "$run_db_downgrades" = true ]]; then
    # downgrade the DB until the initial DB version
    while [[ "$new_db_version" != "$old_db_version" ]]; do
        ./db/admin --env=production down || { echo "DB downgrade failed!"; exit 1; }
        new_db_version=$(get_current_db_version)
        [[ "$new_db_version" =~ ^failed ]] && { echo "get_current_db_version failed: $new_db_version"; exit 1; }
    done
fi

# test full restoration of the initial DB dump
for d in $(get_db_dumps); do
    echo "testing restoration of DB dump: $d"
    pg_restore --verbose --clean --if-exists --create -h $DB_SERVER -p $DB_PORT -U postgres < "$d" > /dev/null || { echo "DB restoration failed: $d"; exit 1; }
done

