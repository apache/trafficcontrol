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

ARG OS_VERSION=8 \
	OS_DISTRO=rockylinux

FROM ${OS_DISTRO}:${OS_VERSION}
MAINTAINER dev@trafficcontrol.apache.org
EXPOSE 80 443

RUN yum install -y epel-release && \
	yum install -y initscripts git jq gcc lua nmap-ncat && \
	yum clean all && \
	# Add skip git's safe.directory check
	git config --global --add safe.directory '*'

# Note that if more than one t3c RPM matches this wildcard, this Dockerfile will
# break because this will create a directory instead of an RPM file, which it
# will then fail to install.
ADD ort_test/trafficcontrol-cache-config*.rpm /trafficcontrol-cache-config.rpm
RUN yum install -y /trafficcontrol-cache-config.rpm && yum clean all

RUN sed -i 's/HOME\/bin/HOME\/bin:\/usr\/local\/go\/bin:/g' /root/.bash_profile &&\
	echo "GOPATH=/root/go; export GOPATH" >> /root/.bash_profile &&\
	echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config &&\
	mkdir /root/go

ADD ort_test/Ort-test.repo /etc/yum.repos.d
ADD ort_test/run.sh ort_test/systemctl.sh /

ENTRYPOINT /run.sh
