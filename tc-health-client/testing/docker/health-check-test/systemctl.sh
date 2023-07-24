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
#

# This is a work around for testing t3c which uses systemctl.
# systemctl does not work in a container very well so this script
# replaces systemctl in the container and always returns a 
# sucessful result to t3c.

USAGE="\nsystemctl COMMAND UNIT\n"

if [ -z $1 ]; then
  echo -e $USAGE
  exit 0
else
  COMMAND=$1
  UNIT=$2
fi

if [ "$COMMAND" != "daemon-reload" ]; then
  if [ "$UNIT" != "trafficserver" ] && [ "$UNIT" != "tc-health-client" ]; then
    echo -e "\nFailed to start ${UNIT}: Unit not found"
    exit 0
  fi
fi

case $COMMAND in 
  daemon-reload)
    ;;
  enable)
    ;;
  restart)
    case $UNIT in
      trafficserver) 
        /opt/trafficserver/bin/trafficserver restart
        ;;
      tc-health-client)
        kill `cat /var/run/tc-health-client.pid`
        tc-health-client -vvv &
        ;;
    esac
    ;;
  status)
    if [ "${UNIT}" = "trafficserver" ]; then
      /opt/trafficserver/bin/trafficserver status
    fi
    ;;
  start)
    case $UNIT in
      trafficserver) 
        /opt/trafficserver/bin/trafficserver start
        ;;
      tc-health-client)
        nohup tc-health-client -vvv &
        ;;
    esac
    ;;
  stop)
    if [ "${UNIT}" = "trafficserver" ]; then
      /opt/trafficserver/bin/trafficserver stop
    fi
    case $UNIT in
      trafficserver) 
        /opt/trafficserver/bin/trafficserver start
        ;;
      tc-health-client)
        kill `cat /var/run/tc-health-client.pid`
        ;;
    esac
    ;;
esac

exit $?
