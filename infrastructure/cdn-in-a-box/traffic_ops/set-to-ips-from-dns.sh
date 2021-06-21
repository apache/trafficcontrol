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
#

base_data_dir="/traffic_ops_data"
servers_dir="${base_data_dir}/servers"
profiles_dir="${base_data_dir}/profiles"

service_names='db trafficops trafficportal trafficmonitor trafficvault trafficrouter enroller dns'

service_domain='infra.ciab.test'

gateway_ip="$(ip route | grep default | cut -d' ' -f3)"
gateway_ip6="$(ip -6 route | grep default | cut -d' ' -f3)"

while [ -z "${gateway_ip}" ]; do
	printf "setting ips from dns: service gateway ip not found! Trying again in 1s\n"
	sleep 1
	gateway_ip="$(ip route | grep default | cut -d' ' -f3)"
	gateway_ip6="$(ip -6 route | grep default | cut -d' ' -f3)"
done

service_ips="${gateway_ip}"
service_ip6s="${gateway_ip6}"
INTERFACE=$(ip link | awk '/\<UP\>/ && !/LOOPBACK/ {sub(/@.*/, "", $2); print $2}')
NETMASK=$(route | awk -v INTERFACE=$INTERFACE '$8 ~ INTERFACE && $1 !~ "default"  {print $3}')
DIG_IP_RETRY=10

for service_name in $service_names; do
	service_fqdn="${service_name}.${service_domain}"

	if [[ ! -e /shared/SKIP_DIG_IP ]]; then
		for (( i=1; i<=DIG_IP_RETRY; i++ )); do
			service_ip="$(dig +short ${service_fqdn} A)"
			if [ -z "${service_ip}" ]; then
				printf "service \"${service_fqdn}\" not found in dns, count=$i, waiting ...\n"
				sleep 3
			else
				break
			fi
		done
	fi

  #
	# TODO add a way to determine if a service wasn't built in the Compose,
	#      so it's possible to Compose only e.g. TO and not everything. Ideas:
	#      1. only wait so long, e.g. 30s. Not ideal, slow, inaccurate
	#      2. dig the Docker DNS name, not the FQDN
	#      3. run this in a cron, with the cron somehow also managing the enroller/init
	#
	if [ -z "${service_ip}" ]; then
		# TODO sleep and try again? Up to n times?
		printf "setting ips from dns: service \"${service_fqdn}\" not found in dns, skipping!\n"
		continue
	fi

	service_ip6="$(dig +short $service_name AAAA)"

	service_ips="${service_ips} ${service_ip}"
	if [ -n "${service_ip6}" ]; then
		service_ip6s="${service_ip6s} ${service_ip6}"
	fi

	# not all services have server files
	printf "setting ips from dns: checking file for dir '${servers_dir}' service '${service_name}'\n"
	service_file="$(ls ${servers_dir}/*-${service_name}* 2>/dev/null)"
	printf "setting ips from dns: trying service file '${service_file}'\n"
	if [ -n "${service_file}" ]; then
		printf "setting ips from dns: service file '${service_file}' exists, adding IPs\n"
		cat "${service_file}" | jq '. + {"ipAddress":"'"${service_ip}"'"}' > "${service_file}.tmp" && mv "${service_file}.tmp" "${service_file}"
		cat "${service_file}" | jq '. + {"ipGateway":"'"${gateway_ip}"'"}' > "${service_file}.tmp" && mv "${service_file}.tmp" "${service_file}"
		cat "${service_file}" | jq '. + {"ipNetmask":"'"${NETMASK}"'"}' > "${service_file}.tmp" && mv "${service_file}.tmp" "${service_file}"
		if [ -n "${service_ip6}" ]; then
			cat "${service_file}" | jq '. + {"ip6Address":"'"${service_ip6}"'"}' > "${service_file}.tmp" && mv "${service_file}.tmp" "${service_file}"
		fi
		if [ -n "${gateway_ip6}" ]; then
			cat "${service_file}" | jq '. + {"ip6Gateway":"'"${gateway_ip6}"'"}' > "${service_file}.tmp" && mv "${service_file}.tmp" "${service_file}"
		fi

		rm -rf "${service_file}.tmp"
	fi
done

ats_profile_type="ATS_PROFILE"

service_ips="$(echo "${service_ips}" | sed 's/^[[:blank:]]*//;s/[[:blank:]]*$//')" # trim
service_ip6s="$(echo "${service_ip6s}" | sed 's/^[[:blank:]]*//;s/[[:blank:]]*$//')" # trim

for profile_file in ${profiles_dir}/*.json; do
	profile_type="$(cat ${profile_file} | jq -r '.type')"
	if [ "${profile_type}" != "${ats_profile_type}" ]; then
		continue
	fi

	# get existing allow_ip, as space-separated
	existing_allow_ips="$(cat ${profile_file} | jq -r '.params | map(select(.name == "allow_ip")) | .[] | .value' 2>/dev/null | tr ',' ' ')"

	new_allow_ips="${existing_allow_ips} ${service_ips}"
	new_allow_ips="$(echo "${new_allow_ips}" | sed 's/^[[:blank:]]*//;s/[[:blank:]]*$//')" # trim
	new_allow_ips="$(echo "${new_allow_ips}" | tr -s ' ' | tr ' ' ',')" # replace spaces with commas, like ATS needs

	# delete existing allow_ip, and add new one
	cat ${profile_file} | jq '. + {params: (.params | map(select(.name != "allow_ip")))} | .params += [{configFile: "astats.config", name: "allow_ip", secure: false, value: "'"${new_allow_ips}"'"}]' > "${profile_file}.tmp" && mv "${profile_file}.tmp" "${profile_file}"


	# get existing allow_ip6, as space-separated
	existing_allow_ip6s="$(cat ${profile_file} | jq -r '.params | map(select(.name == "allow_ip6")) | .[] | .value' 2>/dev/null | tr ',' ' ')"

	new_allow_ip6s="${existing_allow_ip6s} ${service_ip6s}"
	new_allow_ip6s="$(echo "${new_allow_ip6s}" | sed 's/^[[:blank:]]*//;s/[[:blank:]]*$//')" # trim
	new_allow_ip6s="$(echo "${new_allow_ip6s}" | tr -s ' ' | tr ' ' ',')" # replace spaces with commas, like ATS needs

	# delete existing allow_ip, and add new one
	cat ${profile_file} | jq '. + {params: (.params | map(select(.name != "allow_ip6")))} | .params += [{configFile: "astats.config", name: "allow_ip6", secure: false, value: "'"${new_allow_ips}"'"}]' > "${profile_file}.tmp" && mv "${profile_file}.tmp" "${profile_file}"

	rm -rf "${profile_file}.tmp"
done
