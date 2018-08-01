#!/bin/bash

set -ex

d=/docker-entrypoint-initdb.d
for dump in "$d"/*.dump; do
    [[ -f $dump ]] || break
    t=$(mktemp XXX.sql)
    # convert to sql -- can't load a dump until db initialized,  but sql works
    echo "Restoring from $dump"
    pg_restore -f "$t" $dump
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <"$t"
done
