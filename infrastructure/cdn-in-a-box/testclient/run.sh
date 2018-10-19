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

################################################################################
# Wait on SSL certificate generation
until [ -f "$X509_CA_DONE_FILE" ] 
do
  echo "Waiting on Shared SSL certificate generation"
  sleep 3
done

# Source the CIAB-CA shared SSL environment
source $X509_CA_ENV_FILE

# Trust the CIAB-CA at the System level
cp $X509_CA_CERT_FILE /etc/pki/ca-trust/source/anchors
update-ca-trust extract
################################################################################

VNC_DEPTH=${VNC_DEPTH:-32}
VNC_RESOLUTION=${VNC_RESOLUTION:-1440x900}

su -c "vncserver :0 -depth $VNC_DEPTH -geometry $VNC_RESOLUTION" - "$VNC_USER" && tail -F /home/dev/.vnc/vnc:0.log
