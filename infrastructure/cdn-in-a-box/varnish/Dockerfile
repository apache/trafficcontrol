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

ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION} AS common-varnish-cache-config-layers
ARG RHEL_VERSION=8
# Makes RHEL_VERSION available at runtime
ENV RHEL_VERSION="$RHEL_VERSION"

RUN dnf install -y 'dnf-command(config-manager)'
RUN yum config-manager --set-enabled powertools
RUN yum install -y diffutils python3-sphinx

RUN yum install -y epel-release && yum install -y \
    make \
    autoconf \
    automake \
    jemalloc-devel \
    libedit-devel \
    libtool \
    libunwind-devel \
    ncurses-devel \
    pcre2-devel \
    pkgconfig \
    python3-docutils \
    cpio \
    git \
    perl \
    jq \
    gettext

RUN dnf install -y bind-utils kyotocabinet-libs initscripts iproute net-tools nmap-ncat gettext autoconf automake libtool gcc-c++ cronie glibc-devel openssl-devel && \
    dnf install -y logrotate && \
    dnf clean all

RUN curl -O https://varnish-cache.org/downloads/varnish-7.3.0.tgz && tar xf varnish-7.3.0.tgz
RUN cd varnish-7.3.0 && sh autogen.sh && sh configure && make && make install
RUN rm -rf varnish*

COPY infrastructure/cdn-in-a-box/varnish/run.sh infrastructure/cdn-in-a-box/traffic_ops/to-access.sh infrastructure/cdn-in-a-box/enroller/server_template.json /

COPY infrastructure/cdn-in-a-box/dns/set-dns.sh \
    infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
    /usr/local/sbin/
    

ARG ORT_RPM=infrastructure/cdn-in-a-box/cache/trafficcontrol-cache-config.rpm
COPY $ORT_RPM /
RUN rpm -Uvh /$(basename $ORT_RPM) &&\
    rm /$(basename $ORT_RPM)

CMD /run.sh

FROM common-varnish-cache-config-layers AS mid
ENV CACHE_TYPE=mid
COPY infrastructure/cdn-in-a-box/mid/init.d/ /opt/init.d/

FROM common-varnish-cache-config-layers AS edge
ENV CACHE_TYPE=edge
COPY infrastructure/cdn-in-a-box/edge/init.d/ /opt/init.d/


