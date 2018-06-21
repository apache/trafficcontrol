#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
FROM golang:1.10.1
MAINTAINER Dewayne Richardson <dewrich@apache.org>

RUN pwd
ARG DB_NAME
ARG DB_SERVER
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD

ENV DB_NAME=$DB_SERVER
ENV DB_SERVER=$DB_SERVER
ENV DB_PORT=$DB_PORT
ENV DB_USER=$DB_SERVER
ENV DB_PASSWORD=$DB_PASSWORD

ARG SRC_DIR=/go/src
ARG TC_DIR=$SRC_DIR/github.com/apache/trafficcontrol
ARG TESTING_DIR=$TC_DIR/traffic_ops/testing
ARG TESTING_API_DIR=$TESTING_DIR/api
 
#COPY --from=trafficops-perl /opt/traffic_ops/app/bin/traffic_ops_golang /opt/traffic_ops/app/bin/traffic_ops_golang
#COPY --from=trafficops-perl /usr/local/go /usr/local/go

WORKDIR $TC_DIR

# Add the lib/go-tc dir
ADD lib lib

# Add the dependent dirs (without the cruft)
ADD traffic_ops/traffic_ops_golang traffic_ops/traffic_ops_golang
ADD traffic_ops/vendor traffic_ops/vendor
ADD traffic_ops/testing/api traffic_ops/testing/api
ADD traffic_ops/client traffic_ops/client

#ADD conf .
WORKDIR $TESTING_API_DIR/v13
CMD ../docker/run_api_tests.sh

# vi:syntax=Dockerfile
