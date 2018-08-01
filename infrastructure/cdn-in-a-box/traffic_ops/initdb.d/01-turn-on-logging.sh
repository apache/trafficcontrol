#!/bin/sh

cd "$PGDATA"
mkdir -p pg_log
T=`mktemp`
sed "s/#logging_collector = off/logging_collector = on/; s/#log_dest/log_dest/; s/#log_statement = 'none'/log_statement = 'all'/ " <postgresql.conf >"$T"

mv "$T" postgresql.conf
