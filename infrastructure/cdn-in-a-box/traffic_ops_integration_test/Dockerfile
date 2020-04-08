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

FROM golang:1.14.2-alpine3.11 AS integration-builder

RUN apk add --no-cache --update git gcc

RUN go get -u golang.org/x/net/publicsuffix golang.org/x/crypto/scrypt
RUN go get -u github.com/jstemmer/go-junit-report
RUN cd /go/src/github.com/jstemmer/go-junit-report && go build

# integration source and dependencies
COPY ./infrastructure/cdn-in-a-box/ /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/
COPY ./vendor/ /go/src/github.com/apache/trafficcontrol/vendor/
COPY ./traffic_ops/vendor/ /go/src/github.com/apache/trafficcontrol/traffic_ops/vendor/
COPY ./lib/ /go/src/github.com/apache/trafficcontrol/lib/
COPY ./traffic_ops/client/ /go/src/github.com/apache/trafficcontrol/traffic_ops/client/
COPY ./traffic_ops/v1-client/ /go/src/github.com/apache/trafficcontrol/traffic_ops/v1-client/
COPY ./traffic_ops/testing/api /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/api
COPY ./traffic_ops/traffic_ops_golang /go/src/github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang

# if we end up using a different versioning convention, the compile command will need to change
WORKDIR /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/api
RUN CGO_ENABLED=0 go test -c ./v1* -ldflags="-w -s" -o traffic_ops_v1_integration_test
RUN CGO_ENABLED=0 go test -c ./v2* -ldflags="-w -s" -o traffic_ops_v2_integration_test

FROM alpine:3.11

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
COPY ./traffic_ops/testing/api/v1/tc-fixtures.json /opt/integration/app/tc-fixtures-v1.json
COPY ./traffic_ops/testing/api/v2/tc-fixtures.json /opt/integration/app/tc-fixtures-v2.json
COPY --from=integration-builder \
    /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/api/traffic_ops_v1_integration_test \
    /opt/integration/app/
COPY --from=integration-builder \
    /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/api/traffic_ops_v2_integration_test \
    /opt/integration/app/

COPY --from=integration-builder \
    /go/bin/go-junit-report \
    /opt/integration/app/

VOLUME ["/junit"]

WORKDIR /opt/integration/app
CMD ./run.sh
