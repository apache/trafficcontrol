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

FROM golang:1.14.1-alpine3.11 AS base

RUN apk --no-cache add gcc musl-dev git

ARG DIR=github.com/apache/trafficcontrol

WORKDIR /go/src/github.com/apache/trafficcontrol

RUN go get golang.org/x/crypto/ed25519\
           golang.org/x/crypto/scrypt\
           golang.org/x/net/ipv4\
           golang.org/x/net/ipv6\
           golang.org/x/net/publicsuffix\
           golang.org/x/sys/unix

FROM base AS lint

RUN apk --no-cache add curl

RUN curl -L https://github.com/golangci/golangci-lint/releases/download/v1.21.0/golangci-lint-1.21.0-linux-amd64.tar.gz \
	| tar -zxC /tmp \
	&& mv /tmp/golangci-lint-1.21.0-linux-amd64/golangci-lint /usr/local/bin/golangci-lint \
	&& rm -rf /tmp/golangci-lint-1.21.0-linux-amd64

CMD ["golangci-lint", "run", "./..."]

FROM base AS unit

VOLUME ["/junit"]

RUN go get github.com/wadey/gocovmerge &&\
    go get -u github.com/jstemmer/go-junit-report

CMD /go/src/github.com/apache/trafficcontrol/tools/golang/unit-tests.sh
