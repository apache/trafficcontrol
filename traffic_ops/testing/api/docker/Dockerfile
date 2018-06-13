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
FROM trafficops-perl:latest
MAINTAINER Dewayne Richardson <dewrich@apache.org>

RUN yum -y install vim

ARG TO_DIR=traffic_ops
ARG API_DOCKER_DIR=.

ARG TEST_DIR

ENV TEST_DIR=$TEST_DIR

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

#COPY --from=trafficops-perl /opt/traffic_ops /opt/traffic_ops
COPY --from=trafficops-perl /usr/local/go /usr/local/go

ADD $API_DOCKER_DIR/$TO_RPM /

ENV PERL5LIB="/opt/traffic_ops/app/lib:/opt/traffic_ops/app/local/lib/perl5"

WORKDIR /opt/traffic_ops/app
ADD ./run_admin_pl.sh .
RUN ./run_admin_pl.sh

#EXPOSE 60443
