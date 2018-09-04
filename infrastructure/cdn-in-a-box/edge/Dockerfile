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
# Dockerfile to build Edge-Tier Cache container images for
# Apache Traffic Control
# Based on CentOS 7.2
############################################################

FROM centos:7

EXPOSE 80

ARG RPM=https://ci.trafficserver.apache.org/RPMS/CentOS7/trafficserver-7.1.2-4.el7.centos.x86_64.rpm

RUN yum install -y /trafficserver.rpm iproute net-tools nmap-ncat epel-release
RUN yum install -y jq python34-psutil python34-typing python34-setuptools python34-pip && yum clean all
RUN pip3 install --upgrade pip && pip3 install requests urllib3 distro

RUN mkdir -p /var/trafficserver/cache /opt/ort && chown -R ats:ats /etc/trafficserver/ /var/trafficserver/ /opt/ort
ADD edge/parameters.json edge/profile.json edge/server.json edge/run.sh traffic_ops/to-access.sh /

ADD edge/traffic_ops_ort.py /opt/ort

# You need to do this because the RPM in the ATS archives is just all kinds of messed-up
RUN chmod 755 /usr/lib64/trafficserver /etc/trafficserver/body_factory /etc/trafficserver/body_factory/default

CMD /run.sh

