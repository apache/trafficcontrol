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
set -x 

# Enable Security
$RIAK_ADMIN security enable
$RIAK_ADMIN security add-group admins
$RIAK_ADMIN security add-group keysusers

# Add users
$RIAK_ADMIN security add-user "$TV_ADMIN_USER" password="$TV_ADMIN_PASSWORD" groups=admins
$RIAK_ADMIN security add-user "$TV_RIAK_USER" password="$TV_RIAK_PASSWORD" groups=keysusers
$RIAK_ADMIN security add-source "$TV_ADMIN_USER" 0.0.0.0/0 password
$RIAK_ADMIN security add-source "$TV_RIAK_USER" 0.0.0.0/0 password

# Grant privs to admins for everything
$RIAK_ADMIN security grant riak_kv.list_buckets,riak_kv.list_keys,riak_kv.get,riak_kv.put,riak_kv.delete on any to admins

# Grant privs to keysuser for ssl, dnssec, and url_sig_keys buckets only
$RIAK_ADMIN security grant riak_kv.get,riak_kv.put,riak_kv.delete on default ssl to keysusers
$RIAK_ADMIN security grant riak_kv.get,riak_kv.put,riak_kv.delete on default dnssec to keysusers
$RIAK_ADMIN security grant riak_kv.get,riak_kv.put,riak_kv.delete on default url_sig_keys to keysusers
$RIAK_ADMIN security grant riak_kv.get,riak_kv.put,riak_kv.delete on default cdn_uri_sig_keys  to keysusers
