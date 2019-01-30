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

set -eu

domain="ciab.test"
shared_dns_dir="/shared/dns"
dns_key_file_name="K${domain}.private"

my_host="$(hostname -s)"
my_interface="$(ifconfig | egrep "^\w.+ " | grep -v lo | awk '{print $1}')"
my_interface=${my_interface%":"}
my_ip="$(ifconfig $my_interface | grep 'inet ' | awk '{print $2}')" # TODO determine if this should be 'hostname -I'
my_ip=${my_ip#"addr:"}
my_ip6="$(ifconfig $my_interface | egrep 'inet6 .+<global>' | awk '{print $2}')"

full_sub_domain="infra.${domain}"
my_fqdn="${my_host}.${full_sub_domain}"

ttl="86400"

nsupdate_remove_cmd="update delete ${my_fqdn} A"
nsupdate_remove_cmd6="update delete ${my_fqdn} AAAA"

nsupdate_cmd="update add ${my_fqdn} ${ttl} A ${my_ip}"

nsupdate_cmd6=
if [ -n "$my_ip6" ] ; then
	nsupdate_cmd6="update add ${my_fqdn} ${ttl} AAAA ${my_ip6}"
fi

dns_key_path="$(ls ${shared_dns_dir}/*.private || true)"
while [ -z "${dns_key_path}" ]; do
	printf "insert-self-into-dns waiting for dns server to place key\n"
	sleep 1
	dns_key_path="$(ls ${shared_dns_dir}/*private || true)"
done

printf "insert-self-into-dns domain $domain dns_key_path $dns_key_path my_host $my_host my_ip $my_ip my_fqdn $my_fqdn cmd '$nsupdate_cmd'\n"

nsupdate -v -k "${dns_key_path}" << EOF
${nsupdate_remove_cmd}
${nsupdate_remove_cmd6}
${nsupdate_cmd}
${nsupdate_cmd6}
show
send
EOF
