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

.. _to-api-cdns-name-configs-routing:

*********************************
``cdns/{{name}}/configs/routing``
*********************************
.. caution:: This API route is currently broken, see :issue:`2941` for more information.

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
