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
FROM alpine:latest AS trafficrouter-dev

ENV TC=/root/go/src/github.com/apache/trafficcontrol
VOLUME /root/go/src/github.com/apache/trafficcontrol
ENV JPDA_OPTS="-agentlib:jdwp=transport=dt_socket,address=5005,server=y,suspend=n" JAVA_HOME=/usr/lib/jvm/java-8-openjdk M2_HOME=/root/.m2/
EXPOSE 3053:53/tcp 3053:53/udp 3080:80 3443:443 3333:3333 2222:3443 5005:5005

RUN apk add --no-cache openjdk8 inotify-tools maven && ln -s /usr/lib/jvm/java-1.8-openjdk/bin/jdb /bin/jdb

ADD https://dlcdn.apache.org/tomcat/tomcat-9/v9.0.54/bin/apache-tomcat-9.0.54.tar.gz /opt/tomcat.tgz
RUN cd /opt && tar -xf tomcat.tgz && mv apache-tomcat-9.0.54 tomcat && rm /usr/share/java/maven-3/conf/settings.xml

COPY settings.xml /usr/share/java/maven-3/conf/settings.xml

ENTRYPOINT /root/go/src/github.com/apache/trafficcontrol/dev/traffic_router/run.sh
