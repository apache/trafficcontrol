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
#FROM maven:3.5.4-jdk-8-alpine
FROM centos/systemd
ARG DIR=github.com/apache/trafficcontrol

ADD traffic_router /go/src/$DIR/traffic_router
VOLUME ["/junit"]

WORKDIR /go/src/$DIR/traffic_router

RUN yum update -y
RUN yum install -y java-1.6.0-openjdk
RUN yum install -y maven
RUN yum install -y epel-release
RUN yum install -y tomcat-native

CMD bash -c 'mvn test -Djava.library.path=/usr/share/java -DoutputDirectory=/junit 2>&1 && mv core/target/surefire-reports/* /junit'

#
# vi:syntax=Dockerfile
