#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
FROM centos:7
MAINTAINER dev@trafficcontrol.apache.org

RUN yum -y install \
        https://download.postgresql.org/pub/repos/yum/9.6/redhat/rhel-7-x86_64/pgdg-centos96-9.6-3.noarch.rpm

RUN yum -y install \
        vim \
        cpanminus \
        expat-devel \
        gcc-c++ \
        libcurl \
        libcurl-devel \
        libidn-devel \
        libpcap-devel \
        mkisofs \
        nmap-ncat \
        openssl-devel \
        perl \
        perl-App-cpanminus \
        perl-Crypt-ScryptKDF \
        perl-DBD-Pg \
        perl-DBI \
        perl-Digest-SHA1 \
        perl-JSON \
        perl-TermReadKey \
        perl-Test-CPAN-Meta \
        perl-WWW-Curl \
        perl-core \
        perl-libwww-perl \
        postgresql96 \
        postgresql96-devel && \
        yum clean all

#RUN cpanm MIYAGAWA/Carton-v1.0.26.tar.gz
RUN cpanm -n Carton

ADD app /opt/traffic_ops/app
WORKDIR /opt/traffic_ops/app
#RUN carton
RUN POSTGRES_HOME=/usr/pgsql-9.6 carton

ADD install/bin/install_goose.sh /
ADD install/bin/install_go.sh /
RUN /install_go.sh
RUN /install_goose.sh


# ignore this if it fails
#RUN rm -rf /root/.cpan* 2>/dev/null || true

ADD app/bin/tests/runtests.sh /
ARG TESTDIR
ARG TESTENV
ENV TESTDIR=$TESTDIR
ENV TESTENV=$TESTENV
ARG DBHOST
ARG DBPORT
ENV DBHOST=$DBHOST
ENV DBPORT=$DBPORT

ENTRYPOINT /runtests.sh $DBHOST $DBPORT
CMD $TESTENV $TESTDIR

#
# vi:syntax=Dockerfile
