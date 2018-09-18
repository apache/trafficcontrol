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

ADD https://ci.trafficserver.apache.org/RPMS/CentOS7/trafficserver-7.1.4-1.el7.x86_64.rpm /trafficserver.rpm
ADD https://ci.trafficserver.apache.org/RPMS/CentOS7/trafficserver-devel-7.1.4-1.el7.x86_64.rpm /trafficserver-devel.rpm

RUN yum install -y kyotocabinet-libs epel-release rpmdevtools iproute net-tools nmap-ncat gettext autoconf automake libtool gcc-c++ glibc-devel openssl-devel
RUN yum install -y /trafficserver.rpm /trafficserver-devel.rpm jq python34-psutil python34-typing python34-setuptools python34-pip && yum clean all
RUN pip3 install --upgrade pip && pip3 install requests urllib3 distro

ADD traffic_server/plugins/astats_over_http/astats_over_http.c traffic_server/plugins/astats_over_http/Makefile.am /

RUN tsxs -v -c astats_over_http.c -o astats_over_http.so
RUN mkdir -p /usr/libexec/trafficserver && tsxs -v -o astats_over_http.so -i

RUN yum remove -y gcc-c++ glibc-devel autoconf automake libtool && rm -f /astats_over_http.c /Makefile.am

CMD exit
