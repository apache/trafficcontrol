#!/bin/sh
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#

set -o errexit

cd "$TC/traffic_portal"

user=trafficportal
uid="$(stat -c%u "$TC")"
gid="$(stat -c%g "$TC")"
if [[ "$(id -u)" != "$uid" ]]; then
	if ! adduser -Du"$uid" "$user"; then
		user="$(cat /etc/passwd | grep :x:1000: | cut -d: -f1)"
	fi
	sed -Ei "s/^(${user}:.*:)[0-9]+(:)$/\1${gid}\2/" /etc/group
	chown "${uid}:${gid}" /usr/bin
	exec su "$user" -- "$0"
fi

npm ci
./node_modules/.bin/grunt
