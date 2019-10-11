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
# Based on CentOS 7.2
############################################################


FROM centos:7

EXPOSE 443
ENV MOJO_MODE production

RUN rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7 && \
    rpm --import https://dl.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-7 && \
    rpm --import https://download.postgresql.org/pub/repos/yum/RPM-GPG-KEY-PGDG && \
    yum -y update ca-certificates && \
    yum-config-manager --add-repo https://download.postgresql.org/pub/repos/yum/9.6/redhat/rhel-7-x86_64 && \
    yum install -y \
        epel-release \
        pgdg-redhat-latest && \
    yum install -y \
        bind-utils \
        cpanminus \
        expat-devel \
        gcc-c++ \
        gettext \
        git \
        golang \
        iproute \
        jq \
        libcurl \
        libcurl-devel \
        libidn-devel \
        libpcap-devel \
        mkisofs \
        net-tools \
        nmap-ncat \
        openssl \
        openssl-devel \
        perl \
        perl-core \
        perl-Crypt-ScryptKDF \
        perl-DBD-Pg \
        perl-DBI \
        perl-DBIx-Connector \
        perl-Digest-SHA1 \
        perl-JSON \
        perl-JSON-PP \
        perl-libwww-perl \
        perl-TermReadKey \
        perl-Test-CPAN-Meta \
        perl-WWW-Curl \
        postgresql96 \
        postgresql96-devel \
        postgresql96-libs \
        tar &&\
    yum -y clean all

ADD https://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz /

WORKDIR /opt/traffic_ops/app
ADD traffic_ops/app/cpanfile traffic_ops/install/bin/install_goose.sh ./
RUN cpanm -l ./local Carton && \
    POSTGRES_HOME=/usr/pgsql-9.6 PERL5LIB=$(pwd)/local/lib/perl5 ./local/bin/carton  && \
    rm -rf $HOME/.cpan* /tmp/Dockerfile /tmp/local.tar.gz ./cpanfile
RUN ./install_goose.sh

# Override TRAFFIC_OPS_RPM arg to use a different one using --build-arg TRAFFIC_OPS_RPM=...  Can be local file or http://...
ARG TRAFFIC_OPS_RPM=infrastructure/cdn-in-a-box/traffic_ops/traffic_ops.rpm
ADD $TRAFFIC_OPS_RPM /
RUN yum install -y /$(basename $TRAFFIC_OPS_RPM) && \
    rm /$(basename $TRAFFIC_OPS_RPM) && \
    yum -y clean all

# Run carton again, in case the cpanfile included in the RPM differs from the one used earlier in the
# build (should never happen - Perl is supposed to be going away)
RUN POSTGRES_HOME=/usr/pgsql-9.6 PERL5LIB=$(pwd)/local/lib/perl5 ./local/bin/carton && \
    rm -rf $HOME/.cpan* /tmp/Dockerfile /tmp/local.tar.gz

ADD infrastructure/cdn-in-a-box/enroller/server_template.json \
    infrastructure/cdn-in-a-box/traffic_ops/run.sh \
    infrastructure/cdn-in-a-box/traffic_ops/config.sh \
    infrastructure/cdn-in-a-box/traffic_ops/adduser.pl \
    infrastructure/cdn-in-a-box/traffic_ops/to-access.sh \
    infrastructure/cdn-in-a-box/traffic_ops/generate-certs.sh \
    infrastructure/cdn-in-a-box/traffic_ops/trafficops-init.sh \
    infrastructure/cdn-in-a-box/variables.env \
    /

ADD infrastructure/cdn-in-a-box/traffic_ops_data /traffic_ops_data

CMD /run.sh
