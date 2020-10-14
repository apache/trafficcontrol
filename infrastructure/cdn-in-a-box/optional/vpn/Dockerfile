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

FROM ubuntu:18.04
RUN apt-get update && apt-get install -y vim tree wget iptables openvpn iputils-ping net-tools dnsutils
ENV EASYRSA_BIN https://github.com/OpenVPN/easy-rsa/releases/download/v3.0.5/EasyRSA-nix-3.0.5.tgz
RUN cd /root && \
    wget $EASYRSA_BIN && \
    tar -zxf $(basename $EASYRSA_BIN) && \
    rm $(basename $EASYRSA_BIN) && \
    mkdir /vpnca
ADD ./optional/vpn/server.conf /etc/openvpn/server.conf
ADD ./optional/vpn/run.sh ./optional/vpn/vars /

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

ENTRYPOINT /run.sh
EXPOSE 443
