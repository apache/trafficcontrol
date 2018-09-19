#!/bin/bash 
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
