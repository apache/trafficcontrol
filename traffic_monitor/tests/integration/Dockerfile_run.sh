#!/usr/bin/env bash
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

# The following environment variables must be set (ordinarily by `docker run -e` arguments):
# TO_URI
# TO_USER
# TO_PASS
# CDN

# Check that env vars are set
envvars=( TESTTO_URI TESTTO_PORT TESTCACHES_URI TESTCACHES_PORT_START TM_URI )
for v in $envvars
do
	if [[ -z $$v ]]; then echo "$v is unset"; exit 1; fi
done


CFG_FILE=/traffic-monitor-integration-test.cfg

start() {
	printf "DEBUG traffic_monitor_integration starting\n"

	exec /traffic_monitor_integration_test -test.v -cfg $CFG_FILE
}

init() {
  wait_for_to

	curl -Lvsk ${TESTTO_URI}/api/1.2/cdns/fake/snapshot -X POST -d '
{
  "config": {
    "api.cache-control.max_age": "30",
    "consistent.dns.routing": "true",
    "coveragezone.polling.interval": "30",
    "coveragezone.polling.url": "30",
    "dnssec.dynamic.response.expiration": "60",
    "dnssec.enabled": "false",
    "domain_name": "monitor-integration.test",
    "federationmapping.polling.interval": "60",
    "federationmapping.polling.url": "foo",
    "geolocation.polling.interval": "30",
    "geolocation.polling.url": "foo",
    "keystore.maintenance.interval": "30",
    "neustar.polling.interval": "30",
    "neustar.polling.url": "foo",
    "soa": {

    },
    "dnssec.inception": "0",
    "ttls": {
      "admin":   "30",
      "expire":  "30",
      "minimum": "30",
      "refresh": "30",
      "retry":   "30"
		},
    "weight": "1",
    "zonemanager.cache.maintenance.interval": "30",
    "zonemanager.threadpool.scale": "1"
	},
	"contentServers": {
  	"server0": {
      "cacheGroup": "cg0",
      "profile": "Edge0",
      "fqdn": "server0.monitor-integration.test",
      "hashCount": 1,
      "hashId": "server0",
      "httpsPort" : null,
      "ip": "testcaches",
      "ip6": null,
      "locationId": "",
      "port" : 30000,
      "status": "REPORTED",
      "type": "EDGE",
      "deliveryServices": {"ds0":["ds0.monitor-integration.test"]},
      "routingDisabled": 0
      },
  	"server1": {
      "cacheGroup": "cg0",
      "profile": "Edge0",
      "fqdn": "server1.monitor-integration.test",
      "hashCount": 1,
      "hashId": "server1",
      "httpsPort" : null,
      "ip": "testcaches",
      "ip6": null,
      "locationId": "",
      "port" : 30001,
      "status": "REPORTED",
      "type": "EDGE",
      "deliveryServices": {"ds0":["ds0.monitor-integration.test"]},
      "routingDisabled": 0
      }
	},
	"deliveryServices": {
    "ds0": {
      "anonymousBlockingEnabled": false,
      "consistentHashQueryParams": [],
      "consistentHashRegex": "",
      "coverageZoneOnly": false,
      "dispersion": {
			  "limit": 1,
				"shuffled": false
			},
      "domains": ["ds0.monitor-integration.test"],
      "geolocationProvider": null,
      "matchsets": [
			  {
				  "protocol": "HTTP",
					"matchlist": [
					  {
						  "regex": "\\.*ds0\\.*",
							"match-type": "regex"
						}
					]
				}
			],
      "missLocation": {"lat": 0, "lon": 0},
      "protocol": {
        "acceptHttp": true,
        "acceptHttps": false,
        "redirectToHttps": false
      },
      "regionalGeoBlocking": "false",
      "responseHeaders": {},
      "requestHeaders": [],
      "soa": {
        "admin": "60",
        "expire": "60",
        "minimum": "60",
        "refresh": "60",
        "retry": "60"
			},
      "sslEnabled": false,
      "ttl": 60,
      "ttls": {
        "A": "60",
        "AAAA": "60",
        "DNSKEY": "60",
        "DS": "60",
        "NS": "60",
        "SOA": "60"
			},
      "maxDnsIpsForLocation": 3,
      "ip6RoutingEnabled": false,
      "routingName": "ccr",
      "bypassDestination": null,
      "deepCachingType": null,
      "geoEnabled": false,
      "geoLimitRedirectURL": null,
      "staticDnsEntries": []
    }
	},
	"edgeLocations": {
	  "cg0": {"latitude":0, "longitude":0}
	},
	"trafficRouterLocations": {
	  "tr0": {"latitude":0, "longitude":0}
	},
	"monitors": {
    "trafficmonitor": {
        "fqdn": "trafficmonitor.monitor-integration.test",
        "httpsPort": null,
        "ip": "trafficmonitor",
        "ip6": null,
        "location": "cg0",
        "port": 80,
        "profile": "Monitor0",
        "status": "REPORTED"
    }
	},
	"stats": {
    "CDN_name": "fake",
    "date": 1561000000,
    "tm_host": "testto",
    "tm_path": "/fake",
    "tm_user": "fake",
    "tm_version": "integrationtest/0.fake"
	}
}
'

	curl -Lvsk ${TESTTO_URI}/api/1.2/cdns/fake/configs/monitoring.json -X POST -d '
{
  "trafficServers": [
  	{
      "profile": "Edge0",
      "ip": "testcaches",
      "status": "REPORTED",
      "cacheGroup": "cg0",
      "ip6": null,
      "port": 30000,
      "httpsPort": null,
      "hostName": "server0",
      "fqdn": "server0.monitor-integration.test",
      "interfaceName": "bond0",
      "type": "EDGE",
      "hashId": "server0",
      "deliveryServices": {"ds0":["ds0.monitor-integration.test"]}
	  },
  	{
      "profile": "Edge0",
      "ip": "testcaches",
      "status": "REPORTED",
      "cacheGroup": "cg0",
      "ip6": null,
      "port": 30001,
      "httpsPort": null,
      "hostName": "server1",
      "fqdn": "server1.monitor-integration.test",
      "interfaceName": "bond0",
      "type": "EDGE",
      "hashId": "server1",
      "deliveryServices": {"ds0":["ds0.monitor-integration.test"]}
	  }
  ],
	"cacheGroups": [
    {
      "cg0": {
    		"name": "cg0",
    		"coordinates": {"latitude": 0, "longitude": 0}
      }
    }
  ],
	"config": {
    "peers.polling.interval": 30,
    "health.polling.interval": 2000,
    "heartbeat.polling.interval": 2000,
    "tm.polling.interval": 30
  },
	"trafficMonitors": [
	  {
      "port": 80,
      "ip6": "",
      "ip": "trafficmonitor",
      "hostName": "trafficmonitor",
      "fqdn": "trafficmonitor.traffic-monitor-integration.test",
      "profile": "Monitor0",
      "location": "cg0",
      "status": "REPORTED"
		}
  ],
	"deliveryServices": [
      {
  		  "xmlId": "ds0",
				"TotalTpsThreshold": 1000000,
				"status": "Available",
				"TotalKbpsThreshold": 10000000
      }
  ],
	"profiles": [
    {
      "parameters": {
        "health.connection.timeout": 10,
        "health.polling.url": "http://${hostname}/_astats?application=plugin.remap",
        "health.polling.format": "",
        "health.polling.type": "",
        "history.count": 0,
        "MinFreeKbps": 20000,
        "health_threshold": {}
      },
      "name": "Edge0",
      "type": "EDGE"
    },
    {
      "parameters": {
        "health.connection.timeout": 10,
        "health.polling.url": "",
        "health.polling.format": "",
        "health.polling.type": "",
        "history.count": 5,
        "MinFreeKbps": 20000,
        "health_threshold": {}
      },
      "name": "Monitor0",
      "type": "RASCAL"
    }
	]
}
'

	curl -Lvsk ${TESTTO_URI}/api/1.2/servers -X POST -d '
[
  {
    "cachegroup": "foo",
    "cachegroupId": 0,
    "cdnId": 1,
    "cdnName": "fake",
    "deliveryServices": null,
		"fqdn": "trafficmonitor.traffic-monitor-integration.test",
    "guid": "foo",
    "hostName": "trafficmonitor",
    "httpsPort": null,
    "id": 1,
    "iloIpAddress": null,
    "iloIpGateway": null,
    "iloIpNetmask": null,
    "iloPassword": null,
    "iloUsername": null,
    "interfaceMtu": null,
    "interfaceName": "bond0",
    "ip6Address": null,
    "ip6Gateway": null,
    "ipAddress": "trafficmonitor",
    "ipGateway": "192.0.0.1",
    "ipNetmask": "255.255.255.0",
    "lastUpdated": "2019",
    "mgmtIpAddress": null,
    "mgmtIpGateway": null,
    "mgmtIpNetmask": null,
    "offlineReason": "none",
    "physLocation": "",
    "physLocationId": 0,
    "profile": "Monitor0",
    "profileDesc": "nodesc",
    "profileId": 0,
    "rack": "",
    "revalPending": false,
    "routerHostName": "",
    "routerPortName": "",
    "status": "REPORTED",
    "statusId": 0,
    "tcpPort": 80,
    "type": "RASCAL",
    "typeId": 0,
    "updPending": false,
    "xmppId": "",
    "xmppPasswd": ""
  }
]
'

	# DEBUG
	printf "\n\ntestto:\n"
	curl -Lk ${TESTTO_URI}/api/1.2/cdns/foo/snapshot | head -5
	printf "\n\ntestcaches:\n"
	curl -Lk ${TESTCACHES_URI}:${TESTCACHES_PORT_START}/_astats | head -5
	printf "\n\ntraffic_monitor:\n"
	curl -Lk http://trafficmonitor | head -5


	cat > $CFG_FILE <<- EOF
{
  "trafficMonitor": {
    "url": "$TM_URI"
  },
  "default": {
    "session": {
      "timeoutInSecs": 30
    },
    "log": {
      "debug": "stdout",
      "event": "stdout",
      "info": "stdout",
      "error": "stdout",
      "warning": "stdout"
    }
  }
}
EOF

	echo "INITIALIZED=1" >> /etc/environment
}

wait_for_to() {
  while true; do
    curl -Lvsk ${TESTTO_URI}/api/1.2/servers 2>&1 1> /dev/null | grep -E '< HTTP/[0-9]\.?[0-9]* 200'
    RC=$?;
    if [[ $RC -eq 0 ]] ; then
			break;
		fi
		printf "Waiting for Traffic Ops to return a 200 OK\n";
		sleep 2;
	done
}

source /etc/environment
if [ -z "$INITIALIZED" ]; then init; fi
start
