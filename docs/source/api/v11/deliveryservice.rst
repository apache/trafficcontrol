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

.. _to-api-deliveryservices:

*****************************
``/api/1.x/deliveryservices``
*****************************

``GET``
=======
Retrieves all Delivery Services

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| Name            | Required | Description                                                                                                                |
	+=================+==========+============================================================================================================================+
	| ``cdn``         | no       | Show only the Delivery Services belonging to the CDN identified by this integral, unique identifier                        |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| ``id``          | no       | Show only the Delivery Service that has this integral, unique identifier                                                   |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| ``logsEnabled`` | no       | If true, return only Delivery Services with logging enabled, otherwise return only Delivery Services with logging disabled |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| ``profile``     | no       | Return only Delivery Services using the profile identified by this integral, unique identifier                             |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| ``tenant``      | no       | Show only the Delivery Services belonging to the tenant identified by this integral, unique identifier                     |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| ``type``        | no       | Return only Delivery Services of the Delivery Service type identified by this integral, unique identifier                  |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+

Response Structure
------------------
:active:                   ``true`` if the Delivery Service is active, ``false`` otherwise
:anonymousBlockingEnabled: ``true`` if :ref:`Anonymous Blocking <anonymous_blocking-qht>` has been configured for the Delivery Service, ``false`` otherwise
:cacheurl:                 A setting for a deprecated feature of now-unsupported Trafficserver versions
:ccrDnsTtl:                The Time To Live (TTL) of the DNS response for A or AAAA record queries requesting the IP address of the Traffic Router - named "ccrDnsTtl" for legacy reasons
:cdnId:                    The integral, unique identifier of the CDN to which the Delivery Service belongs
:cdnName:                  Name of the CDN to which the Delivery Service belongs
:checkPath:                The path portion of the URL to check connections to this Delivery Service's origin server
:displayName:              The display name of the Delivery Service
:dnsBypassCname:           Domain name to overflow requests for HTTP Delivery Services - bypass starts when the traffic on this Delivery Service exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the Delivery Service
:dnsBypassIp:              The IPv4 IP to use for bypass on a DNS Delivery Service - bypass starts when the traffic on this Delivery Service exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the Delivery Service
:dnsBypassIp6:             The IPv6 IP to use for bypass on a DNS Delivery Service - bypass starts when the traffic on this Delivery Service exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the Delivery Service
:dnsBypassTtl:             The time for which a DNS bypass of this Delivery Service shall remain active
:dscp:                     The Differentiated Services Code Point (DSCP) with which to mark traffic as it leaves the CDN and reaches clients
:edgeHeaderRewrite:        Rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite Apache Trafficserver plugin
:fqPacingRate:             The Fair-Queuing Pacing Rate in Bytes per second set on the all TCP connection sockets in the Delivery Service (see ``man tc-fc_codel`` for more information) - Linux only
:geoLimit:                 The setting that determines how content is geographically limited - this is an integer on the interval [0-2] where the values have these meanings:
:geoLimitCountries:        A string containing a comma-separated list of country codes (e.g. "US,AU") which are allowed to request content through this Delivery Service
:geoLimitRedirectUrl:      A URL to which clients blocked by :ref:`Regional Geographic Blocking <regionalgeo-qht>` or the ``geoLimit`` settings will be re-directed

	0
		None - no limitations
	1
		Only route when the client's IP is found in the Coverage Zone File (CZF)
	2
		Only route when the client's IP is found in the CZF, or when the client can be determined to be from the United States of America

	.. warning:: This does not prevent access to content or make content secure; it merely prevents routing to the content through Traffic Router

:geoProvider:        An integer that represents the provider of a database for mapping IPs to geographic locations; currently only ``0``  - which represents MaxMind - is supported
:globalMaxMbps:      The maximum global bandwidth allowed on this Delivery Service. If exceeded, traffic will be routed to ``dnsBypassIp`` (or ``dnsBypassIp6`` for IPv6 traffic) for DNS Delivery Services and to ``httpBypassFqdn`` for HTTP Delivery Services
:globalMaxTps:       The maximum global transactions per second allowed on this Delivery Service. When this is exceeded traffic will be sent to the dnsByPassIp* for DNS Delivery Services and to the httpBypassFqdn for HTTP Delivery Services
:httpBypassFqdn:     The HTTP destination to use for bypass on an HTTP Delivery Service - bypass starts when the traffic on this Delivery Service exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the Delivery Service
:id:                 An integral, unique identifier for this Delivery Service
:infoUrl:            This is a string which is expected to contain at least one URL pointing to more information about the Delivery Service. Historically, this has been used to link relevant JIRA tickets
:initialDispersion:  The number of caches between which traffic requesting the same object will be randomly split - meaning that if 4 clients all request the same object (one after another), then if this is above 4 there is a possibility that all 4 are cache misses. For most se-cases, this should be 1
:ipv6RoutingEnabled: If ``true``, clients that connect to Traffic Router using IPv6 will be given the IPv6 address of a suitable Edge-tier cache; if ``false`` all addresses will be IPv4, regardless of the client connection\ [2]_
:lastUpdated:        The date and time at which this Delivery Service was last updated, in a ``ctime``-like format
:logsEnabled:        If ``true``, logging is enabled for this Delivery Service, otherwise it is disabled
:longDesc:           A description of the Delivery Service
:longDesc1:          A field used when more detailed information that that provided by ``longDesc`` is desired
:longDesc2:          A field used when even more detailed information that that provided by either ``longDesc`` or ``longDesc1`` is desired
:matchList:          An array TODO: wat?

	:pattern:   A regular expression
	:setNumber: The set Number of the matchList
	:type:      The type of MatchList

:maxDnsAnswers:      The maximum number of IPs to put in a A/AAAA response for a DNS Delivery Service (0 means all available)
:midHeaderRewrite:   Rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite Apache Trafficserver plugin
:missLat:            The latitude to use when the client cannot be found in the CZF or a geographic IP lookup
:missLong:           The longitude to use when the client cannot be found in the CZF or a geographic IP lookup
:multiSiteOrigin:    ``true`` if the Multi Site Origin feature is enabled for this Delivery Service, ``false`` otherwise\ [3]_
:originShield:       An "origin shield" is a forward proxy that sits between Mid-tier caches and the origin and performs further caching beyond what's offered by a standard CDN. This field is a string of FQDNs to use as origin shields, delimited by ``|``
:orgServerFqdn:      The origin server's Fully Qualified Domain Name (FQDN) - including the protocol (e.g. http:// or https://) - for use in retrieving content from the origin server
:profileDescription: The description of the Traffic Router Profile with which this Delivery Service is associated
:profileId:          The integral, unique identifier for the Traffic Router profile with which this Delivery Service is associated
:profileName:        The name of the Traffic Router Profile with which this Delivery Service is associated
:protocol:           The protocol which clients will use to communicate with Edge-tier cache servers\ [2]_ - this is an integer on the interval [0-2] where the values have these meanings:

	0
		HTTP
	1
		HTTPS
	2
		Both HTTP and HTTPS

:qstringIgnore: Tells caches whether or not to consider URLs with different query parameter strings to be distinct - this is an integer on the interval [0-2] where the values have these meanings:

	0
		URLs with different query parameter strings will be considered distinct for caching purposes, and query strings will be passed upstream to the origin
	1
		URLs with different query parameter strings will be considered identical for caching purposes, and query strings will be passed upstream to the origin
	2
		Query strings are stripped out by Edge-tier caches, and thus are neither taken into consideration for caching purposes, nor passed upstream in requests to the origin

:rangeRequestHandling: Tells caches how to handle range requests\ [2]_ - this is an integer on the interval [0-2] where the values have these meanings:

	0
		Range requests will not be cached, but range requests that request ranges of content already cached will be served from the cache
	1
		Use the `background_fetch plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/background_fetch.en.html>`_ to service the range request while caching the whole object
	2
		Use the `experimental cache_range_requests plugin <https://github.com/apache/trafficserver/tree/master/plugins/experimental/cache_range_requests>`_ to treat unique ranges as unique objects

:regexRemap: A regular expression remap rule to apply to this Delivery Service at the Edge tier

	.. seealso:: `The Apache Trafficserver documentation for the Regex Remap plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_remap.en.html>`_

:regionalGeoBlocking: ``true`` if Regional Geo Blocking is in use within this Delivery Service, ``false`` otherwise - see :ref:`regionalgeo-qht` for more information
:remapText:           Additional, raw text to add to the remap line for caches

	.. seealso:: `The Apache Trafficserver documentation for the Regex Remap plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_remap.en.html>`_

:signed:           ``true`` if token-based authentication is enabled for this Delivery Service, ``false`` otherwise
:signingAlgorithm: Type of URL signing method to sign the URLs, basically comes down to one of two plugins or ``null``:

	``null``
		Token-based authentication is not enabled for this Delivery Service
	url_sig:
		URL Signing token-based authentication is enabled for this Delivery Service
	uri_signing
		URI Signing token-based authentication is enabled for this Delivery Service

	.. seealso:: `The Apache Trafficserver documentation for the url_sig plugin <https://docs.trafficserver.apache.org/en/8.0.x/admin-guide/plugins/url_sig.en.html>`_ and `the draft RFC for uri_signing <https://tools.ietf.org/html/draft-ietf-cdni-uri-signing-16>`_


:sslKeyVersion:       TODO: wat?
:tenantId:            The integral, unique identifier of the tenant who owns this Delivery Service
:trRequestHeaders:    If defined, this takes the form of a string of HTTP headers to be included in Traffic Router access logs for requests - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_
:trResponseHeaders:   If defined, this takes the form of a string of HTTP headers to be included in Traffic Router responses - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_
:type:                The name of the routing type of this Delivery Service e.g. "HTTP"
:typeId:              The integral, unique identifier of the routing type of this Delivery Service
:xmlId:               A unique string that describes this Delivery Service - exists for legacy reasons

.. code-block:: json
	:caption: Response Example

	{ "response": [{
		"active": true,
		"anonymousBlockingEnabled": false,
		"cacheurl": null,
		"ccrDnsTtl": null,
		"cdnId": 2,
		"cdnName": "CDN-in-a-Box",
		"checkPath": null,
		"displayName": "Demo 1",
		"dnsBypassCname": null,
		"dnsBypassIp": null,
		"dnsBypassIp6": null,
		"dnsBypassTtl": null,
		"dscp": 0,
		"edgeHeaderRewrite": null,
		"fqPacingRate": null,
		"geoLimit": 0,
		"geoLimitCountries": null,
		"geoLimitRedirectURL": null,
		"geoProvider": 0,
		"globalMaxMbps": null,
		"globalMaxTps": null,
		"httpBypassFqdn": null,
		"id": 1,
		"infoUrl": null,
		"initialDispersion": 1,
		"ipv6RoutingEnabled": true,
		"lastUpdated": "2018-10-24 16:07:05+00",
		"logsEnabled": true,
		"longDesc": "Apachecon North America 2018",
		"longDesc1": null,
		"longDesc2": null,
		"matchList": [
			{
				"type": "HOST_REGEXP",
				"setNumber": 0,
				"pattern": ".*\\.demo1\\..*"
			}
		],
		"maxDnsAnswers": null,
		"midHeaderRewrite": null,
		"missLat": 42,
		"missLong": -88,
		"multiSiteOrigin": false,
		"originShield": null,
		"orgServerFqdn": "http://origin.infra.ciab.test",
		"profileDescription": null,
		"profileId": null,
		"profileName": null,
		"protocol": 0,
		"qstringIgnore": 0,
		"rangeRequestHandling": 0,
		"regexRemap": null,
		"regionalGeoBlocking": false,
		"remapText": null,
		"routingName": "video",
		"signed": false,
		"sslKeyVersion": null,
		"tenantId": 1,
		"type": "HTTP",
		"typeId": 1,
		"xmlId": "demo1",
		"exampleURLs": [
			"http://video.demo1.mycdn.ciab.test"
		]
	}]}

.. [1] Users with the roles "admin" and/or "operations" will be able to see *all* Delivery Services, whereas any other user will only see the Delivery Services their Tenant is allowed to see.
.. [2] This only applies to HTTP Delivery Services
.. [3] See :ref:`multi-site-origin`

	Retrieves the health of all locations (cache groups) for a delivery service.

	Authentication Required: Yes

	Role(s) Required: None

	**Response Properties**

	+------------------+--------+-------------------------------------------------+
	|    Parameter     |  Type  |                   Description                   |
	+==================+========+=================================================+
	| ``totalOnline``  | int    | Total number of online caches across all CDNs.  |
	+------------------+--------+-------------------------------------------------+
	| ``totalOffline`` | int    | Total number of offline caches across all CDNs. |
	+------------------+--------+-------------------------------------------------+
	| ``cachegroups``  | array  | A collection of cache groups.                   |
	+------------------+--------+-------------------------------------------------+
	| ``>online``      | int    | The number of online caches for the cache group |
	+------------------+--------+-------------------------------------------------+
	| ``>offline``     | int    | The number of offline caches for the cache      |
	|                  |        | group.                                          |
	+------------------+--------+-------------------------------------------------+
	| ``>name``        | string | Cache group name.                               |
	+------------------+--------+-------------------------------------------------+

	**Response Example** ::

		{
		 "response": {
				"totalOnline": 148,
				"totalOffline": 0,
				"cachegroups": [
					 {
							"online": 8,
							"offline": 0,
							"name": "us-co-denver"
					 },
					 {
							"online": 7,
							"offline": 0,
							"name": "us-de-newcastle"
					 }
				]
		 }
		}


|

	Retrieves the capacity percentages of a delivery service.

	Authentication Required: Yes

	Role(s) Required: None

	**Request Route Parameters**

	+-----------------+----------+---------------------------------------------------+
	| Name            | Required | Description                                       |
	+=================+==========+===================================================+
	|id               | yes      | delivery service id.                              |
	+-----------------+----------+---------------------------------------------------+

	**Response Properties**

	+------------------------+--------+---------------------------------------------------+
	|       Parameter        |  Type  |                    Description                    |
	+========================+========+===================================================+
	| ``availablePercent``   | number | The percentage of server capacity assigned to     |
	|                        |        | the delivery service that is available.           |
	+------------------------+--------+---------------------------------------------------+
	| ``unavailablePercent`` | number | The percentage of server capacity assigned to the |
	|                        |        | delivery service that is unavailable.             |
	+------------------------+--------+---------------------------------------------------+
	| ``utilizedPercent``    | number | The percentage of server capacity assigned to the |
	|                        |        | delivery service being used.                      |
	+------------------------+--------+---------------------------------------------------+
	| ``maintenancePercent`` | number | The percentage of server capacity assigned to the |
	|                        |        | delivery service that is down for maintenance.    |
	+------------------------+--------+---------------------------------------------------+

	**Response Example** ::

		{
		 "response": {
				"availablePercent": 89.0939840205533,
				"unavailablePercent": 0,
				"utilizedPercent": 10.9060020300395,
				"maintenancePercent": 0.0000139494071146245
		 },
		}


|

	Retrieves the routing method percentages of a delivery service.

	Authentication Required: Yes

	Role(s) Required: None

	**Request Route Parameters**

	+-----------------+----------+---------------------------------------------------+
	| Name            | Required | Description                                       |
	+=================+==========+===================================================+
	|id               | yes      | delivery service id.                              |
	+-----------------+----------+---------------------------------------------------+

	**Response Properties**

	+-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
	|    Parameter    |  Type  |                                                         Description                                                         |
	+=================+========+=============================================================================================================================+
	| ``staticRoute`` | number | The percentage of Traffic Router responses for this deliveryservice satisfied with pre-configured DNS entries.              |
	+-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
	| ``miss``        | number | The percentage of Traffic Router responses for this deliveryservice that were a miss (no location available for client IP). |
	+-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
	| ``geo``         | number | The percentage of Traffic Router responses for this deliveryservice satisfied using 3rd party geo-IP mapping.               |
	+-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
	| ``err``         | number | The percentage of Traffic Router requests for this deliveryservice resulting in an error.                                   |
	+-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
	| ``cz``          | number | The percentage of Traffic Router requests for this deliveryservice satisfied by a CZF hit.                                  |
	+-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+
	| ``dsr``         | number | The percentage of Traffic Router requests for this deliveryservice satisfied by sending the                                 |
	|                 |        | client to the overflow CDN.                                                                                                 |
	+-----------------+--------+-----------------------------------------------------------------------------------------------------------------------------+

	**Response Example** ::

		{
		 "response": {
				"staticRoute": 0,
				"miss": 0,
				"geo": 37.8855391018869,
				"err": 0,
				"cz": 62.1144608981131,
				"dsr": 0
		 },
		}

|

.. _to-api-v11-ds-metrics:

Metrics
+++++++

**GET /api/1.1/deliveryservices/:id/server_types/:type/metric_types/start_date/:start/end_date/:end.json**

	Retrieves detailed and summary metrics for MIDs or EDGEs for a delivery service.

	Authentication Required: Yes

	Role(s) Required: None

	**Request Route Parameters**

	+------------------+----------+-----------------------------------------------------------------------------+
	|       Name       | Required |                                 Description                                 |
	+==================+==========+=============================================================================+
	| ``id``           | yes      | The delivery service id.                                                    |
	+------------------+----------+-----------------------------------------------------------------------------+
	| ``server_types`` | yes      | EDGE or MID.                                                                |
	+------------------+----------+-----------------------------------------------------------------------------+
	| ``metric_types`` | yes      | One of the following: "kbps", "tps", "tps_2xx", "tps_3xx", "tps_4xx",       |
	|                  |          | "tps_5xx".                                                                  |
	+------------------+----------+-----------------------------------------------------------------------------+
	| ``start_date``   | yes      | UNIX time                                                                   |
	+------------------+----------+-----------------------------------------------------------------------------+
	| ``end_date``     | yes      | UNIX time                                                                   |
	+------------------+----------+-----------------------------------------------------------------------------+

	**Request Query Parameters**

	+------------------+----------+-----------------------------------------------------------------------------+
	|       Name       | Required |                                 Description                                 |
	+==================+==========+=============================================================================+
	| ``stats``        | no       | Flag used to return only summary metrics                                    |
	+------------------+----------+-----------------------------------------------------------------------------+

	**Response Properties**

	+----------------------+--------+-------------+
	|      Parameter       |  Type  | Description |
	+======================+========+=============+
	| ``stats``            | hash   |             |
	+----------------------+--------+-------------+
	| ``>>count``          | int    |             |
	+----------------------+--------+-------------+
	| ``>>98thPercentile`` | number |             |
	+----------------------+--------+-------------+
	| ``>>min``            | number |             |
	+----------------------+--------+-------------+
	| ``>>max``            | number |             |
	+----------------------+--------+-------------+
	| ``>>5thPercentile``  | number |             |
	+----------------------+--------+-------------+
	| ``>>95thPercentile`` | number |             |
	+----------------------+--------+-------------+
	| ``>>median``         | number |             |
	+----------------------+--------+-------------+
	| ``>>mean``           | number |             |
	+----------------------+--------+-------------+
	| ``>>stddev``         | number |             |
	+----------------------+--------+-------------+
	| ``>>sum``            | number |             |
	+----------------------+--------+-------------+
	| ``data``             | array  |             |
	+----------------------+--------+-------------+
	| ``>>item``           | array  |             |
	+----------------------+--------+-------------+
	| ``>>time``           | number |             |
	+----------------------+--------+-------------+
	| ``>>value``          | number |             |
	+----------------------+--------+-------------+
	| ``label``            | string |             |
	+----------------------+--------+-------------+

	**Response Example** ::

		{
		 "response": [
				{
					 "stats": {
							"count": 988,
							"98thPercentile": 16589105.55958,
							"min": 3185442.975,
							"max": 17124754.257,
							"5thPercentile": 3901253.95445,
							"95thPercentile": 16013210.034,
							"median": 8816895.576,
							"mean": 8995846.31741194,
							"stddev": 3941169.83683573,
							"sum": 333296106.060112
					 },
					 "data": [
							[
								 1414303200000,
								 12923518.466
							],
							[
								 1414303500000,
								 12625139.65
							]
					 ],
					 "label": "MID Kbps"
				}
		 ],
		}


.. _to-api-v11-ds-server:

Server
++++++

**GET /api/1.1/deliveryserviceserver.json**

	Authentication Required: Yes

	Role(s) Required: Yes

	**Request Query Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	| ``page``  | no       | The page number for use in pagination. |
	+-----------+----------+----------------------------------------+
	| ``limit`` | no       | For use in limiting the result set.    |
	+-----------+----------+----------------------------------------+

	**Response Properties**

	+----------------------+--------+------------------------------------------------+
	| Parameter            | Type   | Description                                    |
	+======================+========+================================================+
	|``lastUpdated``       | array  |                                                |
	+----------------------+--------+------------------------------------------------+
	|``server``            | string |                                                |
	+----------------------+--------+------------------------------------------------+
	|``deliveryService``   | string |                                                |
	+----------------------+--------+------------------------------------------------+

	**Response Example** ::

		{
		 "page": 2,
		 "orderby": "deliveryservice",
		 "response": [
				{
					 "lastUpdated": "2014-09-26 17:53:43",
					 "server": "20",
					 "deliveryService": "1"
				},
				{
					 "lastUpdated": "2014-09-26 17:53:44",
					 "server": "21",
					 "deliveryService": "1"
				},
		 ],
		 "limit": 2
		}

|

.. _to-api-v11-ds-sslkeys:

SSL Keys
+++++++++

**GET /api/1.1/deliveryservices/xmlId/:xmlid/sslkeys.json**

	Authentication Required: Yes

	Role(s) Required: None

	**Request Route Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	| ``xmlId`` | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

	**Request Query Parameters**

	+-------------+----------+--------------------------------+
	|     Name    | Required |          Description           |
	+=============+==========+================================+
	| ``version`` | no       | The version number to retrieve |
	+-------------+----------+--------------------------------+

	**Response Properties**

	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	|    Parameter     |  Type  |                                                               Description                                                               |
	+==================+========+=========================================================================================================================================+
	| ``crt``          | string | base64 encoded crt file for delivery service                                                                                            |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``csr``          | string | base64 encoded csr file for delivery service                                                                                            |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``key``          | string | base64 encoded private key file for delivery service                                                                                    |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``businessUnit`` | string | The business unit entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``city``         | string | The city entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response          |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``organization`` | string | The organization entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response  |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``hostname``     | string | The hostname entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response      |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``country``      | string | The country entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response       |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``state``        | string | The state entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response         |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``version``      | string | The version of the certificate record in Riak                                                                                           |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+

	**Response Example** ::

		{
			"response": {
				"certificate": {
					"crt": "crt",
					"key": "key",
					"csr": "csr"
				},
				"businessUnit": "CDN_Eng",
				"city": "Denver",
				"organization": "KableTown",
				"hostname": "foober.com",
				"country": "US",
				"state": "Colorado",
				"version": "1"
			}
		}

|

**GET /api/1.1/deliveryservices/hostname/:hostname/sslkeys.json**

	Authentication Required: Yes

	Role Required: Admin

	**Request Route Parameters**

	+--------------+----------+---------------------------------------------------+
	|     Name     | Required |                    Description                    |
	+==============+==========+===================================================+
	| ``hostname`` | yes      | pristine hostname of the desired delivery service |
	+--------------+----------+---------------------------------------------------+

	**Request Query Parameters**

	+-------------+----------+--------------------------------+
	|     Name    | Required |          Description           |
	+=============+==========+================================+
	| ``version`` | no       | The version number to retrieve |
	+-------------+----------+--------------------------------+

	**Response Properties**

	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	|    Parameter     |  Type  |                                                               Description                                                               |
	+==================+========+=========================================================================================================================================+
	| ``crt``          | string | base64 encoded crt file for delivery service                                                                                            |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``csr``          | string | base64 encoded csr file for delivery service                                                                                            |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``key``          | string | base64 encoded private key file for delivery service                                                                                    |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``businessUnit`` | string | The business unit entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``city``         | string | The city entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response          |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``organization`` | string | The organization entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response  |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``hostname``     | string | The hostname entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response      |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``country``      | string | The country entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response       |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``state``        | string | The state entered by the user when generating certs.  Field is optional and if not provided by the user will not be in response         |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| ``version``      | string | The version of the certificate record in Riak                                                                                           |
	+------------------+--------+-----------------------------------------------------------------------------------------------------------------------------------------+

	**Response Example** ::

		{
			"response": {
				"certificate": {
					"crt": "crt",
					"key": "key",
					"csr": "csr"
				},
				"businessUnit": "CDN_Eng",
				"city": "Denver",
				"organization": "KableTown",
				"hostname": "foober.com",
				"country": "US",
				"state": "Colorado",
				"version": "1"
			}
		}

|

**GET /api/1.1/deliveryservices/xmlId/:xmlid/sslkeys/delete.json**

	Authentication Required: Yes

	Role Required: Operations

	**Request Route Parameters**

	+-----------+----------+----------------------------------------+
	|    Name   | Required |              Description               |
	+===========+==========+========================================+
	| ``xmlId`` | yes      | xml_id of the desired delivery service |
	+-----------+----------+----------------------------------------+

	**Request Query Parameters**

	+-------------+----------+--------------------------------+
	|     Name    | Required |          Description           |
	+=============+==========+================================+
	| ``version`` | no       | The version number to retrieve |
	+-------------+----------+--------------------------------+

	**Response Properties**

	+--------------+--------+------------------+
	|  Parameter   |  Type  |   Description    |
	+==============+========+==================+
	| ``response`` | string | success response |
	+--------------+--------+------------------+

	**Response Example** ::

		{
			"response": "Successfully deleted ssl keys for <xml_id>"
		}

|

**POST /api/1.1/deliveryservices/sslkeys/generate**

	Generates SSL crt, csr, and private key for a delivery service

	Authentication Required: Yes

	Role(s) Required: Operations

	**Request Properties**

	+--------------+---------+-------------------------------------------------+
	|  Parameter   |   Type  |                   Description                   |
	+==============+=========+=================================================+
	| ``key``      | string  | xml_id of the delivery service                  |
	+--------------+---------+-------------------------------------------------+
	| ``version``  | string  | version of the keys being generated             |
	+--------------+---------+-------------------------------------------------+
	| ``hostname`` | string  | the *pristine hostname* of the delivery service |
	+--------------+---------+-------------------------------------------------+
	| ``country``  | string  |                                                 |
	+--------------+---------+-------------------------------------------------+
	| ``state``    | string  |                                                 |
	+--------------+---------+-------------------------------------------------+
	| ``city``     | string  |                                                 |
	+--------------+---------+-------------------------------------------------+
	| ``org``      | string  |                                                 |
	+--------------+---------+-------------------------------------------------+
	| ``unit``     | boolean |                                                 |
	+--------------+---------+-------------------------------------------------+

	**Request Example** ::

		{
			"key": "ds-01",
			"businessUnit": "CDN Engineering",
			"version": "3",
			"hostname": "tr.ds-01.ott.kabletown.com",
			"certificate": {
				"key": "some_key",
				"csr": "some_csr",
				"crt": "some_crt"
			},
			"country": "US",
			"organization": "Kabletown",
			"city": "Denver",
			"state": "Colorado"
		}

|

	**Response Properties**

	+--------------+--------+-----------------+
	|  Parameter   |  Type  |   Description   |
	+==============+========+=================+
	| ``response`` | string | response string |
	+--------------+--------+-----------------+
	| ``version``  | string | API version     |
	+--------------+--------+-----------------+

	**Response Example** ::

		{
			"response": "Successfully created ssl keys for ds-01"
		}

|

**POST /api/1.1/deliveryservices/sslkeys/add**

	Allows user to add SSL crt, csr, and private key for a delivery service.

	Authentication Required: Yes

	Role(s) Required: Operations

	**Request Properties**

	+-------------+--------+-------------------------------------+
	|  Parameter  |  Type  |             Description             |
	+=============+========+=====================================+
	| ``key``     | string | xml_id of the delivery service      |
	+-------------+--------+-------------------------------------+
	| ``version`` | string | version of the keys being generated |
	+-------------+--------+-------------------------------------+
	| ``csr``     | string |                                     |
	+-------------+--------+-------------------------------------+
	| ``crt``     | string |                                     |
	+-------------+--------+-------------------------------------+
	| ``key``     | string |                                     |
	+-------------+--------+-------------------------------------+

	**Request Example** ::

		{
			"key": "ds-01",
			"version": "1",
			"certificate": {
				"key": "some_key",
				"csr": "some_csr",
				"crt": "some_crt"
			}
		}

|

	**Response Properties**

	+--------------+--------+-----------------+
	|  Parameter   |  Type  |   Description   |
	+==============+========+=================+
	| ``response`` | string | response string |
	+--------------+--------+-----------------+
	| ``version``  | string | API version     |
	+--------------+--------+-----------------+

	**Response Example** ::

		{
			"response": "Successfully added ssl keys for ds-01"
		}


|

**POST /api/1.1/deliveryservices/request**

	Allows a user to send delivery service request details to a specified email address.

	Authentication Required: Yes

	Role(s) Required: None

	**Request Properties**

	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	|  Parameter                             |  Type  | Required |           Description                                                                       |
	+========================================+========+==========+=============================================================================================+
	| ``emailTo``                            | string | yes      | The email to which the delivery service request will be sent.                               |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``details``                            | hash   | yes      | Parameters for the delivery service request.                                                |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>customer``                          | string | yes      | Name of the customer to associated with the delivery service.                               |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>deliveryProtocol``                  | string | yes      | Eg. http or http/https                                                                      |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>routingType``                       | string | yes      | Eg. DNS or HTTP Redirect                                                                    |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>serviceDesc``                       | string | yes      | A description of the delivery service.                                                      |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>peakBPSEstimate``                   | string | yes      | Used to manage cache efficiency and plan for capacity.                                      |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>peakTPSEstimate``                   | string | yes      | Used to manage cache efficiency and plan for capacity.                                      |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>maxLibrarySizeEstimate``            | string | yes      | Used to manage cache efficiency and plan for capacity.                                      |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>originURL``                         | string | yes      | The URL path to the origin server.                                                          |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>hasOriginDynamicRemap``             | bool   | yes      | This is a feature which allows services to use multiple origin URLs for the same service.   |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>originTestFile``                    | string | yes      | A URL path to a test file available on the origin server.                                   |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>hasOriginACLWhitelist``             | bool   | yes      | Is access to your origin restricted using an access control list (ACL or whitelist) of Ips? |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>originHeaders``                     | string | no       | Header values that must be passed to requests to your origin.                               |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>otherOriginSecurity``               | string | no       | Other origin security measures that need to be considered for access.                       |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>queryStringHandling``               | string | yes      | How to handle query strings that come with the request.                                     |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>rangeRequestHandling``              | string | yes      | How to handle range requests.                                                               |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>hasSignedURLs``                     | bool   | yes      | Are Urls signed?                                                                            |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>hasNegativeCachingCustomization``   | bool   | yes      | Any customization required for negative caching?                                            |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>negativeCachingCustomizationNote``  | string | yes      | Negative caching customization instructions.                                                |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>serviceAliases``                    | array  | no       | Service aliases which will be used for this service.                                        |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>rateLimitingGBPS``                  | int    | no       | Rate Limiting - Bandwidth (Gigabits per second)                                             |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>rateLimitingTPS``                   | int    | no       | Rate Limiting - Transactions/Second                                                         |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>overflowService``                   | string | no       | An overflow point (URL or IP address) used if rate limits are met.                          |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>headerRewriteEdge``                 | string | no       | Headers can be added or altered at each layer of the CDN.                                   |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>headerRewriteMid``                  | string | no       | Headers can be added or altered at each layer of the CDN.                                   |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>headerRewriteRedirectRouter``       | string | no       | Headers can be added or altered at each layer of the CDN.                                   |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+
	| ``>notes``                             | string | no       | Additional instructions to provide the delivery service provisioning team.                  |
	+----------------------------------------+--------+----------+---------------------------------------------------------------------------------------------+

	**Request Example** ::

		{
			 "emailTo": "foo@bar.com",
			 "details": {
					"customer": "XYZ Corporation",
					"contentType": "video-on-demand",
					"deliveryProtocol": "http",
					"routingType": "dns",
					"serviceDesc": "service description goes here",
					"peakBPSEstimate": "less-than-5-Gbps",
					"peakTPSEstimate": "less-than-1000-TPS",
					"maxLibrarySizeEstimate": "less-than-200-GB",
					"originURL": "http://myorigin.com",
					"hasOriginDynamicRemap": false,
					"originTestFile": "http://myorigin.com/crossdomain.xml",
					"hasOriginACLWhitelist": true,
					"originHeaders": "",
					"otherOriginSecurity": "",
					"queryStringHandling": "ignore-in-cache-key-and-pass-up",
					"rangeRequestHandling": "range-requests-not-used",
					"hasSignedURLs": true,
					"hasNegativeCachingCustomization": true,
					"negativeCachingCustomizationNote": "negative caching instructions",
					"serviceAliases": [
						 "http://alias1.com",
						 "http://alias2.com"
					],
					"rateLimitingGBPS": 50,
					"rateLimitingTPS": 5000,
					"overflowService": "http://overflowcdn.com",
					"headerRewriteEdge": "",
					"headerRewriteMid": "",
					"headerRewriteRedirectRouter": "",
					"notes": ""
			 }
		}

|

	**Response Properties**

	+-------------+--------+----------------------------------+
	|  Parameter  |  Type  |           Description            |
	+=============+========+==================================+
	| ``alerts``  | array  | A collection of alert messages.  |
	+-------------+--------+----------------------------------+
	| ``>level``  | string | Success, info, warning or error. |
	+-------------+--------+----------------------------------+
	| ``>text``   | string | Alert message.                   |
	+-------------+--------+----------------------------------+
	| ``version`` | string |                                  |
	+-------------+--------+----------------------------------+

	**Response Example** ::

		{
			"alerts": [
						{
								"level": "success",
								"text": "Delivery Service request sent to foo@bar.com."
						}
				]
		}

|
