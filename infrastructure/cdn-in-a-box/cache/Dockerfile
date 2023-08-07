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
# Based on Rocky Linux 8
############################################################

    # Change BASE_IMAGE to centos when RHEL_VERSION=7
ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION} AS common-traffic-server-dependencies
ARG RHEL_VERSION=8
# Makes RHEL_VERSION available at runtime
ENV RHEL_VERSION="$RHEL_VERSION"

# The /traffic_ops_data volume is included to set the ATS RPM version in the MID profile.
VOLUME /traffic_ops_data

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        yum -y install dnf || exit 1; \
    fi

EXPOSE 80

RUN dnf -y install epel-release && \
    if [[ "${RHEL_VERSION%%.*}" -lt 8 ]]; then \
        additional_packages='perl(URI)' || \
        exit 1; \
    else \
        # cjose is required by the ATS TSB RPM
        additional_packages='cjose' || \
        exit 1; \
    fi && \
    dnf -y install              \
        brotli                  \
        # find is required by to-access.sh
        findutils               \
        GeoIP                   \
        git                     \
        groff-base              \
        hwloc                   \
        # jansson is required by the ATS TSB RPM
        jansson                 \
        kyotocabinet-libs       \
        libmaxminddb            \
        libtool-ltdl            \
        libunwind               \
        libuuid                 \
        # luajit is required by the ATS TSB RPM
        luajit                  \
        lzo                     \
        make                    \
        numactl-libs            \
        openssl                 \
        procps-ng               \
        perl                    \
        tcl                     \
        xz                      \
        $additional_packages

RUN dnf install -y bind-utils kyotocabinet-libs initscripts iproute net-tools nmap-ncat gettext autoconf automake libtool gcc-c++ cronie glibc-devel openssl-devel && \
    dnf install -y jq logrotate && \
    dnf clean all

FROM common-traffic-server-dependencies AS common-cache-config-layers

COPY infrastructure/cdn-in-a-box/cache/trafficserver.rpm /trafficserver.rpm
COPY traffic_server/plugins/astats_over_http/astats_over_http.c traffic_server/plugins/astats_over_http/Makefile.am /

RUN set -o errexit -o nounset -o xtrace; \
    rpm -Uvh /trafficserver.rpm; \
    if [[ -e /trafficserver-devel.rpm ]]; then \
        rpm -Uvh /trafficserver-devel.rpm; \
    fi; \
    PATH+=:/opt/trafficserver/bin; \
    tsxs -v -c astats_over_http.c -o astats_over_http.so && \
    tsxs -v -o astats_over_http.so -i && \
    dnf remove -y gcc-c++ glibc-devel autoconf automake libtool && \
    rm -f /astats_over_http.c /Makefile.am

RUN set -o errexit -o nounset -o pipefail -o xtrace && \
    # The ssl directory needs to exist or the first sslkey file will be skipped on badass
    etc_trafficserver="$(rpm -qpl trafficserver.rpm | grep '/etc/trafficserver$')"; \
    if [[ ! -e "${etc_trafficserver}/ssl" ]]; then \
        echo 'Creating ssl directory...'; \
        mkdir "${etc_trafficserver}/ssl"; \
    fi; \
    # trafficserver needs different ownership than the trafficserver RPM sets
    chown -R ats:ats "$etc_trafficserver"

WORKDIR /opt

COPY infrastructure/cdn-in-a-box/cache/traffic_ops_ort.crontab /etc/cron.d/traffic_ops_ort-cron-template
COPY infrastructure/cdn-in-a-box/cache/traffic_ops_ort.logrotate /etc/logrotate.d/ort

COPY infrastructure/cdn-in-a-box/cache/run.sh infrastructure/cdn-in-a-box/traffic_ops/to-access.sh infrastructure/cdn-in-a-box/enroller/server_template.json /

COPY infrastructure/cdn-in-a-box/dns/set-dns.sh \
     infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
     /usr/local/sbin/
# Copy systemctl.sh to /usr/bin specifically because t3c runs /bin/systemctl by absolute path
COPY infrastructure/cdn-in-a-box/cache/systemctl.sh /usr/bin/systemctl

ARG ORT_RPM=infrastructure/cdn-in-a-box/cache/trafficcontrol-cache-config.rpm
COPY $ORT_RPM /
RUN rpm -Uvh /$(basename $ORT_RPM) &&\
    rm /$(basename $ORT_RPM)
CMD /run.sh

ARG TCH_RPM=infrastructure/cdn-in-a-box/health/trafficcontrol-health-client.rpm 
COPY $TCH_RPM /
RUN rpm -Uvh /$(basename $TCH_RPM) &&\
    rm /$(basename $TCH_RPM)

ADD infrastructure/cdn-in-a-box/health/tc-health-client.service /usr/lib/systemd/system
COPY infrastructure/cdn-in-a-box/health/tc-health-client.json  /etc/trafficcontrol/tc-health-client.json
COPY infrastructure/cdn-in-a-box/health/to-creds /etc/to-creds

FROM common-traffic-server-dependencies AS get-delve
RUN dnf -y install golang && \
    go install github.com/go-delve/delve/cmd/dlv@latest

FROM common-cache-config-layers AS mid
ENV CACHE_TYPE=mid
COPY infrastructure/cdn-in-a-box/mid/init.d/ /opt/init.d/

FROM mid AS mid-debug
COPY --from=get-delve /root/go/bin /usr/bin
COPY infrastructure/cdn-in-a-box/cache/init-debug-scripts.sh /opt/init.d/

FROM common-cache-config-layers AS edge
ENV CACHE_TYPE=edge
COPY infrastructure/cdn-in-a-box/edge/init.d/ /opt/init.d/

FROM edge AS edge-debug
COPY --from=get-delve /root/go/bin /usr/bin
COPY infrastructure/cdn-in-a-box/cache/init-debug-scripts.sh /opt/init.d/

# Makes the default target skip the edge-debug stage
FROM edge
