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
# Based on CentOS 7.x
############################################################

FROM centos:7
MAINTAINER dev@trafficcontrol.apache.org

# Default values for TOMCAT RPM and RPM -- override with `docker build --build-arg JDK=...'
ARG TRAFFIC_ROUTER_RPM=traffic_router/traffic_router.rpm
ARG TOMCAT_RPM=traffic_router/tomcat.rpm

RUN yum -y install epel-release git rpm-build net-tools iproute nc wget tar unzip \
          perl-JSON perl-WWWCurl which make autoconf automake gcc gcc-c++ apr apr-devel \
          openssl openssl-devel bind-utils net-tools perl-JSON-PP gettext \
          java-1.8.0-openjdk-headless java-1.8.0-openjdk-devel && \
    yum -y clean all && \
    ln -sfv $(find /usr/lib/jvm -type d -name 'java-1.8.0*openjdk*') /opt/java

ADD $TRAFFIC_ROUTER_RPM /traffic_router.rpm
ADD $TOMCAT_RPM /tomcat.rpm

RUN yum -y install /traffic_router.rpm /tomcat.rpm

RUN find /usr/lib* -name libtc\* -exec ln -sfv '{}' /opt/tomcat/lib \;

ADD enroller/server_template.json \
    traffic_router/run.sh \
    traffic_ops/to-access.sh \
    /

EXPOSE 53 80 3333 3443

CMD /run.sh
