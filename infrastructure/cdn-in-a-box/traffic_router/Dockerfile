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
# Dockerfile to build Traffic Router 3.0
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

MAINTAINER dev@trafficcontrol.apache.org

# Default values for TOMCAT RPM and RPM -- override with `docker build --build-arg JDK=...'
ARG TRAFFIC_ROUTER_RPM=traffic_router/traffic_router.rpm
ARG TOMCAT_RPM=traffic_router/tomcat.rpm

RUN dnf -y install epel-release && \
    dnf -y install \
        jq \
        # find is required by to-access.sh
        findutils \
        git \
        rpm-build \
        net-tools \
        iproute \
        nc \
        procps-ng \
        wget \
        tar \
        unzip \

        perl-JSON perl-WWW-Curl which make autoconf automake gcc gcc-c++ apr apr-devel \
        openssl openssl-devel bind-utils net-tools perl-JSON-PP gettext \
        java-11-openjdk-headless tzdata-java java-11-openjdk-devel tomcat-native && \
    dnf -y clean all && \
    ln -sfv $(realpath /usr/lib/jvm/java-11) /opt/java

ADD $TRAFFIC_ROUTER_RPM /traffic_router.rpm
ADD $TOMCAT_RPM /tomcat.rpm

RUN rpm -Uvh /traffic_router.rpm /tomcat.rpm && \
    find /usr/lib* -name libtc\* -exec ln -sfv {} /opt/tomcat/lib \;

ADD enroller/server_template.json \
    traffic_router/run.sh \
    traffic_ops/to-access.sh \
    /

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

EXPOSE 53 80 3333 3443

CMD /run.sh
