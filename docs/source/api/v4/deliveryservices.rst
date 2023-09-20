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

.. _to-api-v4-deliveryservices:

********************
``deliveryservices``
********************

``GET``
=======
Retrieves :term:`Delivery Services`

:Auth. Required: Yes
:Roles Required: None\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE:READ, CDN:READ, TYPE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                                             |
	+===================+==========+=========================================================================================================================================+
	| cdn               | no       | Show only the :term:`Delivery Services` belonging to the :ref:`ds-cdn` identified by this integral, unique identifier                   |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| id                | no       | Show only the :term:`Delivery Service` that has this integral, unique identifier                                                        |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| logsEnabled       | no       | Show only the :term:`Delivery Services` that have :ref:`ds-logs-enabled` set or not based on this boolean                               |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| profile           | no       | Return only :term:`Delivery Services` using the :term:`Profile` that has this :ref:`profile-id`                                         |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| tenant            | no       | Show only the :term:`Delivery Services` belonging to the :term:`Tenant` identified by this integral, unique identifier                  |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| topology          | no       | Show only the :term:`Delivery Services` assigned to the :term:`Topology` identified by this unique name                                 |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| type              | no       | Return only :term:`Delivery Services` of the :term:`Delivery Service` :ref:`ds-types` identified by this integral, unique identifier    |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| accessibleTo      | no       | Return the :term:`Delivery Services` accessible from a :term:`Tenant` *or it's children* identified by this integral, unique identifier |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| serviceCategory   | no       | Show only the :term:`Delivery Services` belonging to the :term:`Service Category` that has this name                                    |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| xmlId             | no       | Show only the :term:`Delivery Service` that has this text-based, unique identifier                                                      |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| orderby           | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response``                           |
	|                   |          | array                                                                                                                                   |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| sortOrder         | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                                                |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| limit             | no       | Choose the maximum number of results to return                                                                                          |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| offset            | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit                                    |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| page              | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and the first page is 1.    |
	|                   |          | If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to make use of ``page``.                       |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+
	| active            | no       | Show only the :term:`Delivery Services` that have :ref:`ds-active` set or not based on this boolean (whether or not they are active)    |
	+-------------------+----------+-----------------------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/4.1/deliveryservices?xmlId=demo2 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: python-requests/2.24.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:active:                   A boolean that defines :ref:`ds-active`.
:anonymousBlockingEnabled: A boolean that defines :ref:`ds-anonymous-blocking`
:ccrDnsTtl:                 The :ref:`ds-dns-ttl` - named "ccrDnsTtl" for legacy reasons
:cdnId:                     The integral, unique identifier of the :ref:`ds-cdn` to which the :term:`Delivery Service` belongs
:cdnName:                   Name of the :ref:`ds-cdn` to which the :term:`Delivery Service` belongs
:checkPath:                 A :ref:`ds-check-path`
:consistentHashRegex:       A :ref:`ds-consistent-hashing-regex`
:consistentHashQueryParams: An array of :ref:`ds-consistent-hashing-qparams`
:deepCachingType:           The :ref:`ds-deep-caching` setting for this :term:`Delivery Service`
:displayName:               The :ref:`ds-display-name`
:dnsBypassCname:            A :ref:`ds-dns-bypass-cname`
:dnsBypassIp:               A :ref:`ds-dns-bypass-ip`
:dnsBypassIp6:              A :ref:`ds-dns-bypass-ipv6`
:dnsBypassTtl:              The :ref:`ds-dns-bypass-ttl`
:dscp:                      A :ref:`ds-dscp` to be used within the :term:`Delivery Service`
:ecsEnabled:                A boolean that defines the :ref:`ds-ecs` setting on this :term:`Delivery Service`
:edgeHeaderRewrite:         A set of :ref:`ds-edge-header-rw-rules`
:exampleURLs:               An array of :ref:`ds-example-urls`
:firstHeaderRewrite:        A set of :ref:`ds-first-header-rw-rules`
:fqPacingRate:              The :ref:`ds-fqpr`
:geoLimit:                  An integer that defines the :ref:`ds-geo-limit`
:geoLimitCountries:         An array of strings defining the :ref:`ds-geo-limit-countries`
:geoLimitRedirectUrl:       A :ref:`ds-geo-limit-redirect-url`
:geoProvider:               The :ref:`ds-geo-provider`
:globalMaxMbps:             The :ref:`ds-global-max-mbps`
:globalMaxTps:              The :ref:`ds-global-max-tps`
:httpBypassFqdn:            A :ref:`ds-http-bypass-fqdn`
:id:                        An integral, unique identifier for this :term:`Delivery Service`
:infoUrl:                   An :ref:`ds-info-url`
:initialDispersion:         The :ref:`ds-initial-dispersion`
:innerHeaderRewrite:        A set of :ref:`ds-inner-header-rw-rules`
:ipv6RoutingEnabled:        A boolean that defines the :ref:`ds-ipv6-routing` setting on this :term:`Delivery Service`
:lastHeaderRewrite:         A set of :ref:`ds-last-header-rw-rules`
:lastUpdated:               The date and time at which this :term:`Delivery Service` was last updated, in :RFC:`3339` format

	.. versionchanged:: 4.0
		Prior to API version 4.0, this property used :ref:`non-rfc-datetime`.

:logsEnabled: A boolean that defines the :ref:`ds-logs-enabled` setting on this :term:`Delivery Service`
:longDesc:    The :ref:`ds-longdesc` of this :term:`Delivery Service`
:matchList:   The :term:`Delivery Service`'s :ref:`ds-matchlist`

	:pattern:   A regular expression - the use of this pattern is dependent on the ``type`` field (backslashes are escaped)
	:setNumber: An integer that provides explicit ordering of :ref:`ds-matchlist` items - this is used as a priority ranking by Traffic Router, and is not guaranteed to correspond to the ordering of items in the array.
	:type:      The type of match performed using ``pattern``.

:maxDnsAnswers:         The :ref:`ds-max-dns-answers` allowed for this :term:`Delivery Service`
:maxOriginConnections:  The :ref:`ds-max-origin-connections`
:maxRequestHeaderBytes: The :ref:`ds-max-request-header-bytes`
:midHeaderRewrite:      A set of :ref:`ds-mid-header-rw-rules`
:missLat:               The :ref:`ds-geo-miss-default-latitude` used by this :term:`Delivery Service`
:missLong:              The :ref:`ds-geo-miss-default-longitude` used by this :term:`Delivery Service`
:multiSiteOrigin:       A boolean that defines the use of :ref:`ds-multi-site-origin` by this :term:`Delivery Service`
:orgServerFqdn:         The :ref:`ds-origin-url`
:originShield:          A :ref:`ds-origin-shield` string
:profileDescription:    The :ref:`profile-description` of the :ref:`ds-profile` with which this :term:`Delivery Service` is associated
:profileId:             The :ref:`profile-id` of the :ref:`ds-profile` with which this :term:`Delivery Service` is associated
:profileName:           The :ref:`profile-name` of the :ref:`ds-profile` with which this :term:`Delivery Service` is associated
:protocol:              An integral, unique identifier that corresponds to the :ref:`ds-protocol` used by this :term:`Delivery Service`
:qstringIgnore:         An integral, unique identifier that corresponds to the :ref:`ds-qstring-handling` setting on this :term:`Delivery Service`
:rangeRequestHandling:  An integral, unique identifier that corresponds to the :ref:`ds-range-request-handling` setting on this :term:`Delivery Service`
:regexRemap:            A :ref:`ds-regex-remap`
:regional:              A boolean value defining the :ref:`ds-regional` setting on this :term:`Delivery Service`
:regionalGeoBlocking:   A boolean defining the :ref:`ds-regionalgeo` setting on this :term:`Delivery Service`
:remapText:             :ref:`ds-raw-remap`
:requiredCapabilities:  An array of the capabilities that this Delivery Service requires.

	.. versionadded:: 4.1

:serviceCategory:       The name of the :ref:`ds-service-category` with which the :term:`Delivery Service` is associated
:signed:                ``true`` if  and only if ``signingAlgorithm`` is not ``null``, ``false`` otherwise
:signingAlgorithm:      Either a :ref:`ds-signing-algorithm` or ``null`` to indicate URL/URI signing is not implemented on this :term:`Delivery Service`
:rangeSliceBlockSize:   An integer that defines the byte block size for the ATS Slice Plugin. It can only and must be set if ``rangeRequestHandling`` is set to 3.
:sslKeyVersion:         This integer indicates the :ref:`ds-ssl-key-version`
:tenantId:              The integral, unique identifier of the :ref:`ds-tenant` who owns this :term:`Delivery Service`
:tlsVersions:           A list of explicitly supported :ref:`ds-tls-versions`

	.. versionadded:: 4.0

:topology:          The unique name of the :term:`Topology` that this :term:`Delivery Service` is assigned to
:trRequestHeaders:  If defined, this defines the :ref:`ds-tr-req-headers` used by Traffic Router for this :term:`Delivery Service`
:trResponseHeaders: If defined, this defines the :ref:`ds-tr-resp-headers` used by Traffic Router for this :term:`Delivery Service`
:type:              The :ref:`ds-types` of this :term:`Delivery Service`
:typeId:            The integral, unique identifier of the :ref:`ds-types` of this :term:`Delivery Service`
:xmlId:             This :term:`Delivery Service`'s :ref:`ds-xmlid`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 07 Jun 2021 22:52:20 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Jun 2021 21:52:20 GMT
	Content-Length: 847

	{ "response": [
		{
			"active": true,
			"anonymousBlockingEnabled": false,
			"ccrDnsTtl": null,
			"cdnId": 2,
			"cdnName": "CDN-in-a-Box",
			"checkPath": null,
			"consistentHashQueryParams": [],
			"consistentHashRegex": null,
			"deepCachingType": "NEVER",
			"displayName": "Demo 2",
			"dnsBypassCname": null,
			"dnsBypassIp": null,
			"dnsBypassIp6": null,
			"dnsBypassTtl": null,
			"dscp": 0,
			"ecsEnabled": false,
			"edgeHeaderRewrite": null,
			"exampleURLs": [
				"http://video.demo2.mycdn.ciab.test",
				"https://video.demo2.mycdn.ciab.test"
			],
			"firstHeaderRewrite": null,
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
			"innerHeaderRewrite": null,
			"ipv6RoutingEnabled": true,
			"lastHeaderRewrite": null,
			"lastUpdated": "2021-06-07T21:50:03.009954Z",
			"logsEnabled": true,
			"longDesc": "DNS Delivery Service for use with a Federation",
			"matchList": [
				{
					"type": "HOST_REGEXP",
					"setNumber": 0,
					"pattern": ".*\\.demo2\\..*"
				}
			],
			"maxDnsAnswers": null,
			"maxOriginConnections": 0,
			"maxRequestHeaderBytes": 0,
			"midHeaderRewrite": null,
			"missLat": 42,
			"missLong": -88,
			"multiSiteOrigin": true,
			"originShield": null,
			"orgServerFqdn": "http://origin.infra.ciab.test",
			"profileDescription": null,
			"profileId": null,
			"profileName": null,
			"protocol": 2,
			"qstringIgnore": 0,
			"rangeRequestHandling": 0,
			"rangeSliceBlockSize": null,
			"regexRemap": null,
			"regional": false,
			"regionalGeoBlocking": false,
			"remapText": null,
			"requiredCapabilities": [],
			"routingName": "video",
			"serviceCategory": null,
			"signed": false,
			"signingAlgorithm": null,
			"sslKeyVersion": null,
			"tenant": "root",
			"tenantId": 1,
			"tlsVersions": null,
			"topology": "demo1-top",
			"trResponseHeaders": null,
			"trRequestHeaders": null,
			"type": "DNS",
			"typeId": 5,
			"xmlId": "demo2"
		}
	]}


``POST``
========
Allows users to create :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE:CREATE, DELIVERY-SERVICE:READ, CDN:READ, TYPE:READ
:Response Type:  Array

Request Structure
-----------------
:active:                   A boolean that defines :ref:`ds-active`.
:anonymousBlockingEnabled: A boolean that defines :ref:`ds-anonymous-blocking`
:ccrDnsTtl:                 The :ref:`ds-dns-ttl` - named "ccrDnsTtl" for legacy reasons
:cdnId:                     The integral, unique identifier of the :ref:`ds-cdn` to which the :term:`Delivery Service` belongs
:checkPath:                 A :ref:`ds-check-path`
:consistentHashRegex:       A :ref:`ds-consistent-hashing-regex`
:consistentHashQueryParams: An array of :ref:`ds-consistent-hashing-qparams`
:deepCachingType:           The :ref:`ds-deep-caching` setting for this :term:`Delivery Service`
:displayName:               The :ref:`ds-display-name`
:dnsBypassCname:            A :ref:`ds-dns-bypass-cname`
:dnsBypassIp:               A :ref:`ds-dns-bypass-ip`
:dnsBypassIp6:              A :ref:`ds-dns-bypass-ipv6`
:dnsBypassTtl:              The :ref:`ds-dns-bypass-ttl`
:dscp:                      A :ref:`ds-dscp` to be used within the :term:`Delivery Service`
:ecsEnabled:                A boolean that defines the :ref:`ds-ecs` setting on this :term:`Delivery Service`
:edgeHeaderRewrite:         A set of :ref:`ds-edge-header-rw-rules`
:firstHeaderRewrite:        A set of :ref:`ds-first-header-rw-rules`
:fqPacingRate:              The :ref:`ds-fqpr`
:geoLimit:                  An integer that defines the :ref:`ds-geo-limit`
:geoLimitCountries:         A string containing a comma-separated list, or an array of strings defining the :ref:`ds-geo-limit-countries`\ [#geolimit]_
:geoLimitRedirectUrl:       A :ref:`ds-geo-limit-redirect-url`\ [#geolimit]_
:geoProvider:               The :ref:`ds-geo-provider`
:globalMaxMbps:             The :ref:`ds-global-max-mbps`
:globalMaxTps:              The :ref:`ds-global-max-tps`
:httpBypassFqdn:            A :ref:`ds-http-bypass-fqdn`
:infoUrl:                   An :ref:`ds-info-url`
:initialDispersion:         The :ref:`ds-initial-dispersion`
:innerHeaderRewrite:        A set of :ref:`ds-inner-header-rw-rules`
:ipv6RoutingEnabled:        A boolean that defines the :ref:`ds-ipv6-routing` setting on this :term:`Delivery Service`
:lastHeaderRewrite:         A set of :ref:`ds-last-header-rw-rules`
:logsEnabled:               A boolean that defines the :ref:`ds-logs-enabled` setting on this :term:`Delivery Service`
:longDesc:                  The :ref:`ds-longdesc` of this :term:`Delivery Service`
:maxDnsAnswers:             The :ref:`ds-max-dns-answers` allowed for this :term:`Delivery Service`
:maxOriginConnections:      The :ref:`ds-max-origin-connections`
:maxRequestHeaderBytes:     The :ref:`ds-max-request-header-bytes`
:midHeaderRewrite:          A set of :ref:`ds-mid-header-rw-rules`
:missLat:                   The :ref:`ds-geo-miss-default-latitude` used by this :term:`Delivery Service`
:missLong:                  The :ref:`ds-geo-miss-default-longitude` used by this :term:`Delivery Service`
:multiSiteOrigin:           A boolean that defines the use of :ref:`ds-multi-site-origin` by this :term:`Delivery Service`
:orgServerFqdn:             The :ref:`ds-origin-url`
:originShield:              A :ref:`ds-origin-shield` string
:profileId:                 An optional :ref:`profile-id` of a :ref:`ds-profile` with which this :term:`Delivery Service` shall be associated
:protocol:                  An integral, unique identifier that corresponds to the :ref:`ds-protocol` used by this :term:`Delivery Service`
:qstringIgnore:             An integral, unique identifier that corresponds to the :ref:`ds-qstring-handling` setting on this :term:`Delivery Service`
:rangeRequestHandling:      An integral, unique identifier that corresponds to the :ref:`ds-range-request-handling` setting on this :term:`Delivery Service`
:regexRemap:                A :ref:`ds-regex-remap`
:regional:                  A boolean value defining the :ref:`ds-regional` setting on this :term:`Delivery Service`
:regionalGeoBlocking:       A boolean defining the :ref:`ds-regionalgeo` setting on this :term:`Delivery Service`
:remapText:                 :ref:`ds-raw-remap`
:requiredCapabilities:      An array of the capabilities that this Delivery Service requires.

	.. versionadded:: 4.1

:serviceCategory:           The name of the :ref:`ds-service-category` with which the :term:`Delivery Service` is associated - or ``null`` if there is to be no such category
:signed:                    ``true`` if  and only if ``signingAlgorithm`` is not ``null``, ``false`` otherwise
:signingAlgorithm:          Either a :ref:`ds-signing-algorithm` or ``null`` to indicate URL/URI signing is not implemented on this :term:`Delivery Service`
:rangeSliceBlockSize:       An integer that defines the byte block size for the ATS Slice Plugin. It can only and must be set if ``rangeRequestHandling`` is set to 3. It can only be between (inclusive) 262144 (256KB) - 33554432 (32MB).
:sslKeyVersion:             This integer indicates the :ref:`ds-ssl-key-version`
:tenantId:                  The integral, unique identifier of the :ref:`ds-tenant` who owns this :term:`Delivery Service`
:tlsVersions:               An array of explicitly supported :ref:`ds-tls-versions`

	.. versionadded:: 4.0

:topology:          The unique name of the :term:`Topology` that this :term:`Delivery Service` is assigned to
:trRequestHeaders:  If defined, this defines the :ref:`ds-tr-req-headers` used by Traffic Router for this :term:`Delivery Service`
:trResponseHeaders: If defined, this defines the :ref:`ds-tr-resp-headers` used by Traffic Router for this :term:`Delivery Service`
:type:              The :ref:`ds-types` of this :term:`Delivery Service`
:typeId:            The integral, unique identifier of the :ref:`ds-types` of this :term:`Delivery Service`
:xmlId:             This :term:`Delivery Service`'s :ref:`ds-xmlid`

.. code-block:: http
	:caption: Request Example

	POST /api/4.1/deliveryservices HTTP/1.1
	User-Agent: python-requests/2.24.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 1602
	Content-Type: application/json
	Host: trafficops.infra.ciab.test

	{
		"active": false,
		"anonymousBlockingEnabled": false,
		"ccrDnsTtl": null,
		"cdnId": 2,
		"checkPath": null,
		"consistentHashRegex": null,
		"consistentHashQueryParams": [],
		"deepCachingType": "NEVER",
		"displayName": "test",
		"dnsBypassCname": null,
		"dnsBypassIp": null,
		"dnsBypassIp6": null,
		"dnsBypassTtl": null,
		"dscp": 0,
		"ecsEnabled": true,
		"edgeHeaderRewrite": null,
		"firstHeaderRewrite": null,
		"fqPacingRate": null,
		"geoLimit": 0,
		"geoLimitCountries": null,
		"geoLimitRedirectUrl": null,
		"geoProvider": 0,
		"globalMaxMbps": null,
		"globalMaxTps": null,
		"httpBypassFqdn": null,
		"infoUrl": null,
		"initialDispersion": 1,
		"innerHeaderRewrite": null,
		"ipv6RoutingEnabled": false,
		"lastHeaderRewrite": null,
		"logsEnabled": true,
		"longDesc": "A Delivery Service created expressly for API documentation examples",
		"maxDnsAnswers": null,
		"missLat": 0,
		"missLong": 0,
		"maxOriginConnections": 0,
		"maxRequestHeaderBytes": 131072,
		"midHeaderRewrite": null,
		"multiSiteOrigin": false,
		"orgServerFqdn": "http://origin.infra.ciab.test",
		"originShield": null,
		"profileId": null,
		"protocol": 0,
		"qstringIgnore": 0,
		"rangeRequestHandling": 0,
		"regexRemap": null,
		"regional": false,
		"regionalGeoBlocking": false,
		"requiredCapabilities": [],
		"routingName": "test",
		"serviceCategory": null,
		"signed": false,
		"signingAlgorithm": null,
		"rangeSliceBlockSize": null,
		"sslKeyVersion": null,
		"tenant": "root",
		"tenantId": 1,
		"tlsVersions": [
			"1.2",
			"1.3"
		],
		"topology": null,
		"trRequestHeaders": null,
		"trResponseHeaders": null,
		"type": "HTTP",
		"typeId": 1,
		"xmlId": "test"
	}


Response Structure
------------------
:active:                   A boolean that defines :ref:`ds-active`.
:anonymousBlockingEnabled: A boolean that defines :ref:`ds-anonymous-blocking`
:ccrDnsTtl:                 The :ref:`ds-dns-ttl` - named "ccrDnsTtl" for legacy reasons
:cdnId:                     The integral, unique identifier of the :ref:`ds-cdn` to which the :term:`Delivery Service` belongs
:cdnName:                   Name of the :ref:`ds-cdn` to which the :term:`Delivery Service` belongs
:checkPath:                 A :ref:`ds-check-path`
:consistentHashRegex:       A :ref:`ds-consistent-hashing-regex`
:consistentHashQueryParams: An array of :ref:`ds-consistent-hashing-qparams`
:deepCachingType:           The :ref:`ds-deep-caching` setting for this :term:`Delivery Service`
:displayName:               The :ref:`ds-display-name`
:dnsBypassCname:            A :ref:`ds-dns-bypass-cname`
:dnsBypassIp:               A :ref:`ds-dns-bypass-ip`
:dnsBypassIp6:              A :ref:`ds-dns-bypass-ipv6`
:dnsBypassTtl:              The :ref:`ds-dns-bypass-ttl`
:dscp:                      A :ref:`ds-dscp` to be used within the :term:`Delivery Service`
:ecsEnabled:                A boolean that defines the :ref:`ds-ecs` setting on this :term:`Delivery Service`
:edgeHeaderRewrite:         A set of :ref:`ds-edge-header-rw-rules`
:exampleURLs:               An array of :ref:`ds-example-urls`
:firstHeaderRewrite:        A set of :ref:`ds-first-header-rw-rules`
:fqPacingRate:              The :ref:`ds-fqpr`
:geoLimit:                  An integer that defines the :ref:`ds-geo-limit`
:geoLimitCountries:         An array of strings defining the :ref:`ds-geo-limit-countries`
:geoLimitRedirectUrl:       A :ref:`ds-geo-limit-redirect-url`
:geoProvider:               The :ref:`ds-geo-provider`
:globalMaxMbps:             The :ref:`ds-global-max-mbps`
:globalMaxTps:              The :ref:`ds-global-max-tps`
:httpBypassFqdn:            A :ref:`ds-http-bypass-fqdn`
:id:                        An integral, unique identifier for this :term:`Delivery Service`
:infoUrl:                   An :ref:`ds-info-url`
:initialDispersion:         The :ref:`ds-initial-dispersion`
:innerHeaderRewrite:        A set of :ref:`ds-inner-header-rw-rules`
:ipv6RoutingEnabled:        A boolean that defines the :ref:`ds-ipv6-routing` setting on this :term:`Delivery Service`
:lastHeaderRewrite:         A set of :ref:`ds-last-header-rw-rules`
:lastUpdated:               The date and time at which this :term:`Delivery Service` was last updated, in :RFC:`3339` format

	.. versionchanged:: 4.0
		Prior to API version 4.0, this property used :ref:`non-rfc-datetime`.

:logsEnabled: A boolean that defines the :ref:`ds-logs-enabled` setting on this :term:`Delivery Service`
:longDesc:    The :ref:`ds-longdesc` of this :term:`Delivery Service`
:matchList:   The :term:`Delivery Service`'s :ref:`ds-matchlist`

	:pattern:   A regular expression - the use of this pattern is dependent on the ``type`` field (backslashes are escaped)
	:setNumber: An integer that provides explicit ordering of :ref:`ds-matchlist` items - this is used as a priority ranking by Traffic Router, and is not guaranteed to correspond to the ordering of items in the array.
	:type:      The type of match performed using ``pattern``.

:maxDnsAnswers:         The :ref:`ds-max-dns-answers` allowed for this :term:`Delivery Service`
:maxOriginConnections:  The :ref:`ds-max-origin-connections`
:maxRequestHeaderBytes: The :ref:`ds-max-request-header-bytes`
:midHeaderRewrite:      A set of :ref:`ds-mid-header-rw-rules`
:missLat:               The :ref:`ds-geo-miss-default-latitude` used by this :term:`Delivery Service`
:missLong:              The :ref:`ds-geo-miss-default-longitude` used by this :term:`Delivery Service`
:multiSiteOrigin:       A boolean that defines the use of :ref:`ds-multi-site-origin` by this :term:`Delivery Service`
:orgServerFqdn:         The :ref:`ds-origin-url`
:originShield:          A :ref:`ds-origin-shield` string
:profileDescription:    The :ref:`profile-description` of the :ref:`ds-profile` with which this :term:`Delivery Service` is associated
:profileId:             The :ref:`profile-id` of the :ref:`ds-profile` with which this :term:`Delivery Service` is associated
:profileName:           The :ref:`profile-name` of the :ref:`ds-profile` with which this :term:`Delivery Service` is associated
:protocol:              An integral, unique identifier that corresponds to the :ref:`ds-protocol` used by this :term:`Delivery Service`
:qstringIgnore:         An integral, unique identifier that corresponds to the :ref:`ds-qstring-handling` setting on this :term:`Delivery Service`
:rangeRequestHandling:  An integral, unique identifier that corresponds to the :ref:`ds-range-request-handling` setting on this :term:`Delivery Service`
:regexRemap:            A :ref:`ds-regex-remap`
:regional:              A boolean value defining the :ref:`ds-regional` setting on this :term:`Delivery Service`
:regionalGeoBlocking:   A boolean defining the :ref:`ds-regionalgeo` setting on this :term:`Delivery Service`
:remapText:             :ref:`ds-raw-remap`
:requiredCapabilities:  An array of the capabilities that this Delivery Service requires.

	.. versionadded:: 4.1

:serviceCategory:       The name of the :ref:`ds-service-category` with which the :term:`Delivery Service` is associated
:signed:                ``true`` if  and only if ``signingAlgorithm`` is not ``null``, ``false`` otherwise
:signingAlgorithm:      Either a :ref:`ds-signing-algorithm` or ``null`` to indicate URL/URI signing is not implemented on this :term:`Delivery Service`
:rangeSliceBlockSize:   An integer that defines the byte block size for the ATS Slice Plugin. It can only and must be set if ``rangeRequestHandling`` is set to 3.
:sslKeyVersion:         This integer indicates the :ref:`ds-ssl-key-version`
:tenantId:              The integral, unique identifier of the :ref:`ds-tenant` who owns this :term:`Delivery Service`
:tlsVersions:           An array of explicitly supported :ref:`ds-tls-versions`

	.. versionadded:: 4.0

:topology:          The unique name of the :term:`Topology` that this :term:`Delivery Service` is assigned to
:trRequestHeaders:  If defined, this defines the :ref:`ds-tr-req-headers` used by Traffic Router for this :term:`Delivery Service`
:trResponseHeaders: If defined, this defines the :ref:`ds-tr-resp-headers` used by Traffic Router for this :term:`Delivery Service`
:type:              The :ref:`ds-types` of this :term:`Delivery Service`
:typeId:            The integral, unique identifier of the :ref:`ds-types` of this :term:`Delivery Service`
:xmlId:             This :term:`Delivery Service`'s :ref:`ds-xmlid`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 201 Created
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Location: /api/4.1/deliveryservices?id=6
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 07 Jun 2021 23:37:37 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Jun 2021 22:37:37 GMT
	Content-Length: 903

	{ "alerts": [
		{
			"text": "tlsVersions has no effect on 'HTTP' Delivery Services",
			"level": "warning"
		},
		{
			"text": "Delivery Service creation was successful",
			"level": "success"
		}
	],
	"response": [{
		"active": false,
		"anonymousBlockingEnabled": false,
		"ccrDnsTtl": null,
		"cdnId": 2,
		"cdnName": null,
		"checkPath": null,
		"consistentHashQueryParams": [],
		"consistentHashRegex": null,
		"deepCachingType": "NEVER",
		"displayName": "test",
		"dnsBypassCname": null,
		"dnsBypassIp": null,
		"dnsBypassIp6": null,
		"dnsBypassTtl": null,
		"dscp": 0,
		"ecsEnabled": true,
		"edgeHeaderRewrite": null,
		"exampleURLs": [
			"http://test.test.mycdn.ciab.test"
		],
		"firstHeaderRewrite": null,
		"fqPacingRate": null,
		"geoLimit": 0,
		"geoLimitCountries": null,
		"geoLimitRedirectURL": null,
		"geoProvider": 0,
		"globalMaxMbps": null,
		"globalMaxTps": null,
		"httpBypassFqdn": null,
		"id": 6,
		"infoUrl": null,
		"initialDispersion": 1,
		"innerHeaderRewrite": null,
		"ipv6RoutingEnabled": false,
		"lastHeaderRewrite": null,
		"lastUpdated": "2021-06-07T22:37:37.187822Z",
		"logsEnabled": true,
		"longDesc": "A Delivery Service created expressly for API documentation examples",
		"matchList": [
			{
				"type": "HOST_REGEXP",
				"setNumber": 0,
				"pattern": ".*\\.test\\..*"
			}
		],
		"maxDnsAnswers": null,
		"maxOriginConnections": 0,
		"maxRequestHeaderBytes": 131072,
		"midHeaderRewrite": null,
		"missLat": 0,
		"missLong": 0,
		"multiSiteOrigin": false,
		"originShield": null,
		"orgServerFqdn": "http://origin.infra.ciab.test",
		"profileDescription": null,
		"profileId": null,
		"profileName": null,
		"protocol": 0,
		"qstringIgnore": 0,
		"rangeRequestHandling": 0,
		"rangeSliceBlockSize": null,
		"regexRemap": null,
		"regional": false,
		"regionalGeoBlocking": false,
		"remapText": null,
		"requiredCapabilities": [],
		"routingName": "test",
		"serviceCategory": null,
		"signed": false,
		"signingAlgorithm": null,
		"sslKeyVersion": null,
		"tenant": "root",
		"tenantId": 1,
		"tlsVersions": [
			"1.2",
			"1.3"
		],
		"topology": null,
		"trResponseHeaders": null,
		"trRequestHeaders": null,
		"type": "HTTP",
		"typeId": 1,
		"xmlId": "test"
	}]}


.. [#tenancy] Only those :term:`Delivery Services` assigned to :term:`Tenants` that are the requesting user's :term:`Tenant` or children thereof will appear in the output of a ``GET`` request, and the same constraints are placed on the allowed values of the ``tenantId`` field of a ``POST`` request to create a new :term:`Delivery Service`
.. [#geoLimit] These fields must be defined if and only if ``geoLimit`` is non-zero
