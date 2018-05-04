#!/usr/bin/env bash

while ! nc $DB_SERVER $DB_PORT </dev/null; do # &>/dev/null; do
        echo "waiting for $DB_SERVER:$DB_PORT"
        sleep 3
done
psql -h $DB_SERVER -U postgres -c "CREATE USER $DB_USERNAME WITH ENCRYPTED PASSWORD '$DB_USER_PASS'"
createdb $DB_NAME -h $DB_SERVER -U postgres --owner $DB_USERNAME
