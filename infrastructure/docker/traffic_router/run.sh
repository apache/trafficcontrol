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

# Script for running the Dockerfile for Traffic Router.
# The Dockerfile sets up a Docker image which can be used for any new container;
# This script, which should be run when the container is run (it's the ENTRYPOINT), will configure the container.
#
# The following environment variables must be set (ordinarily by `docker run -e` arguments):
# TRAFFIC_OPS_URI
# TRAFFIC_OPS_USER
# TRAFFIC_OPS_PASS
# TRAFFIC_MONITORS # list of semicolon-delimited FQDN:port monitors. E.g. `monitor.foo.com:80;monitor2.bar.org:80`
# ORIGIN_URI # origin server (e.g. hotair), used to create a delivery service

start() {
	./starttr.sh
	touch /opt/traffic_router/var/log/traffic_router.log
	exec tail -f /opt/traffic_router/var/log/traffic_router.log
}

init() {
	TMP_TO_COOKIE="$(curl -v -s -k -X POST --data '{ "u":"'"$TRAFFIC_OPS_USER"'", "p":"'"$TRAFFIC_OPS_PASS"'" }' $TRAFFIC_OPS_URI/api/4.0/user/login 2>&1 | grep 'Set-Cookie' | sed -e 's/.*mojolicious=\(.*\); expires.*/\1/')"
	echo "Got Cookie: $TMP_TO_COOKIE"

#  TMP_IP="$(ifconfig | grep -oE "inet addr:[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}" | grep -v '127.0.0.1' | cut -c11-)"
#	TMP_DOMAIN="$(grep -E "[:blank:]+.+\..+$" /etc/hosts | head -1 | rev | cut -d'.' -f2- --complement | rev)"
#	TMP_GATEWAY="$(route -n | grep -E "^0\.0\.0\.0[[:space:]]" | cut -f1 -d" " --complement | sed -e 's/^[ \t]*//' | cut -f1 -d" ")"
	TMP_IP=$IP
	TMP_DOMAIN=$DOMAIN
	TMP_GATEWAY=$GATEWAY

	TMP_CACHEGROUP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/cachegroups.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="mid-east"]; print match[0]')"
	echo "Got cachegroup ID: $TMP_CACHEGROUP_ID"

	TMP_SERVER_TYPE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="CCR"]; print match[0]')"
	echo "Got server type ID: $TMP_SERVER_TYPE_ID"

	TMP_SERVER_PROFILE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/profiles.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="CCR_CDN"]; print match[0]')"
	echo "Got server profile ID: $TMP_SERVER_PROFILE_ID"

	TMP_PHYS_LOCATION_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/phys_locations.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="plocation-nyc-1"]; print match[0]')"
	echo "Got phys location ID: $TMP_PHYS_LOCATION_ID"

	TMP_CDN_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/cdns.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="cdn"]; print match[0]')"
	echo "Got cdn ID: $TMP_CDN_ID"

	curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "host_name=$HOSTNAME" --data-urlencode "domain_name=$TMP_DOMAIN" --data-urlencode "interface_name=eth0" --data-urlencode "ip_address=$TMP_IP" --data-urlencode "ip_netmask=255.255.0.0" --data-urlencode "ip_gateway=$TMP_GATEWAY" --data-urlencode "interface_mtu=9000" --data-urlencode "cdn=$TMP_CDN_ID" --data-urlencode "cachegroup=$TMP_CACHEGROUP_ID" --data-urlencode "phys_location=$TMP_PHYS_LOCATION_ID" --data-urlencode "type=$TMP_SERVER_TYPE_ID" --data-urlencode "profile=$TMP_SERVER_PROFILE_ID" --data-urlencode "tcp_port=80" $TRAFFIC_OPS_URI/server/create

	TMP_SERVER_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/servers.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["hostName"]=="'"$HOSTNAME"'"]; print match[0]')"
	echo "Got server ID: $TMP_SERVER_ID"

	curl -v -k -H "Cookie: mojolicious=$TMP_TO_COOKIE" -X POST --data-urlencode "id=$TMP_SERVER_ID" --data-urlencode "status=ONLINE" $TRAFFIC_OPS_URI/server/updatestatus

	TMP_DELIVERY_SERVICE_TYPE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="DNS"]; print match[0]')"
	echo "Got delivery service type ID: $TMP_DELIVERY_SERVICE_TYPE_ID"

	curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "ds.xml_id=c2-service" --data-urlencode "ds.display_name=C2 Service" --data-urlencode "ds.cdn_id=$TMP_CDN_ID" --data-urlencode "ds.type=$TMP_DELIVERY_SERVICE_TYPE_ID"  --data-urlencode "ds.protocol=0" --data-urlencode "ds.dscp=0" --data-urlencode "ds.signed=0" --data-urlencode "ds.qstring_ignore=0" --data-urlencode "ds.geo_limit=0" --data-urlencode "ds.http_bypass_fqdn=" --data-urlencode "ds.initial_dispersion=1" --data-urlencode "ds.ipv6_routing_enabled=0" --data-urlencode "ds.range_request_handling=0" --data-urlencode "ds.dns_bypass_ip=" --data-urlencode "ds.dns_bypass_ip6=" --data-urlencode "ds.dns_bypass_cname=" --data-urlencode "ds.dns_bypass_ttl=30" --data-urlencode "ds.max_dns_answers=" --data-urlencode "ds.ccr_dns_ttl=30" --data-urlencode "ds.org_server_fqdn=$ORIGIN_URI" --data-urlencode "ds.multi_site_origin=0" --data-urlencode "ds.profile=$TMP_SERVER_PROFILE_ID" --data-urlencode "ds.global_max_mbps=" --data-urlencode "ds.global_max_tps=" --data-urlencode "ds.miss_lat=41.881944" --data-urlencode "ds.miss_long=-87.627778" --data-urlencode "ds.edge_header_rewrite=" --data-urlencode "ds.mid_header_rewrite=" --data-urlencode "ds.tr_response_headers=" --data-urlencode "ds.regex_remap=" --data-urlencode "ds.remap_text=" --data-urlencode "ds.long_desc=" --data-urlencode "ds.long_desc_1=" --data-urlencode "ds.long_desc_2=" --data-urlencode "ds.info_url=" --data-urlencode "ds.check_path=" --data-urlencode "ds.origin_shield=" --data-urlencode "ds.active=1" --data-urlencode "ds.regex=_____" --data-urlencode "re_type_0=HOST_REGEXP" --data-urlencode "re_order_0=0" --data-urlencode "re_re_0=.*\\.ds1\\..*" $TRAFFIC_OPS_URI/ds/create

	TMP_DELIVERY_SERVICE_LIVE_TYPE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="DNS_LIVE_NATNL"]; print match[0]')"
	echo "Got delivery service live type ID: $TMP_DELIVERY_SERVICE_LIVE_TYPE_ID"

	curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "ds.xml_id=ds2-live" --data-urlencode "ds.display_name=DS2 Live" --data-urlencode "ds.cdn_id=$TMP_CDN_ID" --data-urlencode "ds.type=$TMP_DELIVERY_SERVICE_LIVE_TYPE_ID"  --data-urlencode "ds.protocol=0" --data-urlencode "ds.dscp=0" --data-urlencode "ds.signed=0" --data-urlencode "ds.qstring_ignore=0" --data-urlencode "ds.geo_limit=0" --data-urlencode "ds.http_bypass_fqdn=" --data-urlencode "ds.initial_dispersion=1" --data-urlencode "ds.ipv6_routing_enabled=0" --data-urlencode "ds.range_request_handling=0" --data-urlencode "ds.dns_bypass_ip=" --data-urlencode "ds.dns_bypass_ip6=" --data-urlencode "ds.dns_bypass_cname=" --data-urlencode "ds.dns_bypass_ttl=30" --data-urlencode "ds.max_dns_answers=" --data-urlencode "ds.ccr_dns_ttl=30" --data-urlencode "ds.org_server_fqdn=$ORIGIN_URI" --data-urlencode "ds.multi_site_origin=0" --data-urlencode "ds.profile=$TMP_SERVER_PROFILE_ID" --data-urlencode "ds.global_max_mbps=" --data-urlencode "ds.global_max_tps=" --data-urlencode "ds.miss_lat=41.881944" --data-urlencode "ds.miss_long=-87.627778" --data-urlencode "ds.edge_header_rewrite=" --data-urlencode "ds.mid_header_rewrite=" --data-urlencode "ds.tr_response_headers=" --data-urlencode "ds.regex_remap=" --data-urlencode "ds.remap_text=" --data-urlencode "ds.long_desc=" --data-urlencode "ds.long_desc_1=" --data-urlencode "ds.long_desc_2=" --data-urlencode "ds.info_url=" --data-urlencode "ds.check_path=" --data-urlencode "ds.origin_shield=" --data-urlencode "ds.active=1" --data-urlencode "ds.regex=_____" --data-urlencode "re_type_0=HOST_REGEXP" --data-urlencode "re_order_0=0" --data-urlencode "re_re_0=.*\\.ds2-live\\..*" $TRAFFIC_OPS_URI/ds/create

	TMP_DELIVERY_SERVICE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/deliveryservices.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["xmlId"]=="c2-service"]; print match[0]')"
	echo "Got delivery ID: $TMP_DELIVERY_SERVICE_ID"

	TMP_DELIVERY_SERVICE_HTTP_TYPE_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/types.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["name"]=="HTTP"]; print match[0]')"
	echo "Got delivery service http type ID: $TMP_DELIVERY_SERVICE_HTTP_TYPE_ID"

	curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data-urlencode "ds.xml_id=c3-service" --data-urlencode "ds.display_name=C3 Service" --data-urlencode "ds.cdn_id=$TMP_CDN_ID" --data-urlencode "ds.type=$TMP_DELIVERY_SERVICE_HTTP_TYPE_ID"  --data-urlencode "ds.protocol=0" --data-urlencode "ds.dscp=0" --data-urlencode "ds.signed=0" --data-urlencode "ds.qstring_ignore=0" --data-urlencode "ds.geo_limit=0" --data-urlencode "ds.http_bypass_fqdn=" --data-urlencode "ds.initial_dispersion=1" --data-urlencode "ds.ipv6_routing_enabled=0" --data-urlencode "ds.range_request_handling=0" --data-urlencode "ds.dns_bypass_ip=" --data-urlencode "ds.dns_bypass_ip6=" --data-urlencode "ds.dns_bypass_cname=" --data-urlencode "ds.dns_bypass_ttl=30" --data-urlencode "ds.max_dns_answers=" --data-urlencode "ds.ccr_dns_ttl=30" --data-urlencode "ds.org_server_fqdn=$ORIGIN_URI" --data-urlencode "ds.multi_site_origin=0" --data-urlencode "ds.profile=$TMP_SERVER_PROFILE_ID" --data-urlencode "ds.global_max_mbps=" --data-urlencode "ds.global_max_tps=" --data-urlencode "ds.miss_lat=41.881944" --data-urlencode "ds.miss_long=-87.627778" --data-urlencode "ds.edge_header_rewrite=" --data-urlencode "ds.mid_header_rewrite=" --data-urlencode "ds.tr_response_headers=" --data-urlencode "ds.regex_remap=" --data-urlencode "ds.remap_text=" --data-urlencode "ds.long_desc=" --data-urlencode "ds.long_desc_1=" --data-urlencode "ds.long_desc_2=" --data-urlencode "ds.info_url=" --data-urlencode "ds.check_path=" --data-urlencode "ds.origin_shield=" --data-urlencode "ds.active=1" --data-urlencode "ds.regex=_____" --data-urlencode "re_type_1=HOST_REGEXP" --data-urlencode "re_order_1=0" --data-urlencode "re_re_1=.*\\.ds2\\..*" $TRAFFIC_OPS_URI/ds/create

	TMP_DELIVERY_SERVICE_HTTP_ID="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/deliveryservices.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["xmlId"]=="c3-service"]; print match[0]')"
	echo "Got delivery http ID: $TMP_DELIVERY_SERVICE_HTTP_ID"


	TMP_EDGES="$(curl -s -k -X GET -H "Cookie: mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/servers.json | python -c 'import json,sys;obj=json.load(sys.stdin);match=[x["id"] for x in obj["response"] if x["type"]=="EDGE"]; print "\n".join(match)')"
	echo "Got edges: $TMP_EDGES"
	TMP_EDGE_STR="$(echo "$TMP_EDGES" | xargs -I {} echo 'serverid_{}=on&' | tr -d '\n' | rev | cut -c 2- | rev)"
	echo "Got edge str: $TMP_EDGE_STR"

	curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data $TMP_EDGE_STR $TRAFFIC_OPS_URI/dss/$TMP_DELIVERY_SERVICE_ID/update
	curl -v -k -X POST -H "Cookie: mojolicious=$TMP_TO_COOKIE" --data $TMP_EDGE_STR $TRAFFIC_OPS_URI/dss/$TMP_DELIVERY_SERVICE_HTTP_ID/update

	sed -i -- "s/# traffic_monitor.bootstrap.hosts=some-traffic_monitor.company.net:80;/traffic_monitor.bootstrap.hosts=$TRAFFIC_MONITORS/g" /opt/traffic_router/conf/traffic_monitor.properties

	sed -i -- "s/traffic_ops.username=admin/traffic_ops.username=$TRAFFIC_OPS_USER/g" /opt/traffic_router/conf/traffic_ops.properties
	sed -i -- "s/traffic_ops.password=FIXME/traffic_ops.password=$TRAFFIC_OPS_PASS/g" /opt/traffic_router/conf/traffic_ops.properties

	curl -k -v -X PUT --cookie "mojolicious=$TMP_TO_COOKIE" $TRAFFIC_OPS_URI/api/4.0/snapshot/cdn

	echo "INITIALIZED=1" >> /etc/environment
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
