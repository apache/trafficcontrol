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
:Roles Required: "admin" or "operations"\ [#tenancy]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------------------+
	| Name |                      Description                                               |
	+======+================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service` being modified |
	+------+--------------------------------------------------------------------------------+

:displayName: The :ref:`ds-display-name`
:infoUrl:     An :ref:`ds-info-url`
:longDesc:    The :ref:`ds-longdesc` of this :term:`Delivery Service`
:longDesc1:   The :ref:`ds-longdesc2` of this :term:`Delivery Service`
:longDesc2:   The :ref:`ds-longdesc3` of this :term:`Delivery Service`

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
		"longDesc": "A Delivery Service created for the CDN-in-a-Box project",
		"longDesc1": null,
		"longDesc2": null
	}


Response Structure
------------------
:active:                   A boolean that defines :ref:`ds-active`.
:anonymousBlockingEnabled: A boolean that defines :ref:`ds-anonymous-blocking`
:cacheurl:                 A :ref:`ds-cacheurl`

	.. deprecated:: ATCv3.0
		This field has been deprecated in Traffic Control 3.x and is subject to removal in Traffic Control 4.x or later

:ccrDnsTtl:           The :ref:`ds-dns-ttl` - named "ccrDnsTtl" for legacy reasons
:cdnId:               The integral, unique identifier of the :ref:`ds-cdn` to which the :term:`Delivery Service` belongs
:cdnName:             Name of the :ref:`ds-cdn` to which the :term:`Delivery Service` belongs
:checkPath:           A :ref:`ds-check-path`
:consistentHashRegex: A :ref:`ds-consistent-hashing-regex`

	.. versionadded:: 1.4

:consistentHashQueryParams: An array of :ref:`ds-consistent-hashing-qparams`

	.. versionadded:: 1.4

:deepCachingType:     The :ref:`ds-deep-caching` setting for this :term:`Delivery Service`

	.. versionadded:: 1.3

:displayName:       The :ref:`ds-display-name`
:dnsBypassCname:    A :ref:`ds-dns-bypass-cname`
:dnsBypassIp:       A :ref:`ds-dns-bypass-ip`
:dnsBypassIp6:      A :ref:`ds-dns-bypass-ipv6`
:dnsBypassTtl:      The :ref:`ds-dns-bypass-ttl`
:dscp:              A :ref:`ds-dscp` to be used within the :term:`Delivery Service`
:edgeHeaderRewrite: A set of :ref:`ds-edge-header-rw-rules`
:exampleURLs:       An array of :ref:`ds-example-urls`
:fqPacingRate:      The :ref:`ds-fqpr`

	.. versionadded:: 1.3

:geoLimit:            An integer that defines the :ref:`ds-geo-limit`
:geoLimitCountries:   A string containing a comma-separated list defining the :ref:`ds-geo-limit-countries`
:geoLimitRedirectUrl: A :ref:`ds-geo-limit-redirect-url`
:geoProvider:         The :ref:`ds-geo-provider`
:globalMaxMbps:       The :ref:`ds-global-max-mbps`
:globalMaxTps:        The :ref:`ds-global-max-tps`
:httpBypassFqdn:      A :ref:`ds-http-bypass-fqdn`
:id:                  An integral, unique identifier for this :term:`Delivery Service`
:infoUrl:             An :ref:`ds-info-url`
:initialDispersion:  The :ref:`ds-initial-dispersion`
:ipv6RoutingEnabled: A boolean that defines the :ref:`ds-ipv6-routing` setting on this :term:`Delivery Service`
:lastUpdated:        The date and time at which this :term:`Delivery Service` was last updated, in :rfc:`3339` format
:logsEnabled:        A boolean that defines the :ref:`ds-logs-enabled` setting on this :term:`Delivery Service`
:longDesc:           The :ref:`ds-longdesc` of this :term:`Delivery Service`
:longDesc1:          The :ref:`ds-longdesc2` of this :term:`Delivery Service`
:longDesc2:          The :ref:`ds-longdesc3` of this :term:`Delivery Service`
:matchList:          The :term:`Delivery Service`'s :ref:`ds-matchlist`

	:pattern:   A regular expression - the use of this pattern is dependent on the ``type`` field (backslashes are escaped)
	:setNumber: An integer that provides explicit ordering of :ref:`ds-matchlist` items - this is used as a priority ranking by Traffic Router, and is not guaranteed to correspond to the ordering of items in the array.
	:type:      The type of match performed using ``pattern``.

:maxDnsAnswers:        The :ref:`ds-max-dns-answers` allowed for this :term:`Delivery Service`
:maxOriginConnections: The :ref:`ds-max-origin-connections`

	.. versionadded:: 1.4

:midHeaderRewrite:     A set of :ref:`ds-mid-header-rw-rules`
:missLat:              The :ref:`ds-geo-miss-default-latitude` used by this :term:`Delivery Service`
:missLong:             The :ref:`ds-geo-miss-default-longitude` used by this :term:`Delivery Service`
:multiSiteOrigin:      A boolean that defines the use of :ref:`ds-multi-site-origin` by this :term:`Delivery Service`
:orgServerFqdn:        The :ref:`ds-origin-url`
:originShield:         A :ref:`ds-origin-shield` string
:profileDescription:   The :ref:`profile-description` of the :ref:`ds-profile` with which this :term:`Delivery Service` is associated
:profileId:            The :ref:`profile-id` of the :ref:`ds-profile` with which this :term:`Delivery Service` is associated
:profileName:          The :ref:`profile-name` of the :ref:`ds-profile` with which this :term:`Delivery Service` is associated
:protocol:             An integral, unique identifier that corresponds to the :ref:`ds-protocol` used by this :term:`Delivery Service`
:qstringIgnore:        An integral, unique identifier that corresponds to the :ref:`ds-qstring-handling` setting on this :term:`Delivery Service`
:rangeRequestHandling: An integral, unique identifier that corresponds to the :ref:`ds-range-request-handling` setting on this :term:`Delivery Service`
:regexRemap:           A :ref:`ds-regex-remap`
:regionalGeoBlocking:  A boolean defining the :ref:`ds-regionalgeo` setting on this :term:`Delivery Service`
:remapText:            :ref:`ds-raw-remap`
:signed:               ``true`` if  and only if ``signingAlgorithm`` is not ``null``, ``false`` otherwise
:signingAlgorithm:     Either a :ref:`ds-signing-algorithm` or ``null`` to indicate URL/URI signing is not implemented on this :term:`Delivery Service`

	.. versionadded:: 1.3

:sslKeyVersion: This integer indicates the :ref:`ds-ssl-key-version`
:tenantId:      The integral, unique identifier of the :ref:`ds-tenant` who owns this :term:`Delivery Service`

	.. versionadded:: 1.3

:trRequestHeaders: If defined, this defines the :ref:`ds-tr-req-headers` used by Traffic Router for this :term:`Delivery Service`

	.. versionadded:: 1.3

:trResponseHeaders: If defined, this defines the :ref:`ds-tr-resp-headers` used by Traffic Router for this :term:`Delivery Service`

	.. versionadded:: 1.3

:type:   The :ref:`ds-types` of this :term:`Delivery Service`
:typeId: The integral, unique identifier of the :ref:`ds-types` of this :term:`Delivery Service`
:xmlId:  This :term:`Delivery Service`'s :ref:`ds-xmlid`

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

	{ "response": [{
		"active": true,
		"anonymousBlockingEnabled": false,
		"cacheurl": null,
		"ccrDnsTtl": null,
		"cdnId": 2,
		"cdnName": "CDN-in-a-Box",
		"checkPath": null,
		"displayName": "demo",
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
		"infoUrl": "www.info.com",
		"initialDispersion": 1,
		"ipv6RoutingEnabled": true,
		"lastUpdated": "2019-05-15 14:32:05+00",
		"logsEnabled": true,
		"longDesc": "A Delivery Service created for the CDN-in-a-Box project",
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
		"protocol": 2,
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
			"http://video.demo1.mycdn.ciab.test",
			"https://video.demo1.mycdn.ciab.test"
		],
		"deepCachingType": "NEVER",
		"fqPacingRate": null,
		"signingAlgorithm": null,
		"tenant": "root",
		"trResponseHeaders": null,
		"trRequestHeaders": null,
		"consistentHashRegex": null,
		"consistentHashQueryParams": [
			"abc",
			"pdq",
			"xxx",
			"zyx"
		],
		"maxOriginConnections": 0
	}]}


.. [#tenancy] Only those :term:`Delivery Services` assigned to :term:`Tenants` that are the requesting user's :term:`Tenant` or children thereof may be modified with this endpoint.
