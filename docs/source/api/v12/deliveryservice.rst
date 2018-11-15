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

.. _to-api-deliveryservices-id-servers:

***********************************
``deliveryservices/{{ID}}/servers``
***********************************

``GET``
=======
Retrieves properties of Edge-Tier servers assigned to a Delivery Service.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [1]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------------------+
	| Name | Description                                                                                 |
	+======+=============================================================================================+
	| ID   | The integral, unique identifier of the Delivery service for which servers will be displayed |
	+------+---------------------------------------------------------------------------------------------+

Response Structure
------------------
:cachegroup:     The name of the Cache Group to which the server belongs
:cachegroupId:   An integral, unique identifier for the Cache Group to which the server belongs
:cdnId:          An integral, unique identifier the CDN to which the server belongs
:cdnName:        The name of the CDN to which the server belongs
:domainName:     The domain name part of the Fully Qualified Domain Name (FQDN) of the server
:guid:           Optionally represents an identifier used to uniquely identify the server
:hostName:       The (short) hostname of the server
:httpsPort:      The port on which the server listens for incoming HTTPS requests - 443 in most cases
:id:             An integral, unique identifier for the server
:iloIpAddress:   The IPv4 address of the lights-out-management port\ [2]_
:iloIpGateway:   The IPv4 gateway address of the lights-out-management port\ [2]_
:iloIpNetmask:   The IPv4 subnet mask of the lights-out-management port\ [2]_
:iloPassword:    The password of the of the lights-out-management user - displays as ``******`` unless the requesting user has the 'admin' role)\ [2]_
:iloUsername:    The user name for lights-out-management\ [2]_
:interfaceMtu:   The Maximum Transmission Unit (MTU) to configure for ``interfaceName``

	.. seealso:: `The Wikipedia article on Maximum Transmission Unit <https://en.wikipedia.org/wiki/Maximum_transmission_unit>`_

:interfaceName:  The network interface name used by the server
:ip6Address:     The IPv6 address and subnet mask of the server - applicable for the interface ``interfaceName``
:ip6Gateway:     The IPv6 gateway address of the server - applicable for the interface ``interfaceName``
:ipAddress:      The IPv4 address of the server- applicable for the interface ``interfaceName``
:ipGateway:      The IPv4 gateway of the server- applicable for the interface ``interfaceName``
:ipNetmask:      The IPv4 subnet mask of the server- applicable for the interface ``interfaceName``
:lastUpdated:    The time and date at which this server was last updated, in an ISO-like format
:mgmtIpAddress:  The IPv4 address of the server's management port
:mgmtIpGateway:  The IPv4 gateway of the server's management port
:mgmtIpNetmask:  The IPv4 subnet mask of the server's management port
:offlineReason:  A user-entered reason why the server is in ADMIN_DOWN or OFFLINE status (will be empty if not offline)
:physLocation:   The name of the physical location at which the server resides
:physLocationId: An integral, unique identifier for the physical location at which the server resides
:profile:        The name of the profile assigned to this server
:profileDesc:    A description of the profile assigned to this server
:profileId:      An integral, unique identifier for the profile assigned to this server
:rack:           A string indicating "rack" location
:routerHostName: The human-readable name of the router
:routerPortName: The human-readable name of the router port
:status:         The Status of the server

	.. seealso:: :ref:`health-proto`

:statusId:       An integral, unique identifier for the status of the server

	.. seealso:: :ref:`health-proto`

:tcpPort:        The default port on which the main application listens for incoming TCP connections - 80 in most cases
:type:           The name of the type of this server
:typeId:         An integral, unique identifier for the type of this server
:updPending:     ``true`` if the server has updates pending, ``false`` otherwise

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: MaIvaO8OSjysr4bCkuXFEMf3o6mOqga1aM4IHN/tcP2aa1iXEmA5IrHB7DaqNX/2vGHLXvN+01FEAR/lRNqr1w==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 21:28:23 GMT
	Content-Length: 891

	{ "response": [
		{
			"cachegroup": "CDN_in_a_Box_Edge",
			"cachegroupId": 7,
			"cdnId": 2,
			"cdnName": "CDN-in-a-Box",
			"domainName": "infra.ciab.test",
			"guid": null,
			"hostName": "edge",
			"httpsPort": 443,
			"id": 10,
			"iloIpAddress": "",
			"iloIpGateway": "",
			"iloIpNetmask": "",
			"iloPassword": "",
			"iloUsername": "",
			"interfaceMtu": 1500,
			"interfaceName": "eth0",
			"ip6Address": "fc01:9400:1000:8::100",
			"ip6Gateway": "fc01:9400:1000:8::1",
			"ipAddress": "172.16.239.100",
			"ipGateway": "172.16.239.1",
			"ipNetmask": "255.255.255.0",
			"lastUpdated": "2018-11-14 21:08:44+00",
			"mgmtIpAddress": "",
			"mgmtIpGateway": "",
			"mgmtIpNetmask": "",
			"offlineReason": "",
			"physLocation": "Apachecon North America 2018",
			"physLocationId": 1,
			"profile": "ATS_EDGE_TIER_CACHE",
			"profileDesc": "Edge Cache - Apache Traffic Server",
			"profileId": 9,
			"rack": "",
			"routerHostName": "",
			"routerPortName": "",
			"status": "REPORTED",
			"statusId": 3,
			"tcpPort": 80,
			"type": "EDGE",
			"typeId": 11,
			"updPending": false
		}
	]}


.. [1] Users with the roles "admin" and/or "operations" will be able to the see servers associated with *any* Delivery Services, whereas any other user will only be able to see the servers associated with Delivery Services their Tenant is allowed to see.
.. [2] See `the Wikipedia article on Out-of-Band Management <https://en.wikipedia.org/wiki/Out-of-band_management>`_ for more information.

**PUT /api/1.2/deliveryservices/{:id}**

	Allows user to edit a delivery service.

	Authentication Required: Yes

	Role(s) Required:  admin or oper

	**Request Route Parameters**

	+-----------------+----------+---------------------------------------------------+
	| Name            | Required | Description                                       |
	+=================+==========+===================================================+
	|id               | yes      | delivery service id.                              |
	+-----------------+----------+---------------------------------------------------+

	**Request Properties**

	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| Parameter                | Required | Description                                                                                             |
	+==========================+==========+=========================================================================================================+
	| active                   | yes      | true if active, false if inactive.                                                                      |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| anonymousBlockingEnabled | no       | - true: enable blocking clients with anonymous ips                                                      |
	|                          |          | - false: disabled                                                                                       |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| cacheurl                 | no       | Cache URL rule to apply to this delivery service.                                                       |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| ccrDnsTtl                | no       | The TTL of the DNS response for A or AAAA queries requesting the IP address of the tr.host.             |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| cdnId                    | yes      | cdn id                                                                                                  |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| checkPath                | no       | The path portion of the URL to check this deliveryservice for health.                                   |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| deepCachingType          | no       | When to do Deep Caching for this Delivery Service:                                                      |
	|                          |          |                                                                                                         |
	|                          |          | - NEVER (default)                                                                                       |
	|                          |          | - ALWAYS                                                                                                |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| displayName              | yes      | Display name                                                                                            |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| dnsBypassCname           | no       | Bypass CNAME                                                                                            |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| dnsBypassIp              | no       | The IPv4 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the       |
	|                          |          | globalMaxMbps traffic on this deliveryservice.                                                          |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| dnsBypassIp6             | no       | The IPv6 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the       |
	|                          |          | globalMaxMbps traffic on this deliveryservice.                                                          |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| dnsBypassTtl             | no       | The TTL of the DNS bypass response.                                                                     |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| dscp                     | yes      | The Differentiated Services Code Point (DSCP) with which to mark downstream (EDGE -> customer) traffic. |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| edgeHeaderRewrite        | no       | The EDGE header rewrite actions to perform.                                                             |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| fqPacingRate             | no       | The maximum rate in bytes per second for each TCP connection in this delivery service. If exceeded,     |
	|                          |          | will be rate limited by the Linux kernel. A default value of 0 disables this feature                    |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| geoLimitRedirectURL      | no       | This is the URL Traffic Router will redirect to when Geo Limit Failure.                                 |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| geoLimit                 | yes      | - 0: None - no limitations                                                                              |
	|                          |          | - 1: Only route on CZF file hit                                                                         |
	|                          |          | - 2: Only route on CZF hit or when from geo limit countries                                             |
	|                          |          |                                                                                                         |
	|                          |          | Note that this does not prevent access to content or makes content secure; it just prevents             |
	|                          |          | routing to the content by Traffic Router.                                                               |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| geoLimitCountries        | no       | The geo limit countries.                                                                                |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| geoProvider              | yes      | - 0: Maxmind(default)                                                                                   |
	|                          |          | - 1: Neustar                                                                                            |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| globalMaxMbps            | no       | The maximum global bandwidth allowed on this deliveryservice. If exceeded, the traffic routes to the    |
	|                          |          | dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for HTTP deliveryservices.              |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| globalMaxTps             | no       | The maximum global transactions per second allowed on this deliveryservice. When this is exceeded       |
	|                          |          | traffic will be sent to the dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for         |
	|                          |          | HTTP deliveryservices                                                                                   |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| httpBypassFqdn           | no       | The HTTP destination to use for bypass on an HTTP deliveryservice - bypass starts when serving more     |
	|                          |          | than the globalMaxMbps traffic on this deliveryservice.                                                 |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| infoUrl                  | no       | Use this to add a URL that points to more information about that deliveryservice.                       |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| initialDispersion        | yes|no   | Initial dispersion. Required for HTTP* delivery services.                                               |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| ipv6RoutingEnabled       | yes|no   | false: send IPv4 address of Traffic Router to client on HTTP type del.                                  |
	|                          |          | Required for DNS*, HTTP* and STEERING* delivery services.                                               |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| logsEnabled              | yes      | - false: No                                                                                             |
	|                          |          | - true: Yes                                                                                             |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| longDesc                 | no       | Description field.                                                                                      |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| longDesc1                | no       | Description field 1.                                                                                    |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| longDesc2                | no       | Description field 2.                                                                                    |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| maxDnsAnswers            | no       | The maximum number of IPs to put in a A/AAAA response for a DNS deliveryservice (0 means all            |
	|                          |          | available).                                                                                             |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| midHeaderRewrite         | no       | The MID header rewrite actions to perform.                                                              |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| missLat                  | yes|no   | The latitude as decimal degrees to use when the client cannot be found in the CZF or the Geo lookup.    |
	|                          |          | e.g. 39.7391500 or null. Required for DNS* and HTTP* delivery services.                                 |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| missLong                 | yes|no   | The longitude as decimal degrees to use when the client cannot be found in the CZF or the Geo lookup.   |
	|                          |          | e.g. -104.9847000 or null. Required for DNS* and HTTP* delivery services.                               |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| multiSiteOrigin          | yes|no   | true if enabled, false if disabled. Required for DNS* and HTTP* delivery services.                      |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| orgServerFqdn            | yes|no   | The origin server base URL (FQDN when used in this instance, includes the                               |
	|                          |          | protocol (http:// or https://) for use in retrieving content from the origin server. This field is      |
	|                          |          | required if type is DNS* or HTTP*.                                                                      |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| originShield             | no       | Origin shield                                                                                           |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| profileId                | no       | DS profile ID                                                                                           |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| protocol                 | yes|no   | - 0: serve with http:// at EDGE                                                                         |
	|                          |          | - 1: serve with https:// at EDGE                                                                        |
	|                          |          | - 2: serve with both http:// and https:// at EDGE                                                       |
	|                          |          |                                                                                                         |
	|                          |          | Required for DNS*, HTTP* or *STEERING* delivery services.                                               |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| qstringIgnore            | yes|no   | - 0: no special query string handling; it is for use in the cache-key and pass up to origin.            |
	|                          |          | - 1: ignore query string in cache-key, but pass it up to parent and or origin.                          |
	|                          |          | - 2: drop query string at edge, and do not use it in the cache-key.                                     |
	|                          |          |                                                                                                         |
	|                          |          | Required for DNS* and HTTP* delivery services.                                                          |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| rangeRequestHandling     | yes|no   | How to treat range requests (required for DNS* and HTTP* delivery services):                            |
	|                          |          | - 0 Do not cache (ranges requested from files taht are already cached due to a non range request will   |
	|                          |          | be a HIT)                                                                                               |
	|                          |          | - 1 Use the background_fetch plugin.                                                                    |
	|                          |          | - 2 Use the cache_range_requests plugin.                                                                |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| regexRemap               | no       | Regex Remap rule to apply to this delivery service at the Edge tier.                                    |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| regionalGeoBlocking      | yes      | Is the Regional Geo Blocking feature enabled.                                                           |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| remapText                | no       | Additional raw remap line text.                                                                         |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| routingName              | yes      | The routing name of this deliveryservice, e.g. <routingName>.<xmlId>.cdn.com.                           |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| signed                   | no       | - false: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.          |
	|                          |          | - true: token based auth is enabled for this deliveryservice.                                           |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| signingAlgorithm         | no       | - null: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.           |
	|                          |          | - "url_sig": URL Sign token based auth is enabled for this deliveryservice.                             |
	|                          |          | - "uri_signing": URI Signing token based auth is enabled for this deliveryservice.                      |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| sslKeyVersion            | no       | SSL key version                                                                                         |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| tenantId                 | No       | Owning tenant ID                                                                                        |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| trRequestHeaders         | no       | Traffic router log request headers                                                                      |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| trResponseHeaders        | no       | Traffic router additional response headers                                                              |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| typeId                   | yes      | The type of this deliveryservice (one of :ref:to-api-v12-types use_in_table='deliveryservice').         |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| xmlId                    | yes      | Unique string that describes this deliveryservice. This value cannot be changed on update.              |
	+--------------------------+----------+---------------------------------------------------------------------------------------------------------+


	**Request Example** ::

		{
				"xmlId": "my_ds_1",
				"displayName": "my_ds_displayname_1",
				"tenantId": 1,
				"protocol": 1,
				"orgServerFqdn": "http://10.75.168.91",
				"cdnId": 2,
				"typeId": 42,
				"active": false,
				"dscp": 10,
				"geoLimit": 0,
				"geoProvider": 0,
				"initialDispersion": 1,
				"ipv6RoutingEnabled": false,
				"logsEnabled": false,
				"multiSiteOrigin": false,
				"missLat": 39.7391500,
				"missLong": -104.9847000,
				"qstringIgnore": 0,
				"rangeRequestHandling": 0,
				"regionalGeoBlocking": false,
				"anonymousBlockingEnabled": false,
				"signed": false,
				"signingAlgorithm": null
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
						"infoUrl": "",
						"initialDispersion": "1",
						"ipv6RoutingEnabled": true,
						"lastUpdated": "2016-01-26 08:49:35",
						"logsEnabled": false,
						"longDesc": "",
						"longDesc1": "",
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

|

**PUT /api/1.2/deliveryservices/{:id}/safe**

	Allows a user to edit limited fields of an assigned delivery service.

	Authentication Required: Yes

	Role(s) Required:  users with the delivery service assigned or ops and above

	**Request Route Parameters**

	+-----------------+----------+---------------------------------------------------+
	| Name            | Required | Description                                       |
	+=================+==========+===================================================+
	|id               | yes      | delivery service id.                              |
	+-----------------+----------+---------------------------------------------------+

	**Request Properties**

	+------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| Parameter              | Required | Description                                                                                             |
	+========================+==========+=========================================================================================================+
	| displayName            | no       | Display name                                                                                            |
	+------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| infoUrl                | no       | Use this to add a URL that points to more information about that deliveryservice.                       |
	+------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| longDesc               | no       | Description field.                                                                                      |
	+------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| longDesc1              | no       | Description field 1.                                                                                    |
	+------------------------+----------+---------------------------------------------------------------------------------------------------------+
	| all other fields       | n/a      | All other fields will be silently ignored                                                               |
	+------------------------+----------+---------------------------------------------------------------------------------------------------------+


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

|

**DELETE /api/1.2/deliveryservices/{:id}**

	Allows user to delete a delivery service.

	Authentication Required: Yes

	Role(s) Required:  admin or oper

	**Request Route Parameters**

	+-----------------+----------+---------------------------------------------------+
	| Name            | Required | Description                                       |
	+=================+==========+===================================================+
	| id              | yes      | delivery service id.                              |
	+-----------------+----------+---------------------------------------------------+

	 **Response Example** ::

		{
					 "alerts": [
										 {
														 "level": "success",
														 "text": "Delivery service was deleted."
										 }
						 ],
		}

|

**POST /api/1.2/deliveryservices/:xml_id/servers**

	Assign caches to a delivery service.

	Authentication Required: Yes

	Role(s) Required:  admin or oper

	**Request Route Parameters**

	+--------+----------+-----------------------------------+
	| Name   | Required | Description                       |
	+========+==========+===================================+
	| xml_id | yes      | the xml_id of the deliveryservice |
	+--------+----------+-----------------------------------+

	**Request Properties**

	+--------------+----------+-------------------------------------------------------------------------------------------------------------+
	| Parameter    | Required | Description                                                                                                 |
	+==============+==========+=============================================================================================================+
	| serverNames  | yes      | array of hostname of cache servers to assign to this deliveryservice, for example: [ "server1", "server2" ] |
	+--------------+----------+-------------------------------------------------------------------------------------------------------------+

	**Request Example** ::

		{
				"serverNames": [
						"tc1_ats1"
				]
		}

	**Response Properties**

	+--------------+--------+-------------------------------------------------------------------------------------------------------------+
	| Parameter    | Type   | Description                                                                                                 |
	+==============+========+=============================================================================================================+
	| xml_id       | string | Unique string that describes this delivery service.                                                         |
	+--------------+--------+-------------------------------------------------------------------------------------------------------------+
	| serverNames  | string | array of hostname of cache servers to assign to this deliveryservice, for example: [ "server1", "server2" ] |
	+--------------+--------+-------------------------------------------------------------------------------------------------------------+


	 **Response Example** ::

		{
				"response":{
						"serverNames":[
								"tc1_ats1"
						],
						"xmlId":"my_ds_1"
				}
		}

|

URI Signing Keys
++++++++++++++++

**DELETE /api/1.2/deliveryservices/:xml_id/urisignkeys**

	Deletes URISigning objects for a delivery service.

	Authentication Required: Yes

	Role(s) Required: admin

	**Request Route Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	| xml_id    | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

**GET /api/1.2/deliveryservices/:xml_id/urisignkeys**

	Retrieves one or more URISigning objects for a delivery service.

	Authentication Required: Yes

	Role(s) Required: admin

	**Request Route Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	| xml_id    | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

	**Response Properties**

	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	|    Parameter        |  Type  |                                                               Description                                                               |
	+=====================+========+=========================================================================================================================================+
	| ``Issuer``          | string | a string describing the issuer of the URI signing object. Multiple URISigning objects may be returned in a response, see example        |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``renewal_kid``     | string | a string naming the jwt key used for renewals.                                                                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``keys``            | string | json array of jwt symmetric keys                                                             .                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``alg``             | string | this parameter repeats for each jwt key in the array and specifies the jwa encryption algorithm to use with this key, RFC 7518.         |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kid``             | string | this parameter repeats for each jwt key in the array and specifies the unique id for the key as defined in RFC 7516.                    |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kty``             | string | this parameter repeats for each jwt key in the array and specifies the key type as defined in RFC 7516.                                 |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``k``               | string | this parameter repeats for each jwt key in the array and specifies the base64 encoded symmetric key see RFC 7516.                       |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+

	**Response Example** ::

		{
			"Kabletown URI Authority": {
				"renewal_kid": "Second Key",
				"keys": [
					{
						"alg": "HS256",
						"kid": "First Key",
						"kty": "oct",
						"k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
					},
					{
						"alg": "HS256",
						"kid": "Second Key",
						"kty": "oct",
						"k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
					}
				]
			}
		}


**POST /api/1.2/deliveryservices/:xml_id/urisignkeys**

	Assigns URISigning objects to a delivery service.

	Authentication Required: Yes

	Role(s) Required: admin

	**Request Route Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	|   xml_id  | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

	**Request Properties**

	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	|    Parameter        |  Type  |                                                               Description                                                               |
	+=====================+========+=========================================================================================================================================+
	| ``Issuer``          | string | a string describing the issuer of the URI signing object. Multiple URISigning objects may be returned in a response, see example        |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``renewal_kid``     | string | a string naming the jwt key used for renewals.                                                                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``keys``            | string | json array of jwt symmetric keys                                                             .                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``alg``             | string | this parameter repeats for each jwt key in the array and specifies the jwa encryption algorithm to use with this key, RFC 7518.         |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kid``             | string | this parameter repeats for each jwt key in the array and specifies the unique id for the key as defined in RFC 7516.                    |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kty``             | string | this parameter repeats for each jwt key in the array and specifies the key type as defined in RFC 7516.                                 |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``k``               | string | this parameter repeats for each jwt key in the array and specifies the base64 encoded symmetric key see RFC 7516.                       |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+

	**Request Example** ::

		{
			"Kabletown URI Authority": {
				"renewal_kid": "Second Key",
				"keys": [
					{
						"alg": "HS256",
						"kid": "First Key",
						"kty": "oct",
						"k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
					},
					{
						"alg": "HS256",
						"kid": "Second Key",
						"kty": "oct",
						"k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
					}
				]
			}
		}

**PUT /api/1.2/deliveryservices/:xml_id/urisignkeys**

	updates URISigning objects on a delivery service.

	Authentication Required: Yes

	Role(s) Required: admin

	**Request Route Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	|  xml_id   | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

	**Request Properties**

	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	|    Parameter        |  Type  |                                                               Description                                                               |
	+=====================+========+=========================================================================================================================================+
	| ``Issuer``          | string | a string describing the issuer of the URI signing object. Multiple URISigning objects may be returned in a response, see example        |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``renewal_kid``     | string | a string naming the jwt key used for renewals.                                                                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``keys``            | string | json array of jwt symmetric keys                                                             .                                          |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``alg``             | string | this parameter repeats for each jwt key in the array and specifies the jwa encryption algorithm to use with this key, RFC 7518.         |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kid``             | string | this parameter repeats for each jwt key in the array and specifies the unique id for the key as defined in RFC 7516.                    |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``kty``             | string | this parameter repeats for each jwt key in the array and specifies the key type as defined in RFC 7516.                                 |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``k``               | string | this parameter repeats for each jwt key in the array and specifies the base64 encoded symmetric key see RFC 7516.                       |
	+---------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+

	**Request Example** ::

		{
			"Kabletown URI Authority": {
				"renewal_kid": "Second Key",
				"keys": [
					{
						"alg": "HS256",
						"kid": "First Key",
						"kty": "oct",
						"k": "Kh_RkUMj-fzbD37qBnDf_3e_RvQ3RP9PaSmVEpE24AM"
					},
					{
						"alg": "HS256",
						"kid": "Second Key",
						"kty": "oct",
						"k": "fZBpDBNbk2GqhwoB_DGBAsBxqQZVix04rIoLJ7p_RlE"
					}
				]
			}
		}

|

