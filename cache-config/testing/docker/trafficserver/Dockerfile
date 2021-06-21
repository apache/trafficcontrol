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

###############################################################
# Dockerfile to build Traffic Server RPM
###############################################################

ARG RHEL_VERSION=8
FROM centos:${RHEL_VERSION}
ARG RHEL_VERSION=8
# Makes RHEL_VERSION available in later layers without needing to specify it again
ENV RHEL_VERSION=$RHEL_VERSION

MAINTAINER dev@trafficcontrol.apache.org

VOLUME /atsbuild

RUN echo "Image Version: ${RHEL_VERSION}"

### Common for all sub-component builds
RUN set -o errexit; \
    if [[ ${RHEL_VERSION%%.*} -ge 8 ]]; then \
		rpm_gpg_key=RPM-GPG-KEY-centosofficial; \
		yum install -y 'dnf-command(config-manager)'; \
		yum config-manager --set-enabled powertools; \
	else \
		rpm_gpg_key="RPM-GPG-KEY-CentOS-${RHEL_VERSION%%.*}"; \
		yum install -y deltarpm centos-release-scl-rh; \
		yum-config-manager --enable rhel-server-rhscl-7-rpms; \
	fi && \
	rpm --import "/etc/pki/rpm-gpg/${rpm_gpg_key}" && \
	rpm --import "https://dl.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-${RHEL_VERSION%%.*}" && \
	yum -y clean all && \
	yum -y update ca-certificates && \
	yum -y install \
		git \
		rpm-build \
		rsync \
		epel-release && \
	yum -y clean all

RUN if [[ ${RHEL_VERSION%%.*} -ge 8 ]]; then \
		os_pkgs=( \
			brotli \
			brotli-devel \
			curl \
			gcc-toolset-9 \
			gcc-toolset-9-runtime \
			jansson \
			jansson-devel \
			libmaxminddb \
			libmaxminddb-devel); \
		os_toolset="gcc-toolset-9"; \
	else \
		os_pkgs=(devtoolset-9); \
	fi \
	&& yum install -y \
		${os_pkgs[*]} \
		autoconf \
		automake \
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

ADD  traffic_server/plugins/astats_over_http /astats_over_http

COPY cache-config/testing/docker/trafficserver/traffic_server_jemalloc \
     cache-config/testing/docker/trafficserver/trafficserver.spec \
     cache-config/testing/docker/trafficserver/cjose.pic.patch  \
     cache-config/testing/docker/trafficserver/jansson.pic.patch \
     cache-config/testing/docker/trafficserver/run.sh \ 
    / 

CMD /run.sh
