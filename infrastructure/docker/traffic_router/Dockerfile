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

FROM centos:6.6
MAINTAINER Dan Kirkwood

# Default values for JDK and RPM -- override with `docker build --build-arg JDK=...'
ARG JDK=http://download.oracle.com/otn-pub/java/jdk/8u92-b14/jdk-8u92-linux-x64.rpm
ARG RPM=http://traffic-control-cdn.net/downloads/1.6.0/RELEASE-1.6.0/traffic_router-1.6.0-3503.4899d302.x86_64.rpm

RUN yum install -y wget tar unzip perl-JSON perl-WWW-Curl which

# NOTE: Can't just use ADD because of the required cookie
RUN curl -O -L -b "oraclelicense=a" $JDK
ADD $RPM /

RUN yum install -y /$(basename $JDK) /$(basename $RPM)
RUN rm /$(basename $JDK) /$(basename $RPM)

EXPOSE 53 80 3333
ADD run.sh /
ENTRYPOINT /run.sh
