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
# Based on CentOS 8
############################################################

ARG RHEL_VERSION=8
FROM centos:${RHEL_VERSION} AS common-traffic-server-dependencies
ARG RHEL_VERSION=8
# Makes RHEL_VERSION available at runtime
ENV RHEL_VERSION="$RHEL_VERSION"

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        yum -y install dnf || exit 1; \
    fi

EXPOSE 80

RUN dnf -y install epel-release && \
    if [[ "${RHEL_VERSION%%.*}" -ge 8 ]]; then \
        additional_packages='compat-openssl10 pkgconf-pkg-config' || \
        exit 1; \
    else \
        additional_packages=openssl || \
        exit 1; \
    fi && \
    dnf -y install              \
        GeoIP                   \
        git                     \
        groff-base              \
        hwloc                   \
        hwloc-libs              \
        kyotocabinet-libs       \
        libtool-ltdl            \
        libunwind               \
        lzo                     \
        make                    \
        numactl-libs            \
        perl                    \
        tcl                     \
        $additional_packages && \
    if [[ "${RHEL_VERSION%%.*}" -eq 8 ]]; then \
        # Pretend that we have the right library versions.
        # TODO: Use a proper CentOS 7 or 8 RPM once trafficserver
        # is in EPEL again (see apache/trafficserver#6855)
        cd /usr/lib64 && \
        ln -s libtcl8.6.so        libtcl8.5.so && \
        ln -s libncursesw.so.6    libncursesw.so.5 && \
        ln -s libtinfo.so.6       libtinfo.so.5  || \
            exit 1; \
    fi

RUN dnf install -y bind-utils kyotocabinet-libs initscripts iproute net-tools nmap-ncat gettext autoconf automake libtool gcc-c++ cronie glibc-devel openssl-devel && \
    dnf install -y jq logrotate && \
    dnf clean all

FROM common-traffic-server-dependencies AS common-cache-config-layers

ADD https://ci.trafficserver.apache.org/RPMS/CentOS7/trafficserver-7.1.4-2.el7.x86_64.rpm /trafficserver.rpm
ADD https://ci.trafficserver.apache.org/RPMS/CentOS7/trafficserver-devel-7.1.4-2.el7.x86_64.rpm /trafficserver-devel.rpm

COPY traffic_server/plugins/astats_over_http/astats_over_http.c traffic_server/plugins/astats_over_http/Makefile.am /

RUN rpm -Uvh --nodeps /trafficserver.rpm /trafficserver-devel.rpm && \
    tsxs -v -c astats_over_http.c -o astats_over_http.so

# The symbolic links here are used because trafficserver-7.1.4-2.el7.x86_64.rpm has / as TS_ROOT,
# but the Cache Configs assume that /opt/trafficserver is TS_ROOT. The symlinks can be removed
# once a Traffic Server RPM using /opt/trafficserver as TS_ROOT is used.
RUN ln -s / /opt/trafficserver && \
    mkdir /libexec /etc/trafficserver/ssl && \
    ln -s /usr/lib64/trafficserver/plugins /libexec/trafficserver && \
    tsxs -v -o astats_over_http.so -i

RUN dnf remove -y gcc-c++ glibc-devel autoconf automake libtool && \
    rm -f /astats_over_http.c /Makefile.am

# traffic_server needs different ownership than trafficserver-7.1.4-2.el7.x86_64.rpm sets
RUN chown -R ats:ats /etc/trafficserver/ /usr/lib64/trafficserver/


WORKDIR /opt

COPY infrastructure/cdn-in-a-box/cache/traffic_ops_ort.crontab /etc/cron.d/traffic_ops_ort-cron-template
COPY infrastructure/cdn-in-a-box/cache/traffic_ops_ort.logrotate /etc/logrotate.d/ort

COPY infrastructure/cdn-in-a-box/cache/run.sh infrastructure/cdn-in-a-box/traffic_ops/to-access.sh infrastructure/cdn-in-a-box/enroller/server_template.json /

COPY infrastructure/cdn-in-a-box/dns/set-dns.sh \
     infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
     /usr/local/sbin/
# Copy systemctl.sh to /usr/bin specifically because t3c runs /bin/systemctl by absolute path
COPY infrastructure/cdn-in-a-box/cache/init-debug-scripts.sh /usr/bin/systemctl

ARG ORT_RPM=infrastructure/cdn-in-a-box/cache/trafficcontrol-cache-config.rpm
COPY $ORT_RPM /
RUN rpm -Uvh /$(basename $ORT_RPM) &&\
    rm /$(basename $ORT_RPM)
CMD /run.sh

FROM common-traffic-server-dependencies AS get-delve
RUN dnf -y install golang && \
    go get -u github.com/go-delve/delve/cmd/dlv

FROM common-cache-config-layers AS mid
COPY infrastructure/cdn-in-a-box/mid/init.d/ /opt/init.d/

FROM mid AS mid-debug
COPY --from=get-delve /root/go/bin /usr/bin
COPY infrastructure/cdn-in-a-box/cache/init-debug-scripts.sh /opt/init.d/

FROM common-cache-config-layers AS edge
COPY infrastructure/cdn-in-a-box/edge/init.d/ /opt/init.d/

FROM edge AS edge-debug
COPY --from=get-delve /root/go/bin /usr/bin
COPY infrastructure/cdn-in-a-box/cache/init-debug-scripts.sh /opt/init.d/

# Makes the default target skip the edge-debug stage
FROM edge
