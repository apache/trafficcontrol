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

set -o errexit -o nounset

cd "$TC/experimental/traffic-portal"

user=tpv2
uid="$(stat -c%u "$TC")"
gid="$(stat -c%g "$TC")"
if [[ "$(id -u)" != "$uid" ]]; then
	for dir in "${TC}/.npm" .angular node_modules; do
		if [[ -e "$dir" ]] && [[ "$(stat -c%u "$dir")" -ne "$uid" || "$(stat -c%g "$dir")" -ne "$gid" ]] ; then
			chown -R "${uid}:${gid}" "$dir"
		fi
	done

	if ! adduser --disabled-password -u "$uid" "$user"; then
		user="$(cat /etc/passwd | grep :x:${uid}: | cut -d: -f1)"
	fi
	sed -Ei "s/^(${user}:.*:)[0-9]+(:)$/\1${gid}\2/" /etc/group
	chown "${uid}:${gid}" /usr/bin
	exec su "$user" -- "$0"
fi

npm i
./node_modules/.bin/ng serve --ssl --ssl-cert /server.crt --ssl-key /server.key --watch --proxy-config "$TC/dev/tpv2/proxy.json" --port 443 --host "::0" --live-reload
