#!/bin/bash
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

# Script for running the Dockerfile for Traffic Router.
# The Dockerfile sets up a Docker image which can be used for any new container;
# This script simulates the systemd unit file that is used to start traffic router on 
# servers in the real world, but in Docker containers systemd is disabled. 
# Therefore it is important to keep this script up to date with any changes that are
# made to traffic_router/build/build_rpm.sh and traffic_router/build/pom.xml

export JAVA_HOME="$(command -v java | xargs realpath | xargs dirname)/.."
export CATALINA_PID=/opt/traffic_router/temp/tomcat.pid
export CATALINA_HOME=/opt/tomcat
export CATALINA_BASE=/opt/traffic_router
export CATALINA_OUT=/opt/tomcat/logs/catalina.log
export CATALINA_OPTS="\
  -server -Xms512m -Xmx1g \
  -Dlog4j.configuration=$CATALINA_BASE/conf/log4j.properties \
  -Djava.library.path=/usr/lib64 \
  -Dorg.apache.catalina.connector.Response.ENFORCE_ENCODING_IN_GET_WRITER=false \
  -XX:+UseG1GC \
  -XX:+UnlockExperimentalVMOptions \
  -XX:InitiatingHeapOccupancyPercent=30"
export JAVA_OPTS="\
  -Djava.awt.headless=true \
  -Djava.security.egd=file:/dev/./urandom"

ulimit -c unlimited
/opt/tomcat/bin/startup.sh
