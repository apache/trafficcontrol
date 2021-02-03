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

# Example Build and Run:
# docker network create cdnet
# docker build --rm --tag traffic_ops:1.7.0 --build-arg=RPM=http://traffic-control-cdn.net/downloads/1.7.0/RELEASE-1.7.0/traffic_ops-1.7.0-3908.5b77f60f.x86_64.rpm traffic_ops
#
# docker run --name my-traffic-ops-mysql --hostname my-traffic-ops-mysql --net cdnet --env MYSQL_ROOT_PASSWORD=secretrootpass --detach mysql:5.5
#
# docker run --name my-traffic-ops --hostname my-traffic-ops --net cdnet --publish 443:443 --env MYSQL_IP=my-traffic-ops-mysql --env MYSQL_PORT=3306 --env MYSQL_ROOT_PASS=secretrootpass --env MYSQL_TRAFFIC_OPS_PASS=supersecretpassword --env ADMIN_USER=superroot --env ADMIN_PASS=supersecreterpassward --env CERT_COUNTRY=US --env CERT_STATE=Colorado --env CERT_CITY=Denver --env CERT_COMPANY=NotComcast --env TRAFFIC_VAULT_PASS=marginallylesssecret --env DOMAIN=cdnet --detach traffic_ops:1.5.1

FROM centos/systemd
MAINTAINER dev@trafficcontrol.apache.org

RUN yum install -y \
  https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm \
  initscripts epel-release perl-Crypt-ScryptKDF perl cpanminus perl-Test-CPAN-Meta perl-DBIx-Connector

RUN cpanm Carton

# Override RPM arg to use a different one using --build-arg RPM=...  Can be local file or http://...
ARG RPM=traffic_ops.rpm
ADD traffic_ops/$RPM /
RUN yum install -y /$(basename $RPM)

# once installed, remove rpm to lower image size
RUN rm /$(basename $RPM)

RUN POSTGRES_HOME=/usr/pgsql-9.6 cd /opt/traffic_ops/app && carton

RUN export PERL5LIB=/opt/traffic_ops/app/local/lib/perl5/:/opt/traffic_ops/install/lib/ \
	&& export TERM=xterm \
	&& export USER=root 

# fixes an 'Invalid Argument' bug; TODO diagnose    , fix, & remove
RUN cp /opt/traffic_ops/app/bin/traffic_ops_golang{,.new} && mv /opt/traffic_ops/app/bin/traffic_ops_golang{.new,} 

EXPOSE 443
WORKDIR /opt/traffic_ops/app
ENV MOJO_MODE production
ADD traffic_ops/profile.origin.traffic_ops /
ADD traffic_ops/run.sh /
CMD /run.sh
