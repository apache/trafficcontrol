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

function longer_dns_timeout() {
  local day_in_ms dns_properties;
  day_in_ms=$(( 1000 * 60 * 60 * 24 )); # Timing out debugging after 1 day seems fair
  dns_properties=/opt/traffic_router/conf/dns.properties;
  <<-DNS_CONFIG_LINES cat >> $dns_properties;
		dns.tcp.timeout.task=$(( day_in_ms ))
		dns.udp.timeout.task=$(( day_in_ms ))
DNS_CONFIG_LINES
}

set-dns.sh
insert-self-into-dns.sh

# Global Vars for FQDNs, ports, etc
source /to-access.sh

CATALINA_HOME="/opt/tomcat"
CATALINA_BASE="/opt/traffic_router"
CATALINA_OUT="$CATALINA_HOME/logs/catalina.log"
CATALINA_LOG="$CATALINA_HOME/logs/catalina.$(date +%Y-%m-%d).log"
CATALINA_PID="$CATALINA_BASE/temp/tomcat.pid"

CATALINA_OPTS="\
  -server -Xms2g -Xmx8g \
  -Djava.library.path=/usr/lib64:$CATALINA_BASE/lib:$CATALINA_HOME/lib \
  -Dlog4j.configurationFile=$CATALINA_BASE/conf/log4j2.xml \
  -Dorg.apache.catalina.connector.Response.ENFORCE_ENCODING_IN_GET_WRITER=false \
  -XX:+UseG1GC \
  -XX:+UnlockExperimentalVMOptions \
  -XX:InitiatingHeapOccupancyPercent=30"

if [[ "$TR_DEBUG_ENABLE" == true ]]; then
    export JPDA_OPTS="-agentlib:jdwp=transport=dt_socket,address=*:5005,server=y,suspend=n";
    longer_dns_timeout;
fi;

JAVA_HOME=/opt/java
JAVA_OPTS="\
  -Djava.library.path=/usr/lib64 \
  -Dcache.config.json.refresh.period=5000 \
  -Djava.awt.headless=true \
  -Djava.security.egd=file:/dev/./urandom"

TO_PROPERTIES="$CATALINA_BASE/conf/traffic_ops.properties"
TM_PROPERTIES="$CATALINA_BASE/conf/traffic_monitor.properties"
LOGFILE="/var/log/traffic_router/traffic_router.log"
ACCESSLOG="/var/log/traffic_router/access.log"

export JAVA_HOME JAVA_OPTS
export TO_PROPERTIES TM_PROPERTIES
export CATALINA_HOME CATALINA_BASE CATALINA_OPTS CATALINA_OUT CATALINA_PID

# Wait on SSL certificate generation
until [[ -f "$X509_CA_ENV_FILE" ]]
do
  echo "Waiting on Shared SSL certificate generation"
  sleep 3
done

# Source the CIAB-CA shared SSL environment
until [[ -n "$X509_GENERATION_COMPLETE" ]]
do
  echo "Waiting on X509 vars to be defined"
  sleep 1
  source "$X509_CA_ENV_FILE"
done

# Copy the CIAB-CA certificate to the traffic_router conf so it can be added to the trust store
cp $X509_CA_ROOT_CERT_FILE $CATALINA_BASE/conf
cp $X509_CA_INTR_CERT_FILE $CATALINA_BASE/conf
cp $X509_CA_CERT_FULL_CHAIN_FILE /etc/pki/ca-trust/source/anchors
update-ca-trust extract

# Add traffic
for crtfile in $(find $CATALINA_BASE/conf -name \*.crt -type f)
do
  alias=$(echo $crtfile |sed -e 's/.crt//g' |tr [:upper:] [:lower:]);
  cacerts=$(find $JAVA_HOME -follow -name cacerts); echo $cacerts;
  keytool=$JAVA_HOME/bin/keytool;

  $keytool -list -alias $alias -keystore $cacerts -storepass changeit -noprompt > /dev/null;

  if [ $? -ne 0 ]; then
     echo "Installing certificate ${crtfile}..";
     $keytool -import -trustcacerts -file $crtfile -alias $alias -keystore $cacerts -storepass changeit -noprompt;
  fi;
done

/opt/traffic_router/conf/generatingCerts.sh

# Configure TO properties
# File: /opt/traffic_router/conf/traffic_ops.properties
echo "" > $TO_PROPERTIES
echo "traffic_ops.username=$TO_ADMIN_USER" >> $TO_PROPERTIES
echo "traffic_ops.password=$TO_ADMIN_PASSWORD" >> $TO_PROPERTIES

# Configure TM properties
# File: /opt/traffic_router/conf/traffic_monitor.properties
echo "traffic_monitor.bootstrap.hosts=$TM_FQDN:$TM_PORT;" >> $TM_PROPERTIES
echo "traffic_monitor.properties.reload.period=60000" >> $TM_PROPERTIES

# Enroll Traffic Router
to-enroll tr || (while true; do echo "enroll failed."; sleep 3 ; done)

# Wait for traffic monitor
until nc $TM_FQDN $TM_PORT </dev/null >/dev/null 2>&1; do
  echo "Waiting for Traffic Monitor to start..."
  sleep 3
done

touch $LOGFILE $ACCESSLOG
if [[ "$TR_DEBUG_ENABLE" == true ]]; then
    exec /opt/tomcat/bin/catalina.sh jpda start &
else
    exec /opt/tomcat/bin/catalina.sh run &
fi;

tail -F $CATALINA_OUT $CATALINA_LOG $LOGFILE $ACCESSLOG
