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
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set (ordinarily by `docker run -e` arguments):
# TRAFFIC_OPS_URI
# TRAFFIC_OPS_USER
# TRAFFIC_OPS_PASS
# TRAFFIC_MONITORS # list of semicolon-delimited FQDN:port monitors. E.g. `monitor.foo.com:80;monitor2.bar.org:80`
# ORIGIN_URI # origin server (e.g. hotair), used to create a delivery service

export JAVA_HOME="$(command -v java | xargs realpath | xargs dirname)/.."
export CATALINA_PID=/opt/traffic_router/temp/tomcat.pid
export CATALINA_HOME=/opt/tomcat
export CATALINA_BASE=/opt/traffic_router
export CATALINA_OUT=/opt/tomcat/logs/catalina.log
source /opt/traffic_router/conf/startup.properties
/opt/tomcat/bin/shutdown.sh
