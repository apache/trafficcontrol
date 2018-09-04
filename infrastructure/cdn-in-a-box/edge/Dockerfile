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

FROM tccache:latest


RUN mkdir -p /var/trafficserver/cache /opt/ort && chown -R ats:ats /etc/trafficserver/ /var/trafficserver/ /opt/ort
ADD edge/run.sh traffic_ops/to-access.sh enroller/server_template.json /
ADD edge/traffic_ops_ort.py /opt/ort

# You need to do this because the RPM in the ATS archives is just all kinds of messed-up
RUN chmod 755 /usr/lib64/trafficserver /etc/trafficserver/body_factory /etc/trafficserver/body_factory/default

CMD /run.sh

