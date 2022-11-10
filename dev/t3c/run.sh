#!/bin/sh

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

set -o errexit
trap '[ $? -eq 0 ] && exit 0 || echo "Error on line ${LINENO} of ${0}"; exit 1' EXIT

cd "$TC/tc-health-client"
go build --gcflags "all=-N -l" .

cd "$TC/cache-config"

# Build area may contain non-debug binaries
make clean && make -j debug

for component in "t3c t3c-apply t3c-check t3c-check-refs t3c-check-reload t3c-diff t3c-generate t3c-preprocess t3c-request t3c-update"; do
	if [[ ! -f "/usr/bin/$component" ]]; then
		ln -s "$TC/cache-config/$component/$component" /usr/bin
	fi
done

if [[ ! -f /usr/bin/tc-health-client ]]; then
	ln -s "$TC/tc-health-client/tc-health-client" /usr/bin/
fi

su -c traffic_server ats &

while inotifywait --exclude '.*(\.md|\.json|\.pl|\.rst|_test\.go|\.gitignore|__debug_bin|-logrotate|.service)$|^\./(build|t3c-check-refs/test-files|testing|t3util/testing|tm-health-client/(config|tmagent)/test_files)/.*' -e modify -r . ; do
	T3C_PID="$(ps | grep t3c | grep -v grep | grep -v inotifywait | grep -v run.sh | tr -s ' ' | cut -d ' ' -f2)"
	if [[ ! -z "$T3"]]; then
		echo "$T3C_PID" | xargs kill;
	fi
	# TODO: is it even necessary to restart ATS?
	if [[ -f /var/trafficserver/server.lock ]]; then
		rm /var/trafficserver/server.lock;
	fi
	ps | grep traffic_server | grep -v grep | tr -s ' ' | cut -d ' ' -f2 | xargs kill
	su -c traffic_server ats &
	# for whatever reason, without this the repeated call to inotifywait will
	# sometimes lose track of th current directory. It spits out:
	# Couldn't watch .: No such file or directory
	# which is a bit odd.
	sleep 0.5
done;
