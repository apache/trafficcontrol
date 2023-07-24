#!/usr/bin/env bash

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

export_environment() {
	while IFS= read -r line; do
		export "$line";
	done < <(sed '/#/d' /ciab.env)
}

maybe_debug() {
	set -o errexit -o nounset
	hostname="$1"
	shift
	debug_port="$1"
	shift
	actual_binary="$1"
	shift
	hostname="${hostname//-/_}" # replace - with _
	hostname="${hostname^^}" # uppercase
	debug_variable_name="T3C_DEBUG_COMPONENT_${hostname}"
	if [[ "${!debug_variable_name}" == "${actual_binary%.actual}" ]]; then
		command=(dlv --listen=":${debug_port}" --headless=true --api-version=2 exec "/usr/bin/${actual_binary}" --)
	else
		command=("$actual_binary")
	fi
	exec "${command[@]}" "$@"
}

hostname="$(hostname --short)"
for t3c_tool in $(compgen -c t3c | sort | uniq); do
	(
		path="$(type -p "$t3c_tool")"
		cd "$(dirname "$path")"
		dlv_script="${t3c_tool}.debug"
		actual_binary="${t3c_tool}.actual"
		<<-DLV_SCRIPT cat > "$dlv_script"
		#!/usr/bin/env bash
		$(type export_environment | tail -n+2)
		$(type maybe_debug | tail -n+2)
		export_environment
		maybe_debug "${hostname}" "${DEBUG_PORT}" "${actual_binary}" "\$@"
		DLV_SCRIPT
		chmod +x "$dlv_script"
		mv "$t3c_tool" "$actual_binary"
		ln -s "$dlv_script" "$t3c_tool"
	)
done
