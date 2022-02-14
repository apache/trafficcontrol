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

cd "$TC/traffic_router"

mvn -Dmaven.test.skip=true compile -P \!rpm-build
mvn -Dmaven.test.skip=true package -P \!rpm-build

chmod -R a+rw "$TC/dev/traffic_router/"
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
