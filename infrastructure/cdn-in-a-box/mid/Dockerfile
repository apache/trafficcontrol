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
# Dockerfile to build Mid-Tier Cache container images for
# Apache Traffic Control
# Based on CentOS 7.2
############################################################

FROM centos:7

ADD https://ci.trafficserver.apache.org/RPMS/CentOS7/trafficserver-7.1.2-4.el7.centos.x86_64.rpm /trafficserver-7.1.2-4.el7.centos.x86_64.rpm

RUN yum install -y /trafficserver-7.1.2-4.el7.centos.x86_64.rpm iproute net-tools nmap-ncat && yum clean all

RUN chown -R ats:ats /etc/trafficserver/
RUN mkdir -p /var/trafficserver/cache
RUN chown -R ats:ats /var/trafficserver/
RUN mkdir -p /opt/ort
RUN chown ats:ats /opt/ort
ADD mid_parameters.json /mid_parameters.json
ADD mid_profile.json /mid_profile.json
ADD server.json /server.json
ADD setup.sh /setup.sh

EXPOSE 80

CMD /setup.sh
