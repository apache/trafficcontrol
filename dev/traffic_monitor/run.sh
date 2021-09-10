#!/bin/sh

set -o errexit
trap '[ $? -eq 0 ] && exit 0 || echo "Error on line ${LINENO} of ${0}"; exit 1' EXIT

TC=/root/go/src/github.com/apache/trafficcontrol

cd "$TC/traffic_monitor"

dlv --accept-multiclient --continue --listen=:81 --headless --api-version=2 debug -- --opsCfg="$TC/dev/traffic_monitor/ops.config.json" --config="$TC/dev/traffic_monitor/tm.config.json" &

# "static" files need to be watched since TM caches their contents, so it needs
# to be restarted to apply changes to them.
while inotifywait --exclude '.*(\.md|\.csv|\.json|\.cfg|_test\.go|\.gitignore|__debug_bin)$|^\./(build|tools|tests|tmclient)/.*' -e modify -r . ; do
	kill "$(netstat -nlp | grep ':80' | grep __debug_bin | head -n1 | tr -s ' ' | cut -d ' ' -f7 | cut -d '/' -f1)"
	kill "$(netstat -nlp | grep ':81' | grep dlv | head -n1 | tr -s ' ' | cut -d ' ' -f7 | cut -d '/' -f1)"
	dlv --accept-multiclient --continue --listen=:81 --headless --api-version=2 debug -- --opsCfg="$TC/dev/traffic_monitor/ops.config.json" --config="$TC/dev/traffic_monitor/tm.config.json" &
	# for whatever reason, without this the repeated call to inotifywait will
	# sometimes lose track of th current directory. It spits out:
	# Couldn't watch .: No such file or directory
	# which is a bit odd.
	sleep 0.5
done;
