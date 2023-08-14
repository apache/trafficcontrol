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
FROM alpine:3.13 AS base

COPY --from=get-go /usr/local/go /usr/local/go
ENV PATH=/usr/local/go/bin:${PATH} \
    CGO_ENABLED=0 \
    GOPATH=/go

RUN set -o errexit; \
    architecture=$(uname -m); \
    mkdir lib64; \
    # Use musl libc where the go binary expects glibc
    # Less-generalized: ln -s /lib/ld-musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
    ln -s /lib/ld-musl-${architecture}.so.[0-9] /lib64/ld-linux-${architecture//_/-}.so.2;

RUN apk --no-cache add gcc musl-dev git

ARG DIR=github.com/apache/trafficcontrol

WORKDIR /go/src/github.com/apache/trafficcontrol

FROM base AS lint

RUN apk --no-cache add curl

RUN version=1.32.0 && \
    curl -L https://github.com/golangci/golangci-lint/releases/download/v${version}/golangci-lint-${version}-linux-amd64.tar.gz \
	| tar -zxC /tmp \
	&& mv /tmp/golangci-lint-${version}-linux-amd64/golangci-lint /usr/local/bin/golangci-lint \
	&& rm -rf /tmp/golangci-lint-${version}-linux-amd64

CMD go mod vendor -v && \
    golangci-lint run ./...

FROM base AS unit

VOLUME ["/junit"]

RUN go install github.com/wadey/gocovmerge@latest &&\
    go install github.com/jstemmer/go-junit-report@latest

CMD /go/src/github.com/apache/trafficcontrol/tools/golang/unit-tests.sh
