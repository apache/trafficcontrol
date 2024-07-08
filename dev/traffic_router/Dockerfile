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
ARG TOMCAT_VERSION

ENV TC=/go/src/github.com/apache/trafficcontrol
VOLUME "$TC"
ENV JPDA_OPTS="-agentlib:jdwp=transport=dt_socket,address=*:5005,server=y,suspend=n" \
	JAVA_HOME=/usr/lib/jvm/java-11-openjdk M2_HOME=${TC}/trafficcontrol/.m2 \
	CATALINA_BASE=/opt/traffic_router \
	TRAFFIC_MONITOR_HOSTS=trafficmonitor \
	CATALINA_OPTS=-Dlog4j.configurationFile=/opt/traffic_router/conf/log4j2.xml
EXPOSE 3053:53/tcp \
	3053:53/udp \
	3080:80 \
	3443:443 \
	3333:3333 \
	2222:3443 \
	5005:5005

RUN apk add --no-cache openjdk11 inotify-tools maven tomcat-native openssl && ln -s /usr/lib/jvm/java-11-openjdk/bin/jdb /bin/jdb

ADD https://archive.apache.org/dist/tomcat/tomcat-9/v${TOMCAT_VERSION}/bin/apache-tomcat-${TOMCAT_VERSION}.tar.gz /opt/tomcat.tgz

RUN cd /opt && \
	tar -xf tomcat.tgz && \
	mv apache-tomcat-* tomcat && \
	rm -r tomcat.tgz tomcat/webapps/* /usr/share/java/maven-3/conf/settings.xml && \
	ln -s "${TC}/dev/traffic_router" /opt/traffic_router

COPY settings.xml /usr/share/java/maven-3/conf/settings.xml

CMD "${TC}/dev/traffic_router/run.sh"
