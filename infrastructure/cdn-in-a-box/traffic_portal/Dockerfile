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
# Dockerfile to build Traffic Portal container images
# Based on Rocky Linux 8
############################################################

    # Change BASE_IMAGE to centos when RHEL_VERSION=7
ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION}
ARG RHEL_VERSION=8

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        yum -y install dnf || exit 1; \
    fi

RUN set -o pipefail && \
    curl -sLf https://rpm.nodesource.com/setup_16.x | bash

# Override TRAFFIC_PORTAL_RPM arg to use a different one using --build-arg TRAFFIC_PORTAL_RPM=...  Can be local file or http://...
ARG TRAFFIC_PORTAL_RPM=traffic_portal/traffic_portal.rpm

ARG TO_HOST=$TO_HOST

# Install and delete the TRAFFIC_PORTAL_RPM when finished
RUN dnf install -y epel-release && \
    dnf install -y \
      jq \
      nodejs \
      openssl \
      procps-ng \
      # find is required by to-access.sh
      findutils \
      gettext \
      bind-utils \
      net-tools && \
    dnf clean all

ADD $TRAFFIC_PORTAL_RPM /
RUN rpm -Uvh /$(basename $TRAFFIC_PORTAL_RPM) && \
    rm /$(basename $TRAFFIC_PORTAL_RPM)

ADD enroller/server_template.json \
    traffic_portal/run.sh \
    traffic_ops/to-access.sh /

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

EXPOSE 443

CMD /run.sh
