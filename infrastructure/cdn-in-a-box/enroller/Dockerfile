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
FROM golang:${GO_VERSION} AS get-go
FROM debian:bullseye AS enroller-dependencies

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        git \
        wget && \
    apt-get clean

COPY --from=get-go /usr/local/go /usr/local/go
ENV PATH=/usr/local/go/bin:${PATH} \
    GOPATH=/go \
    CGO_ENABLED=0

FROM enroller-dependencies AS enroller-builder
ARG ENROLLER_DEBUG_BUILD=false

# enroller source and dependencies
COPY ./lib/ /go/src/github.com/apache/trafficcontrol/lib/
COPY ./go.mod ./go.sum /go/src/github.com/apache/trafficcontrol/
COPY ./vendor/ /go/src/github.com/apache/trafficcontrol/vendor/
COPY ./traffic_ops/toclientlib/ /go/src/github.com/apache/trafficcontrol/traffic_ops/toclientlib/
COPY ./traffic_ops/v4-client/ /go/src/github.com/apache/trafficcontrol/traffic_ops/v4-client/
COPY ./traffic_ops/v5-client/ /go/src/github.com/apache/trafficcontrol/traffic_ops/v5-client/
COPY ./infrastructure/cdn-in-a-box/ /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/

# varnishcfg requires t3c for ToData struct and not needed for enroller
RUN rm -rf /go/src/github.com/apache/trafficcontrol/lib/varnishcfg

WORKDIR /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/enroller
RUN set -o errexit -o nounset; \
    go clean; \
    go mod vendor -v; \
    gcflags= ldflags=; \
    if [ "$ENROLLER_DEBUG_BUILD" = true ]; then \
        apt-get install -y --no-install-recommends gcc libstdc++-9-dev; \
        echo 'Building Enroller without optimization or inlining'; \
        gcflags='all=-N -l'; \
    else \
        echo 'Optimizing Enroller build'; \
        ldflags='-s -w'; \
    fi; \
    go build -ldflags "$ldflags" -gcflags "$gcflags"

FROM enroller-dependencies as get-delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM debian:bullseye AS enroller

RUN apt-get update && apt-get install -y \
        netcat curl dnsutils net-tools \
        #to-access dependencies
        jq gettext && \
    apt-get clean
COPY --from=enroller-builder \
   /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/enroller/enroller \
   /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/enroller/run.sh \
   /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/enroller/server_template.json \
   /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/traffic_ops/to-access.sh \
   /

COPY infrastructure/cdn-in-a-box/dns/set-dns.sh \
     infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
     /usr/local/sbin/

WORKDIR /shared/enroller
CMD /run.sh

FROM enroller AS enroller-debug
COPY --from=get-delve /go/bin /usr/bin

# Makes the default target skip the enroller-debug stage
FROM enroller
