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
FROM alpine:3.13 AS integration-builder

COPY --from=get-go /usr/local/go /usr/local/go
ENV PATH=/usr/local/go/bin:${PATH} \
    GOPATH=/go \
    CGO_ENABLED=0
ENV PATH="${PATH}:${GOPATH}/bin"

RUN apk add --no-cache --update git gcc

RUN go install github.com/jstemmer/go-junit-report@latest

COPY ./vendor/ /go/src/github.com/apache/trafficcontrol/vendor/

# integration source and dependencies
COPY ./cache-config/t3cutil/ /go/src/github.com/apache/trafficcontrol/cache-config/t3cutil/
COPY ./lib/ /go/src/github.com/apache/trafficcontrol/lib/
COPY ./traffic_ops/toclientlib/ /go/src/github.com/apache/trafficcontrol/traffic_ops/toclientlib/
COPY ./traffic_ops/v3-client/ /go/src/github.com/apache/trafficcontrol/traffic_ops/v3-client/
COPY ./traffic_ops/v4-client/ /go/src/github.com/apache/trafficcontrol/traffic_ops/v4-client/
COPY ./traffic_ops/v5-client/ /go/src/github.com/apache/trafficcontrol/traffic_ops/v5-client/
COPY ./traffic_ops/testing/api /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/api
COPY ./traffic_ops/traffic_ops_golang /go/src/github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang

COPY ./go.mod ./go.sum /go/src/github.com/apache/trafficcontrol/
# if we end up using a different versioning convention, the compile command will need to change
WORKDIR /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/api
RUN go mod vendor -v

RUN set -o errexit -o xtrace;\
    for api_version in v3 v4 v5; do\
        go test -c ./$api_version -ldflags="-w -s" -o traffic_ops_${api_version}_integration_test;\
    done

FROM alpine:3.13

RUN apk add --no-cache --update \
  curl \
  wget \
  bash \
  nmap \
  bind-tools \
  net-tools

# MANIFEST
# run.sh                        (wait on TO, then run bin)
# to-access.sh                  (sourced, get to-ping and env vars)
# config.sh                     (sourced, and creates config for bin file)
# tc-fixtures.json              (test data to run tests with)
# traffic_ops_integration_test  (main bin, from integration-builder)
COPY ./infrastructure/cdn-in-a-box/traffic_ops_integration_test/run.sh    /opt/integration/app/
COPY ./infrastructure/cdn-in-a-box/traffic_ops/to-access.sh /opt/integration/app/
COPY ./infrastructure/cdn-in-a-box/traffic_ops_integration_test/config.sh /opt/integration/app/
COPY ./traffic_ops/testing/api/v3/tc-fixtures.json /opt/integration/app/tc-fixtures-v3.json
COPY ./traffic_ops/testing/api/v4/tc-fixtures.json /opt/integration/app/tc-fixtures-v4.json
COPY ./traffic_ops/testing/api/v5/tc-fixtures.json /opt/integration/app/tc-fixtures-v5.json
COPY ./traffic_ops/testing/api/v5/client-intermediate-chain.crt.pem  ./traffic_ops/testing/api/v5/client.key.pem /opt/integration/app/
COPY --from=integration-builder \
    /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/api/traffic_ops_v3_integration_test \
    /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/api/traffic_ops_v4_integration_test \
    /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/api/traffic_ops_v5_integration_test \
    /opt/integration/app/

COPY --from=integration-builder \
    /go/bin/go-junit-report \
    /opt/integration/app/

COPY infrastructure/cdn-in-a-box/dns/set-dns.sh \
     infrastructure/cdn-in-a-box/dns/insert-self-into-dns.sh \
     /usr/local/sbin/

VOLUME ["/junit"]

WORKDIR /opt/integration/app
CMD ./run.sh
