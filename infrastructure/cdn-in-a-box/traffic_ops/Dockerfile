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

RUN rpm --import /etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7 && \
    rpm --import http://dl.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-7 && \
    rpm --import https://download.postgresql.org/pub/repos/yum/RPM-GPG-KEY-PGDG && \
    yum -y update ca-certificates && \
    yum-config-manager --add-repo https://download.postgresql.org/pub/repos/yum/9.6/redhat/rhel-7-x86_64 && \
    yum install -y \
        epel-release \
        pgdg-redhat-latest && \
    yum -y clean all

RUN yum install -y \
        cpanminus \
        bind-utils \
        net-tools \
        gettext \
        golang \
        nmap-ncat \
        openssl \
        perl \
        perl-Crypt-ScryptKDF \
        perl-Test-CPAN-Meta \
        perl-JSON-PP \
        git \
        iproute \
        jq && \
    yum -y clean all

RUN yum install -y perl-DBIx-Connector

# Override TRAFFIC_OPS_RPM arg to use a different one using --build-arg TRAFFIC_OPS_RPM=...  Can be local file or http://...
ARG TRAFFIC_OPS_RPM=traffic_ops/traffic_ops.rpm
ADD $TRAFFIC_OPS_RPM /
RUN yum install -y \
        /$(basename $TRAFFIC_OPS_RPM) \
        rm /$(basename $TRAFFIC_OPS_RPM) && \
    yum clean all

# copy any dir structure in overrides to TO -- allows modification of the install and shortcut to get perl modules/goose installed
WORKDIR /opt/traffic_ops/app
COPY traffic_ops/overrides/ /opt/traffic_ops/.
RUN cpanm -l ./local Carton

# run carton whether or not local dir was installed
RUN POSTGRES_HOME=/usr/pgsql-9.6 PERL5LIB=$(pwd)/local/lib/perl5 ./local/bin/carton && \
     rm -rf $HOME/.cpan* /tmp/Dockerfile /tmp/local.tar.gz

RUN /opt/traffic_ops/install/bin/install_goose.sh

RUN export PERL5LIB=/opt/traffic_ops/app/local/lib/perl5/:/opt/traffic_ops/install/lib/ \
	&& export TERM=xterm \
	&& export USER=root \
	&& /opt/traffic_ops/install/bin/download_web_deps -i || \
	true # keep a failure here from failing all..

ADD http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz /

EXPOSE 443
WORKDIR /opt/traffic_ops/app
ENV MOJO_MODE production

ADD enroller/server_template.json \
    traffic_ops/run.sh \
    traffic_ops/config.sh \
    traffic_ops/adduser.pl \
    traffic_ops/to-access.sh \
    traffic_ops/generate-certs.sh \
    traffic_ops/trafficops-init.sh \
    variables.env \
    /

ADD traffic_ops_data /traffic_ops_data

RUN chown -R trafops:trafops /opt/traffic_ops
CMD /run.sh
