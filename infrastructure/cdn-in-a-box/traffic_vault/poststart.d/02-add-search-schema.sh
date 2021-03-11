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
source /to-access.sh

curl -kvs -XPUT -H 'Content-Type:application/xml' "https://$TV_ADMIN_USER:$TV_ADMIN_PASSWORD@$TV_FQDN:$TV_HTTPS_PORT/search/schema/sslkeys" -d @/sslkeys.xml || cat /var/log/riak/*

curl -kvs -XPUT -H 'Content-Type:application/json' "https://$TV_ADMIN_USER:$TV_ADMIN_PASSWORD@$TV_FQDN:$TV_HTTPS_PORT/search/index/sslkeys" -d '{"schema":"sslkeys"}'

curl -kvs -XPUT -H 'Content-Type:application/json' "https://$TV_ADMIN_USER:$TV_ADMIN_PASSWORD@$TV_FQDN:$TV_HTTPS_PORT/buckets/ssl/props" -d'{"props":{"search_index":"sslkeys"}}'

to-enroll "tv" ALL "" "8088" "8088" || (while true; do echo "enroll failed."; sleep 3 ; done)
