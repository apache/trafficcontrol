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

set -ex
env

envvars=( DB_SERVER POSTGRES_PASSWORD DB_PORT DB_NAME POSTGRES_USER DB_USER DB_PASSWORD TO_HOST TO_EMAIL TP_HOST TP_EMAIL TO_SECRET RIAK_USER RIAK_PASSWORD )
for v in ${envvars[*]}
do
    val=${!v}
    [[ -z $val ]] && echo $v is unset && exit 1
done

export TO=/opt/traffic_ops/app

# TODO: change sslmode=require when enabled in db
cat >$TO/db/dbconf.yml <<-DBCONF
version: "1.0"
name: dbconf.yml
production:
  driver: postgres
  open: host=$DB_SERVER port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable
DBCONF

cat >$TO/conf/production/database.conf <<-DATABASECONF
{
	"description": "PostgreSQL database on $DB_SERVER port $TO_PORT",
	"dbname": "$DB_NAME",
	"hostname": "$DB_SERVER",
	"user": "$DB_USER",
	"password": "$DB_PASSWORD",
	"port": "$DB_PORT",
	"ssl": false,
	"type": "Pg"
}
DATABASECONF

crt=$TO/conf/trafficops.crt
key=$TO/conf/trafficops.key
openssl req -newkey rsa:2048 -nodes -keyout $key -x509 -days 365 -out $crt -subj "/C=$CERT_COUNTRY/ST=$CERT_STATE/L=$CERT_CITY/O=$CERT_COMPANY"

cat >$TO/conf/cdn.conf <<-CDNCONF
{
    "hypnotoad" : {
        "listen" : [
            "https://trafficops-perl:60443?cert=$crt&key=$key&verify=0x00&ciphers=AES128-GCM-SHA256:HIGH:!RC4:!MD5:!aNULL:!EDH:!ED"
        ],
        "user" : "trafops",
        "group" : "trafops",
        "heartbeat_timeout" : 20,
        "pid_file" : "/var/run/traffic_ops.pid",
        "workers" : 12
    },
    "traffic_ops_golang" : {
        "port" : "443",
        "proxy_timeout" : 60,
        "proxy_keep_alive" : 60,
        "proxy_tls_timeout" : 60,
        "proxy_read_header_timeout" : 60,
        "read_timeout" : 60,
        "read_header_timeout" : 60,
        "write_timeout" : 60,
        "idle_timeout" : 60,
        "log_location_error": "/var/log/traffic_ops/error.log",
        "log_location_warning": "",
        "log_location_info": "",
        "log_location_debug": "",
        "log_location_event": "/var/log/traffic_ops/access.log",
        "max_db_connections": 20,
        "backend_max_connections": {
            "mojolicious": 4
        }
    },
    "cors" : {
        "access_control_allow_origin" : "*"
    },
    "to" : {
        "base_url" : "https://$TO_HOST",
        "email_from" : "no-reply@traffic-ops-domain.com",
        "no_account_found_msg" : "A Traffic Ops user account is required for access. Please contact your Traffic Ops user administrator."
    },
    "portal" : {
        "base_url" : "http://localhost:8080/!#/",
        "email_from" : "no-reply@traffic-portal-domain.com",
        "pass_reset_path" : "user",
        "user_register_path" : "user"
    },
    "secrets" : [
        "$TO_SECRET"
    ],
    "geniso" : {
        "iso_root_path" : "$TO/public"
    },
    "inactivity_timeout" : 60
}
CDNCONF

# ensure db is available for DB_USER
export PGPASSWORD="$DB_PASSWORD"
until psql -h "$DB_SERVER" -U "$DB_USER" -p "$DB_PORT" -c '\l'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

sleep 3 # give it a few extra seconds

psql -h "$DB_SERVER" -U "$DB_USER" -p "$DB_PORT" -c '\du'

>&2 echo "Postgres is up"


export PERL5LIB=$TO/lib:$TO/local/lib/perl5
cd $TO || (echo "NO $TO found" && exit 1)

chown -R trafops:trafops .

# needed for goose
export GOPATH=/opt/traffic_ops/go
export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin

./db/admin.pl -env production upgrade
./db/adduser.pl $TO_ADMIN_USER $TO_ADMIN_PASSWORD admin | tee /adduser.sql | psql -h "$DB_SERVER" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME" -e
./local/bin/hypnotoad script/cdn

while [[ ! -f /var/log/traffic_ops/traffic_ops.log ]]; do
    echo waiting for /var/log/traffic_ops/traffic_ops.log
    ps -ef | grep traffic
    sleep 1
done

touch $TO/conf/ready

exec tail -f /var/log/traffic_ops/traffic_ops.log
