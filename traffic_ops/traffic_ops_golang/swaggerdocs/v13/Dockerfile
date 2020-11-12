# 
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at

#  http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
 
#  Swagger Spec Server
FROM debian:buster AS swagger-server

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        git \
        wget && \
    apt-get clean

COPY GO_VERSION /
RUN go_version=$(cat /GO_VERSION) && \
    wget -O go.tar.gz https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz && \
    tar -C /usr/local -xvzf go.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/bin/go && \
    rm go.tar.gz
ENV GOPATH=/go


COPY traffic_ops/traffic_ops_golang/swaggerdocs/v13/swaggerspec .
COPY traffic_ops/traffic_ops_golang/swaggerdocs/v13/swaggerspec-server /usr/src/swaggerspec-server
WORKDIR /usr/src/swaggerspec-server

RUN go build 
ENTRYPOINT ["./swaggerspec-server"]
