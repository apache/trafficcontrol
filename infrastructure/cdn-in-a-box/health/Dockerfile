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
# Change BASE_IMAGE to centos when RHEL_VERSION=7

ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION} as common-dependencies
ARG RHEL_VERSION=8
# Makes RHEL_VERSION available in later layers without needing to specify it again
ENV RHEL_VERSION="$RHEL_VERSION"

MAINTAINER dev@trafficcontrol.apache.org

# top level of trafficcontrol directory must be mounted as a volume:
# docker run --volume /trafficcontrol:$(pwd) ...
VOLUME /trafficcontrol

# CentOS 7 has pandoc in epel (which we installed above), 8 has it in [pP]ower[tT]ools
RUN if [[ ${RHEL_VERSION%%.*} -ge 8 ]]; then \
        yum -y install dnf-plugins-core && \
        (yum -y config-manager --set-enabled PowerTools || yum -y config-manager --set-enabled powertools) \
    fi && \
    yum -y install pandoc

### tc-health-client specific requirements
FROM common-dependencies AS tchealthclient

### Common for all sub-component builds
RUN rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-* && \
    rpm --import "https://dl.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-${RHEL_VERSION%%.*}" && \
    yum -y update ca-certificates && \
    yum -y install \
        jq \
        nodejs \
        openssl \
        procps-ng \
        git \
        rpm-build \
        rsync \
        gzip \
        # find is required by to-access.sh
        findutils \
        gettext \
        bind-utils \
        net-tools \
        epel-release && \
    yum -y clean all

ARG TCH_RPM=infrastructure/cdn-in-a-box/health/trafficcontrol-health-client.rpm
COPY $TCH_RPM /
RUN rpm -Uvh /$(basename $TCH_RPM) &&\
    rm /$(basename $TCH_RPM)

ADD enroller/server_template.json \
    health/run.sh \
    traffic_ops/to-access.sh /

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

RUN ["chmod", "+x", "/run.sh"]
CMD [ "/run.sh" ]
