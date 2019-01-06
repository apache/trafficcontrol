#!/usr/bin/env bash

# TODO fix to work with docker-compose up --scale
# TODO fix duplication with docker-compose.yml
domain='infra.ciab.test'

# service_names must match the docker-compose service name, or hostname.
service_names='db trafficops trafficops-perl trafficportal trafficmonitor trafficvault trafficrouter edge mid origin enroller socksproxy client vnc dns'

bind_zone_dir='/data/bind/etc/'
bind_zone_file='zone.ciab.test'

bind_zone_file_path="${bind_zone_dir}/${bind_zone_file}"

origin="${domain}."
origin_line="\$ORIGIN ${origin}"

changed_zone_file="false"

function add_zone_entry {
	host="$1"
	ip="$2"
	record="$3"
	entry="${host}                IN ${record}    ${ip}"

	grep -E "^${host}\s+IN\s+${record}\s+${ip}" $bind_zone_file_path > /dev/null
	if [ $? -eq 0 ]; then
		: # echo "host \"$host\" ip \"$ip\" record \"$record\" already in zonefile"
	else
		echo "host \"$host\" ip \"$ip\" record \"$record\" not in zonefile, inserting"
		sed -i "s/${origin_line}/${origin_line}\n\n${entry}/" ${bind_zone_file_path}
		changed_zone_file="true"
	fi
}

gateway_ip="$(ip route | grep default | cut -d' ' -f3)"
gateway_ip6="$(ip -6 route | grep default | cut -d' ' -f3)"

if [ -z $gateway_ip ]; then
	echo "service gateway ip not found yet"
else
	add_zone_entry gw $gateway_ip A
fi

if [ -z $gateway_ip6 ]; then
	echo "service gateway ip6 not found yet"
else
	add_zone_entry gw $gateway_ip6 AAAA
fi

for service_name in $service_names; do
	service_ip="$(dig +short $service_name A)"
	if [ -z $service_ip ]; then
		echo "service $service_name ip not found yet"
	else
		add_zone_entry $service_name $service_ip A
	fi

	service_ip6="$(dig +short $service_name AAAA)"
	if [ -z $service_ip6 ]; then
		echo "service $service_name ip6 not found yet"
	else
		add_zone_entry $service_name $service_ip6 AAAA
	fi
done

if [ "$changed_zone_file" = "true" ] ; then
	/usr/sbin/rndc reload 2>&1 >> /rndc.log
	rndc_exit=$?
	while [ $rndc_exit -ne 0 ]; do
		printf "bind reload exit $rndc_exit failed, trying again in 2s\n" >> /rndc.log
		sleep 2
		/usr/sbin/rndc reload 2>&1 >> /rndc.log
		rndc_exit=$?
	done
fi
