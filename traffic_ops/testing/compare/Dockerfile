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

############################################################
# Dockerfile to build Docker images for testing via the 'compare' tool
# Based on Alpine Linux
############################################################

FROM golang:alpine

ENV GOPATH /go/
RUN mkdir -p /go/src/github.com/apache/trafficcontrol/traffic_control/clients/python\
	/go/src/github.com/apache/trafficcontrol/lib\
	/go/src/github.com/apache/trafficcontrol/traffic_ops/testing/compare\
	/go/src/github.com/apache/trafficcontrol/traffic_ops/vendor/github.com/kelseyhightower\
	/artifacts

RUN apk update
RUN apk add python3 git
ADD traffic_control/clients/python /go/src/github.com/apache/trafficcontrol/traffic_control/clients/python/
RUN python3 -m ensurepip && python3 -m pip install --upgrade pip && python3 -m pip install /go/src/github.com/apache/trafficcontrol/traffic_control/clients/python/

ADD lib /go/src/github.com/apache/trafficcontrol/lib
ADD vendor /go/src/github.com/apache/trafficcontrol/vendor
ADD traffic_ops/vendor/github.com/kelseyhightower /go/src/github.com/apache/trafficcontrol/traffic_ops/vendor/github.com/kelseyhightower
ADD traffic_ops/testing/compare /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/compare

WORKDIR /go/src/github.com/apache/trafficcontrol/traffic_ops/testing/compare/
RUN go get -v ./...
RUN go build .
RUN cp testroutes.txt /artifacts/

ARG MODE="-s"
ENV mode=$MODE

CMD ./genConfigRoutes.py $mode -k --refURL=$TO_URL --testURL=$TEST_URL --refUser=$TO_USER --refPasswd=$TO_PASSWORD --testUser=$TEST_USER --testPasswd=$TEST_PASSWORD -l INFO 2>&1 >>/artifacts/testroutes.txt | tee /artifacts/genRoutesConfig.log &&\
	./compare --ref_url=$TO_URL --test_url=$TEST_URL --ref_user=$TO_USER --ref_passwd=$TO_PASSWORD --test_user=$TEST_USER --test_passwd=$TEST_PASSWORD -r /artifacts </artifacts/testroutes.txt
