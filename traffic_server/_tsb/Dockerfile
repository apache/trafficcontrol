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
    # Change BASE_IMAGE to centos when RHEL_VERSION=7
ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION} as ats-common-dependencies
ARG RHEL_VERSION=8

MAINTAINER dev@trafficcontrol.apache.org

### Common for all sub-component builds
RUN if [[ ${RHEL_VERSION%%.*} -ge 8 ]]; then \
		yum install -y 'dnf-command(config-manager)'; \
		yum config-manager --set-enabled powertools; \
	else \
		yum install -y deltarpm centos-release-scl-rh; \
		yum-config-manager --enable rhel-server-rhscl-7-rpms; \
	fi && \
	rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-* && \
	rpm --import "https://dl.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-${RHEL_VERSION%%.*}" && \
	yum -y clean all && \
	yum -y update ca-certificates && \
	yum -y install \
		git \
		rpm-build \
		rsync \
		epel-release && \
	yum -y clean all

### ats specific requirements
FROM ats-common-dependencies AS build-ats-specific
ARG RHEL_VERSION=8
# Makes RHEL_VERSION accessible to CMD
ENV RHEL_VERSION="$RHEL_VERSION"

RUN if [[ ${RHEL_VERSION%%.*} -ge 8 ]]; then \
		os_pkgs=( \
			brotli \
			brotli-devel \
			curl \
			gcc-toolset-11 \
			gcc-toolset-11-runtime \
			jansson \
			jansson-devel \
			libmaxminddb \
			libmaxminddb-devel); \
		os_toolset="gcc-toolset-11"; \
	else \
		os_pkgs=(devtoolset-11); \
		os_toolset="devtoolset-11"; \
	fi \
	&& yum install -y \
		${os_pkgs[*]} \
		autoconf \
		automake \
		ed \
		expat-devel \
		flex \
		gcc-c++ \
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
		sudo \
		tcl-devel \
		zlib \
		zlib-devel \
	&& yum clean all
COPY	jansson.pic.patch /opt/src/
COPY	cjose.pic.patch /opt/src/
RUN	pip3 install --user Sphinx
COPY	run.sh /run.sh
COPY	trafficserver.spec /rpmbuilddir/SPECS/trafficserver.spec
COPY	traffic_server_jemalloc /rpmbuilddir/SOURCES/traffic_server_jemalloc
COPY	trafficserver.env /rpmbuilddir/SOURCES/trafficserver.env
RUN	/usr/sbin/useradd -u 176 -r ats -s /sbin/nologin -d /
CMD if [[ ${RHEL_VERSION%%.*} -ge 8 ]]; then \
		os_toolset=gcc-toolset-11; \
		openssl_included='--without_openssl'; \
	else \
		os_toolset=devtoolset-7; \
		openssl_included='--with_openssl'; \
	fi \
	&& set -o pipefail; scl enable ${os_toolset} "./run.sh ${openssl_included}" 2>&1 | tee /rpmbuilddir/RPMS/build-trafficserver.log
