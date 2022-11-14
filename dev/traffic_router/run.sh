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

set -o errexit -o nounset

cd "$TC/traffic_router"

user=trafficrouter
uid="$(stat -c%u "$TC")"
gid="$(stat -c%g "$TC")"
if [[ "$(id -u)" != "$uid" ]]; then
	for dir in "${TC}/.m2"  */target; do
		if [[ -e "$dir" ]] && [[ "$(stat -c%u "$dir")" -ne "$uid" || "$(stat -c%g "$dir")" -ne "$gid" ]] ; then
			chown -R "${uid}:${gid}" "$dir"
		fi
	done

	adduser -Du"$uid" "$user"
	sed -Ei "s/^(${user}:.*:)[0-9]+(:)$/\1${gid}\2/" /etc/group
	chown -R "${uid}:${gid}" /opt
	exec su "$user" -- "$0"
fi

mvn -Dmaven.test.skip=true compile package -P \!rpm-build

cd "$TC/dev/traffic_router"
exec /opt/tomcat/bin/catalina.sh jpda run
