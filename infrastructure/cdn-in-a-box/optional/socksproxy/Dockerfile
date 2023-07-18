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

############################################################
# Dockerfile to build optional CiaB Socks Proxy
# Based on Rocky Linux 8
############################################################
    # Change BASE_IMAGE to centos when RHEL_VERSION=7
ARG BASE_IMAGE=centos \
    RHEL_VERSION=7
FROM ${BASE_IMAGE}:${RHEL_VERSION}
ARG RHEL_VERSION=7

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        yum -y install dnf || exit 1; \
    fi

ARG DANTES_SRC=https://www.inet.no/dante/files/dante-1.4.3.tar.gz

RUN dnf install -y net-tools bind-utils iproute wget curl automake autoconf gcc make && \
    curl -Ls -o /tmp/dante.tar.gz $DANTES_SRC && \
    tar -C /usr/src -zxvpf /tmp/dante.tar.gz && \
    cd /usr/src/dante* && \
    ./configure --prefix=/usr && \
    make -j 4 && \
    make install && \
    groupadd -g 8062 sockd  && \
    useradd -m -u 8062 -g sockd sockd && \
    dnf remove -y automake autoconf gcc make && \
    dnf clean all && \
    rm -rf /tmp/*  

COPY optional/socksproxy/sockd.conf /etc
COPY optional/socksproxy/run.sh /

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

EXPOSE 1080

CMD ["/run.sh"]
