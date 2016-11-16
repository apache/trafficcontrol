#!/bin/bash -x

# wait for traffic_ops.sql file to appear
while [[ ! -f /docker-entrypoint-initdb.d/traffic_ops.sql ]]; do
	sleep 1
done
