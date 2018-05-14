#!/bin/bash

if [[ ! -x ./bin/traffic_server ]]; then
	echo "Missing 'traffic_server' executable!" >&2
	exit 1
fi

./bin/traffic_server 2>&1 >/var/log/traffic_server/ts.log &
disown
exec tail -f /var/log/traffic_server/ts.log
