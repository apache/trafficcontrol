#!/bin/bash -x

set -x

waiting=/sync/waiting-for-pgloader
touch $waiting

# Wait for pgloader to finish
while [[ -f $waiting ]]; do
    ls -l $waiting
    sleep 3
done

echo "Looks like pgloader is finished..  Converting.."

# Load required conversion of booleans
psql postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST/$POSTGRES_DB < ./convert_bools.sql
