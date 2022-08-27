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

# For cache, you may either use (RAM or disk) block devices or disk directories
# To use RAM block devices, pass them as /dev/ram0 and /dev/ram1 via `docker run --device`
# To use disk directories, simply don't pass devices, and the container will configure Traffic Server for directories

# Block devices may be created on the native machine with, for example, `modprobe brd`.
# The recommended minimum size for each block devices is 1G.
# For example, `sudo modprobe brd rd_size=1048576 rd_nr=2`

ARG OS_DISTRO=rockylinux
ARG OS_VERSION=8
FROM ${OS_DISTRO}:${OS_VERSION}
ARG OS_DISTRO
ARG OS_VERSION
# Makes OS_VERSION available in later layers without needing to specify it again
ENV OS_VERSION=$OS_VERSION
ENV OS_DISTRO=$OS_DISTRO
MAINTAINER dev@trafficcontrol.apache.org
EXPOSE 80 443

# these packages load for both centos 7 and rockylinux 8
# no checks required at this time.
RUN yum install -y epel-release && yum repolist && \
  yum install -y brotli initscripts jansson jansson-devel \
    git gcc hwloc jq lua luajit man-db tcl \
    # libmaxminddb is required by trafficserver
    libmaxminddb \
    && \
  yum clean all

RUN echo "Using Image ${OS_DISTRO}:${OS_VERSION}"

ADD health-check-test/systemctl.sh /
RUN cp /usr/bin/systemctl /usr/bin/systemctl.save
RUN cp /systemctl.sh /usr/bin/systemctl && chmod 0755 /usr/bin/systemctl

# Note that if more than one t3c RPM matches this wildcard, this Dockerfile will
# break because this will create a directory instead of an RPM file, which it
# will then fail to install.
ADD health-check-test/trafficserver-[0-9]*.rpm /trafficserver.rpm
ADD health-check-test/trafficcontrol-health-client*x86_64.rpm /trafficcontrol-health-client.rpm
RUN rpm -i /trafficserver.rpm /trafficcontrol-health-client.rpm 
ADD health-check-test/tc-health-client.json /etc/trafficcontrol/

RUN sed -i 's/HOME\/bin/HOME\/bin:\/usr\/local\/go\/bin:/g' /root/.bash_profile &&\
  echo "GOPATH=/root/go; export GOPATH" >> /root/.bash_profile &&\
  echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config &&\
  mkdir /root/go

ADD variables.env /
ADD health-check-test/run.sh /
ADD health-check-test/parent.config /
ADD health-check-test/strategies.yaml /
RUN chmod +x /run.sh

ENTRYPOINT /run.sh
