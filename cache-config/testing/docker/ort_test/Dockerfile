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
# Dockerfile to build Traffic Server container images
#   as Edges for Traffic Control 1.4
# Based on CentOS 6.6
############################################################

# For cache, you may either use (RAM or disk) block devices or disk directories
# To use RAM block devices, pass them as /dev/ram0 and /dev/ram1 via `docker run --device`
# To use disk directories, simply don't pass devices, and the container will configure Traffic Server for directories

# Block devices may be created on the native machine with, for example, `modprobe brd`.
# The recommended minimum size for each block devices is 1G.
# For example, `sudo modprobe brd rd_size=1048576 rd_nr=2`

FROM centos:7
MAINTAINER dev@trafficcontrol.apache.org

RUN yum install -y \
  https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm \
  epel-release initscripts postgresql13.x86_64 git gcc lua-5.1.4-15.el7 lua-devel-5.1.4-15.el7 \
  ImageMagick-c++-devel && \
  # jq is used in run.sh to update tc-fixtures.json with the ATS RPM version
  yum install -y jq

ADD ort_test/trafficcontrol-cache-config*.rpm /trafficcontrol-cache-config.rpm
RUN yum install -y /trafficcontrol-cache-config.rpm

RUN sed -i 's/HOME\/bin/HOME\/bin:\/usr\/local\/go\/bin:/g' /root/.bash_profile
RUN echo "GOPATH=/root/go; export GOPATH" >> /root/.bash_profile
RUN echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config
RUN mkdir /root/go
RUN mkdir -p /var/lib/trafficcontrol-cache-config/status

EXPOSE 80 443
ADD ort_test/run.sh /
ADD ort_test/Ort-test.repo /etc/yum.repos.d
ADD ort_test/systemctl.sh /

ENTRYPOINT /run.sh
