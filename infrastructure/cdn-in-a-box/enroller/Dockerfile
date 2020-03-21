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

FROM golang:1.14.1 AS enroller-builder

# enroller source and dependencies
COPY ./lib/ /go/src/github.com/apache/trafficcontrol/lib/
COPY ./vendor/ /go/src/github.com/apache/trafficcontrol/vendor/
COPY ./traffic_ops/client/ /go/src/github.com/apache/trafficcontrol/traffic_ops/client/
COPY ./traffic_ops/vendor/ /go/src/github.com/apache/trafficcontrol/traffic_ops/vendor/
COPY ./infrastructure/cdn-in-a-box/ /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/

WORKDIR /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/enroller
RUN go clean && go get -v && go build



FROM debian:stretch

RUN apt-get update && apt-get install -y netcat curl dnsutils net-tools && apt-get clean
COPY --from=enroller-builder \
   /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/enroller/enroller \
   /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/enroller/run.sh \
   /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/enroller/server_template.json \
   /go/src/github.com/apache/trafficcontrol/infrastructure/cdn-in-a-box/traffic_ops/to-access.sh \
   /

WORKDIR /shared/enroller
CMD /run.sh
