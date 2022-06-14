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
# Dockerfile to build Traffic Ops container images
# Based on CentOS 8
############################################################

ARG OS_DISTRO=rockylinux
ARG OS_VERSION=8
FROM ${OS_DISTRO}:${OS_VERSION}
ARG OS_DISTRO
ARG OS_VERSION
ENV OS_DISTRO=${OS_DISTRO}
ENV OS_VERSION=${OS_VERSION}

RUN if [[ "${OS_VERSION%%.*}" -eq 7 ]]; then \
		yum -y install dnf || exit 1; \
	fi

RUN set -o nounset -o errexit && \
	mkdir -p /etc/cron.d; \
	if [[ "${OS_VERSION%%.*}" -eq 7 ]]; then \
		use_repo=''; \
		enable_repo=''; \
		# needed for llvm-toolset-7-clang, which is needed for postgresql13-devel-13.2-1PGDG, required by TO rpm
		dnf -y install gcc centos-release-scl-rh; \
	else \
		use_repo='--repo=pgdg13'; \
		enable_repo='--enablerepo=powertools'; \
	fi; \
	dnf -y install "https://download.postgresql.org/pub/repos/yum/reporpms/EL-${OS_VERSION%%.*}-x86_64/pgdg-redhat-repo-latest.noarch.rpm"; \
	# libicu required by postgresql13
	dnf -y install libicu; \
	dnf -y $use_repo -- install postgresql13; \
	dnf -y install epel-release; \
	dnf -y $enable_repo install      \
		bind-utils           \
		gettext              \
		# ip commands is used in set-to-ips-from-dns.sh
		iproute              \
		isomd5sum            \
		jq                   \
		libidn-devel         \
		libpcap-devel        \
		mkisofs              \
		net-tools            \
		nmap-ncat            \
		openssl              \
		perl-Crypt-ScryptKDF \
		perl-Digest-SHA1     \
		perl-JSON-PP         \
		python3              \
		# rsync is used to copy certs in "Shared SSL certificate generation" step
		rsync;               \
	dnf clean all

EXPOSE 443

WORKDIR /opt/traffic_ops/app

RUN yum -y install procps
ADD traffic_ops/traffic_ops*x86_64.rpm /traffic_ops.rpm
RUN yum -y install /traffic_ops.rpm && \
	rm /traffic_ops.rpm

ADD traffic_ops/run.sh /

EXPOSE 443
CMD /run.sh
