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
# Dockerfile to build Traffic Router 1.4 container images
# Based on CentOS 6.6
############################################################

# Example Build and Run:
# docker build --rm --tag traffic_router:1.4 Traffic_Router
# docker build --rm --build-arg JDK=http://download.oracle.com/<path to jdk rpm> --build-arg RPM=<path to traffic_router rpm> --tag traffic_router:<version> traffic_router

# docker run --name my-traffic-router --hostname my-traffic-router --net cdnet --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --env TRAFFIC_MONITORS="my-traffic-monitor-0:80;my-traffic-monitor-1:80" --env ORIGIN_URI="http://my-hotair.cdnet" --detach traffic_router:1.4

FROM centos:7
MAINTAINER dev@trafficcontrol.apache.org

# Default values for TMCAT RPM and RPM -- override with `docker build --build-arg JDK=...'
ARG RPM=traffic_router.rpm
ARG TMCAT=tomcat.rpm
ARG TC_REPO=traffic-control.repo
ADD $TMCAT /
ADD $RPM /
ADD $TC_REPO /etc/yum.repos.d/
ADD starttr.sh /
ADD shutdowntr.sh /

### Common for all sub-component builds
RUN	yum -y install \
		epel-release \
		git \
		rpm-build && \
	yum -y clean all


RUN yum install -y wget tar unzip perl-JSON perl-WWW-Curl which
RUN yum search tomcat-native
RUN yum search jdk


RUN yum install -y /$(basename $TMCAT) /$(basename $RPM)
#RUN rm /$(basename $JDK) /$(basename $RPM)

EXPOSE 53 80 3333
ADD run.sh /
ENTRYPOINT /run.sh
