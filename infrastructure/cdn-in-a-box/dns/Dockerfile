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

ENV BIND_USER=bind \
    BIND_VERSION=1:9.11.3 \
    DATA_DIR=/data

RUN echo 'APT::Install-Recommends 0;' >> /etc/apt/apt.conf.d/01norecommends \
 && echo 'APT::Install-Suggests 0;' >> /etc/apt/apt.conf.d/01norecommends \
 && apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y vim.tiny wget net-tools sudo net-tools ca-certificates unzip apt-transport-https \
 && rm -rf /var/lib/apt/lists/* && rm -rf /etc/apt/apt.conf.d/docker-gzip-indexes \
 && apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y \
        #to-access dependencies
        jq gettext \
        bind9=${BIND_VERSION}* bind9-host=${BIND_VERSION}* dnsutils \
 && rm -rf /var/lib/apt/lists/*

COPY dns/entrypoint.sh /sbin/entrypoint.sh
COPY dns/named.conf.local /etc/bind
COPY dns/named.conf.options /etc/bind
COPY dns/zone.ciab.test /etc/bind
COPY dns/zone.ip4.arpa /etc/bind
COPY dns/zone.ip6.arpa /etc/bind
COPY traffic_ops/to-access.sh /
COPY enroller/server_template.json /

COPY dns/set-self-dns.sh \
     dns/set-dns-update.sh \
     /usr/local/sbin/

RUN chmod 755 /sbin/entrypoint.sh

EXPOSE 53/udp 53/tcp
ENTRYPOINT ["/sbin/entrypoint.sh"]
CMD ["/usr/sbin/named"]
