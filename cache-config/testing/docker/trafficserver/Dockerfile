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

###############################################################
# Dockerfile to build Traffic Server RPM
# Based on CentOS 7 for ATS RPM ot match the ort_test 
# container running CentOS 7.
###############################################################

ARG RHEL_VERSION=7
FROM centos:${RHEL_VERSION}
ARG RHEL_VERSION=7
ENV RHEL_VERSION="$RHEL_VERSION"
VOLUME /atsbuild

# install the required toolsets
RUN set -o errexit; \
  if [[ "${RHEL_VERSION%%.*}" -le 7 ]]; then \
    yum -y install centos-release-scl; \
    enable_repo=''; \
    toolset_rpms='devtoolset-9'; \
  else \
    enable_repo='--enablerepo=powertools'; \
    toolset_rpms='gcc-toolset-9-gcc gcc-toolset-9-gcc-c++'; \
  fi; \
  yum -y $enable_repo install \
    autoconf \
    automake \
    curl \
    flex \
    git \
    hwloc \
    hwloc-devel \
    libtool \
    libcap \
    libcap-devel \
    lua \
    make \
    openssl \
    openssl-devel \
    pcre \
    pcre-devel \
    pkgconfig \
    python3 \
    redhat-rpm-config \
    rpm-build \
    tcl \
    tcl-devel \
    $toolset_rpms

COPY cache-config/testing/docker/trafficserver/trafficserver-8.spec \
     cache-config/testing/docker/trafficserver/trafficserver-9.spec \
     cache-config/testing/docker/trafficserver/run.sh \
     /

CMD /run.sh
