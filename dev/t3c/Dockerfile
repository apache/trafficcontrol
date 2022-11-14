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
ARG ATS_VERSION
ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine AS get-go
FROM ghcr.io/apache/trafficcontrol/ci/trafficserver-alpine:${ATS_VERSION}
COPY --from=get-go /usr/local/go /usr/local/go
ENV PATH=/usr/local/go/bin:${PATH} \
	GOPATH=/go
ENV PATH=${GOPATH}/bin:${PATH}

ENV TC="/go/src/github.com/apache/trafficcontrol" GOFLAGS="--buildvcs=false"
VOLUME $TC
EXPOSE 80 8081

USER root
RUN apk add --no-cache \
	# inotify-tools is used for inotifywait in run.sh
	inotify-tools \
	# make is used for the t3c Makefile
	make \
	# gcc and musl-dev are used to build packages using CGO
	gcc musl-dev && \
	go install github.com/go-delve/delve/cmd/dlv@latest && \
	rm -rf $GOPATH/pkg/*

RUN echo "stats_over_http.so" >> /etc/trafficserver/plugin.config && echo "system_stats.so" >> /etc/trafficserver/plugin.config

CMD ${TC}/dev/t3c/run.sh
