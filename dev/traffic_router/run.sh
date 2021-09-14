#!/bin/sh

set -o errexit
trap '[ $? -eq 0 ] && exit 0 || echo "Error on line ${LINENO} of ${0}"; exit 1' EXIT

cd "$TC/traffic_router"

mvn -Dmaven.test.skip=true compile -P \!rpm-build
mvn -Dmaven.test.skip=true package -P \!rpm-build

/opt/tomcat/bin/catalina.sh jpda run
# java -agentlib:jdwp=transport=dt_socket,address=5005,server=y,suspend=n StartTrafficRouter

# while inotifywait --exclude '.*(\.md|_test\.go|\.gitignore|__debug_bin)$' -e modify -r . ; do
# 	kill "$(netstat -nlp | grep ':443' | grep __debug_bin | head -n1 | tr -s ' ' | cut -d ' ' -f7 | cut -d '/' -f1)"
# 	kill "$(netstat -nlp | grep ':6444' | grep dlv | head -n1 | tr -s ' ' | cut -d ' ' -f7 | cut -d '/' -f1)"
# 	dlv --accept-multiclient --continue --listen=:6444 --headless --api-version=2 debug -- --cfg=../../dev/traffic_ops/cdn.json --dbcfg=../../dev/traffic_ops/db.config.json &
# 	# for whatever reason, without this the repeated call to inotifywait will
# 	# sometimes lose track of th current directory. It spits out:
# 	# Couldn't watch .: No such file or directory
# 	# which is a bit odd.
# 	sleep 0.5
# done;
