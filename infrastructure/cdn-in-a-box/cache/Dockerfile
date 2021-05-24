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
FROM centos:${RHEL_VERSION} AS common-cache-server-layers
ARG RHEL_VERSION=8

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
        perl-Carp               \
        perl-constant           \
        perl-Data-Dumper        \
        perl-Encode             \
        perl-Exporter           \
        perl-File-Path          \
        perl-File-Temp          \
        perl-Filter             \
        perl-Getopt-Long        \
        perl-HTTP-Tiny          \
        perl-libs               \
        perl-macros             \
        perl-parent             \
        perl-PathTools          \
        perl-Pod-Escapes        \
        perl-podlators          \
        perl-Pod-Perldoc        \
        perl-Pod-Simple         \
        perl-Pod-Usage          \
        perl-Scalar-List-Utils  \
        perl-Socket             \
        perl-Storable           \
        perl-Text-ParseWords    \
        perl-threads            \
        perl-threads-shared     \
        perl-Time-HiRes         \
        perl-Time-Local         \
        perl-URI                \
        tcl                     \
        $additional_packages && \
    if [[ "${RHEL_VERSION%%.*}" -eq 8 ]]; then \
        set -- \
            # Pretend that we have the right library versions.
            # TODO: Use a proper CentOS 7 or 8 RPM once trafficserver
            # is in EPEL again (see apache/trafficserver#6855)
            libtcl8.6.so        libtcl8.5.so     \
            libncursesw.so.6    libncursesw.so.5 \
            libtinfo.so.6       libtinfo.so.5    \
            || exit 1; \
    fi && \
    cd /usr/lib64 && \
    while [[ $# -gt 0 ]]; do \
        source="$1" && \
        shift && \
        target="$1" && \
        shift && \
        ln -s "$source" "$target" || exit 1; \
    done

ADD https://ci.trafficserver.apache.org/RPMS/CentOS7/trafficserver-7.1.4-2.el7.x86_64.rpm /trafficserver.rpm
ADD https://ci.trafficserver.apache.org/RPMS/CentOS7/trafficserver-devel-7.1.4-2.el7.x86_64.rpm /trafficserver-devel.rpm

RUN rpm -Uvh --nodeps /trafficserver.rpm /trafficserver-devel.rpm && \
    dnf install -y jq python3-psutil python3-setuptools python3-pip logrotate && \
    dnf clean all

RUN dnf install -y bind-utils kyotocabinet-libs initscripts iproute net-tools nmap-ncat gettext autoconf automake libtool gcc-c++ cronie glibc-devel openssl-devel

RUN python3 -m pip install --upgrade pip && python3 -m pip install requests urllib3 distro

ADD traffic_server/plugins/astats_over_http/astats_over_http.c traffic_server/plugins/astats_over_http/Makefile.am /

RUN tsxs -v -c astats_over_http.c -o astats_over_http.so

# The symbolic link here is a shim for broken atstccfg behavior - remove when it's fixed.
RUN mkdir -p /usr/libexec/trafficserver /opt/ort /opt/trafficserver/etc/trafficserver/ /opt/init.d && ln -s /opt/trafficserver/etc/trafficserver/ssl /etc/trafficserver/ssl && tsxs -v -o astats_over_http.so -i

RUN dnf remove -y gcc-c++ glibc-devel autoconf automake libtool && rm -f /astats_over_http.c /Makefile.am

# You need to do this because the RPM in the ATS archives is just all kinds of messed-up
RUN chmod 755 /usr/lib64/trafficserver /etc/trafficserver/body_factory /etc/trafficserver/body_factory/default
RUN mkdir -p /var/trafficserver /opt/ort && \
    dd if=/dev/zero bs=1M count=1000 of=/var/trafficserver/cache && \
    chown -R ats:ats /etc/trafficserver/ /var/trafficserver/ /opt/ort /usr/lib64/trafficserver/ && \
    sed -i 's/STRING 8080 8080:ipv6/STRING 80 80:ipv6/' /etc/trafficserver/records.config

RUN setcap CAP_NET_BIND_SERVICE=+eip /bin/traffic_server && setcap CAP_NET_BIND_SERVICE=+eip /bin/traffic_manager && setcap CAP_NET_BIND_SERVICE=+eip /bin/trafficserver && setcap CAP_NET_BIND_SERVICE=+eip /bin/traffic_cop

WORKDIR /opt

ADD infrastructure/cdn-in-a-box/ort /opt/ort/
ADD traffic_control/clients/python /opt/Apache-TrafficControl/

RUN touch /var/log/ort.log && \
	pip3 install ./Apache-TrafficControl && \
	pip3 install ./ort && \
	cp ort/traffic_ops_ort.crontab /etc/cron.d/traffic_ops_ort-cron-template && \
	cp ort/traffic_ops_ort.logrotate /etc/logrotate.d/ort

ADD infrastructure/cdn-in-a-box/cache/run.sh infrastructure/cdn-in-a-box/traffic_ops/to-access.sh infrastructure/cdn-in-a-box/enroller/server_template.json /

COPY infrastructure/cdn-in-a-box/dns/set-dns.sh \
     infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
     /usr/local/sbin/

ARG ORT_RPM=infrastructure/cdn-in-a-box/cache/trafficcontrol-cache-config.rpm
ADD $ORT_RPM /
RUN rpm -Uvh --nodeps /$(basename $ORT_RPM) &&\
    rm /$(basename $ORT_RPM)
CMD /run.sh

FROM common-cache-server-layers AS mid
COPY infrastructure/cdn-in-a-box/mid/init.d/ /opt/init.d/

FROM common-cache-server-layers AS edge
COPY infrastructure/cdn-in-a-box/edge/init.d/ /opt/init.d/
