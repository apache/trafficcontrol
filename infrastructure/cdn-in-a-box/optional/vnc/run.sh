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

################################################################################
# Wait on SSL certificate generation
set +x 
set +e 
set +m

[[ -f "/usr/local/sbin/set-dns.sh" ]] && /usr/local/sbin/set-dns.sh
[[ -f "/usr/local/sbin/insert-self-into-dns.sh" ]] && /usr/local/sbin/insert-self-into-dns.sh

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

# Trust the CIAB-CA at the System level
cp $X509_CA_CERT_FULL_CHAIN_FILE /etc/pki/ca-trust/source/anchors
update-ca-trust extract
################################################################################

VNC_DEPTH=${VNC_DEPTH:-32}
VNC_RESOLUTION=${VNC_RESOLUTION:-1440x900}
VNC_INSTANCE_NUM=9

if [[ -z $VNC_PASSWD ]] ; then 
   echo "WARNING: no \$VNC_PASSWD environment has been set."
   sleep 10
else 
   echo "Changing vnc console password to: $VNC_PASSWD" 
   echo -en "$VNC_PASSWD" | vncpasswd -f > /home/$VNC_USER/.vnc/passwd
   sync 
   sleep 1
fi

su -c "vncserver :$VNC_INSTANCE_NUM -depth $VNC_DEPTH -geometry $VNC_RESOLUTION" - "$VNC_USER" 

tail -F -- /home/ciabuser/.vnc/vnc*:9.log
