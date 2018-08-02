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
# Dockerfile to build Traffic Vault image
# Based on CentOS 7.2
############################################################
FROM centos:7

ARG ADMIN_PASS
ARG USER_PASS

# download and install Riak
RUN curl -s https://packagecloud.io/install/repositories/basho/riak/script.rpm.sh | bash
RUN yum -y install \
        epel-release && \
    yum -y install \
        initscripts \
        jq \
        nmap-ncat \
        riak && \
    yum -y clean all && \
    rm -rf /var/yum/cache

EXPOSE 8098 8087 8088
ADD to-access.sh /
ADD traffic_vault/run.sh /
ENTRYPOINT [ "./run.sh", "/to-access.sh" ]
