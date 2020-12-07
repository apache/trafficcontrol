# Licensed to the Apache Software Fou:qndation (ASF) under one
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
#
# SPDX-License-Identifier: Apache-2.0
FROM centos:7

RUN	yum clean all \
	&& yum install -y deltarpm epel-release centos-release-scl-rh \
	&& yum-config-manager --enable rhel-server-rhscl-7-rpms \
	&& yum clean all \
	&& yum install -y \
		autoconf \
		automake \
		devtoolset-7 \
		ed \
		expat-devel \
		flex \
		gcc-c++ \
		git \
		glibc-devel \
		hwloc \
		hwloc-devel \
		libcap-devel \
		libcurl-devel \
		libtool \
		libuuid-devel \
		lua-devel \
		luajit-devel \
		make \
		man \
		nano \
		ncurses-devel \
		nmap-ncat \
		openssl \
		openssl-devel \
		pcre \
		pcre-devel \
		perl-Digest-SHA \
		perl-ExtUtils-MakeMaker \
		perl-URI \
		pkgconfig \
		python3 \
		rpm-build \
		sudo \
		tcl-devel \
		zlib \
		zlib-devel \
	&& yum clean all

COPY	jansson.pic.patch /opt/src/
COPY	cjose.pic.patch /opt/src/
ADD	https://bootstrap.pypa.io/get-pip.py /
RUN	python get-pip.py
RUN	pip install --user Sphinx
COPY	run.sh /
COPY	trafficserver.spec /rpmbuilddir/SPECS/trafficserver.spec
COPY	traffic_server_jemalloc /rpmbuilddir/SOURCES/traffic_server_jemalloc
RUN	/usr/sbin/useradd -u 176 -r ats -s /sbin/nologin -d /
CMD	set -o pipefail; scl enable devtoolset-7 ./run.sh 2>&1 | tee /rpmbuilddir/RPMS/x86_64/build.log
