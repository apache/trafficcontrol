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
# Dockerfile to initialized Traffic Ops Database container 
# Based on CentOS 7.2
############################################################

FROM centos/systemd

RUN yum install -y https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm

RUN yum -y install \
  postgresql13 \
  nmap-ncat \
  cpanminus && \
  yum clean all

ENV POSTGRES_HOME $POSTGRES_HOME
ENV PGPASSWORD $PGPASSWORD 
ENV DB_USERNAME $DB_USERNAME
ENV DB_NAME $DB_NAME
ENV DB_USER_PASS $DB_USER_PASS 
ENV DB_SERVER $DB_SERVER
ENV DB_PORT $DB_PORT

ADD db_init/dbInit.sh /
CMD /dbInit.sh
