#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.
#
ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine AS trafficmonitor-dev

ENV TC=/go/src/github.com/apache/trafficcontrol GOFLAGS="--buildvcs=false"
VOLUME $TC
EXPOSE 80 81

RUN apk add --no-cache inotify-tools gcc libc-dev && \
	go install github.com/go-delve/delve/cmd/dlv@latest && \
	rm -rf $GOPATH/pkg/*
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

CMD ${TC}/dev/traffic_monitor/run.sh
