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


FROM trafficops-go

EXPOSE 443
ENV MOJO_MODE production

RUN yum -y install \
        iproute \
        libpcap-devel \
        mkisofs \
        perl-DBIx-Connector && \
    yum -y clean all && mkdir -p /opt/traffic_ops/app/public

ADD traffic_router/core/src/test/resources/geo/GeoLite2-City.mmdb.gz /opt/traffic_ops/app/public/

# Override TRAFFIC_OPS_RPM arg in Dockerfile-go to use a different RPM using --build-arg TRAFFIC_OPS_RPM=...  Can be local file or http://...
RUN yum -y reinstall /$(basename $TRAFFIC_OPS_RPM) && \
    rm /$(basename $TRAFFIC_OPS_RPM) && \
    yum -y clean all && \
    rm /opt/traffic_ops/app/bin/traffic_ops_golang

WORKDIR /opt/traffic_ops/app
RUN cpanm -l ./local Carton && \
    POSTGRES_HOME=/usr/pgsql-9.6 PERL5LIB=$(pwd)/local/lib/perl5 ./local/bin/carton  && \
    rm -rf $HOME/.cpan* /tmp/Dockerfile /tmp/local.tar.gz
RUN /opt/traffic_ops/install/bin/install_goose.sh

COPY infrastructure/cdn-in-a-box/traffic_ops/run.sh \
    infrastructure/cdn-in-a-box/traffic_ops/adduser.pl \
    infrastructure/cdn-in-a-box/traffic_ops/generate-certs.sh \
    infrastructure/cdn-in-a-box/traffic_ops/trafficops-init.sh \
    infrastructure/cdn-in-a-box/variables.env \
    /

COPY infrastructure/cdn-in-a-box/traffic_ops_data /traffic_ops_data

CMD /run.sh
