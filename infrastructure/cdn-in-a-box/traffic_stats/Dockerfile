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

FROM centos:7

# Default values for RPM -- override with `docker build --build-arg RPM=...'
ARG TRAFFIC_TS_RPM=traffic_stats/traffic_stats.rpm

RUN yum install -y epel-release && \
    yum install -y \
        jq \
        nmap-ncat \
        net-tools \
        gettext \
        bind-utils \
        openssl && \
    yum clean all

ADD $TRAFFIC_TS_RPM /
RUN yum install -y /$(basename $TRAFFIC_TS_RPM) \
    rm /$(basename $TRAFFIC_TS_RPM)

ADD enroller/server_template.json \
    traffic_ops/to-access.sh \
    traffic_stats/run.sh \
    /

ENTRYPOINT /run.sh
