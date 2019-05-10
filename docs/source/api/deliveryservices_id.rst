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

***************************
``deliveryservices/{{ID}}``
***************************
.. deprecated:: 1.1
	Use the ``id`` query parameter of :ref:`to-api-deliveryservices` instead

``GET``
=======
Retrieves a specific :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------+----------+------------------------------------------------------------------------------------------------------------------------------------------------+
	| Name        | Required | Description                                                                                                                                    |
	+=============+==========+================================================================================================================================================+
	| cdn         | no       | Show only the :term:`Delivery Service`\ s belonging to the CDN identified by this integral, unique identifier                                  |
	+-------------+----------+------------------------------------------------------------------------------------------------------------------------------------------------+
	| logsEnabled | no       | If true, return only :term:`Delivery Service`\ s with logging enabled, otherwise return only :term:`Delivery Service`\ s with logging disabled |
	+-------------+----------+------------------------------------------------------------------------------------------------------------------------------------------------+
	| profile     | no       | Return only :term:`Delivery Service`\ s using the profile identified by this integral, unique identifier                                       |
	+-------------+----------+------------------------------------------------------------------------------------------------------------------------------------------------+
	| tenant      | no       | Show only the :term:`Delivery Service`\ s belonging to the tenant identified by this integral, unique identifier                               |
	+-------------+----------+------------------------------------------------------------------------------------------------------------------------------------------------+
	| type        | no       | Return only :term:`Delivery Service`\ s of the :term:`Delivery Service` type identified by this integral, unique identifier                    |
	+-------------+----------+------------------------------------------------------------------------------------------------------------------------------------------------+

.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------+
	| Name | Description                                                                     |
	+======+=================================================================================+
	| ID   | The integral, unique identifier of the :term:`Delivery Service` to be retrieved |
	+------+---------------------------------------------------------------------------------+


Response Structure
------------------
.. versionchanged:: 1.5
	Added ``cacheassignmentgroups`` field

.. versionchanged:: 1.4
	Added ``consistentHashRegex`` and ``maxOriginConnection`` fields

.. versionchanged:: 1.3
	Added fields: ``deepCachingType``, ``signingAlgorithm``, ``tenant``, ``fqPacingRate``, ``trResponseHeaders`` and ``trRequestHeaders``.

:active:                   ``true`` if the :term:`Delivery Service` is active, ``false`` otherwise
:anonymousBlockingEnabled: ``true`` if :ref:`Anonymous Blocking <anonymous_blocking-qht>` has been configured for the :term:`Delivery Service`, ``false`` otherwise
:cacheassignmentgroups:    Array of cache assignment group IDs that are assigned to this delivery service

	.. versionadded:: 1.5

:cacheurl:                 A setting for a deprecated feature of now-unsupported Trafficserver versions

	.. deprecated:: ATCv3.0
		This field has been deprecated in Traffic Control 3.x and is subject to removal in Traffic Control 4.x or later

:ccrDnsTtl:                The Time To Live (TTL) of the DNS response for A or AAAA record queries requesting the IP address of the Traffic Router - named "ccrDnsTtl" for legacy reasons
:cdnId:                    The integral, unique identifier of the CDN to which the :term:`Delivery Service` belongs
:cdnName:                  Name of the CDN to which the :term:`Delivery Service` belongs
:checkPath:                The path portion of the URL to check connections to this :term:`Delivery Service`'s origin server
:consistentHashRegex:      If defined, this is a regex used for the Pattern-Based Consistent Hashing feature. It is only applicable for HTTP and Steering Delivery Services

	.. versionadded:: 1.4

:deepCachingType:          A string that describes when "Deep Caching" will be used by this :term:`Delivery Service` - one of:

	ALWAYS
		"Deep Caching" will always be used with this :term:`Delivery Service`
	NEVER
		"Deep Caching" will never be used with this :term:`Delivery Service`

	.. versionadded:: 1.3

:displayName:              The display name of the :term:`Delivery Service`
:dnsBypassCname:           Domain name to overflow requests for HTTP :term:`Delivery Service`\ s - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`
:dnsBypassIp:              The IPv4 IP to use for bypass on a DNS :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`
:dnsBypassIp6:             The IPv6 IP to use for bypass on a DNS :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`
:dnsBypassTtl:             The time for which a DNS bypass of this :term:`Delivery Service`\ shall remain active
:dscp:                     The :abbr:`FQDN (Differentiated Services Code Point)` with which to mark traffic as it leaves the CDN and reaches clients
:edgeHeaderRewrite:        Rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite :abbr:`ATS (Apache Traffic Server)` plugin
:fqPacingRate:             The Fair-Queuing Pacing Rate in Bytes per second set on the all TCP connection sockets in the :term:`Delivery Service` (see ``man tc-fc_codel`` for more information) - Linux only

	.. versionadded:: 1.3

:geoLimit:                 The setting that determines how content is geographically limited - this is an integer on the interval [0-2] where the values have these meanings:
:geoLimitCountries:        A string containing a comma-separated list of country codes (e.g. "US,AU") which are allowed to request content through this :term:`Delivery Service`
:geoLimitRedirectUrl:      A URL to which clients blocked by :ref:`Regional Geographic Blocking <regionalgeo-qht>` or the ``geoLimit`` settings will be re-directed

	0
		None - no limitations
	1
		Only route when the client's IP is found in the :term:`Coverage Zone File`
	2
		Only route when the client's IP is found in the :term:`Coverage Zone File`, or when the client can be determined to be from the United States of America

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

:maxDnsAnswers:      The maximum number of IPs to put in a A/AAAA response for a DNS :term:`Delivery Service` (0 means all available)
:midHeaderRewrite:   Rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite Apache Trafficserver plugin
:missLat:            The latitude to use when the client cannot be found in the CZF or a geographic IP lookup
:missLong:           The longitude to use when the client cannot be found in the CZF or a geographic IP lookup
:multiSiteOrigin:    ``true`` if the Multi Site Origin feature is enabled for this :term:`Delivery Service`, ``false`` otherwise\ [3]_
:originShield:       An "origin shield" is a forward proxy that sits between Mid-tier caches and the origin and performs further caching beyond what's offered by a standard CDN. This field is a string of FQDNs to use as origin shields, delimited by ``|``
:orgServerFqdn:      The origin server's Fully Qualified Domain Name (FQDN) - including the protocol (e.g. http:// or https://) - for use in retrieving content from the origin server
:profileDescription: The description of the Traffic Router Profile with which this :term:`Delivery Service` is associated
:profileId:          The integral, unique identifier for the Traffic Router profile with which this :term:`Delivery Service` is associated
:profileName:        The name of the Traffic Router Profile with which this :term:`Delivery Service` is associated
:protocol:           The protocol which clients will use to communicate with Edge-tier :term:`cache server` s\ [2]_ - this is an integer on the interval [0-2] where the values have these meanings:

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

	.. versionadded:: 1.3

:trResponseHeaders:   If defined, this takes the form of a string of HTTP headers to be included in Traffic Router responses - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_

	.. versionadded:: 1.3

:type:                The name of the routing type of this :term:`Delivery Service` e.g. "HTTP"
:typeId:              The integral, unique identifier of the routing type of this :term:`Delivery Service`
:xmlId:               A unique string that describes this :term:`Delivery Service` - exists for legacy reasons

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: Mw4ZsiNKfnxZvN+LsfAzxIZjgGTzcBLcZK24mMdhN1XMRBtwEj9VI3ExNvWKv3dp0f3HRRCUTx6C+ST8bRL9jA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 21:43:36 GMT
	Content-Length: 1290

	{ "response": [
		{
			"active": true,
			"anonymousBlockingEnabled": false,
			"cacheassignmentgroups": [1, 2, 3],
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
		}
	]}


.. [1] Users with the roles "admin" and/or "operation" will be able to see *all* :term:`Delivery Service`\ s, whereas any other user will only see the :term:`Delivery Service`\ s their Tenant is allowed to see.
.. [2] This only applies to HTTP :term:`Delivery Service`\ s
.. [3] See :ref:`ds-multi-site-origin`
.. [4] This only applies to DNS-routed :term:`Delivery Service`\ s

``PUT``
=======
Allows users to edit an existing :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [10]_
:Response Type:  **NOT PRESENT** - Despite returning a ``200 OK`` response (rather than e.g. a ``204 NO CONTENT`` response), this endpoint does **not** return a representation of the modified resource in its payload, and instead returns nothing - not even a success message.

Request Structure
-----------------
:active:                   If ``true``, the :term:`Delivery Service` will immediately become active and serves traffic
:anonymousBlockingEnabled: An optional field which, if defined and ``true`` will cause :ref:`Anonymous Blocking <anonymous_blocking-qht>` to be used with the new :term:`Delivery Service`
:cacheassignmentgroups:    Array of cache assignment group IDs to be assigned to this delivery service
	.. versionadded:: 1.5

:cacheurl:                 An optional setting for a deprecated feature of now-unsupported Trafficserver versions (read: "Don't use this")

	.. deprecated:: ATCv3.0
		This field has been deprecated in Traffic Control 3.x and is subject to removal in Traffic Control 4.x or later

:ccrDnsTtl:                The Time To Live (TTL) in seconds of the DNS response for A or AAAA record queries requesting the IP address of the Traffic Router - named "ccrDnsTtl" for legacy reasons
:cdnId:                    The integral, unique identifier for the CDN to which this :term:`Delivery Service`\ shall be assigned
:checkPath:                The path portion of the URL which will be used to check connections to this :term:`Delivery Service`'s origin server
:consistentHashRegex:      If defined, this is a regex used for the Pattern-Based Consistent Hashing feature. It is only applicable for HTTP and Steering Delivery Services

	.. versionadded:: 1.4

:deepCachingType:          A string describing when to do Deep Caching for this :term:`Delivery Service`:

	NEVER
		Deep Caching will never be used by this :term:`Delivery Service` (default)
	ALWAYS
		Deep Caching will always be used by this :term:`Delivery Service`

:displayName:       The human-friendly name for this :term:`Delivery Service`
:dnsBypassCname:    Domain name to overflow requests for HTTP :term:`Delivery Service`\ s - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`
:dnsBypassIp:       The IPv4 IP to use for bypass on a DNS :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`
:dnsBypassIp6:      The IPv6 IP to use for bypass on a DNS :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`
:dnsBypassTtl:      The time for which a DNS bypass of this :term:`Delivery Service`\ shall remain active
:dscp:              The Differentiated Services Code Point (DSCP) with which to mark downstream (EDGE -> customer) traffic. This should be zero in most cases
:edgeHeaderRewrite: An optional string which, if present, defines rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite Apache Trafficserver plugin
:fqPacingRate:      An optional integer which, if present, sets the Fair-Queuing Pacing Rate in bytes per second set on the all TCP connection sockets in the :term:`Delivery Service` (see ``man tc-fc_codel`` for more information) - Linux only, defaults to 0 meaning "disabled"
:geoLimit:          The setting that determines how content is geographically limited - this is an integer on the interval [0-2] where the values have these meanings:

	0
		None - no limitations
	1
		Only route when the client's IP is found in the Coverage Zone File (CZF)
	2
		Only route when the client's IP is found in the CZF, or when the client can be determined to be from the United States of America

	.. warning:: This does not prevent access to content or make content secure; it merely prevents routing to the content through Traffic Router

:geoLimitCountries:   A string containing a comma-separated list of country codes (e.g. "US,AU") which are allowed to request content through this :term:`Delivery Service`\ [5]_
:geoLimitRedirectUrl: A URL to which clients blocked by :ref:`Regional Geographic Blocking <regionalgeo-qht>` or the ``geoLimit`` settings will be re-directed\ [5]_
:geoProvider:         An integer that represents the provider of a database for mapping IPs to geographic locations; currently only the following values are supported:

	0
		The "Maxmind" GeoIP2 database (default)
	1
		Neustar

:globalMaxMbps:      An optional integer that will set the maximum global bandwidth allowed on this :term:`Delivery Service`. If exceeded, traffic will be routed to ``dnsBypassIp`` (or ``dnsBypassIp6`` for IPv6 traffic) for DNS :term:`Delivery Service`\ s and to ``httpBypassFqdn`` for HTTP :term:`Delivery Service`\ s
:globalMaxTps:       An optional integer that will set the maximum global transactions per second allowed on this :term:`Delivery Service`. When this is exceeded traffic will be sent to the ``dnsBpassIp`` (and/or ``dnsBypassIp6``)for DNS :term:`Delivery Service`\ s and to the ``httpBypassFqdn`` for HTTP :term:`Delivery Service`\ s
:httpBypassFqdn:     An optional Fully Qualified Domain Name (FQDN) to use for bypass on an HTTP :term:`Delivery Service` - bypass starts when the traffic on this :term:`Delivery Service` exceeds ``globalMaxMbps``, or when more than ``globalMaxTps`` is being exceeded within the :term:`Delivery Service`\ [2]_
:infoUrl:            An optional string which, if present, is expected to contain at least one URL pointing to more information about the :term:`Delivery Service`. Historically, this has been used to link relevant JIRA tickets
:initialDispersion:  The number of caches between which traffic requesting the same object will be randomly split - meaning that if 4 clients all request the same object (one after another), then if this is above 4 there is a possibility that all 4 are cache misses. For most use-cases, this should be 1\ [2]_\ [6]_
:ipv6RoutingEnabled: If ``true``, clients that connect to Traffic Router using IPv6 will be given the IPv6 address of a suitable Edge-tier cache; if ``false`` all addresses will be IPv4, regardless of the client connection - optional for ANY_MAP :term:`Delivery Service`\ s
:logsEnabled:        If ``true``, logging is enabled for this :term:`Delivery Service`, otherwise it is disabled
:longDesc:           An optional description of the :term:`Delivery Service`
:longDesc1:          An optional field used when more detailed information that that provided by ``longDesc`` is desired
:longDesc2:          An optional field used when even more detailed information that that provided by either ``longDesc`` or ``longDesc1`` is desired
:maxDnsAnswers:      An optional field which, when present, specifies the maximum number of IPs to put in responses to A/AAAA DNS record requests - defaults to 0, meaning "no limit"\ [4]_
:midHeaderRewrite:   An optional string containing rewrite operations to be performed on TCP headers at the Edge-tier cache level - used by the Header Rewrite Apache Trafficserver plugin
:missLat:            The latitude to use when the client cannot be found in the CZF or a geographic IP lookup\ [7]_
:missLong:           The longitude to use when the client cannot be found in the CZF or a geographic IP lookup\ [7]_
:multiSiteOrigin:    ``true`` if the Multi Site Origin feature is enabled for this :term:`Delivery Service`, ``false`` otherwise\ [3]_\ [7]_
:orgServerFqdn:      The URL of the :term:`Delivery Service`'s origin server for use in retrieving content from the origin server\ [7]_

	.. note:: Despite the field name, this must truly be a full URL - including the protocol (e.g. ``http://`` or ``https://``) - **NOT** merely the server's Fully Qualified Domain Name (FQDN)

:originShield: An "origin shield" is a forward proxy that sits between Mid-tier caches and the origin and performs further caching beyond what's offered by a standard CDN. This optional field is a string of FQDNs to use as origin shields, delimited by ``|``
:profileId:    An optional, integral, unique identifier for the Traffic Router profile with which this :term:`Delivery Service`\ shall be associated
:protocol:     The protocol which clients will use to communicate with Edge-tier :term:`cache server` s - this is an (optional for ANY_MAP :term:`Delivery Service`\ s) integer on the interval [0,2] where the values have these meanings:

	0
		HTTP
	1
		HTTPS
	2
		Both HTTP and HTTPS

:qstringIgnore: Tells caches whether or not to consider URLs with different query parameter strings to be distinct\ [7]_ - this is an integer on the interval [0-2] where the values have these meanings:

	0
		URLs with different query parameter strings will be considered distinct for caching purposes, and query strings will be passed upstream to the origin
	1
		URLs with different query parameter strings will be considered identical for caching purposes, and query strings will be passed upstream to the origin
	2
		Query strings are stripped out by Edge-tier caches, and thus are neither taken into consideration for caching purposes, nor passed upstream in requests to the origin

:rangeRequestHandling: Tells caches how to handle range requests\ [7]_ - this is an integer on the interval [0,2] where the values have these meanings:

	0
		Range requests will not be cached, but range requests that request ranges of content already cached will be served from the cache
	1
		Use the `background_fetch plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/background_fetch.en.html>`_ to service the range request while caching the whole object
	2
		Use the `experimental cache_range_requests plugin <https://github.com/apache/trafficserver/tree/master/plugins/experimental/cache_range_requests>`_ to treat unique ranges as unique objects

:regexRemap: An optional, regular expression remap rule to apply to this :term:`Delivery Service` at the Edge tier

	.. seealso:: `The Apache Trafficserver documentation for the Regex Remap plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_remap.en.html>`_

:regionalGeoBlocking: ``true`` if Regional Geo Blocking is in use within this :term:`Delivery Service`, ``false`` otherwise - see :ref:`regionalgeo-qht` for more information
:remapText:           Optional, raw text to add to the remap line for caches

	.. seealso:: `The Apache Trafficserver documentation for the Regex Remap plugin <https://docs.trafficserver.apache.org/en/latest/admin-guide/plugins/regex_remap.en.html>`_

:routingName:      The routing name of this :term:`Delivery Service`, used as the top-level part of the FQDN used by clients to request content from the :term:`Delivery Service` e.g. ``routingName.xml_id.CDNName.com``
:signed:           An optional field which should be ``true`` if token-based authentication\ [8]_ will be enabled for this :term:`Delivery Service`, ``false`` (default) otherwise
:signingAlgorithm: Type of URL signing method to sign the URLs\ [8]_, basically comes down to one of two plugins or ``null``:

	``null``
		Token-based authentication is not enabled for this :term:`Delivery Service`
	url_sig:
		URL Signing token-based authentication is enabled for this :term:`Delivery Service`
	uri_signing
		URI Signing token-based authentication is enabled for this :term:`Delivery Service`

	.. seealso:: `The Apache Trafficserver documentation for the url_sig plugin <https://docs.trafficserver.apache.org/en/8.0.x/admin-guide/plugins/url_sig.en.html>`_ and `the draft RFC for uri_signing <https://tools.ietf.org/html/draft-ietf-cdni-uri-signing-16>`_ - note, however that the current implementation of uri_signing uses Draft 12 of that RFC document, **NOT** the latest

:sslKeyVersion: This optional integer indicates the generation of keys to be used by the :term:`Delivery Service` - if any - and is incremented by the Traffic Portal client whenever new keys are generated

	.. warning:: This number will not be correct if keys are manually replaced using the API, as the key generation API does not increment it!

:tenantId:            An optional, integral, unique identifier of the tenant who will own this :term:`Delivery Service`
:trRequestHeaders:    If defined, this takes the form of a string of HTTP headers to be included in Traffic Router access logs for requests - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_
:trResponseHeaders:   If defined, this takes the form of a string of HTTP headers to be included in Traffic Router responses - it's a template where ``__RETURN__`` translates to a carriage return and line feed (``\r\n``)\ [2]_
:typeId:              The integral, unique identifier for the routing type of this :term:`Delivery Service`
:xmlId:               A unique string that describes this :term:`Delivery Service` - exists for legacy reasons

	.. note:: While this field **must** be present, it is **not** allowed to change; this must be the same as the ``xml_id`` the :term:`Delivery Service` already has. This should almost never be different from the :term:`Delivery Service`'s ``displayName``.


.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/deliveryservices/1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 761
	Content-Type: application/json

	{
		"active": true,
		"anonymousBlockingEnabled": false,
		"cacheassignmentgroups": [],
		"cdnId": 2,
		"cdnName": "CDN-in-a-Box",
		"deepCachingType": "NEVER",
		"displayName": "demo",
		"exampleURLs": [
			"http://video.demo.mycdn.ciab.test"
		],
		"dscp": 0,
		"geoLimit": 0,
		"geoProvider": 0,
		"initialDispersion": 1,
		"ipv6RoutingEnabled": false,
		"lastUpdated": "2018-11-14 18:21:17+00",
		"logsEnabled": true,
		"longDesc": "A :term:`Delivery Service` created expressly for API documentation examples",
		"missLat": -1,
		"missLong": -1,
		"multiSiteOrigin": false,
		"orgServerFqdn": "http://origin.infra.ciab.test",
		"protocol": 0,
		"qstringIgnore": 0,
		"rangeRequestHandling": 0,
		"regionalGeoBlocking": false,
		"routingName": "video",
		"signed": false,
		"tenant": "root",
		"tenantId": 1,
		"typeId": 1,
		"xmlId": "demo1"
	}

.. [5] These fields must be defined if and only if ``geoLimit`` is non-zero
.. [6] These fields are required for HTTP-routed :term:`Delivery Service`\ s, and optional for all others
.. [7] These fields are required for HTTP-routed and DNS-routed :term:`Delivery Service`\ s, but are optional for (and in fact may have no effect on) STEERING and ANY_MAP :term:`Delivery Service`\ s
.. [8] See "token-based-auth" TODO --- wat for more information

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 20 Nov 2018 14:12:25 GMT
	Content-Length: 0
	Content-Type: text/plain; charset=utf-8


.. [10] Users with the roles "admin" and/or "operation" will be able to edit *all* :term:`Delivery Service`\ s, whereas any other user will only be able to edit the :term:`Delivery Service`\ s their Tenant is allowed to edit.

``DELETE``
==========
Deletes the target :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [11]_
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------+
	| Name | Description                                                                     |
	+======+=================================================================================+
	| ID   | The integral, unique identifier of the :term:`Delivery Service` to be retrieved |
	+------+---------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/deliveryservices/2 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...


Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: w9NlQpJJEl56r6iYq/fk8o5WfAXeUS5XR9yDHvKUgPO8lYEo8YyftaSF0MPFseeOk60dk6kQo+MLYTDIAhhRxw==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 20 Nov 2018 14:56:37 GMT
	Content-Length: 57

	{ "alerts": [
		{
			"text": "ds was deleted.",
			"level": "success"
		}
	]}

.. [11] Users with the roles "admin" and/or "operation" will be able to delete *all* :term:`Delivery Service`\ s, whereas any other user will only be able to delete the :term:`Delivery Service`\ s their Tenant is allowed to delete.
