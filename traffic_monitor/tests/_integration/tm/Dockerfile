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

# This is a very simple Dockerfile.
# All it does is install and start the Traffic Monitor, given a Traffic Ops to point it to.
# It doesn't do any of the complex things the Dockerfiles in infrastructure/docker or infrastructure/cdn-in-a-box do, like inserting itself into Traffic Ops.
# It is designed for a very simple use case, where the complex orchestration of other Traffic Control components is done elsewhere (or manually).

FROM rockylinux:8
MAINTAINER dev@trafficcontrol.apache.org

ARG RPM=traffic_monitor.rpm
ADD $RPM /

RUN yum install -y initscripts jq /$(basename $RPM) && rm /$(basename $RPM)

ADD Dockerfile_run.sh /
ENTRYPOINT /Dockerfile_run.sh
