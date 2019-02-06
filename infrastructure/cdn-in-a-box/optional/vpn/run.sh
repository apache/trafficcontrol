#!/bin/bash
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

set-dns.sh
insert-self-into-dns.sh

set -x
set -e

INTERFACE=$(find /sys/class/net -type l '!' -name lo -exec basename '{}' \;)
NETWORK=$(route |grep -v default |grep $INTERFACE |awk '{print $1}')
NETMASK=$(route |grep -v default |grep $INTERFACE |awk '{print $3}')
DNSADDR=$(dig +short dns)

if [[ -z "$PRIVATE_NETWORK" ]] || [[ -z "$PRIVATE_NETMASK" ]]; then
  if [[ $NETWORK == 192.* ]] || [[ $HOST == 10.* ]]; then
    PRIVATE_NETWORK="172.16.127.0"
  else
    PRIVATE_NETWORK="10.16.127.0"
  fi
  PRIVATE_NETMASK="255.255.255.240"
fi



# Check if vpn ca existed
if [ ! -f "/vpnca/completed" ]; then
  cp /vars /root/EasyRSA-*
  cd /root/EasyRSA-*
  ./easyrsa init-pki
  ./easyrsa build-ca nopass
  ./easyrsa gen-dh
  ./easyrsa gen-req vpnserver nopass
  ./easyrsa sign-req server vpnserver
  ./easyrsa build-client-full vpnclient01 nopass
  cd pki
  cp ca.crt dh.pem issued/vpn* private/vpn* /vpnca
  openvpn --genkey --secret /vpnca/tls.key
  cat <<-EOF > /vpnca/client.ovpn
client
proto tcp
dev tun
remote REALHOSTIP REALPORT
resolv-retry infinite
nobind
persist-key
persist-tun
ns-cert-type server
comp-lzo
verb 3

<ca>
$(cat /vpnca/ca.crt)
</ca>
<cert>
$(cat /vpnca/vpnclient01.crt)
</cert>
<key>
$(cat /vpnca/vpnclient01.key)
</key>
key-direction 1
<tls-auth>
$(cat /vpnca/tls.key)
</tls-auth>
EOF
  touch /vpnca/completed

  # Update the permissions to be read/write for all
  find /vpnca -type d -exec chmod a+rwx '{}' \;
  find /vpnca -type f -exec chmod a+rw '{}' \;

  # Update the ownership to be nobody:nogroup
  chown -R nobody:nogroup /vpnca
fi

echo 1 > /proc/sys/net/ipv4/ip_forward
iptables -t nat -I POSTROUTING 1 -s $PRIVATE_NETWORK/$PRIVATE_NETMASK -o $INTERFACE -j MASQUERADE

/usr/sbin/openvpn --status /server.status 10 --cd /etc/openvpn --script-security 2 --config /etc/openvpn/server.conf --push "dhcp-option DNS $DNSADDR" --push "route $NETWORK $NETMASK" --server $PRIVATE_NETWORK $PRIVATE_NETMASK
