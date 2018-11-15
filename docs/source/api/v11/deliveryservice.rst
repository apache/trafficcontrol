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

********************
``deliveryservices``
********************

``GET``
=======
Retrieves all Delivery Services

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| Name        | Required | Description                                                                                                                |
	+=============+==========+============================================================================================================================+
	| cdn         | no       | Show only the Delivery Services belonging to the CDN identified by this integral, unique identifier                        |
	+-------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| id          | no       | Show only the Delivery Service that has this integral, unique identifier                                                   |
	+-------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| logsEnabled | no       | If true, return only Delivery Services with logging enabled, otherwise return only Delivery Services with logging disabled |
	+-------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| profile     | no       | Return only Delivery Services using the profile identified by this integral, unique identifier                             |
	+-------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| tenant      | no       | Show only the Delivery Services belonging to the tenant identified by this integral, unique identifier                     |
	+-------------+----------+----------------------------------------------------------------------------------------------------------------------------+
	| type        | no       | Return only Delivery Services of the Delivery Service type identified by this integral, unique identifier                  |
	+-------------+----------+----------------------------------------------------------------------------------------------------------------------------+

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
:initialDispersion:  The number of caches between which traffic requesting the same object will be randomly split - meaning that if 4 clients all request the same object (one after another), then if this is above 4 there is a possibility that all 4 are cache misses. For most use-cases, this should be 1
:ipv6RoutingEnabled: If ``true``, clients that connect to Traffic Router using IPv6 will be given the IPv6 address of a suitable Edge-tier cache; if ``false`` all addresses will be IPv4, regardless of the client connection\ [2]_
:lastUpdated:        The date and time at which this Delivery Service was last updated, in a ``ctime``-like format
:logsEnabled:        If ``true``, logging is enabled for this Delivery Service, otherwise it is disabled
:longDesc:           A description of the Delivery Service
:longDesc1:          A field used when more detailed information that that provided by ``longDesc`` is desired
:longDesc2:          A field used when even more detailed information that that provided by either ``longDesc`` or ``longDesc1`` is desired
:matchList:          An array of methods used by Traffic Router to determine whether or not a request can be serviced by this Delivery Service

	:pattern:   A regular expression - the use of this pattern is dependent on the ``type`` field (backslashes are escaped)
	:setNumber: An integral, unique identifier for the set of types to which the ``type`` field belongs
	:type:      The type of match performed using ``pattern`` to determine whether or not to use this Delivery Service

		HOST_REGEXP
			Use the Delivery Service if ``pattern`` matches the ``Host:`` HTTP header of an HTTP request\ [2]_
		HEADER_REGEXP
			Use the Delivery Service if ``pattern`` matches an HTTP header (both the name and value) in an HTTP request\ [2]_
		PATH_REGEXP
			Use the Delivery Service if ``pattern`` matches the request path of this Delivery Service's URL
		STEERING_REGEXP
			Use the Delivery Service if ``pattern`` matches the ``xml_id`` of one of this Delivery Service's "Steering" target Delivery Services

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

	.. seealso:: `The Apache Trafficserver documentation for the url_sig plugin <https://docs.trafficserver.apache.org/en/8.0.x/admin-guide/plugins/url_sig.en.html>`_ and `the draft RFC for uri_signing <https://tools.ietf.org/html/draft-ietf-cdni-uri-signing-16>`_ - note, however that the current implementation of uri_signing uses Draft 12 of that RFC document, NOT the latest.


:sslKeyVersion:       This integer indicates the generation of keys in use by the Delyvery Service - if any - and is incremented by the Traffic Portal client whenever new keys are generated

	.. warning:: This number will not be correct if keys are manually replaced using the API, as the key generation API does not increment it!

:tenantId:            The integral, unique identifier of the tenant who owns this Delivery Service
:trRequestHeaders:    If defined, this takes the form of a string of HTTP headers to be included in Traffic Router access logs for requests - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_
:trResponseHeaders:   If defined, this takes the form of a string of HTTP headers to be included in Traffic Router responses - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_
:type:                The name of the routing type of this Delivery Service e.g. "HTTP"
:typeId:              The integral, unique identifier of the routing type of this Delivery Service
:xmlId:               A unique string that describes this Delivery Service - exists for legacy reasons

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: mCLMjvACRKHNGP/OSx4javkOtxxzyiDdQzsV78IamUhVmvyKyKaCeOKRmpsG69w+nhh3OkPZ6e9MMeJpcJSKcA==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 15 Nov 2018 19:04:29 GMT
	Transfer-Encoding: chunked

	{ "response": [
	{
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
		"lastUpdated": "2018-11-14 18:21:17+00",
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
		],
		"deepCachingType": "NEVER",
		"signingAlgorithm": null,
		"tenant": "root"
	}]}

.. [1] Users with the roles "admin" and/or "operations" will be able to see *all* Delivery Services, whereas any other user will only see the Delivery Services their Tenant is allowed to see.
.. [2] This only applies to HTTP Delivery Services
.. [3] See :ref:`multi-site-origin`

``POST``
========
Allows users to create Delivery Service.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Array

Request Structure
-----------------
:active:                   If ``true``, the Delivery Service will immediately become active and
:anonymousBlockingEnabled:     | no       | - true: enable blocking clients with anonymous ips                                                      |
                                            - false: disabled                                                                                       |
:cacheurl:                     | no       | Cache URL rule to apply to this delivery service.                                                       |
:ccrDnsTtl:                    | no       | The TTL of the DNS response for A or AAAA queries requesting the IP address of the tr.host.             |
:cdnId:                        | yes      | cdn id                                                                                                  |
:checkPath:                    | no       | The path portion of the URL to check this deliveryservice for health.                                   |
:deepCachingType:              | no       | When to do Deep Caching for this Delivery Service:                                                      |

                                            - NEVER (default)                                                                                       |
                                            - ALWAYS                                                                                                |
:displayName:                  | yes      | Display name                                                                                            |
:dnsBypassCname:               | no       | Bypass CNAME                                                                                            |
:dnsBypassIp:                  | no       | The IPv4 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the globalMaxMbps traffic on this deliveryservice.                                                          |
:dnsBypassIp6:                 | no       | The IPv6 IP to use for bypass on a DNS deliveryservice - bypass starts when serving more than the globalMaxMbps traffic on this deliveryservice.                                                          |
:dnsBypassTtl:                 | no       | The TTL of the DNS bypass response.                                                                     |
:dscp:                         | yes      | The Differentiated Services Code Point (DSCP) with which to mark downstream (EDGE -> customer) traffic. |
:edgeHeaderRewrite:            | no       | The EDGE header rewrite actions to perform.                                                             |
:fqPacingRate:                 | no       | The maximum rate in bytes per second for each TCP connection in this delivery service. If exceeded, will be rate limited by the Linux kernel. A default value of 0 disables this feature                    |
:geoLimitRedirectURL:          | no       | This is the URL Traffic Router will redirect to when Geo Limit Failure.                                 |
:geoLimit:                     | yes      | - 0: None - no limitations                                                                              |
                                          | - 1: Only route on CZF file hit                                                                         |
                                          | - 2: Only route on CZF hit or when from geo limit countries                                             |

                                                 Note that this does not prevent access to content or makes content secure; it just prevents routing to the content by Traffic Router.                                                               |
:geoLimitCountries:            | no       | The geo limit countries.                                                                                |
:geoProvider:                  | yes      | - 0: Maxmind(default)                                                                                   |
                                            - 1: Neustar                                                                                            |
:globalMaxMbps:                | no       | The maximum global bandwidth allowed on this deliveryservice. If exceeded, the traffic routes to the dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for HTTP deliveryservices.              |
:globalMaxTps:                 | no       | The maximum global transactions per second allowed on this deliveryservice. When this is exceeded traffic will be sent to the dnsByPassIp* for DNS deliveryservices and to the httpBypassFqdn for HTTP deliveryservices
:httpBypassFqdn:               | no       | The HTTP destination to use for bypass on an HTTP deliveryservice - bypass starts when serving more than the globalMaxMbps traffic on this deliveryservice.                                                 |
:infoUrl:                      | no       | Use this to add a URL that points to more information about that deliveryservice.                       |
:initialDispersion:            | yes|no   | Initial dispersion. Required for HTTP* delivery services.                                               |
:ipv6RoutingEnabled:           | yes|no   | false: send IPv4 address of Traffic Router to client on HTTP type del. Required for DNS*, HTTP* and STEERING* delivery services.
:logsEnabled:                  | yes      | - false: No                                                                                             |
                                            - true: Yes                                                                                             |
:longDesc:                     | no       | Description field.                                                                                      |
:longDesc1:                    | no       | Description field 1.                                                                                    |
:longDesc2:                    | no       | Description field 2.                                                                                    |
:maxDnsAnswers:                | no       | The maximum number of IPs to put in a A/AAAA response for a DNS deliveryservice (0 means all available).                                                                                             |
:midHeaderRewrite:             | no       | The MID header rewrite actions to perform.                                                              |
:missLat:                      | yes|no   | The latitude as decimal degrees to use when the client cannot be found in the CZF or the Geo lookup. e.g. 39.7391500 or null. Required for DNS* and HTTP* delivery services.                                 |
:missLong:                     | yes|no   | The longitude as decimal degrees to use when the client cannot be found in the CZF or the Geo lookup. e.g. -104.9847000 or null. Required for DNS* and HTTP* delivery services.                               |
:multiSiteOrigin:              | yes|no   | true if enabled, false if disabled. Required for DNS* and HTTP* delivery services.                      |
:orgServerFqdn:                | yes|no   | The origin server base URL (FQDN when used in this instance, includes the protocol (http:// or https://) for use in retrieving content from the origin server. This field is required if type is DNS* or HTTP*
:originShield:                 | no       | Origin shield                                                                                           |
:profileId:                    | no       | DS profile ID                                                                                           |
:protocol:                     | yes|no   | - 0: serve with http:// at EDGE                                                                         |
                                            - 1: serve with https:// at EDGE                                                                        |
                                            - 2: serve with both http:// and https:// at EDGE                                                       |
                                                                                                                                                    |
                                                 Required for DNS*, HTTP* or *STEERING* delivery services.                                               |
:qstringIgnore:                | yes|no   | - 0: no special query string handling; it is for use in the cache-key and pass up to origin.            |
                                            - 1: ignore query string in cache-key, but pass it up to parent and or origin.                          |
                                            - 2: drop query string at edge, and do not use it in the cache-key.                                     |

                                                   Required for DNS* and HTTP* delivery services.                                                          |
:rangeRequestHandling:         | yes|no   | How to treat range requests (required for DNS* and HTTP* delivery services):                            |
                                            - 0 Do not cache (ranges requested from files taht are already cached due to a non range request will be a HIT)                                                                                               |
                                            - 1 Use the background_fetch plugin.                                                                    |
                                            - 2 Use the cache_range_requests plugin.                                                                |
:regexRemap:                   | no       | Regex Remap rule to apply to this delivery service at the Edge tier.                                    |
:regionalGeoBlocking:          | yes      | Is the Regional Geo Blocking feature enabled.                                                           |
:remapText:                    | no       | Additional raw remap line text.                                                                         |
:routingName:                  | yes      | The routing name of this deliveryservice, e.g. <routingName>.<xmlId>.cdn.com.                           |
:signed:                       | no       | - false: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.          |
                                            - true: token based auth is enabled for this deliveryservice.                                           |
:signingAlgorithm:             | no       | - null: token based auth (see :ref:token-based-auth) is not enabled for this deliveryservice.           |
                                            - "url_sig": URL Sign token based auth is enabled for this deliveryservice.                             |
                                            - "uri_signing": URI Signing token based auth is enabled for this deliveryservice.                      |
:sslKeyVersion:                | no       | SSL key version                                                                                         |
:tenantId:                     | No       | Owning tenant ID                                                                                        |
:trRequestHeaders:             | no       | Traffic router log request headers                                                                      |
:trResponseHeaders:            | no       | Traffic router additional response headers                                                              |
:typeId:                       | yes      | The type of this deliveryservice (one of :ref:to-api-v12-types use_in_table='deliveryservice').         |
:xmlId:                        | yes      | Unique string that describes this deliveryservice.                                                      |


.. code-block:: http
	:caption: Request Example

	POST /api/1.4/deliveryservices HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 761
	Content-Type: application/json

	{
		"active": false,
		"anonymousBlockingEnabled": false,
		"cdnId": 2,
		"cdnName": "CDN-in-a-Box",
		"deepCachingType": "NEVER",
		"displayName": "test",
		"exampleURLs": [
			"http://test.test.mycdn.ciab.test"
		],
		"dscp": 0,
		"geoLimit": 0,
		"geoProvider": 0,
		"initialDispersion": 1,
		"ipv6RoutingEnabled": false,
		"lastUpdated": "2018-11-14 18:21:17+00",
		"logsEnabled": true,
		"longDesc": "A Delivery Service created expressly for API documentation examples",
		"missLat": -1,
		"missLong": -1,
		"multiSiteOrigin": false,
		"orgServerFqdn": "http://origin.infra.ciab.test",
		"protocol": 0,
		"qstringIgnore": 0,
		"rangeRequestHandling": 0,
		"regionalGeoBlocking": false,
		"routingName": "test",
		"signed": false,
		"tenant": "root",
		"tenantId": 1,
		"typeId": 1,
		"xmlId": "test"
	}



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
:initialDispersion:  The number of caches between which traffic requesting the same object will be randomly split - meaning that if 4 clients all request the same object (one after another), then if this is above 4 there is a possibility that all 4 are cache misses. For most use-cases, this should be 1
:ipv6RoutingEnabled: If ``true``, clients that connect to Traffic Router using IPv6 will be given the IPv6 address of a suitable Edge-tier cache; if ``false`` all addresses will be IPv4, regardless of the client connection\ [2]_
:lastUpdated:        The date and time at which this Delivery Service was last updated, in a ``ctime``-like format
:logsEnabled:        If ``true``, logging is enabled for this Delivery Service, otherwise it is disabled
:longDesc:           A description of the Delivery Service
:longDesc1:          A field used when more detailed information that that provided by ``longDesc`` is desired
:longDesc2:          A field used when even more detailed information that that provided by either ``longDesc`` or ``longDesc1`` is desired
:matchList:          An array of methods used by Traffic Router to determine whether or not a request can be serviced by this Delivery Service

	:pattern:   A regular expression - the use of this pattern is dependent on the ``type`` field (backslashes are escaped)
	:setNumber: An integral, unique identifier for the set of types to which the ``type`` field belongs
	:type:      The type of match performed using ``pattern`` to determine whether or not to use this Delivery Service

		HOST_REGEXP
			Use the Delivery Service if ``pattern`` matches the ``Host:`` HTTP header of an HTTP request\ [2]_
		HEADER_REGEXP
			Use the Delivery Service if ``pattern`` matches an HTTP header (both the name and value) in an HTTP request\ [2]_
		PATH_REGEXP
			Use the Delivery Service if ``pattern`` matches the request path of this Delivery Service's URL
		STEERING_REGEXP
			Use the Delivery Service if ``pattern`` matches the ``xml_id`` of one of this Delivery Service's "Steering" target Delivery Services

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

	.. seealso:: `The Apache Trafficserver documentation for the url_sig plugin <https://docs.trafficserver.apache.org/en/8.0.x/admin-guide/plugins/url_sig.en.html>`_ and `the draft RFC for uri_signing <https://tools.ietf.org/html/draft-ietf-cdni-uri-signing-16>`_ - note, however that the current implementation of uri_signing uses Draft 12 of that RFC document, NOT the latest.


:sslKeyVersion:       This integer indicates the generation of keys in use by the Delyvery Service - if any - and is incremented by the Traffic Portal client whenever new keys are generated

	.. warning:: This number will not be correct if keys are manually replaced using the API, as the key generation API does not increment it!

:tenantId:            The integral, unique identifier of the tenant who owns this Delivery Service
:trRequestHeaders:    If defined, this takes the form of a string of HTTP headers to be included in Traffic Router access logs for requests - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_
:trResponseHeaders:   If defined, this takes the form of a string of HTTP headers to be included in Traffic Router responses - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_
:type:                The name of the routing type of this Delivery Service e.g. "HTTP"
:typeId:              The integral, unique identifier of the routing type of this Delivery Service
:xmlId:               A unique string that describes this Delivery Service - exists for legacy reasons

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
