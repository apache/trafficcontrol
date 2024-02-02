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
# Dockerfile to build Traffic Stats container images
# Based on CentOS
############################################################

    # Change BASE_IMAGE to centos when RHEL_VERSION=7
ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION} AS trafficstats-dependencies
ARG RHEL_VERSION=8

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        yum -y install dnf || exit 1; \
    fi

RUN dnf install -y epel-release && \
    dnf install -y \
        jq \
        nmap-ncat \
        net-tools \
        # find is required by to-access.sh
        findutils \
        gettext \
        bind-utils \
        openssl \
        procps-ng && \
    dnf clean all

FROM trafficstats-dependencies AS trafficstats

# Default values for RPM -- override with `docker build --build-arg RPM=...'
ARG TRAFFIC_TS_RPM=traffic_stats/traffic_stats.rpm
ADD $TRAFFIC_TS_RPM /
RUN rpm -Uvh /$(basename $TRAFFIC_TS_RPM) && \
    rm /$(basename $TRAFFIC_TS_RPM)

ADD enroller/server_template.json \
    traffic_ops/to-access.sh \
    traffic_stats/run.sh \
    /

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

ENTRYPOINT /run.sh

FROM trafficstats-dependencies AS get-delve

RUN dnf -y install golang && \
    go install github.com/go-delve/delve/cmd/dlv@latest

FROM trafficstats AS trafficstats-debug
COPY --from=get-delve /root/go/bin /usr/bin

# Makes the default target skip the trafficstats-debug stage
FROM trafficstats
