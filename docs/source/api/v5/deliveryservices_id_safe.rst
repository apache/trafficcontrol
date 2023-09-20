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
:Roles Required: None\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE-SAFE:UPDATE, DELIVERY-SERVICE:READ, TYPE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------------------+
	| Name | Description                                                                    |
	+======+================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service` being modified |
	+------+--------------------------------------------------------------------------------+

:displayName: A string that is the :ref:`ds-display-name`
:infoUrl:     An optional\ [#optional]_ string containing the :ref:`ds-info-url`
:longDesc:    An optional\ [#optional]_ string containing the :ref:`ds-longdesc` of this :term:`Delivery Service`

.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/deliveryservices/1/safe HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: access_token=...; mojolicious=...
	Content-Length: 132
	Content-Type: application/json

	{
		"displayName": "test",
		"infoUrl": "this is not even a real URL",
		"longDesc": "this is a description of the delivery service"
	}

Response Structure
------------------
:active:                    The :term:`Delivery Service`'s :ref:`ds-active` state
:anonymousBlockingEnabled:  A boolean that defines :ref:`ds-anonymous-blocking`
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
:logsEnabled:               A boolean that defines the :ref:`ds-logs-enabled` setting on this :term:`Delivery Service`
:longDesc:                  The :ref:`ds-longdesc` of this :term:`Delivery Service`
:matchList:                 The :term:`Delivery Service`'s :ref:`ds-matchlist`

	:pattern:   A regular expression - the use of this pattern is dependent on the ``type`` field (backslashes are escaped)
	:setNumber: An integer that provides explicit ordering of :ref:`ds-matchlist` items - this is used as a priority ranking by Traffic Router, and is not guaranteed to correspond to the ordering of items in the array.
	:type:      The type of match performed using ``pattern``.

:maxDnsAnswers:        The :ref:`ds-max-dns-answers` allowed for this :term:`Delivery Service`
:maxOriginConnections: The :ref:`ds-max-origin-connections`
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
:regional:             A boolean value defining the :ref:`ds-regional` setting on this :term:`Delivery Service`
:regionalGeoBlocking:  A boolean defining the :ref:`ds-regionalgeo` setting on this :term:`Delivery Service`
:remapText:            :ref:`ds-raw-remap`
:signed:               ``true`` if  and only if ``signingAlgorithm`` is not ``null``, ``false`` otherwise
:signingAlgorithm:     Either a :ref:`ds-signing-algorithm` or ``null`` to indicate URL/URI signing is not implemented on this :term:`Delivery Service`
:rangeSliceBlockSize:  An integer that defines the byte block size for the ATS Slice Plugin. It can only and must be set if ``rangeRequestHandling`` is set to 3.
:sslKeyVersion:        This integer indicates the :ref:`ds-ssl-key-version`
:tenantId:             The integral, unique identifier of the :ref:`ds-tenant` who owns this :term:`Delivery Service`
:tlsVersions:           A list of explicitly supported :ref:`ds-tls-versions`
:topology:             The unique name of the :term:`Topology` that this :term:`Delivery Service` is assigned to
:trRequestHeaders:     If defined, this defines the :ref:`ds-tr-req-headers` used by Traffic Router for this :term:`Delivery Service`
:trResponseHeaders:    If defined, this defines the :ref:`ds-tr-resp-headers` used by Traffic Router for this :term:`Delivery Service`
:type:                 The :ref:`ds-types` of this :term:`Delivery Service`
:typeId:               The integral, unique identifier of the :ref:`ds-types` of this :term:`Delivery Service`
:xmlId:                This :term:`Delivery Service`'s :ref:`ds-xmlid`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Thu, 29 Sep 2022 23:30:08 GMT; Max-Age=3600; HttpOnly, access_token=...; Path=/; Expires=Thu, 29 Sep 2022 23:30:08 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 29 Sep 2022 22:30:08 GMT
	Content-Length: 904

	{ "alerts": [{
		"text": "Delivery Service safe update successful.",
		"level": "success"
	}],
	"response": {
		"active": "ACTIVE",
		"anonymousBlockingEnabled": false,
		"ccrDnsTtl": null,
		"cdnId": 2,
		"cdnName": "CDN-in-a-Box",
		"checkPath": null,
		"consistentHashQueryParams": [
			"abc",
			"pdq",
			"xxx",
			"zyx"
		],
		"consistentHashRegex": null,
		"deepCachingType": "NEVER",
		"displayName": "test",
		"dnsBypassCname": null,
		"dnsBypassIp": null,
		"dnsBypassIp6": null,
		"dnsBypassTtl": null,
		"dscp": 0,
		"ecsEnabled": false,
		"edgeHeaderRewrite": null,
		"exampleURLs": [
			"http://video.demo1.mycdn.ciab.test",
			"https://video.demo1.mycdn.ciab.test"
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
		"infoUrl": "this is not even a real URL",
		"initialDispersion": 1,
		"innerHeaderRewrite": null,
		"ipv6RoutingEnabled": true,
		"lastHeaderRewrite": null,
		"lastUpdated": "2022-09-29T22:30:08.574513Z",
		"logsEnabled": true,
		"longDesc": "this is a description of the delivery service",
		"matchList": [
			{
				"type": "HOST_REGEXP",
				"setNumber": 0,
				"pattern": ".*\\.demo1\\..*"
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
		"routingName": "video",
		"serviceCategory": null,
		"signed": false,
		"signingAlgorithm": null,
		"sslKeyVersion": 1,
		"tenant": "root",
		"tenantId": 1,
		"tlsVersions": null,
		"topology": "demo1-top",
		"trResponseHeaders": null,
		"trRequestHeaders": null,
		"type": "HTTP",
		"typeId": 1,
		"xmlId": "demo1"
	}}

.. [#tenancy] Only those :term:`Delivery Services` assigned to :term:`Tenants` that are the requesting user's :term:`Tenant` or children thereof may be modified with this endpoint.
.. [#optional] If ``infoUrl`` is not present in the request body it is *implicitly set to* ``null``. If ``longDesc`` isn't in the request body, it is *implicitly set to* an empty string.
