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
FROM ${BASE_IMAGE}:${RHEL_VERSION} as os-dependencies
ARG RHEL_VERSION=8

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        yum -y install dnf || exit 1; \
    fi

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        utils_package=yum-utils; \
    else \
        utils_package=dnf-utils; \
    fi && \
    dnf install -y \
        bind-utils \
        # chromium and jq are in EPEL
        epel-release \
        # find is required by to-access.sh
        findutils \
        GConf2 \
        git \
        java-1.8.0-openjdk \
        net-tools \
        $utils_package && \
    set -o pipefail && \
    dnf -y install chromium jq && \
    dnf -y clean all

FROM os-dependencies AS node-dependencies
# Download and install node
RUN set -o pipefail; \
    curl -sLf https://rpm.nodesource.com/setup_16.x | bash
RUN dnf -y install nodejs || ( \
        # rpm --rebuilddb after installing the nodejs is necessary on docker 20.10.8,
        # build 3967b7d as of 2023-01-19, though that was not always the case for that
        # docker version
        rpm --rebuilddb && \
        echo 'RPM DB rebuilt' && \
        dnf -y install nodejs \
    )
RUN npm i -g webdriver-manager

FROM node-dependencies

COPY    traffic_portal/ /lang/traffic_portal/
WORKDIR /lang/traffic_portal/test/integration

RUN npm ci

RUN mkdir /portaltestresults

RUN jq ' \
        .capabilities.chromeOptions.args += [ \
            "--disable-extensions", \
            "--ignore-certificate-errors", \
            "--no-sandbox", \
            "--whitelisted-ips=" \
        ] | \
        .getPageTimeout = 120000 | \
        .allScriptsTimeout = 300000 | \
        .defaultTimeoutInterval = 12000 | \
        .params.junitReporter = true' \
        config.json > conf.json.tmp && \
    mv conf.json.tmp config.json

# Install chromium
RUN dnf -y install chromium

RUN set -o pipefail && \
    webdriver-manager clean && \
    chromium-browser --version | \
        grep -o '[0-9.]\+' | \
        xargs webdriver-manager update --gecko false --versions.chrome

COPY infrastructure/cdn-in-a-box/traffic_ops/to-access.sh \
     infrastructure/cdn-in-a-box/traffic_portal_integration_test/run.sh \
     infrastructure/cdn-in-a-box/dns/set-dns.sh \
     infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
     /usr/local/sbin/

CMD run.sh
