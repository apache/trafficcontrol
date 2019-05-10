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
Allows a user to edit metadata fields of a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [1]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------------------+
	| Name |                      Description                                               |
	+======+================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service` being modified |
	+------+--------------------------------------------------------------------------------+

:displayName: The human-friendly name for this :term:`Delivery Service`
:infoUrl:     A string which is expected to contain at least one URL pointing to more information about the :term:`Delivery Service`. Historically, this has been used to link relevant JIRA tickets
:longDesc:    A description of the :term:`Delivery Service`
:longDesc1:   A field used when more detailed information that that provided by ``longDesc`` is desired
:longDesc2:   A field used when even more detailed information that that provided by either ``longDesc`` or ``longDesc1`` is desired

.. note:: All of these fields are optional; this ``PUT`` behaves more like a ``PATCH``

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/deliveryservices/1/safe HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 165
	Content-Type: application/x-www-form-urlencoded

	{
		"displayName": "demo",
		"infoUrl": "www.info.com",
		"longDesc": "A :term:`Delivery Service` created for the CDN-in-a-Box project",
		"longDesc1": null,
		"longDesc2": null
	}


Response Structure
------------------
.. versionchanged:: 1.3
	Removed ``fqPacingRate`` field, added fields: ``deepCachingType``, ``signingAlgorithm``, and ``tenant``.

:active:                   ``true`` if the :term:`Delivery Service` is active, ``false`` otherwise
:anonymousBlockingEnabled: ``true`` if :ref:`Anonymous Blocking <anonymous_blocking-qht>` has been configured for the :term:`Delivery Service`, ``false`` otherwise
:cacheurl:                 A setting for a deprecated feature of now-unsupported Trafficserver versions

	.. deprecated:: ATCv3.0
		This field has been deprecated in Traffic Control 3.x and is subject to removal in Traffic Control 4.x or later

:ccrDnsTtl:                The Time To Live (TTL) of the DNS response for A or AAAA record queries requesting the IP address of the Traffic Router - named "ccrDnsTtl" for legacy reasons
:cdnId:                    The integral, unique identifier of the CDN to which the :term:`Delivery Service` belongs
:cdnName:                  Name of the CDN to which the :term:`Delivery Service` belongs
:checkPath:                The path portion of the URL to check connections to this :term:`Delivery Service`'s origin server
:consistentHashRegex:      If defined, this is a regex used for the Pattern-Based Consistent Hashing feature. It is only applicable for HTTP and Steering Delivery Services

	.. versionadded:: 1.5

:deepCachingType:          A string that describes when "Deep Caching" will be used by this :term:`Delivery Service` - one of:

	ALWAYS
		"Deep Caching" will always be used with this :term:`Delivery Service`
	NEVER
		"Deep Caching" will never be used with this :term:`Delivery Service`

	.. versionadded:: 1.3

:displayName:              The display name of the :term:`Delivery Service`
:dnsBypassCname:           Domain name to overflow requests for HTTP :term:`Delivery Service`\ s - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`\ [4]_
:dnsBypassIp:              The IPv4 IP to use for bypass on a DNS :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`\ [4]_
:dnsBypassIp6:             The IPv6 IP to use for bypass on a DNS :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`\ [4]_
:dnsBypassTtl:             The time for which a DNS bypass of this :term:`Delivery Service`\ shall remain active\ [4]_
:dscp:                     The Differentiated Services Code Point (DSCP) with which to mark traffic as it leaves the CDN and reaches clients
:edgeHeaderRewrite:        Rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite Apache Trafficserver plugin
:fqPacingRate:             The Fair-Queuing Pacing Rate in Bytes per second set on the all TCP connection sockets in the :term:`Delivery Service` (see ``man tc-fc_codel`` for more information) - Linux only

	.. deprecated:: 1.3
		This field is only present/available in API versions 1.2 and lower - it has been removed in API version 1.3

:geoLimit:                 The setting that determines how content is geographically limited - this is an integer on the interval [0-2] where the values have these meanings:
:geoLimitCountries:        A string containing a comma-separated list of country codes (e.g. "US,AU") which are allowed to request content through this :term:`Delivery Service`
:geoLimitRedirectUrl:      A URL to which clients blocked by :ref:`Regional Geographic Blocking <regionalgeo-qht>` or the ``geoLimit`` settings will be re-directed

	0
		None - no limitations
	1
		Only route when the client's IP is found in the Coverage Zone File (CZF)
	2
		Only route when the client's IP is found in the CZF, or when the client can be determined to be from the United States of America

	.. warning:: This does not prevent access to content or make content secure; it merely prevents routing to the content through Traffic Router

:geoProvider:        An integer that represents the provider of a database for mapping IPs to geographic locations; currently only ``0``  - which represents MaxMind - is supported
:globalMaxMbps:      The maximum global bandwidth allowed on this :term:`Delivery Service`. If exceeded, traffic will be routed to ``dnsBypassIp`` (or ``dnsBypassIp6`` for IPv6 traffic) for DNS :term:`Delivery Service`\ s and to ``httpBypassFqdn`` for HTTP :term:`Delivery Service`\ s
:globalMaxTps:       The maximum global transactions per second allowed on this :term:`Delivery Service`. When this is exceeded traffic will be sent to the dnsByPassIp* for DNS :term:`Delivery Service`\ s and to the httpBypassFqdn for HTTP :term:`Delivery Service`\ s
:httpBypassFqdn:     The HTTP destination to use for bypass on an HTTP :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`
:id:                 An integral, unique identifier for this :term:`Delivery Service`
:infoUrl:            This is a string which is expected to contain at least one URL pointing to more information about the :term:`Delivery Service`. Historically, this has been used to link relevant JIRA tickets
:initialDispersion:  The number of caches between which traffic requesting the same object will be randomly split - meaning that if 4 clients all request the same object (one after another), then if this is above 4 there is a possibility that all 4 are cache misses. For most use-cases, this should be 1
:ipv6RoutingEnabled: If ``true``, clients that connect to Traffic Router using IPv6 will be given the IPv6 address of a suitable Edge-tier cache; if ``false`` all addresses will be IPv4, regardless of the client connection\ [2]_
:lastUpdated:        The date and time at which this :term:`Delivery Service` was last updated, in a ``ctime``-like format
:logsEnabled:        If ``true``, logging is enabled for this :term:`Delivery Service`, otherwise it is disabled
:longDesc:           A description of the :term:`Delivery Service`
:longDesc1:          A field used when more detailed information that that provided by ``longDesc`` is desired
:longDesc2:          A field used when even more detailed information that that provided by either ``longDesc`` or ``longDesc1`` is desired
:matchList:          An array of methods used by Traffic Router to determine whether or not a request can be serviced by this :term:`Delivery Service`

	:pattern:   A regular expression - the use of this pattern is dependent on the ``type`` field (backslashes are escaped)
	:setNumber: An integral, unique identifier for the set of types to which the ``type`` field belongs
	:type:      The type of match performed using ``pattern`` to determine whether or not to use this :term:`Delivery Service`

		HOST_REGEXP
			Use the :term:`Delivery Service` if ``pattern`` matches the ``Host:`` HTTP header of an HTTP request\ [2]_
		HEADER_REGEXP
			Use the :term:`Delivery Service` if ``pattern`` matches an HTTP header (both the name and value) in an HTTP request\ [2]_
		PATH_REGEXP
			Use the :term:`Delivery Service` if ``pattern`` matches the request path of this :term:`Delivery Service`'s URL
		STEERING_REGEXP
			Use the :term:`Delivery Service` if ``pattern`` matches the ``xml_id`` of one of this :term:`Delivery Service`'s "Steering" target :term:`Delivery Service`\ s

:maxDnsAnswers:      The maximum number of IPs to put in a A/AAAA response for a DNS :term:`Delivery Service` (0 means all available)\ [4]_
:midHeaderRewrite:   Rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite Apache Trafficserver plugin
:missLat:            The latitude to use when the client cannot be found in the CZF or a geographic IP lookup
:missLong:           The longitude to use when the client cannot be found in the CZF or a geographic IP lookup
:multiSiteOrigin:    ``true`` if the Multi Site Origin feature is enabled for this :term:`Delivery Service`, ``false`` otherwise\ [3]_
:originShield:       An "origin shield" is a forward proxy that sits between Mid-tier caches and the origin and performs further caching beyond what's offered by a standard CDN. This field is a string of FQDNs to use as origin shields, delimited by ``|``
:orgServerFqdn:      The origin server's Fully Qualified Domain Name (FQDN) - including the protocol (e.g. http:// or https://) - for use in retrieving content from the origin server
:profileDescription: The description of the Traffic Router Profile with which this :term:`Delivery Service` is associated
:profileId:          The integral, unique identifier for the Traffic Router profile with which this :term:`Delivery Service` is associated
:profileName:        The name of the Traffic Router Profile with which this :term:`Delivery Service` is associated
:protocol:           The protocol which clients will use to communicate with Edge-tier :term:`cache server`\ s\ [2]_ - this is an integer on the interval [0-2] where the values have these meanings:

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

:rangeRequestHandling: Tells caches how to handle range requests\ [5]_ - this is an integer on the interval [0-2] where the values have these meanings:

	0
		Range requests will not be cached, but range requests that request ranges of content already cached will be served from the cache
	1
		Use the `background_fetch plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/background_fetch.en.html>`_ to service the range request while caching the whole object
	2
		Use the `experimental cache_range_requests plugin <https://github.com/apache/trafficserver/tree/master/plugins/experimental/cache_range_requests>`_ to treat unique ranges as unique objects

:regexRemap: A regular expression remap rule to apply to this :term:`Delivery Service` at the Edge tier

	.. seealso:: `The Apache Trafficserver documentation for the Regex Remap plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_remap.en.html>`_

:regionalGeoBlocking: ``true`` if Regional Geo Blocking is in use within this :term:`Delivery Service`, ``false`` otherwise - see :ref:`regionalgeo-qht` for more information
:remapText:           Additional, raw text to add to the remap line for caches

	.. seealso:: `The Apache Trafficserver documentation for the Regex Remap plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_remap.en.html>`_

:signed:           ``true`` if token-based authentication is enabled for this :term:`Delivery Service`, ``false`` otherwise
:signingAlgorithm: Type of URL signing method to sign the URLs, basically comes down to one of two plugins or ``null``:

	``null``
		Token-based authentication is not enabled for this :term:`Delivery Service`
	url_sig:
		URL Signing token-based authentication is enabled for this :term:`Delivery Service`
	uri_signing
		URI Signing token-based authentication is enabled for this :term:`Delivery Service`

	.. seealso:: `The Apache Trafficserver documentation for the url_sig plugin <https://docs.trafficserver.apache.org/en/8.0.x/admin-guide/plugins/url_sig.en.html>`_ and `the draft RFC for uri_signing <https://tools.ietf.org/html/draft-ietf-cdni-uri-signing-16>`_ - note, however that the current implementation of uri_signing uses Draft 12 of that RFC document, NOT the latest.

	.. versionadded:: 1.3

:sslKeyVersion:       This integer indicates the generation of keys in use by the :term:`Delivery Service` - if any - and is incremented by the Traffic Portal client whenever new keys are generated

	.. warning:: This number will not be correct if keys are manually replaced using the API, as the key generation API does not increment it!

:tenant:            The name of the tenant who owns this :term:`Delivery Service`

	.. versionadded:: 1.3

:tenantId:            The integral, unique identifier of the tenant who owns this :term:`Delivery Service`
:trRequestHeaders:    If defined, this takes the form of a string of HTTP headers to be included in Traffic Router access logs for requests - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_
:trResponseHeaders:   If defined, this takes the form of a string of HTTP headers to be included in Traffic Router responses - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_
:type:                The name of the routing type of this :term:`Delivery Service` e.g. "HTTP"
:typeId:              The integral, unique identifier of the routing type of this :term:`Delivery Service`
:xmlId:               A unique string that describes this :term:`Delivery Service` - exists for legacy reasons

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Mon, 19 Nov 2018 19:29:40 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Mon, 19 Nov 2018 23:29:40 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: wSCPoNQbFTN0FonjXYH13jwTvOwo0ltSD0ACRQ4d/eaWIfzNyAFAD/RapflUP2PIqttb6NlnHkZve0j6ETJ+gw==
	Content-Length: 1439

	{ "alerts": [
		{
			"level": "success",
			"text": "Deliveryservice safe update was successful."
		}
	],
	"response": [
		{
			"profileId": null,
			"protocol": 0,
			"deepCachingType": "NEVER",
			"regionalGeoBlocking": 0,
			"routingName": "video",
			"orgServerFqdn": "http://origin.infra.ciab.test",
			"cdnId": 2,
			"geoProvider": 0,
			"longDesc2": null,
			"globalMaxMbps": null,
			"dnsBypassIp6": null,
			"geoLimit": 0,
			"maxDnsAnswers": null,
			"id": 1,
			"sslKeyVersion": null,
			"midHeaderRewrite": null,
			"geoLimitRedirectURL": null,
			"active": 1,
			"logsEnabled": 1,
			"initialDispersion": 1,
			"regexRemap": null,
			"geoLimitCountries": null,
			"missLat": 42,
			"anonymousBlockingEnabled": 0,
			"longDesc": "A :term:`Delivery Service` created for the CDN-in-a-Box project",
			"matchList": [
				{
					"pattern": ".*\\.demo1\\..*",
					"setNumber": 0,
					"type": "HOST_REGEXP"
				}
			],
			"rangeRequestHandling": 0,
			"profileName": null,
			"dnsBypassCname": null,
			"globalMaxTps": null,
			"type": "HTTP",
			"httpBypassFqdn": null,
			"infoUrl": "www.info.com",
			"signingAlgorithm": null,
			"missLong": -88,
			"trRequestHeaders": null,
			"trResponseHeaders": null,
			"exampleURLs": [
				"http://video.demo1.mycdn.ciab.test"
			],
			"remapText": null,
			"longDesc1": null,
			"displayName": "demo",
			"qstringIgnore": 0,
			"multiSiteOrigin": 0,
			"xmlId": "demo1",
			"lastUpdated": "2018-11-19 16:26:57.310527+00",
			"ipv6RoutingEnabled": 1,
			"ccrDnsTtl": null,
			"dscp": 0,
			"dnsBypassIp": null,
			"dnsBypassTtl": null,
			"originShield": null,
			"cacheurl": null,
			"edgeHeaderRewrite": null,
			"profileDescription": null,
			"typeId": 1,
			"cdnName": "CDN-in-a-Box",
			"signed": false,
			"checkPath": null,
			"fqPacingRate": null
		}
	]}


.. [1] Users with the "admin" or "operations" roles will be able to edit *any*:term:`Delivery Service`, whereas other users will only be able to edit :term:`Delivery Service`\ s that their tenant has permissions to edit.
.. [2] This only applies to HTTP-routed :term:`Delivery Service`\ s
.. [3] See :ref:`ds-multi-site-origin`
.. [4] This only applies to DNS-routed :term:`Delivery Service`\ s
.. [5] These fields are required for HTTP-routed and DNS-routed :term:`Delivery Service`\ s, but are optional for (and in fact may have no effect on) STEERING and ANY_MAP :term:`Delivery Service`\ s
