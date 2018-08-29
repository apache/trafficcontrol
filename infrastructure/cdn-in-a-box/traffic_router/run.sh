#!/usr/bin/env bash
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
NAME="Traffic Router Application"

CATALINA_HOME=/opt/tomcat
CATALINA_BASE=/opt/traffic_router
CATALINA_OUT=$CATALINA_HOME/logs/catalina.log
CATALINA_PID=$CATALINA_BASE/temp/tomcat.pid

CATALINA_OPTS="\
  -server -Xms2g -Xmx8g \
  -Djava.library.path=$CATALINA_HOME/lib \
  -Dlog4j.configuration=file://$CATALINA_BASE/conf/log4j.properties \
  -Dorg.apache.catalina.connector.Response.ENFORCE_ENCODING_IN_GET_WRITER=false \
  -XX:+UseG1GC \
  -XX:+UnlockExperimentalVMOptions \
  -XX:InitiatingHeapOccupancyPercent=30"

JAVA_HOME=/opt/java
JAVA_OPTS="\
  -Djava.awt.headless=true \
  -Djava.security.egd=file:/dev/./urandom"

TO_PROPERTIES="$CATALINA_BASE/conf/traffic_ops.properties"
TM_PROPERTIES="$CATALINA_BASE/conf/traffic_monitor.properties"
LOGFILE="$CATALINA_BASE/var/log/traffic_router.log"
ACCESSLOG="$CATALINA_BASE/var/log/access.log"


export JAVA_HOME JAVA_OPTS
export TO_PROPERTIES TM_PROPERTIES 
export CATALINA_HOME CATALINA_BASE CATALINA_OPTS CATALINA_OUT CATALINA_PID

# Wait for Enroller
export TO_USER=$TO_ADMIN_USER
export TO_PASSWORD=$TO_ADMIN_PASSWORD
source /to-access.sh
to-enroll $(hostname -s)

# Configure TO properties
# File: /opt/traffic_router/conf/traffic_ops.properties
echo "" > $TO_PROPERTIES
echo "traffic_ops.username=$TO_ADMIN_USER" >> $TO_PROPERTIES
echo "traffic_ops.password=$TO_ADMIN_PASSWORD" >> $TO_PROPERTIES

# Configure TM properties
# File: /opt/traffic_router/conf/traffic_monitor.properties
echo "traffic_monitor.bootstrap.hosts=trafficmonitor:80;" >> $TM_PROPERTIES
echo "traffic_monitor.properties.reload.period=60000" >> $TM_PROPERTIES

touch $LOGFILE $ACCESSLOG
tail -F $CATALINA_BASE/var/log/traffic_router.log $CATALINA_BASE/var/log/access.log &  \
  /opt/tomcat/bin/catalina.sh run 
