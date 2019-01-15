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

set -eux

domain="ciab.test"

shared_dns_dir="/shared/dns"

dns_key_file_name="K${domain}.private"

dns_conf_file="/data/bind/etc/named.conf.local"

if cat "${dns_conf_file}" | grep "key \"${domain}.\" {" > /dev/null; then
	printf "set-dns-update: key already exists, not recreating\n"
	exit 0 # if the key exists from a previous docker run, don't recreate it.
fi

printf "set-dns-update: no key exists, creating\n"

# no key exists in the dns conf, but one might in the shared volume, rm just in case
rm -f K${domain}*
rm -f /shared/dns/*.private
rm -f /shared/dns/*.key

dns_key_name="$(dnssec-keygen -C -r /dev/urandom -a HMAC-MD5 -b 512 -n HOST "${domain}")"
dns_key_private="${dns_key_name}.private"
dns_key_key="${dns_key_name}.key"

dns_key_secret="$(cat ${dns_key_private} | grep '^Key' | awk '{print $2}')"

printf "waiting for self to serve dns...\n"
while ! dig +short "@$(hostname -s)" "$(hostname -s)"; do
	printf "waiting for self to serve dns...\n"
	sleep 1
done

cat << EOF >> "${dns_conf_file}"
key "${domain}." {
  algorithm hmac-md5;
  secret "${dns_key_secret}";
};
EOF

# origin_line="zone \"${domain}\" {"
# allow_update_line="  allow-update { key \"${domain}.\"; };"
# sed -i "s/${origin_line}/${origin_line}\n${allow_update_line}/" "${dns_conf_file}"

mkdir -p "${shared_dns_dir}"

/usr/sbin/rndc reload 2>&1 >> /rndc.log

# copy the key after reloading, so by the time clients get the key, it's usable.
cp "${dns_key_private}" "${shared_dns_dir}"
cp "${dns_key_key}" "${shared_dns_dir}"
