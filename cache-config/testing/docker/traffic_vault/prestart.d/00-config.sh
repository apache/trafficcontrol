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

# The following environment variables must be set, see 'variables.env'
# CERT_CITY
# CERT_COMPANY
# CERT_COUNTRY
# CERT_STATE

# Set the Riak certs in the config (this cert+key will be created in the run.sh script
sed -i -- 's/## ssl.certfile = $(platform_etc_dir)\/cert.pem/ssl.certfile = \/etc\/riak\/certs\/server.crt/g' /etc/riak/riak.conf
sed -i -- 's/## ssl.keyfile = $(platform_etc_dir)\/key.pem/ssl.keyfile = \/etc\/riak\/certs\/server.key/g' /etc/riak/riak.conf
sed -i -- 's/## ssl.cacertfile = $(platform_etc_dir)\/cacertfile.pem/ssl.cacertfile = \/etc\/riak\/certs\/ca-bundle.crt/g' /etc/riak/riak.conf
sed -i -- "s/nodename = riak@127.0.0.1/nodename = riak@0.0.0.0/g" /etc/riak/riak.conf
sed -i -- "s/listener.http.internal = 127.0.0.1:8098/listener.http.internal = 0.0.0.0:8098/g" /etc/riak/riak.conf
sed -i -- "s/listener.protobuf.internal = 127.0.0.1:8087/listener.protobuf.internal = 0.0.0.0:8087/g" /etc/riak/riak.conf
sed -i -- "s/## listener.https.internal = 127.0.0.1:8098/listener.https.internal = 0.0.0.0:8088/g" /etc/riak/riak.conf
sed -i -- "s/search = off/search = on/g" /etc/riak/riak.conf

echo "tls_protocols.tlsv1.1 = on" >> /etc/riak/riak.conf

mkdir /etc/riak/certs

openssl req -newkey rsa:2048 -nodes -keyout /etc/riak/certs/server.key -x509 -days 365 -out /etc/riak/certs/server.crt -subj "/C=$CERT_COUNTRY/ST=$CERT_STATE/L=$CERT_CITY/O=$CERT_COMPANY"
cp /etc/riak/certs/server.crt /etc/riak/certs/ca-bundle.crt
