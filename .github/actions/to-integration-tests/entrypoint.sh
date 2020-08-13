#!/bin/sh -l
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

set -ex

GOPATH="$(mktemp -d)"
SRCDIR="$GOPATH/src/github.com/apache"
mkdir -p "$SRCDIR"
ln -s "$PWD" "$SRCDIR/trafficcontrol"

cd "$SRCDIR/trafficcontrol/traffic_ops/app/db"

# /usr/local/go/bin/go get ./...
# /usr/local/go/bin/go build ./admin.go

echo 'version: "1.0"
name: dbconf.yml

test:
  driver: postgres
  open: host=postgres port=5432 user=traffic_ops password=twelve dbname=traffic_ops sslmode=disable

' > dbconf.yml

psql -d postgresql://traffic_ops:twelve@postgres:5432/traffic_ops < ./create_tables.sql
goose --env=test --path="$PWD" up
psql -d postgresql://traffic_ops:twelve@postgres:5432/traffic_ops < ./seeds.sql
psql -d postgresql://traffic_ops:twelve@postgres:5432/traffic_ops < ./patches.sql


# cd ..

# ./db/admin --env="test" reset

cd "$SRCDIR/trafficcontrol/traffic_ops/traffic_ops_golang"


/usr/local/go/bin/go get ./...
/usr/local/go/bin/go build .

openssl req -x509 -out localhost.crt -keyout localhost.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=localhost' -extensions EXT -config <( \
   printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")

cat <<-EOF >cdn.conf
{
    "hypnotoad" : {
        "listen" : [
            "https://not-a-real-host.test:1?cert=$PWD/localhost.crt&key=$PWD/localhost.key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED"
        ],
        "user" : "trafops",
        "group" : "trafops",
        "heartbeat_timeout" : 20,
        "pid_file" : "/var/run/traffic_ops.pid",
        "workers" : 12
    },
    "use_ims": true,
    "traffic_ops_golang" : {
        "insecure": true,
        "port" : "6443",
        "log_location_error": "stderr",
        "log_location_warning": "stderr",
        "log_location_info": "stderr",
        "log_location_debug": "stderr",
        "log_location_event": "stderr",
        "max_db_connections": 20,
        "db_conn_max_lifetime_seconds": 60,
        "db_query_timeout_seconds": 20,
        "supported_ds_metrics": [ "kbps", "tps_total", "tps_2xx", "tps_3xx", "tps_4xx", "tps_5xx" ]
    },
    "cors" : {
        "access_control_allow_origin" : "*"
    },
    "to" : {
        "base_url" : "https://localhost",
        "email_from" : "no-reply@traffic-ops.test",
        "no_account_found_msg" : "A Traffic Ops user account is required for access. Please contact your Traffic Ops user administrator."
    },
    "portal" : {
        "base_url" : "https://not-a-real-host.test/#!/",
        "email_from" : "no-reply@traffic-portal.test",
        "pass_reset_path" : "user",
        "user_register_path" : "user"
    },
    "secrets" : [
        "blahblah"
    ],
    "geniso" : {
        "iso_root_path" : "/opt/traffic_ops/app/public"
    },
    "inactivity_timeout" : 60,
    "smtp" : {
        "enabled" : false
    },
    "InfluxEnabled": false
}
EOF

cat <<-EOF >database.conf
{
        "description": "Local PostgreSQL database on port 5432",
        "dbname": "traffic_ops",
        "hostname": "postgres",
        "user": "traffic_ops",
        "password": "twelve",
        "port": "5432",
        "ssl": false,
        "type": "Pg"
}
EOF



./traffic_ops_golang --cfg ./cdn.conf --dbcfg ./database.conf >out.log 2>err.log &

cd ../testing/api/v1

cat <<-EOF >traffic-ops-test.conf
{
    "default": {
        "logLocations": {
            "debug": "stdout",
            "error": "stdout",
            "event": "stdout",
            "info": "stdout",
            "warning": "stdout"
        },
        "session": {
            "timeoutInSecs": 60
        }
    },
    "trafficOps": {
        "URL": "https://localhost:6443",
        "password": "twelve",
        "users": {
            "disallowed": "disallowed",
            "operations": "operations",
            "admin": "admin",
            "federation": "federation",
            "portal": "portal",
            "readOnly": "readOnly",
            "extension": "extension"
        }
    },
    "trafficOpsDB": {
        "dbname": "traffic_ops",
        "description": "Test database to_test",
        "hostname": "postgres",
        "password": "twelve",
        "port": "5432",
        "type": "Pg",
        "user": "traffic_ops"
    }
}
EOF

go test -v --cfg ./traffic-ops-test

exit $?
