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
# Dockerfile to build Traffic Server container images
#   as Mids for Traffic Control 1.4
# Based on CentOS 6.6
############################################################

# For cache, you may either use (RAM or disk) block devices or disk directories
# To use RAM block devices, pass them as /dev/ram0 and /dev/ram1 via `docker run --device`
# To use disk directories, simply don't pass devices, and the container will configure Traffic Server for directories

# Block devices may be created on the native machine with, for example, `modprobe brd`.
# The recommended minimum size for each block devices is 1G.
# For example, `sudo modprobe brd rd_size=1048576 rd_nr=2`

# Example Build and Run:
#
# docker build --rm --tag traffic_server_mid:1.4 traffic_server_mid
#
# docker run --name my-mid-0 --hostname my-mid-0 --net cdnet --device /dev/ram0:/dev/ram0 --device /dev/ram1:/dev/ram1 --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --detach traffic_server_mid:1.4
#
# OR
#
# docker run --name my-mid-0 --hostname my-mid-0 --net cdnet --env TRAFFIC_OPS_URI=http://my-traffic-ops:3000 --env TRAFFIC_OPS_USER=superroot --env TRAFFIC_OPS_PASS=supersecreterpassward --detach traffic_server_mid:1.4

FROM centos:6.6
MAINTAINER dev@trafficcontrol.apache.org

RUN yum install -y perl-JSON

RUN curl -O http://traffic-control-cdn.net/downloads/trafficserver-5.3.2-599.089d585.el6.x86_64.rpm
RUN yum install -y trafficserver-5.3.2-599.089d585.el6.x86_64.rpm

RUN mkdir /opt/ort
RUN cd /opt/ort && curl -LO https://github.com/apache/trafficcontrol/raw/RELEASE-1.4.0-RC0/traffic_ops/bin/traffic_ops_ort.pl
RUN chmod +x /opt/ort/traffic_ops_ort.pl
RUN yum install -y "perl(JSON)"
RUN curl -O http://traffic-control-cdn.net/downloads/astats_over_http-1.2-8.el6.x86_64.rpm
RUN yum install -y astats_over_http-1.2-8.el6.x86_64.rpm

RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_cop
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_crashlog
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_ctl
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_layout
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_line
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_logcat
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_logstats
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_manager
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_sac
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/trafficserver
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_server
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_top
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/traffic_via
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/tspush
RUN setcap 'cap_net_bind_service=+ep' /opt/trafficserver/bin/tsxs

EXPOSE 80 443
ADD run.sh /
ENTRYPOINT /run.sh
