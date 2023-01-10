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

ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine AS get-go
FROM alpine:3.14 AS load-test-builder

COPY --from=get-go /usr/local/go /usr/local/go
ENV PATH=/usr/local/go/bin:${PATH} \
    GOPATH=/go \
    CGO_ENABLED=0
ENV PATH="${PATH}:${GOPATH}/bin"
ARG TC_DIR=${GOPATH}/src/github.com/apache/trafficcontrol

COPY ./lib/ ${TC_DIR}/lib/
COPY ./traffic_ops/toclientlib/ ${TC_DIR}/traffic_ops/toclientlib/
COPY ./traffic_ops/v4-client/ ${TC_DIR}/traffic_ops/v4-client/
COPY ./go.mod ./go.sum ${TC_DIR}/
COPY ./traffic_router/ultimate-test-harness ${TC_DIR}/traffic_router/ultimate-test-harness

RUN cd ${TC_DIR}/traffic_router/ultimate-test-harness && \
    go test -c

FROM alpine:3.14 AS load-test

RUN apk add --no-cache \
    # for to-access.sh
    bash \
    # contains dig, required by to-access.sh
    bind-tools \
    # to recognize the CDN in a Box CA certificate
    ca-certificates

WORKDIR /opt/load-test/app
COPY --from=load-test-builder /go/src/github.com/apache/trafficcontrol/traffic_router/ultimate-test-harness .

COPY ./infrastructure/cdn-in-a-box/traffic_ops/to-access.sh /
COPY infrastructure/cdn-in-a-box/dns/set-dns.sh \
     infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
     infrastructure/cdn-in-a-box/traffic_router_load_test/run.sh /usr/local/bin/

CMD run.sh
