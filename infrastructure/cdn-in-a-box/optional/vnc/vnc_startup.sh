#!/bin/bash -x
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

VNC_DEPTH=${VNC_DEPTH:-32}
VNC_RESOLUTION=${VNC_DEPTH:-1440x900}

startfluxbox &

sleep 3

vncconfig -iconic &
xterm -bg black -fg white +sb &

until nc 'trafficportal.infra.ciab.test' 443 </dev/null >/dev/null 2>&1; do
  echo "Waiting for Traffic Portal to start" 
  sleep 2
done

firefox 'http://trafficportal.infra.ciab.test' &
