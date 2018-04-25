#!/usr/bin/env bash
#
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

set -ex
env

export TP=/opt/traffic_portal/

envvars=( TO_PORT )
for v in ${envvars[*]}; do
	val=${!v}
	[[ -z $val ]] && echo "$v is unset" && exit 1
done

mkdir -p /etc/traffic_portal/conf/
mkdir -p /etc/pki/tls/private
mkdir -p /etc/pki/tls/certs

cat >/etc/traffic_portal/conf/config.js <<-TPCONF
/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// this is the config that is consumed by /server.js on traffic portal startup (sudo service traffic_portal start)
module.exports = {
    timeout: '120s',
    useSSL: true, // set to true if you plan to use https (self-signed or trusted certs).
    port: 80, // set to http port
    sslPort: 61443, // set to https port
    // if useSSL is true, generate ssl certs and provide the proper locations.
    ssl: {
        key:    '/etc/pki/tls/private/localhost.key',
        cert:   '/etc/pki/tls/certs/localhost.crt',
        ca:     [ '/etc/pki/tls/certs/ca-bundle.crt' ]
    },
    // set api 'base_url' to the traffic ops api url (all api calls made from the traffic portal will be proxied to the api base_url)
    api: {
        base_url: "https://${TO_HOST}/api/"
    },
    // default static files location (this is where the traffic portal html, css and javascript was installed. rpm installs these files at /opt/traffic_portal/public
    // change this to ./app/dist/public/ if you are running locally for development
    files: {
        static: '/opt/traffic_portal/public'
    },
    // default log location (this is where traffic_portal logs are written)
    // change this to ./server/log/access.log if you are running traffic portal locally for development
    log: {
        stream: '/var/log/traffic_portal/access.log'
    },
    reject_unauthorized: 0 // 0 if using self-signed certs, 1 if trusted certs
};

TPCONF

/etc/init.d/traffic_portal start

# Give trafficportal a second to create the logfile
sleep 3

# Print out the status, because trafficportal will fail silently if you let it
/etc/init.d/traffic_portal status

#Fallback
if [[ ! -f /var/log/traffic_portal/traffic_portal.log ]]; then
	touch /var/log/traffic_portal/traffic_portal.log
fi

exec tail -f /var/log/traffic_portal/traffic_portal.log
