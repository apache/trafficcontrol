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

.. _to-api-v4-deliveryservices-id:

***************************
``deliveryservices/{{ID}}``
***************************

``PUT``
=======
Allows users to edit an existing :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE:UPDATE, DELIVERY-SERVICE:READ, CDN:READ, TYPE:READ
:Response Type:  Array (should always have a length of exactly one on success)

Request Structure
-----------------
:active:                   A boolean that defines :ref:`ds-active`.
:anonymousBlockingEnabled: A boolean that defines :ref:`ds-anonymous-blocking`
:ccrDnsTtl:                 The :ref:`ds-dns-ttl` - named "ccrDnsTtl" for legacy reasons
:cdnId:                     The integral, unique identifier of the :ref:`ds-cdn` to which the :term:`Delivery Service` belongs

		.. note:: If the Delivery Service has SSL Keys, then cdnId is not allowed to change as that would invalidate the SSL Key

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
:midHeaderRewrite:          A set of :ref:`ds-mid-header-rw-rules`
:missLat:                   The :ref:`ds-geo-miss-default-latitude` used by this :term:`Delivery Service`
:missLong:                  The :ref:`ds-geo-miss-default-longitude` used by this :term:`Delivery Service`
:multiSiteOrigin:           A boolean that defines the use of :ref:`ds-multi-site-origin` by this :term:`Delivery Service`
:orgServerFqdn:             The :ref:`ds-origin-url`
:originShield:              A :ref:`ds-origin-shield` string
:profileId:                 An optional :ref:`profile-id` of the :ref:`ds-profile` with which this :term:`Delivery Service` will be associated
:protocol:                  An integral, unique identifier that corresponds to the :ref:`ds-protocol` used by this :term:`Delivery Service`
:qstringIgnore:             An integral, unique identifier that corresponds to the :ref:`ds-qstring-handling` setting on this :term:`Delivery Service`
:rangeRequestHandling:      An integral, unique identifier that corresponds to the :ref:`ds-range-request-handling` setting on this :term:`Delivery Service`
:regexRemap:                A :ref:`ds-regex-remap`
:regional:                  A boolean value defining the :ref:`ds-regional` setting on this :term:`Delivery Service`
:regionalGeoBlocking:       A boolean defining the :ref:`ds-regionalgeo` setting on this :term:`Delivery Service`
:remapText:                 :ref:`ds-raw-remap`
:requiredCapabilities:      An array of the capabilities that this Delivery Service requires.

	.. versionadded:: 4.1

:routingName:               The :ref:`ds-routing-name` of this :term:`Delivery Service`

		.. note:: If the Delivery Service has SSL Keys, then ``routingName`` is not allowed to change as that would invalidate the SSL Key

:signed:              ``true`` if  and only if ``signingAlgorithm`` is not ``null``, ``false`` otherwise
:signingAlgorithm:    Either a :ref:`ds-signing-algorithm` or ``null`` to indicate URL/URI signing is not implemented on this :term:`Delivery Service`
:rangeSliceBlockSize: An integer that defines the byte block size for the ATS Slice Plugin. It can only and must be set if ``rangeRequestHandling`` is set to 3. It can only be between (inclusive) 262144 (256KB) - 33554432 (32MB).
:sslKeyVersion:       This integer indicates the :ref:`ds-ssl-key-version`
:tenantId:            The integral, unique identifier of the :ref:`ds-tenant` who owns this :term:`Delivery Service`
:tlsVersions:         An array of explicitly supported :ref:`ds-tls-versions`

	.. versionadded:: 4.0

:topology:          The unique name of the :term:`Topology` that this :term:`Delivery Service` is assigned to
:trRequestHeaders:  If defined, this defines the :ref:`ds-tr-req-headers` used by Traffic Router for this :term:`Delivery Service`
:trResponseHeaders: If defined, this defines the :ref:`ds-tr-resp-headers` used by Traffic Router for this :term:`Delivery Service`
:typeId:            The integral, unique identifier of the :ref:`ds-types` of this :term:`Delivery Service`
:xmlId:             This :term:`Delivery Service`'s :ref:`ds-xmlid`

	.. note:: While this field **must** be present, it is **not** allowed to change; this must be the same as the ``xml_id`` the :term:`Delivery Service` already has. This should almost never be different from the :term:`Delivery Service`'s ``displayName``.


.. code-block:: http
	:caption: Request Example

	PUT /api/4.1/deliveryservices/6 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: python-requests/2.24.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 1585
	Content-Type: application/json

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
		"tlsVersions": null,
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
:requiredCapabilities:  An array of the capabilities that this Delivery Service requires.

	.. versionadded:: 4.1

:remapText:             :ref:`ds-raw-remap`
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

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 08 Jun 2021 00:34:04 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: tTncbRoJR+pyykVbEc6nWyoFnhlJzsbge9hVZfw+WK28rzSGECZ/Q4zXTQtFjHWY5G+0Rk4w9GKrSFK3k+u5Ng==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 07 Jun 2021 23:34:04 GMT
	Content-Length: 840

	{ "alerts": [
		{
			"text": "Delivery Service update was successful",
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
		"exampleURLs": null,
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
		"lastUpdated": "2021-06-07T23:34:04.831215Z",
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
		"tlsVersions": null,
		"topology": null,
		"trResponseHeaders": null,
		"trRequestHeaders": null,
		"type": "HTTP",
		"typeId": 1,
		"xmlId": "test"
	}]}


``DELETE``
==========
Deletes the target :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE:DELETE, DELIVERY-SERVICE:READ, CDN:READ, TYPE:READ
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------------------------+
	| Name | Description                                                                   |
	+======+===============================================================================+
	| ID   | The integral, unique identifier of the :term:`Delivery Service` to be deleted |
	+------+-------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/4.1/deliveryservices/2 HTTP/1.1
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
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
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


.. [#tenancy] Only those :term:`Delivery Services` assigned to :term:`Tenants` that are the requesting user's :term:`Tenant` or children thereof will appear in the output of a ``GET`` request, and the same constraints are placed on the allowed values of the ``tenantId`` field of a ``PUT`` request to update a new :term:`Delivery Service`. Furthermore, the only :term:`Delivery Services` a user may delete are those assigned to a :term:`Tenant` that is either the same :term:`Tenant` as the user's :term:`Tenant`, or a descendant thereof.
.. [#geoLimit] These fields must be defined if and only if ``geoLimit`` is non-zero
