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

.. _to-api-deliveryservices-id:

************************************
``/api/1.x/deliveryservices/{{ID}}``
************************************
.. deprecated:: 1.1
	Use the ``id`` query parameter of :ref:`to-api-deliveryservices` instead

``GET``
=======
Retrieves a specific Delivery Service

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
	| ``logsEnabled`` | no       | If true, return only Delivery Services with logging enabled, otherwise return only Delivery Services with logging disabled |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| ``profile``     | no       | Return only Delivery Services using the profile identified by this integral, unique identifier                             |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| ``tenant``      | no       | Show only the Delivery Services belonging to the tenant identified by this integral, unique identifier                     |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| ``type``        | no       | Return only Delivery Services of the Delivery Service type identified by this integral, unique identifier                  |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+

.. table:: Request Path Parameters

	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| Name            | Required | Description                                                                                                                |
	+=================+==========+============================================================================================================================+
	| ``id``          | yes      | The integral, unique identifier of the Delivery Service to be retrieved                                                    |
	+-----------------+----------+----------------------------------------------------------------------------------------------------------------------------+


Response Structure
------------------
:active:                   ``true`` if the Delivery Service is active, ``false`` otherwise
:anonymousBlockingEnabled: ``true`` if :ref:`Anonymous Blocking <anonymous_blocking-qht>` has been configured for the Delivery Service, ``false`` otherwise
:cacheurl:                 A setting for a deprecated feature of now-unsupported Trafficserver versions
:ccrDnsTtl:                The Time To Live (TTL) of the DNS response for A or AAAA record queries requesting the IP address of the Traffic Router - named "ccrDnsTtl" for legacy reasons
:cdnId:                    The integral, unique identifier of the CDN to which the Delivery Service belongs
:cdnName:                  Name of the CDN to which the Delivery Service belongs
:checkPath:                The path portion of the URL to check this Delivery Service for health - TODO: wat?
:displayName:              The display name of the Delivery Service
:dnsBypassCname:           TODO: wat?
:dnsBypassIp:              The IPv4 IP to use for bypass on a DNS Delivery Service  - bypass starts when the traffic on this Delivery Service exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the Delivery Service
:dnsBypassIp6:             The IPv6 IP to use for bypass on a DNS Delivery Service - bypass starts when the traffic on this Delivery Service exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the Delivery Service
:dnsBypassTtl:             The time for which a DNS bypass of this Delivery Service shall remain active
:dscp:                     The Differentiated Services Code Point (DSCP) with which to mark traffic as it leaves the CDN and reaches clients
:edgeHeaderRewrite:        Rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite Apache Trafficserver plugin
:fqPacingRate:             TODO: wat?
:geoLimitRedirectUrl:      A URL to which clients blocked by :ref:`Regional Geographic Blocking <regionalgeo-qht>` will be re-directed
:geoLimit:                 The setting that determines how content is geographically limited - this is an integer on the interval [0-2] where the values have these meanings:

	0
		None - no limitations
	1
		Only route when the client's IP is found in the Coverage Zone File (CZF)
	2
		Only route when the client's IP is found in the CZF, or when the client can be determined to be from the United States of America

	.. warning:: This does not prevent access to content or make content secure; it merely prevents routing to the content through Traffic Router

:geoLimitCountries:  |
:geoProvider:        |
:globalMaxMbps:      The maximum global bandwidth allowed on this Delivery Service. If exceeded, traffic will be routed to ``dnsBypassIp`` (or ``dnsBypassIp6`` for IPv6 traffic) for DNS Delivery Services and to ``httpBypassFqdn`` for HTTP Delivery Services
:globalMaxTps:       The maximum global transactions per second allowed on this Delivery Service. When this is exceeded traffic will be sent to the dnsByPassIp* for DNS Delivery Services and to the httpBypassFqdn for HTTP Delivery Services
:httpBypassFqdn:     The HTTP destination to use for bypass on an HTTP Delivery Service - bypass starts when the traffic on this Delivery Service exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the Delivery Service
:id:                 An integral, unique identifier for this Delivery Service
:infoUrl:            This is a string which is expected to contain at least one URL pointing to more information about the Delivery Service. Historically, this has been used to link relevant JIRA tickets
:initialDispersion:  |
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

:maxDnsAnswers:        The maximum number of IPs to put in a A/AAAA response for a DNS Delivery Service (0 means all available)
:midHeaderRewrite:     Rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite Apache Trafficserver plugin
:missLat:              The latitude to use when the client cannot be found in the CZF or a geographic IP lookup
:missLong:             The longitude to use when the client cannot be found in the CZF or a geographic IP lookup
:multiSiteOrigin:      ``true`` if the Multi Site Origin feature is enabled for this Delivery Service, ``false`` otherwise\ [3]_
:multiSiteOriginAlgor: TODO: is this ever real? ``true`` the Multi Site Origin feature enabled for this Delivery Service\ [3]_
:originShield:
:orgServerFqdn:        The origin server's Fully Qualified Domain Name (FQDN) - including the protocol (e.g. http:// or https://) - for use in retrieving content from the origin server
:profileDescription:   The description of the Traffic Router Profile with which this Delivery Service is associated
:profileId:            The integral, unique identifier for the Traffic Router profile with which this Delivery Service is associated
:profileName:          The name of the Traffic Router Profile with which this Delivery Service is associated
:protocol:             The protocol which clients will use to communicate with Edge-tier cache servers\ [2]_ - this is an integer on the interval [0-2] where the values have these meanings:

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

:regexRemap:          A regular expression remap rule to apply to this Delivery Service at the Edge tier
:regionalGeoBlocking: ``true`` if Regional Geo Blocking is in use within this Delivery Service, ``false`` otherwise - see :ref:`regionalgeo-qht` for more information
:remapText:           Additional, raw text to add to the remap line for caches
:signed:              ``true`` if token-based authentication is enabled for this Delivery Service, ``false`` otherwise
:signingAlgorithm:    TODO: does this exist for unsigned ds? | string | - null: token based auth (see :ref:token-based-auth) is not enabled for this Delivery Service
                         |        | - "url_sig": URL Sign token based auth is enabled for this Delivery Service
                         |        | - "uri_signing": URI Signing token based auth is enabled for this Delivery Service
:sslKeyVersion:       TODO: wat?
:tenantId:            The integral, unique identifier of the tenant who owns this Delivery Service
:trRequestHeaders:    TODO: do these exist?
:trResponseHeaders:   TODO: do these exist?
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

.. [1] Users with the roles "admin" and/or "operation" will be able to see *all* Delivery Services, whereas any other user will only see the Delivery Services their Tenant is allowed to see.
.. [2] This only applies to HTTP Delivery Services
.. [3] See :ref:`multi-site-origin`
