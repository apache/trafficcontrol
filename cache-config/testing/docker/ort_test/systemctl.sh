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

USAGE="\nsystemctl COMMAND NAME\n"

if [ -z $1 ] || [ -z $2 ]; then
  echo -e $USAGE
  exit 0
else
  COMMAND=$1
  NAME=$2
fi

if [ "$2" != "trafficserver.service" ]; then
  echo -e "\nFailed to start ${NAME}.service: Unit not found.n"
  exit 0
fi

case $COMMAND in 
  enable)
    ;;
  restart)
    /opt/trafficserver/bin/trafficserver restart
    ;;
  status)
    /opt/trafficserver/bin/trafficserver status
    ;;
  start)
    /opt/trafficserver/bin/trafficserver start
    ;;
  stop)
    /opt/trafficserver/bin/trafficserver stop
    ;;
esac

exit $?
