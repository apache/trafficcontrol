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

.. _to-api-deliveryservices-id-safe:

********************************
``deliveryservices/{{ID}}/safe``
********************************

``PUT``
=======
Allows a user to edit metadata fields of a Delivery Service.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [1]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------+
	| Name |                      Description                                       |
	+======+========================================================================+
	|  ID  | The integral, unique identifier of the Delivery Service being modified |
	+------+------------------------------------------------------------------------+

:displayName: The human-friendly name for this Delivery Service
:infoUrl:     A string which is expected to contain at least one URL pointing to more information about the Delivery Service. Historically, this has been used to link relevant JIRA tickets
:longDesc:    A description of the Delivery Service
:longDesc1:   A field used when more detailed information that that provided by ``longDesc`` is desired
:longDesc2:   A field used when even more detailed information that that provided by either ``longDesc`` or ``longDesc1`` is desired

.. note:: All of these fields are optional; this ``PUT`` behaves more like a ``PATCH``

	**Request Example** ::

		{
				"displayName": "My Cool Delivery Service",
				"infoUrl": "www.info.com",
				"longDesc": "some info about the service",
				"longDesc1": "the customer label"
		}


	**Response Properties**

	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| Parameter                    | Type   | Description                                                                                                                          |
	+==============================+========+======================================================================================================================================+
	| ``active``                   | bool   | true if active, false if inactive.                                                                                                   |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``anonymousBlockingEnabled`` | bool   | - true: enable blocking clients with anonymous ips                                                                                   |
	|                              |        | - false: disabled                                                                                                                    |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``cacheurl``                 | string | Cache URL rule to apply to this delivery service.                                                                                    |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``ccrDnsTtl``                | int    | The TTL of the DNS response for A or AAAA queries requesting the IP address of the tr. host.                                         |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``cdnId``                    | int    | Id of the CDN to which the delivery service belongs to.                                                                              |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``cdnName``                  | string | Name of the CDN to which the delivery service belongs to.                                                                            |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``checkPath``                | string | The path portion of the URL to check this deliveryservice for health.                                                                |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``deepCachingType``          | string | When to do Deep Caching for this Delivery Service:                                                                                   |
	|                              |        |                                                                                                                                      |
	|                              |        | - NEVER (default)                                                                                                                    |
	|                              |        | - ALWAYS                                                                                                                             |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``displayName``              | string | The display name of the delivery service.                                                                                            |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``dnsBypassCname``           | string |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``dnsBypassIp``              | string | The IPv4 IP to use for bypass on a DNS deliveryservice  - bypass starts when serving more than the                                   |
	|                              |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``dnsBypassIp6``             | string | The IPv6 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the                                    |
	|                              |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``dnsBypassTtl``             | int    | The TTL of the DNS bypass response.                                                                                                  |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``dscp``                     | int    | The Differentiated Services Code Point (DSCP) with which to mark downstream (EDGE ->  customer) traffic.                             |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``edgeHeaderRewrite``        | string | The EDGE header rewrite actions to perform.                                                                                          |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``exampleURLs``              | array  | Entry points into the CDN for this deliveryservice.                                                                                  |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``fqPacingRate``             |  int   | The maximum rate in bytes per second for each TCP connection in this delivery service. If exceeded,                                  |
	|                              |        | will be rate limited by the Linux kernel. A default value of 0 disables this feature                                                 |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``geoLimitRedirectUrl``      | string |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``geoLimit``                 | int    | - 0: None - no limitations                                                                                                           |
	|                              |        | - 1: Only route on CZF file hit                                                                                                      |
	|                              |        | - 2: Only route on CZF hit or when from USA                                                                                          |
	|                              |        |                                                                                                                                      |
	|                              |        | Note that this does not prevent access to content or makes content secure; it just prevents                                          |
	|                              |        | routing to the content by Traffic Router.                                                                                            |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``geoLimitCountries``        | string |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``geoProvider``              | int    |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``globalMaxMbps``            | int    | The maximum global bandwidth allowed on this deliveryservice. If exceeded, the traffic routes to the                                 |
	|                              |        | dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for HTTP deliveryservices.                                           |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``globalMaxTps``             | int    | The maximum global transactions per second allowed on this deliveryservice. When this is exceeded                                    |
	|                              |        | traffic will be sent to the dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for                                      |
	|                              |        | HTTP deliveryservices                                                                                                                |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``httpBypassFqdn``           | string | The HTTP destination to use for bypass on an HTTP deliveryservice - bypass starts when serving more than the                         |
	|                              |        | globalMaxMbps traffic on this deliveryservice.                                                                                       |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``id``                       | int    | The deliveryservice id (database row number).                                                                                        |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``infoUrl``                  | string | Use this to add a URL that points to more information about that deliveryservice.                                                    |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``initialDispersion``        | int    |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``ipv6RoutingEnabled``       | bool   | false: send IPv4 address of Traffic Router to client on HTTP type del.                                                               |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``lastUpdated``              | string |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``logsEnabled``              | bool   |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``longDesc``                 | string | Description field.                                                                                                                   |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``longDesc1``                | string | Description field 1.                                                                                                                 |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``longDesc2``                | string | Description field 2.                                                                                                                 |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``matchList``                | array  | Array of matchList hashes.                                                                                                           |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``>>type``                   | string | The type of MatchList (one of :ref:to-api-v11-types use_in_table='regex').                                                           |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``>>setNumber``              | string | The set Number of the matchList.                                                                                                     |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``>>pattern``                | string | The regexp for the matchList.                                                                                                        |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``maxDnsAnswers``            | int    | The maximum number of IPs to put in a A/AAAA response for a DNS deliveryservice (0 means all                                         |
	|                              |        | available).                                                                                                                          |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``midHeaderRewrite``         | string | The MID header rewrite actions to perform.                                                                                           |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``missLat``                  | float  | The latitude as decimal degrees to use when the client cannot be found in the CZF or the Geo lookup.                                 |
	|                              |        | - e.g. 39.7391500 or null                                                                                                            |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``missLong``                 | float  | The longitude as decimal degrees to use when the client cannot be found in the CZF or the Geo lookup.                                |
	|                              |        | - e.g. -104.9847000 or null                                                                                                          |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``multiSiteOrigin``          | bool   | Is the Multi Site Origin feature enabled for this delivery service (0=false, 1=true). See :ref:`multi-site-origin`                   |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``orgServerFqdn``            | string | The origin server base URL (FQDN when used in this instance, includes the                                                            |
	|                              |        | protocol (http:// or https://) for use in retrieving content from the origin server.                                                 |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``originShield``             | string |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``profileDescription``       | string | The description of the Traffic Router Profile with which this deliveryservice is associated.                                         |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``profileId``                | int    | The id of the Traffic Router Profile with which this deliveryservice is associated.                                                  |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``profileName``              | string | The name of the Traffic Router Profile with which this deliveryservice is associated.                                                |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``protocol``                 | int    | - 0: serve with http:// at EDGE                                                                                                      |
	|                              |        | - 1: serve with https:// at EDGE                                                                                                     |
	|                              |        | - 2: serve with both http:// and https:// at EDGE                                                                                    |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``qstringIgnore``            | int    | - 0: no special query string handling; it is for use in the cache-key and pass up to origin.                                         |
	|                              |        | - 1: ignore query string in cache-key, but pass it up to parent and or origin.                                                       |
	|                              |        | - 2: drop query string at edge, and do not use it in the cache-key.                                                                  |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``rangeRequestHandling``     | int    | How to treat range requests:                                                                                                         |
	|                              |        | - 0 Do not cache (ranges requested from files taht are already cached due to a non range request will be a HIT)                      |
	|                              |        | - 1 Use the `background_fetch <https://docs.trafficserver.apache.org/en/latest/reference/plugins/background_fetch.en.html>`_ plugin. |
	|                              |        | - 2 Use the cache_range_requests plugin.                                                                                             |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``regexRemap``               | string | Regex Remap rule to apply to this delivery service at the Edge tier.                                                                 |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``regionalGeoBlocking``      | bool   | Regex Remap rule to apply to this delivery service at the Edge tier.                                                                 |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``remapText``                | string | Additional raw remap line text.                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``routingName``              | string | The routing name of this deliveryservice, e.g. <routingName>.<xmlId>.cdn.com.                                                        |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``signed``                   | bool   | - false: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.                                       |
	|                              |        | - true: token based auth is enabled for this deliveryservice.                                                                        |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``signingAlgorithm``         | string | - null: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.                                        |
	|                              |        | - "url_sig": URL Sign token based auth is enabled for this deliveryservice.                                                          |
	|                              |        | - "uri_signing": URI Signing token based auth is enabled for this deliveryservice.                                                   |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``sslKeyVersion``            | int    |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``trRequestHeaders``         | string |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``trResponseHeaders``        | string |                                                                                                                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``typeId``                   | int    | The type of this deliveryservice (one of :ref:to-api-v11-types use_in_table='deliveryservice').                                      |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+
	| ``xmlId``                    | string | Unique string that describes this deliveryservice.                                                                                   |
	+------------------------------+--------+--------------------------------------------------------------------------------------------------------------------------------------+

	**Response Example** ::

		{
			"response": [
				{
						"active": true,
						"anonymousBlockingEnabled": false,
						"cacheurl": null,
						"ccrDnsTtl": "3600",
						"cdnId": "2",
						"cdnName": "over-the-top",
						"checkPath": "",
						"deepCachingType": "NEVER",
						"displayName": "My Cool Delivery Service",
						"dnsBypassCname": "",
						"dnsBypassIp": "",
						"dnsBypassIp6": "",
						"dnsBypassTtl": "30",
						"dscp": "40",
						"edgeHeaderRewrite": null,
						"exampleURLs": [
								"http://foo.foo-ds.foo.bar.net"
						],
						"geoLimit": "0",
						"geoLimitCountries": null,
						"geoLimitRedirectURL": null,
						"geoProvider": "0",
						"globalMaxMbps": null,
						"globalMaxTps": "0",
			"fqPacingRate": "0",
						"httpBypassFqdn": "",
						"id": "442",
						"infoUrl": "www.info.com",
						"initialDispersion": "1",
						"ipv6RoutingEnabled": true,
						"lastUpdated": "2016-01-26 08:49:35",
						"logsEnabled": false,
						"longDesc": "some info about the service",
						"longDesc1": "the customer label",
						"longDesc2": "",
						"matchList": [
								{
										"pattern": ".*\\.foo-ds\\..*",
										"setNumber": "0",
										"type": "HOST_REGEXP"
								}
						],
						"maxDnsAnswers": "0",
						"midHeaderRewrite": null,
						"missLat": "39.7391500",
						"missLong": "-104.9847000",
						"multiSiteOrigin": false,
						"orgServerFqdn": "http://baz.boo.net",
						"originShield": null,
						"profileDescription": "Content Router for over-the-top",
						"profileId": "5",
						"profileName": "ROUTER_TOP",
						"protocol": "0",
						"qstringIgnore": "1",
						"rangeRequestHandling": "0",
						"regexRemap": null,
						"regionalGeoBlocking": false,
						"remapText": null,
						"routingName": "foo",
						"signed": false,
						"signingAlgorithm": null,
						"sslKeyVersion": "0",
						"tenantId": 1,
						"trRequestHeaders": null,
						"trResponseHeaders": "Access-Control-Allow-Origin: *",
						"type": "HTTP",
						"typeId": "8",
						"xmlId": "foo-ds"
				}
			]
		}

.. [1] Users with the "admin" or "operations" roles will be able to edit *any* Delivery Service, whereas other users will only be able to edit Delivery Services that their tenant has permissions to edit.
