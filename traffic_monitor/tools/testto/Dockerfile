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

FROM centos:8
MAINTAINER dev@trafficcontrol.apache.org

RUN dnf -y install golang git
ENV GOPATH=/go

RUN go get -u github.com/apache/trafficcontrol/v6/lib/go-rfc github.com/apache/trafficcontrol/v6/lib/go-tc
ADD . ${GOPATH}/src/github.com/apache/trafficcontrol/traffic_monitor/tools/testto
WORKDIR ${GOPATH}/src/github.com/apache/trafficcontrol/traffic_monitor/tools/testto
RUN go build && cp testto /usr/sbin

CMD ${GOPATH}/src/github.com/apache/trafficcontrol/traffic_monitor/tools/testto/Dockerfile_run.sh
