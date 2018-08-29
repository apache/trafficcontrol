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
# Dockerfile to build Traffic Router 3.0 
# Based on CentOS 7.x
############################################################

FROM centos:7
MAINTAINER Dan Kirkwood

# Default values for TMCAT RPM and RPM -- override with `docker build --build-arg JDK=...'
ARG TRAFFIC_ROUTER_RPM=traffic_router/traffic_router.rpm
ARG TOMCAT_RPM=traffic_router/tomcat.rpm 
ADD $TRAFFIC_ROUTER_RPM /traffic_router.rpm
ADD $TOMCAT_RPM /tomcat.rpm

### Common for all sub-component builds
RUN	yum -y install epel-release git rpm-build net-tools iproute nc wget tar unzip \
          perl-JSON perl-WWWCurl which tcpdump autoconf automake gcc gcc-c++ \
          apr apr-devel java-1.8.0-openjdk java-1.8.0-openjdk-devel make \
					openssl openssl-devel && \
	  yum -y clean all

WORKDIR /opt
RUN curl -L -b "oraclelicense=a" http://download.oracle.com/otn-pub/java/jdk/8u131-b11/d54c1d3a095b4ff2b6607d096fa80163/jdk-8u131-linux-x64.tar.gz | tar xzf - && \
    curl -L https://repo.maven.apache.org/maven2/org/apache/maven/apache-maven/3.3.9/apache-maven-3.3.9-bin.tar.gz | tar xzf -
RUN ln -s jdk* java && ln -s apache-maven* maven

RUN alternatives --install /usr/bin/java java /opt/java/bin/java 2 && \
    alternatives --install /usr/bin/jar jar /opt/java/bin/jar 2 && \
    alternatives --install /usr/bin/javac javac /opt/java/bin/javac 2 && \
    alternatives --set java /opt/java/bin/java && \
    alternatives --set jar /opt/java/bin/jar && \
    alternatives --set javac /opt/java/bin/javac && \
    alternatives --install /usr/bin/mvn mvn /opt/maven/bin/mvn 1

# Install Tomcat and Tomcat Native 
RUN cd / && rpm2cpio /traffic_router.rpm | cpio -ivd && \
    rpm2cpio /tomcat.rpm | cpio -ivd && \ 
 		tar -C /opt -xvpf /opt/tomcat/bin/tomcat-native.tar.gz

RUN	cd $(find /opt -maxdepth 1 -type d | grep tomcat-native )/native && \
    ./configure --prefix=/opt/tomcat/lib --with-apr=`which apr-1-config` --with-java-home=/opt/java --with-ssl && \
		make install
     
EXPOSE 53 80 3333
ADD traffic_router/run.sh traffic_router/to-access.sh /
CMD /run.sh
