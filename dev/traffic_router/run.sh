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
user=trafficrouter
uid="$(stat -c%u .)"
gid="$(stat -c%g .)"
adduser -Du"$uid" "$user"
sed -Ei "s/^(${user}:.*:)[0-9]+(:)$/\1${gid}\2/" /etc/group
chown -R "${uid}:${gid}" /opt

su "$user" -- /usr/bin/mvn -Dmaven.test.skip=true compile package -P \!rpm-build

cd "$TC/dev/traffic_router"
exec su "$user" -- /opt/tomcat/bin/catalina.sh jpda run
