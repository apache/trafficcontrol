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

.. _to-api-servers-id-deliveryservices:

***********************************
``servers/{{ID}}/deliveryservices``
***********************************

``GET``
=======
Retrieves all :term:`Delivery Service`\ s assigned to a specific server.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------------------------------------------------+
	| Name | Description                                                                                                  |
	+======+==============================================================================================================+
	|  ID  | The integral, unique identifier of the server for which assigned :term:`Delivery Service`\ s shall be listed |
	+------+--------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/servers/9/deliveryservices HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:active:                   ``true`` if the :term:`Delivery Service` is active, ``false`` otherwise
:anonymousBlockingEnabled: ``true`` if :ref:`Anonymous Blocking <anonymous_blocking-qht>` has been configured for the :term:`Delivery Service`, ``false`` otherwise
:cacheurl:                 A setting for a deprecated feature of now-unsupported :abbr:`ATS (Apache Traffic Server)` versions

	.. deprecated:: ATCv3.0
		This field has been deprecated in Traffic Control 3.x and is subject to removal in Traffic Control 4.x or later

:ccrDnsTtl:                The :abbr:`TTL (Time To Live)` of the DNS response for A or AAAA record queries requesting the IP address of the Traffic Router - named "ccrDnsTtl" for legacy reasons
:cdnId:                    The integral, unique identifier of the CDN to which the :term:`Delivery Service` belongs
:cdnName:                  Name of the CDN to which the :term:`Delivery Service` belongs
:checkPath:                The path portion of the URL to check connections to this :term:`Delivery Service`'s origin server
:consistentHashRegex:      If defined, this is a regex used for the Pattern-Based Consistent Hashing feature. It is only applicable for HTTP and Steering Delivery Services

	.. versionadded:: 1.5

:displayName:              The display name of the :term:`Delivery Service`
:dnsBypassCname:           Domain name to overflow requests for HTTP :term:`Delivery Service`\ s - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`\ [3]_
:dnsBypassIp:              The IPv4 IP to use for bypass on a DNS :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`\ [3]_
:dnsBypassIp6:             The IPv6 IP to use for bypass on a DNS :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`\ [3]_
:dnsBypassTtl:             The time for which a DNS bypass of this :term:`Delivery Service`\ shall remain active\ [3]_
:dscp:                     The :abbr:`DSCP (Differentiated Services Code Point)` with which to mark traffic as it leaves the CDN and reaches clients
:edgeHeaderRewrite:        Rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite :abbr:`ATS (Apache Traffic Server)` plugin
:fqPacingRate:             The Fair-Queuing Pacing Rate in Bytes per second set on the all TCP connection sockets in the :term:`Delivery Service` (see :manpage:`tc-fq_codel(8)` for more information) - Linux only
:geoLimit:                 The setting that determines how content is geographically limited - this is an integer on the interval [0-2] where the values have these meanings:

	0
		None - no limitations
	1
		Only route when the client's IP is found in the :term:`Coverage Zone File`
	2
		Only route when the client's IP is found in the :term:`Coverage Zone File`, or when the client can be determined to be from the United States of America

	.. warning:: This does not prevent access to content or make content secure; it merely prevents routing to the content through Traffic Router

:geoLimitCountries:   A string containing a comma-separated list of country codes (e.g. "US,AU") which are allowed to request content through this :term:`Delivery Service`
:geoLimitRedirectUrl: A URL to which clients blocked by :ref:`Regional Geographic Blocking <regionalgeo-qht>` or the ``geoLimit`` settings will be re-directed
:geoProvider:         An integer that represents the provider of a database for mapping IPs to geographic locations; currently only the following values are supported:

	0
		`The "Maxmind" GeoIP2 database (default) <https://www.maxmind.com/en/geoip2-databases>`_
	1
		`Neustar GeoPoint IP address database <https://www.security.neustar/digital-performance/ip-intelligence/ip-address-data>`_

		.. warning:: It's not clear whether Neustar databases are actually supported; this is an old option and compatibility may have been broken over time.

:globalMaxMbps:       The maximum global bandwidth allowed on this :term:`Delivery Service`. If exceeded, traffic will be routed to ``dnsBypassIp`` (or ``dnsBypassIp6`` for IPv6 traffic) for DNS :term:`Delivery Service`\ s and to ``httpBypassFqdn`` for HTTP :term:`Delivery Service`\ s
:globalMaxTps:        The maximum global transactions per second allowed on this :term:`Delivery Service`. When this is exceeded traffic will be sent to the ``dnsBypassIp`` (and/or ``dnsBypassIp6``) for DNS :term:`Delivery Service`\ s and to the httpBypassFqdn for HTTP :term:`Delivery Service`\ s
:httpBypassFqdn:      The HTTP destination to use for bypass on an HTTP :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`
:id:                  An integral, unique identifier for this :term:`Delivery Service`
:infoUrl:             This is a string which is expected to contain at least one URL pointing to more information about the :term:`Delivery Service`. Historically, this has been used to link relevant JIRA tickets
:initialDispersion:  The number of :term:`cache server`\ s between which traffic requesting the same object will be randomly split - meaning that if 4 clients all request the same object (one after another), then if this is above 4 there is a possibility that all 4 are cache misses. For most use-cases, this should be 1\ [1]_
:ipv6RoutingEnabled: If ``true``, clients that connect to Traffic Router using IPv6 will be given the IPv6 address of a suitable Edge-tier :term:`cache server`; if ``false`` all addresses will be IPv4, regardless of the client connection\ [1]_
:lastUpdated:        The date and time at which this :term:`Delivery Service` was last updated, in a :manpage:`ctime(3)`-like format
:logsEnabled:        If ``true``, logging is enabled for this :term:`Delivery Service`, otherwise it is disabled
:longDesc:           A description of the :term:`Delivery Service`
:longDesc1:          A field used when more detailed information that that provided by ``longDesc`` is desired
:longDesc2:          A field used when even more detailed information that that provided by either ``longDesc`` or ``longDesc1`` is desired
:matchList:          An array of methods used by Traffic Router to determine whether or not a request can be serviced by this :term:`Delivery Service`

	:pattern:   A regular expression - the use of this pattern is dependent on the ``type`` field (backslashes are escaped)
	:setNumber: An integral, unique identifier for the set of types to which the ``type`` field belongs
	:type:      The :term:`Type` of match performed using ``pattern`` to determine whether or not to use this :term:`Delivery Service`

		HOST_REGEXP
			Use the :term:`Delivery Service` if ``pattern`` matches the ``Host:`` HTTP header of an HTTP request\ [1]_
		HEADER_REGEXP
			Use the :term:`Delivery Service` if ``pattern`` matches an HTTP header (both the name and value) in an HTTP request\ [1]_
		PATH_REGEXP
			Use the :term:`Delivery Service` if ``pattern`` matches the request path of this :term:`Delivery Service`'s URL
		STEERING_REGEXP
			Use the :term:`Delivery Service` if ``pattern`` matches the ``xml_id`` of one of this :term:`Delivery Service`'s "Steering" target :term:`Delivery Service`\ s

:maxDnsAnswers:    The maximum number of IPs to put in responses to A/AAAA DNS record requests (0 means all available)\ [3]_
:midHeaderRewrite: Rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite :abbr:`ATS (Apache Traffic Server)` plugin
:missLat:          The latitude to use when the client cannot be found in the :term:`Coverage Zone File` or a geographic IP lookup
:missLong:         The longitude to use when the client cannot be found in the :term:`Coverage Zone File` or a geographic IP lookup
:multiSiteOrigin:  ``true`` if the Multi Site Origin feature is enabled for this :term:`Delivery Service`, ``false`` otherwise\ [2]_
:orgServerFqdn:    The URL of the :term:`Delivery Service`'s origin server for use in retrieving content from the :term:`origin server`

	.. note:: Despite the field name, this must truly be a full URL - including the protocol (e.g. ``http://`` or ``https://``) - **NOT** merely the server's :abbr:`FQDN (Fully Qualified Domain Name)`

:originShield:       An "origin shield" is a forward proxy that sits between Mid-tier :term:`cache server`\ s and the :term:`origin` and performs further caching beyond what's offered by a standard CDN. This field is a string of :abbr:`FQDN (Fully Qualified Domain Name)`\ s to use as origin shields, delimited by ``|``
:profileDescription: The description of the Traffic Router :term:`Profile` with which this :term:`Delivery Service` is associated
:profileId:          The integral, unique identifier for the Traffic Router :term:`Profile` with which this :term:`Delivery Service` is associated
:profileName:        The name of the Traffic Router :term:`Profile` with which this :term:`Delivery Service` is associated
:protocol:           The protocol which clients will use to communicate with Edge-tier :term:`cache server`\ s\ [1]_ - this is an integer on the interval [0-2] where the values have these meanings:

	0
		HTTP
	1
		HTTPS
	2
		Both HTTP and HTTPS

:qstringIgnore: Tells :term:`cache server`\ s whether or not to consider URLs with different query parameter strings to be distinct - this is an integer on the interval [0-2] where the values have these meanings:

	0
		URLs with different query parameter strings will be considered distinct for caching purposes, and query strings will be passed upstream to the :term:`origin`
	1
		URLs with different query parameter strings will be considered identical for caching purposes, and query strings will be passed upstream to the :term:`origin`
	2
		Query strings are stripped out by Edge-tier :term:`cache server`\ s, and thus are neither taken into consideration for caching purposes, nor passed upstream in requests to the :term:`origin`

:rangeRequestHandling: Tells caches how to handle range requests\ [4]_ - this is an integer on the interval [0,2] where the values have these meanings:

	0
		Range requests will not be cached, but range requests that request ranges of content already cached will be served from the :term:`cache server`
	1
		Use the `background_fetch plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/background_fetch.en.html>`_ to service the range request while caching the whole object
	2
		Use the `experimental cache_range_requests plugin <https://github.com/apache/trafficserver/tree/master/plugins/experimental/cache_range_requests>`_ to treat unique ranges as unique objects

:regexRemap: A regular expression "remap rule" to apply to this :term:`Delivery Service` at the Edge tier

	.. seealso:: `The Apache Trafficserver documentation for the Regex Remap plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_remap.en.html>`_

:regionalGeoBlocking: ``true`` if Regional Geo Blocking is in use within this :term:`Delivery Service`, ``false`` otherwise

	.. seealso:: See :ref:`regionalgeo-qht` for more information

:remapText: Additional, raw text to add to the line for this :term:`Delivery Service` for :term:`cache server`\ s

	.. seealso:: `The Apache Trafficserver documentation for the Regex Remap plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_remap.en.html>`_

:signed:           ``true`` if token-based authentication is enabled for this :term:`Delivery Service`, ``false`` otherwise
:signingAlgorithm: Type of URL signing method to sign the URLs, basically comes down to one of two plugins or ``null``:

	``null``
		Token-based authentication is not enabled for this :term:`Delivery Service`
	url_sig:
		URL Signing token-based authentication is enabled for this :term:`Delivery Service`
	uri_signing
		URI Signing token-based authentication is enabled for this :term:`Delivery Service`

	.. seealso:: `The Apache Trafficserver documentation for the url_sig plugin <https://docs.trafficserver.apache.org/en/8.0.x/admin-guide/plugins/url_sig.en.html>`_ and `the draft RFC for uri_signing <https://tools.ietf.org/html/draft-ietf-cdni-uri-signing-16>`_ - note, however that the current implementation of uri_signing uses Draft 12 of that RFC document, **NOT** the latest

:sslKeyVersion: This integer indicates the generation of keys in use by the :term:`Delivery Service` - if any - and is incremented by the Traffic Portal client whenever new keys are generated

	.. warning:: This number will not be correct if keys are manually replaced using the API, as the key generation API does not increment it!

:tenantId:            The integral, unique identifier of the :term:`Tenant` who owns this :term:`Delivery Service`
:trRequestHeaders:    If defined, this takes the form of a string of HTTP headers to be included in Traffic Router access logs for requests - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [1]_
:trResponseHeaders:   If defined, this takes the form of a string of HTTP headers to be included in Traffic Router responses - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [1]_
:type:                The name of the routing type of this :term:`Delivery Service` e.g. "HTTP"
:typeId:              The integral, unique identifier of the routing type of this :term:`Delivery Service`
:xmlId:               A unique string that describes this :term:`Delivery Service` - exists for legacy reasons, but is used heavily by Traffic Control components

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: heK6DafnKW6KdyqQ7lTJQcStli3ixkWYjnbQ2EzR8ZU6Tlij3Takr6CNr0BcD5yWFVN1D8mvMPcj5XLP3FTt5w==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 16:53:04 GMT
	Content-Length: 1129

	{ "response": [
		{
			"active": true,
			"cacheurl": null,
			"ccrDnsTtl": null,
			"cdnId": 2,
			"checkPath": null,
			"deepCachingType": null,
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
			"lastUpdated": "2018-12-05 17:51:00+00",
			"logsEnabled": true,
			"longDesc": "Apachecon North America 2018",
			"longDesc1": null,
			"longDesc2": null,
			"maxDnsAnswers": null,
			"midHeaderRewrite": null,
			"missLat": 42,
			"missLong": -88,
			"multiSiteOrigin": false,
			"multiSiteOriginAlgo": null,
			"originShield": null,
			"orgServerFqdn": "http://origin.infra.ciab.test",
			"profileDescription": null,
			"profileId": null,
			"protocol": 0,
			"qstringIgnore": 0,
			"rangeRequestHandling": 0,
			"regexRemap": null,
			"regionalGeoBlocking": false,
			"remapText": null,
			"routingName": "video",
			"signingAlgorithm": null,
			"sslKeyVersion": null,
			"trRequestHeaders": null,
			"trResponseHeaders": null,
			"tenantId": 1,
			"typeId": 1,
			"xmlId": "demo1"
		}
	]}

.. [1] This only applies to HTTP-routed :term:`Delivery Service`\ s
.. [2] See :ref:`ds-multi-site-origin`
.. [3] This only applies to DNS-routed :term:`Delivery Service`\ s
.. [4] These fields are required for HTTP-routed and DNS-routed :term:`Delivery Service`\ s, but are optional for (and in fact may have no effect on) STEERING and ANY_MAP :term:`Delivery Service`\ s
