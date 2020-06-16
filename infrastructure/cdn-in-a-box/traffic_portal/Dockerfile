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
# Based on CentOS 7.2
############################################################

FROM centos:7

RUN curl -sL https://rpm.nodesource.com/setup_12.x | bash -

# Override TRAFFIC_PORTAL_RPM arg to use a different one using --build-arg TRAFFIC_PORTAL_RPM=...  Can be local file or http://...
ARG TRAFFIC_PORTAL_RPM=traffic_portal/traffic_portal.rpm

ARG TO_HOST=$TO_HOST

# Install and delete the TRAFFIC_PORTAL_RPM when finished
RUN yum install -y \
    epel-release && \
    yum install -y \
      jq \
      nodejs \
      openssl \
      gettext \
      bind-utils \
      net-tools && \
    yum clean all || \
    echo "ERROR INSTALLING PACKAGES"

ADD $TRAFFIC_PORTAL_RPM /
RUN yum install -y /$(basename $TRAFFIC_PORTAL_RPM) && \
    rm /$(basename $TRAFFIC_PORTAL_RPM)

ADD enroller/server_template.json \
    traffic_portal/run.sh \
    traffic_ops/to-access.sh /

CMD /run.sh
