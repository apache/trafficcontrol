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
# Based on CentOS 7.2
############################################################

FROM centos:7

EXPOSE 80

ADD https://ci.trafficserver.apache.org/RPMS/CentOS7/trafficserver-7.1.4-2.el7.x86_64.rpm /trafficserver.rpm
ADD https://ci.trafficserver.apache.org/RPMS/CentOS7/trafficserver-devel-7.1.4-2.el7.x86_64.rpm /trafficserver-devel.rpm

RUN yum install -y bind-utils kyotocabinet-libs epel-release initscripts iproute net-tools nmap-ncat gettext autoconf automake libtool gcc-c++ cronie glibc-devel openssl-devel
RUN yum install -y /trafficserver.rpm /trafficserver-devel.rpm jq python36-psutil python36-typing python36-setuptools python36-pip logrotate && yum clean all
RUN python3 -m pip install --upgrade pip && python3 -m pip install requests urllib3 distro

ADD traffic_server/plugins/astats_over_http/astats_over_http.c traffic_server/plugins/astats_over_http/Makefile.am /

RUN tsxs -v -c astats_over_http.c -o astats_over_http.so
RUN mkdir -p /usr/libexec/trafficserver /opt/ort && tsxs -v -o astats_over_http.so -i

RUN yum remove -y gcc-c++ glibc-devel autoconf automake libtool && rm -f /astats_over_http.c /Makefile.am

# You need to do this because the RPM in the ATS archives is just all kinds of messed-up
RUN chmod 755 /usr/lib64/trafficserver /etc/trafficserver/body_factory /etc/trafficserver/body_factory/default
RUN mkdir -p /var/trafficserver /opt/ort && \
    dd if=/dev/zero bs=1M count=1000 of=/var/trafficserver/cache && \
    chown -R ats:ats /etc/trafficserver/ /var/trafficserver/ /opt/ort /usr/lib64/trafficserver/ && \
    sed -i 's/STRING 8080 8080:ipv6/STRING 80 80:ipv6/' /etc/trafficserver/records.config

RUN setcap CAP_NET_BIND_SERVICE=+eip /bin/traffic_server && setcap CAP_NET_BIND_SERVICE=+eip /bin/traffic_manager && setcap CAP_NET_BIND_SERVICE=+eip /bin/trafficserver && setcap CAP_NET_BIND_SERVICE=+eip /bin/traffic_cop

ADD infrastructure/cdn-in-a-box/ort /opt/ort/
ADD traffic_control/clients/python /opt/Apache-TrafficControl/

WORKDIR /opt

RUN touch /var/log/ort.log && pip3 install ./Apache-TrafficControl && pip3 install ./ort && cp ort/traffic_ops_ort.crontab /etc/cron.d/traffic_ops_ort-cron-template && mkdir -p /opt/init.d && cp ort/traffic_ops_ort.logrotate /etc/logrotate.d/ort

ADD infrastructure/cdn-in-a-box/cache/run.sh infrastructure/cdn-in-a-box/traffic_ops/to-access.sh infrastructure/cdn-in-a-box/enroller/server_template.json /
CMD /run.sh
