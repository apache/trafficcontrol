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
# Based on Rocky Linux 8
############################################################

    # Change BASE_IMAGE to centos when RHEL_VERSION=7
ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION} as trafficops-dependencies
ARG RHEL_VERSION=8
# Makes RHEL_VERSION available in later layers without needing to specify it again
ENV RHEL_VERSION="$RHEL_VERSION"

VOLUME /traffic_ops_data

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
		yum -y install dnf || exit 1; \
	fi

RUN set -o nounset -o errexit -o xtrace && \
	mkdir -p /etc/cron.d; \
	if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
		use_repo=''; \
		enable_repo=''; \
	else \
		use_repo='--repo=pgdg13'; \
		enable_repo='--enablerepo=powertools'; \
	fi; \
	dnf -y install epel-release libicu; \
	dnf -y install "https://download.postgresql.org/pub/repos/yum/reporpms/EL-${RHEL_VERSION%%.*}-$(rpm --eval %_arch)/pgdg-redhat-repo-latest.noarch.rpm"; \
	dnf -y $use_repo -- install postgresql13; \
	if [[ -n "$enable_repo" ]]; then \
		dnf -y remove pgdg-redhat-repo; \
	fi; \
	dnf -y $enable_repo install \
		bind-utils           \
		# find is required by to-access.sh
		findutils            \
		gettext              \
		# ip commands is used in set-to-ips-from-dns.sh
		iproute              \
		isomd5sum            \
		jq                   \
		net-tools            \
		nmap-ncat            \
		openssl              \
		procps-ng            \
		# rsync is used to copy certs in "Shared SSL certificate generation" step
		rsync                \

		# Traffic Ops dependencies (Not all needed for CDN in a Box, but all
		# required by the Traffic Ops RPM)
		cpanminus            \
		expat-devel          \
		libcap               \
		libcurl-devel        \
		libidn-devel         \
		libpcap-devel        \
		mkisofs              \
		openssl-devel        \
		perl-core            \
		perl-Crypt-ScryptKDF \
		perl-DBD-Pg          \
		perl-DBI             \
		perl-Digest-SHA1     \
		perl-JSON            \
		perl-libwww-perl     \
		perl-TermReadKey     \
		perl-Test-CPAN-Meta  \
		perl-WWW-Curl        \
		python3; \
	dnf clean all

FROM trafficops-dependencies AS trafficops

EXPOSE 443

WORKDIR /opt/traffic_ops/app

ADD infrastructure/cdn-in-a-box/traffic_ops_data /traffic_ops_data

# Override TRAFFIC_OPS_RPM arg to use a different one using --build-arg TRAFFIC_OPS_RPM=...  Can be local file or http://...
#
ARG TRAFFIC_OPS_RPM=infrastructure/cdn-in-a-box/traffic_ops/traffic_ops.rpm

COPY $TRAFFIC_OPS_RPM /traffic_ops.rpm
RUN rpm -Uvh /traffic_ops.rpm && \
	rm /traffic_ops.rpm

COPY .github/actions/to-integration-tests/intermediate.crt.pem \
	.github/actions/to-integration-tests/intermediate.key.pem \
	.github/actions/to-integration-tests/rootca.crt.pem \
	.github/actions/to-integration-tests/rootca.key.pem \
	.github/actions/to-integration-tests/server.crt.pem \
	.github/actions/to-integration-tests/server.key.pem \
	/etc/pki/tls/traffic_ops/

COPY infrastructure/cdn-in-a-box/enroller/server_template.json \
	infrastructure/cdn-in-a-box/traffic_ops/config.sh \
	infrastructure/cdn-in-a-box/traffic_ops/run-go.sh \
	infrastructure/cdn-in-a-box/traffic_ops/to-access.sh \
	infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
	infrastructure/cdn-in-a-box/dns/set-dns.sh \
	infrastructure/cdn-in-a-box/traffic_ops/set-to-ips-from-dns.sh \
	infrastructure/cdn-in-a-box/traffic_ops/generate-certs.sh \
	infrastructure/cdn-in-a-box/traffic_ops/trafficops-init.sh \
	infrastructure/cdn-in-a-box/variables.env \
	/

EXPOSE 443
CMD /run-go.sh
HEALTHCHECK --interval=10s --timeout=1s \
	CMD bash -c 'source /to-access.sh && [[ "$(curl -sk "https://${TO_FQDN}/api/${TO_API_VERSION}/ping" | jq .ping)" == \"pong\" ]]'

FROM trafficops-dependencies as get-delve

RUN dnf -y install golang && \
    go install github.com/go-delve/delve/cmd/dlv@latest

FROM trafficops AS trafficops-debug
COPY --from=get-delve /root/go/bin /usr/bin
COPY infrastructure/cdn-in-a-box/traffic_ops/tv.conf conf/production/

# Makes the default target skip the trafficops-debug stage
FROM trafficops
