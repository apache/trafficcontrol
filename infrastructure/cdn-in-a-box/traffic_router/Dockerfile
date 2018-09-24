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
ARG JDK8_RPM=http://download.oracle.com/otn-pub/java/jdk/8u181-b13/96a7b8442fe848ef90c96a2fad6ed6d1/jdk-8u181-linux-x64.rpm
ARG TRAFFIC_ROUTER_RPM=traffic_router/traffic_router.rpm
ARG TOMCAT_RPM=traffic_router/tomcat.rpm 

### Common for all sub-component builds
RUN	yum -y install epel-release git rpm-build net-tools iproute nc wget tar unzip \
          perl-JSON perl-WWWCurl which make autoconf automake gcc gcc-c++ apr apr-devel \
					openssl openssl-devel bind-utils net-tools perl-JSON-PP gettext && \ 
    yum -y clean all

ADD $TRAFFIC_ROUTER_RPM /traffic_router.rpm
ADD $TOMCAT_RPM /tomcat.rpm

## Install and link java_home to /opt/java
RUN curl -kvsL -o /jdk8.rpm -H  'Cookie: oraclelicense=accept-securebackup-cookie;' "$JDK8_RPM"  && \
    rpm -ivh /jdk8.rpm && \
    ln -sfv $(find /usr/java -mindepth 1 -maxdepth 1 -type d -name jdk\*)   /opt/java

# Install Tomcat and Tomcat Native 
RUN cd / && rpm2cpio /traffic_router.rpm | cpio -ivd && \
    rpm2cpio /tomcat.rpm | cpio -ivd && \ 
 		tar -C /opt -xvpf /opt/tomcat/bin/tomcat-native.tar.gz

# Compile and install Tomcat-native
RUN	cd $(find /opt -maxdepth 1 -type d | grep tomcat-native )/native && \
    ./configure --prefix=/opt/tomcat --with-apr=`which apr-1-config` --with-java-home=/opt/java --with-ssl && \
		make install
     
ADD enroller/server_template.json \ 
    traffic_router/run.sh \ 
    traffic_ops/to-access.sh \ 
    /

EXPOSE 53 80 3333

CMD /run.sh
