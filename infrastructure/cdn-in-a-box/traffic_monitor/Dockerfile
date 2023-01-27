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
# Dockerfile to build Traffic Monitor container images
# Based on CentOS
############################################################

    # Change BASE_IMAGE to centos when RHEL_VERSION=7
ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION} as trafficmonitor-dependencies
ARG RHEL_VERSION=8

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        yum -y install dnf || exit 1; \
    fi

RUN dnf install -y epel-release && \
    dnf install -y \
        jq \
        nmap-ncat \
        iproute \
        net-tools \
        # find is required by to-access.sh
        findutils \
        gettext \
        bind-utils \
        openssl \
        procps-ng \
        initscripts && \
    dnf clean all

FROM trafficmonitor-dependencies AS trafficmonitor

# Default values for RPM -- override with `docker build --build-arg RPM=...'
ARG RPM=traffic_monitor/traffic_monitor.rpm
ADD $RPM /
RUN rpm -Uvh  /$(basename $RPM) && \
    rm /$(basename $RPM)

RUN mkdir -p /opt/traffic_monitor/conf
ADD traffic_monitor/traffic_monitor.cfg /opt/traffic_monitor/conf/traffic_monitor.cfg.template

ADD enroller/server_template.json \
    traffic_ops/to-access.sh \
    /

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

EXPOSE 80
ADD traffic_monitor/run.sh /
CMD /run.sh
HEALTHCHECK --interval=10s --timeout=1s \
    CMD bash -c 'source /to-access.sh && [[ "$(curl -s http://trafficmonitor.infra.ciab.test/api/traffic-ops-uri)" == "$TO_URL" ]]'

FROM trafficmonitor-dependencies as get-delve

RUN dnf -y install golang && \
    go install github.com/go-delve/delve/cmd/dlv@latest

FROM trafficmonitor as trafficmonitor-debug
COPY --from=get-delve /root/go/bin /usr/bin

# Makes the default target skip the trafficmonitor-debug stage
FROM trafficmonitor
