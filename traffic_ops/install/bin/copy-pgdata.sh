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

#
# Copy traffic_ops postgres db from traffic_ops and restore to local postgres.
#
# NOTE!! This script uses `psql` on the local machine to delete and overwrite the database.
#
#      *******  This is destructive to the local database!!  ******
#      *******  This is destructive to the local database!!  ******
#      *******  This is destructive to the local database!!  ******

# Set the following env vars to avoid having prompts for each of these..
# Environment variables used:
#     TO_URL       - Source URL for traffic_ops
#     TO_USER      - User to authenticate (requires admin privileges for db dump)
#     TO_PASSWORD  - Password for above user
#
#     TO_DEST_URL  - Destination URL for traffic_ops on the local machine (the one to be overwritten)
#     TODB_TERMINATE_CONNECTIONS - set to "y" to force current database connections to be forced off

cookie_current() {
    local cookiefile=$1
    [[ -f $cookiefile ]] || return 1

    # get expiration from cookiejar -- compare to current time
    exp=$(awk '/mojolicious/ {print $5}' $cookiefile | tail -n 1)
    cur=$(date +%s)

    # compare expiration with current time
    (( $exp > $cur ))
}

to-auth () {
    [[ -z $TO_URL ]] && read -p 'Traffic Ops URL: ' TO_URL
    [[ -z $TO_USER ]] && read -p 'Traffic Ops user: ' TO_USER
    [[ -z $TO_PASSWORD ]] && read -s -p "Traffic Ops password for $TO_USER: " TO_PASSWORD

    COOKIEJAR=/tmp/cookiejar.$(echo $TO_URL $TO_USER | md5sum | awk '{print $1}')
    cookie_current $COOKIEJAR && return
    local datadir=$(mktemp -d)
    local login="$datadir/login.json"
    local url=$TO_URL/api/4.0/user/login
    local datatype='Accept: application/json'
    cat > "$login"  <<-CREDS
        { "u" : "$TO_USER", "p" : "$TO_PASSWORD" }
CREDS

    res=$(curl -k -H "$datatype" --cookie "$COOKIEJAR" --cookie-jar "$COOKIEJAR" -X POST --data @"$login" "$url")

    # clean up creds
    rm -rf $datadir
    if [[ $res != *"Successfully logged in."* ]]; then
        echo $res
        return 1
    fi
}

to-get () {
    to-auth && curl -L -k -s --cookie "$COOKIEJAR" -X GET "$TO_URL/$1"
}

# Dump the postgres db to a file from traffic ops
dump_source_db() {
    local dumpfile="$1"
    to-get api/4.0/dbdump >"$dumpfile"
}

# Prepare the destination db by terminating any existing connections and then dropping and creating the db.
prep_destination() {
    local target_db=${1:-traffic_ops}
    # export TODB_TERMINATE_CONNECTIONS=y will terminate the connections without prompting
    local ans=$TODB_TERMINATE_CONNECTIONS

    while [[ $ans != y ]]; do
        read -p "Terminating connections to $target_db.  OK? (y/n) " ans
        case $ans in
            n)
                echo "Not terminating connections"
                exit
                ;;
            y)
                break
                ;;
            *)
                echo "Answer y or n"
                ;;
        esac
    done

    # Create the sql to terminate connections -- avoid using single-quotes
    read -d% termsql <<-EOF
    SELECT pg_terminate_backend(pg_stat_activity.pid)
        FROM pg_stat_activity
        WHERE pg_stat_activity.datname = \$\$$target_db\$\$
        AND pid <> pg_backend_pid();
EOF

    # terminate any connections to the destination db -- won't complete otherwise
    echo $termsql | psql -Upostgres -h localhost
    # drop and create destination db
    dropdb -h localhost -Upostgres $target_db && createdb -Upostgres -h localhost --owner traffic_ops $target_db
}

# write cr-config
write_crconfig() {
    to-put "/api/4.0/snapshot/$1" >/dev/null
}


#----------------------------------------------------------------
# main starts here
#
#

usage() {
   fmt <<-USAGE
 $0 [<sql file> ...]

 $0 copies a Traffic Ops postgresql database from an existing Traffic Ops installation
 into a postgresql installation.  This must be run from a server that has access to the
 db using the psql command with the postgres user.

 The script prompts for all needed information unless the corresponding env var is
 set for each piece:

     TO_URL       - Source URL for traffic_ops

     TO_USER      - User to authenticate (requires admin privileges for db dump)

     TO_PASSWORD  - Password for above user

     TO_DEST_URL  - Destination URL for traffic_ops on the local machine (the one to be overwritten)

     TODB_TERMINATE_CONNECTIONS - set to "y" for current database connections to be forced off (NOTE:
      this is required for the database to be replaced)

 To run this script from an automated system (e.g. Jenkins or some other CI system),  set each of the
 above variables.

 Any sql files listed on the command line are executed on the database once copied.  This allows for
 adjusting to a different set of servers, for example..
USAGE
    exit
}

while getopts ":h" opt; do
  case ${opt} in
    h ) usage ;;
  esac
done

target_db=traffic_ops

# Create a tmp path to collect the dump
dumpfile=$(mktemp --tmpdir=/tmp pg-XXX.dump)
cat >$dumpfile </dev/null
cleanup() {
    rm -f $dumpfile
}
trap cleanup EXIT

# Dump db to dump file
echo "Dumping db to $dumpfile"
dump_source_db "$dumpfile"

# Check validity of dump before proceeding
pg_restore -l "$dumpfile" > /dev/null
if [[ $? -ne 0 ]] ; then
	exit 1
fi

echo "Prepping destination $dumpfile"
prep_destination traffic_ops

echo "Restoring from $dumpfile"
pg_restore --verbose --clean --create -h localhost -U postgres -d traffic_ops <$dumpfile

if [[ $# > 0 ]]; then
        for a in "$@"; do
                [[ $a == *.sql ]] || continue
                echo "Loading $a"
                psql -h localhost -U postgres -d traffic_ops -f "$a"
        done
fi
(cd /opt/traffic_ops/app;  PATH=$PATH:/opt/traffic_ops/go/bin ./db/admin -env production upgrade)

TO_URL=$TO_DEST_URL
echo "Snapshotting CRConfigs on $TO_URL"

# get the list of cdns from the copied db
cdns=$(to-get api/4.0/cdns | jq -Sr '.response|.[]|.name' | grep -v ALL)

for c in $cdns; do
    write_crconfig "$c"
done

