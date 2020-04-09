..
..
.. Licensed under the Apache License, Version 2.0 (the "License");
.. you may not use this file except in compliance with the License.
.. You may obtain a copy of the License at
..
..     http://www.apache.org/licenses/LICENSE-2.0
..
.. Unless required by applicable law or agreed to in writing, software
.. distributed under the License is distributed on an "AS IS" BASIS,
.. WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
.. See the License for the specific language governing permissions and
.. limitations under the License.
..

.. _to-api-v1-cdns-name-configs-routing:

*********************************
``cdns/{{name}}/configs/routing``
*********************************
.. deprecated:: ATCv4
	Use the ``GET`` method of :ref:`to-api-v1-cdns-name-snapshot` instead.

``GET``
=======
Retrieves CDN routing information.

:Auth. Required: Yes
:Roles Required: None
:Response Type:

Request Structure
-----------------
.. table:: Request Path Parameters

	+----------+----------+-------------------------------------+
	|   Name   | Required | Description                         |
	+==========+==========+=====================================+
	| ``name`` | yes      | The name of the CDN to be inspected |
	+----------+----------+-------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.5/cdns/CDN-in-a-Box/configs/routing HTTP/1.1
	User-Agent: curl/7.29.0
	Host: trafficops.infra.ciab.test
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cacheGroups: A collection of objects that represent :term:`Cache Groups`.

	:coordinates: An object that represents the geographic location of the :term:`Cache Group`

		:latitude:  number
		:longitude: number

	:name: string

:config: object

	:coveragezone.polling.url:       string
	:domain_name:                    string
	:geolocation.polling.interval:   integer
	:geolocation.polling.url:        string
	:geolocation6.polling.interval:  integer
	:geolocation6.polling.url:       string
	:tcoveragezone.polling.interval: integer
	:tld.soa.admin:                  string
	:tld.soa.expire:                 integer
	:tld.soa.minimum:                integer
	:tld.soa.refresh:                integer
	:tld.soa.retry:                  integer
	:tld.ttls.A:                     integer
	:tld.ttls.AAAA:                  integer
	:tld.ttls.NS:                    integer
	:tld.ttls.SOA:                   integer

:deliveryServices: An array of delivery services.

	:coverageZoneOnly: boolean
	:bypassDestination: object

		:maxDnsIpsForLocation: integer
		:ttl:                  integer
		:type:                 string

	:geoEnabled:       string
	:matchSets:        array

		:protocol:  string
		:matchList: array

			:matchType: string
			:regex:     string

	:missCoordinates: object

		:latitude:  number
		:longitude: number

	:soa: object

		:admin:   string
		:expire:  integer
		:minimum: integer
		:refresh: integer
		:retry:   integer

	:ttl:              integer
	:ttls: object

		:A:    integer
		:AAAA: integer
		:NS:   integer
		:SOA:  integer

	:xmlId:            string

:stats: object

	:cdnName:           string
	:date:              integer
	:trafficOpsHost:    string
	:trafficOpsPath:    string
	:trafficOpsUser:    string
	:trafficOpsVersion: string

:trafficMonitors: An array of Traffic Monitors

	:fqdn:     string
	:hostName: string
	:ip6:      string
	:ip:       string
	:location: string
	:port:     integer
	:profile:  string
	:status:   string

:trafficRouters: object

	:apiPort:  integer
	:fqdn:     string
	:hostName: string
	:ip6:      string
	:ip:       string
	:location: string
	:port:     integer
	:profile:  integer
	:status:   string

:trafficServers: An array of Traffic Servers.

	:cacheGroup:       string
	:deliveryServices: array

		:xmlId:    string
		:remaps:   array
		:hostName: string

	:fqdn:          string
	:hashId:        string
	:interfaceName: string
	:ip:            string
	:ip6:           string
	:port:          integer
	:profile:       string
	:status:        string
	:type:          string

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Mon, 27 Jan 2020 19:20:14 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Mon, 27 Jan 2020 23:20:14 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: Dxgtd9e67IRb9HyhPxG94zijfpCB44mdstlf5ZXokCQoAUKbcPaTu2szPMgineWmNvWxAfgrXo0ZVUnCRqxh7A==
	Transfer-Encoding: chunked

	{
		"alerts": [
			{
				"level": "warning",
				"text": "This endpoint is deprecated, please use 'GET /cdns/{{name}}/snapshot' instead"
			}
		],
		"response": {
			"trafficServers": [
				{
					"profile": "ATS_MID_TIER_CACHE",
					"ip": "172.16.239.5",
					"status": "REPORTED",
					"cacheGroup": "CDN_in_a_Box_Mid",
					"ip6": "fc01:9400:1000:8::5",
					"port": 80,
					"deliveryServices": [],
					"hostName": "mid",
					"fqdn": "mid.infra.ciab.test",
					"interfaceName": "eth0",
					"type": "MID",
					"hashId": "mid"
				},
				{
					"profile": "ATS_EDGE_TIER_CACHE",
					"ip": "172.16.239.4",
					"status": "REPORTED",
					"cacheGroup": "CDN_in_a_Box_Edge",
					"ip6": "fc01:9400:1000:8::4",
					"port": 80,
					"deliveryServices": [],
					"hostName": "edge",
					"fqdn": "edge.infra.ciab.test",
					"interfaceName": "eth0",
					"type": "EDGE",
					"hashId": "edge"
				}
			],
			"stats": {
				"trafficOpsPath": "/api/1.5/cdns/CDN-in-a-Box/configs/routing",
				"cdnName": "CDN-in-a-Box",
				"trafficOpsVersion": "4.0.0-10449.03d91ae3.el7",
				"trafficOpsUser": "admin",
				"date": 1580152814,
				"trafficOpsHost": "trafficops.infra.ciab.test"
			},
			"cacheGroups": [
				{
					"coordinates": {
						"longitude": -77.036574,
						"latitude": 38.897663
					},
					"name": "CDN_in_a_Box_Edge"
				},
				{
					"coordinates": {
						"longitude": -77.036574,
						"latitude": 38.897663
					},
					"name": "CDN_in_a_Box_Mid"
				}
			],
			"config": {
				"tld.soa.admin": "twelve_monkeys",
				"dnssec.dynamic.response.expiration": "300s",
				"api.cache-control.max-age": 10,
				"neustar.polling.url": "https://trafficops.infra.ciab.test:443/neustar.tar.gz",
				"zonemanager.threadpool.scale": "0.50",
				"coveragezone.polling.interval": 3600000,
				"federationmapping.polling.interval": 60000,
				"steeringmapping.polling.interval": 60000,
				"tld.ttls.DNSKEY": 30,
				"geolocation.polling.interval": 86400000,
				"tld.soa.expire": 604800,
				"federationmapping.polling.url": "https://${toHostname}/api/1.5/federations/all",
				"coveragezone.polling.url": "https://trafficops.infra.ciab.test:443/coverage-zone.json",
				"tld.soa.minimum": 30,
				"geolocation.polling.url": "https://trafficops.infra.ciab.test:443/GeoLite2-City.mmdb.gz",
				"keystore.maintenance.interval": 300,
				"zonemanager.cache.maintenance.interval": 300,
				"domain_name": "mycdn.ciab.test",
				"tld.ttls.AAAA": 3600,
				"tld.soa.refresh": 28800,
				"consistent.dns.routing": "true",
				"tld.ttls.SOA": 86400,
				"neustar.polling.interval": 86400000,
				"tld.ttls.NS": 3600,
				"tld.ttls.DS": 30,
				"certificates.polling.interval": 300000,
				"tld.ttls.A": 3600,
				"tld.soa.retry": 7200
			},
			"trafficMonitors": [
				{
					"profile": "RASCAL-Traffic_Monitor",
					"location": "CDN_in_a_Box_Edge",
					"ip": "172.16.239.11",
					"status": "ONLINE",
					"ip6": "fc01:9400:1000:8::b",
					"port": 80,
					"hostName": "trafficmonitor",
					"fqdn": "trafficmonitor.infra.ciab.test"
				}
			],
			"trafficRouters": [
				{
					"profile": "CCR_CIAB",
					"location": "CDN_in_a_Box_Edge",
					"ip": "172.16.239.12",
					"status": "ONLINE",
					"secureApiPort": 3333,
					"ip6": "fc01:9400:1000:8::c",
					"port": 80,
					"hostName": "trafficrouter",
					"fqdn": "trafficrouter.infra.ciab.test",
					"apiPort": 3333
				}
			]
		}
	}
