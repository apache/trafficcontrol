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

ARG RHEL_VERSION=8
FROM centos:${RHEL_VERSION}
ARG RHEL_VERSION=8
# Makes RHEL_VERSION available in later layers without needing to specify it again
ENV RHEL_VERSION=$RHEL_VERSION

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
		yum -y install dnf || exit 1; \
	fi

RUN set -o nounset -o errexit && \
	mkdir -p /etc/cron.d; \
	if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
		use_repo=''; \
		enable_repo=''; \
	else \
		use_repo='--repo=pgdg96'; \
		enable_repo='--enablerepo=powertools'; \
	fi; \
	dnf -y install "https://download.postgresql.org/pub/repos/yum/reporpms/EL-${RHEL_VERSION%%.*}-x86_64/pgdg-redhat-repo-latest.noarch.rpm"; \
	dnf -y $use_repo -- install postgresql96; \
	dnf -y install epel-release; \
	dnf -y $enable_repo install      \
		jq              \
		bind-utils      \
		net-tools       \
		gettext         \
		perl-JSON-PP    \
		perl-Crypt-ScryptKDF  \
		mkisofs         \
		isomd5sum       \
		nmap-ncat       \
		openssl         \
		python3         \
		golang          \
		git             \
		libidn-devel    \
		libpcap-devel   \
		perl-Digest-SHA1\
		# Used to copy certs in "Shared SSL certificate generation" step
		rsync;          \
	dnf clean all

EXPOSE 443

ADD traffic_router/core/src/test/resources/geo/GeoLite2-City.mmdb.gz /opt/traffic_ops/app/public/

WORKDIR /opt/traffic_ops/app
ADD traffic_ops/install/bin/install_goose.sh ./
RUN ./install_goose.sh && rm ./install_goose.sh && dnf -y remove git && dnf clean all

ADD infrastructure/cdn-in-a-box/traffic_ops_data /traffic_ops_data

# Override TRAFFIC_OPS_RPM arg to use a different one using --build-arg TRAFFIC_OPS_RPM=...  Can be local file or http://...
#
ARG TRAFFIC_OPS_RPM=infrastructure/cdn-in-a-box/traffic_ops/traffic_ops.rpm

COPY $TRAFFIC_OPS_RPM /traffic_ops.rpm
RUN yum -y install /traffic_ops.rpm && \
	rm /traffic_ops.rpm

COPY infrastructure/cdn-in-a-box/enroller/server_template.json \
	infrastructure/cdn-in-a-box/traffic_ops/config.sh \
	infrastructure/cdn-in-a-box/traffic_ops/run-go.sh \
	infrastructure/cdn-in-a-box/traffic_ops/to-access.sh \
	infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
	infrastructure/cdn-in-a-box/dns/set-dns.sh \
	infrastructure/cdn-in-a-box/traffic_ops/set-to-ips-from-dns.sh \
	infrastructure/cdn-in-a-box/traffic_ops/generate-certs.sh \
	infrastructure/cdn-in-a-box/traffic_ops/adduser.pl \
	infrastructure/cdn-in-a-box/traffic_ops/trafficops-init.sh \
	infrastructure/cdn-in-a-box/variables.env \
	/

EXPOSE 443
CMD /run-go.sh
HEALTHCHECK --interval=10s --timeout=1s \
	CMD bash -c 'source /to-access.sh && [[ "$(curl -sk "https://${TO_FQDN}/api/${TO_API_VERSION}/ping" | jq .ping)" == \"pong\" ]]'
