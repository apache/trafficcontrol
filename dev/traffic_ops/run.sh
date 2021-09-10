#!/bin/sh

set -o errexit
trap '[ $? -eq 0 ] && exit 0 || echo "Error on line ${LINENO} of ${0}"; exit 1' EXIT

while ! pg_isready -h db -p 5432 -d postgres; do
	echo "waiting for db on postgresql://db:5432/postgres";
	sleep 3;
done

TC=/root/go/src/github.com/apache/trafficcontrol
ADMIN="$TC/traffic_ops/app/db/admin"

cd "$TC"
make
cd "$TC/dev/traffic_ops"

"$ADMIN" -c ./dbconf.yml -s "$TC/traffic_ops/app/db/create_tables.sql" -S "$TC/traffic_ops/app/db/seeds.sql" -p "$TC/traffic_ops/app/db/patches.sql" -m "$TC/traffic_ops/app/db/migrations" reset
"$ADMIN" -c ./dbconf.yml -s "$TC/traffic_ops/app/db/create_tables.sql" -S "$TC/traffic_ops/app/db/seeds.sql" -p "$TC/traffic_ops/app/db/patches.sql" -m "$TC/traffic_ops/app/db/migrations" upgrade
"$ADMIN" -v -c ./traffic.vault.dbconf.yml -s "$TC/traffic_ops/app/db/trafficvault/create_tables.sql" -m "$TC/traffic_ops/app/db/trafficvault/migrations" reset
"$ADMIN" -v -c ./traffic.vault.dbconf.yml -s "$TC/traffic_ops/app/db/trafficvault/create_tables.sql" -m "$TC/traffic_ops/app/db/trafficvault/migrations" upgrade

psql -d 'postgres://traffic_ops:twelve12@db:5432/traffic_ops_development?sslmode=disable' -f ./insert.admin.psql

cd "$TC/traffic_ops/traffic_ops_golang"

dlv --accept-multiclient --continue --listen=:6444 --headless --api-version=2 debug -- --cfg=../../dev/traffic_ops/cdn.json --dbcfg=../../dev/traffic_ops/db.config.json &


while inotifywait --exclude '.*(\.md|_test\.go|\.gitignore|__debug_bin)$' -e modify -r . ; do
	kill "$(netstat -nlp | grep ':443' | grep __debug_bin | head -n1 | tr -s ' ' | cut -d ' ' -f7 | cut -d '/' -f1)"
	kill "$(netstat -nlp | grep ':6444' | grep dlv | head -n1 | tr -s ' ' | cut -d ' ' -f7 | cut -d '/' -f1)"
	dlv --accept-multiclient --continue --listen=:6444 --headless --api-version=2 debug -- --cfg=../../dev/traffic_ops/cdn.json --dbcfg=../../dev/traffic_ops/db.config.json &
	# for whatever reason, without this the repeated call to inotifywait will
	# sometimes lose track of th current directory. It spits out:
	# Couldn't watch .: No such file or directory
	# which is a bit odd.
	sleep 0.5
done;
