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
# Dockerfile to build Traffic Portal container images
# Based on CentOS 7.2
############################################################

FROM centos/systemd

RUN curl -sL https://rpm.nodesource.com/setup_12.x | bash -

ARG RPM=traffic_portal.rpm
ADD $RPM /

RUN yum install -y epel-release
RUN yum install -y jq nodejs openssl gettext bind-utils net-tools /$(basename $RPM)

RUN rm /$(basename $RPM)

RUN yum clean all

ADD run.sh /
CMD /run.sh
