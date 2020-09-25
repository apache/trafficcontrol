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

.. _to-api-deliveryservice_requests:

****************************
``deliveryservice_requests``
****************************

``GET``
=======
Retrieves representations of :term:`Delivery Service Requests`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                              |
	+===========+==========+==========================================================================================+
	| assignee  | no       | Filter for :term:`Delivery Service Requests` that are assigned to the user               |
	|           |          | identified by this username.                                                             |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| assigneeId| no       | Filter for :term:`Delivery Service Requests` that are assigned to the user               |
	|           |          | identified by this integral, unique identifier                                           |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| author    | no       | Filter for :term:`Delivery Service Requests` submitted by the user                       |
	|           |          | identified by this username                                                              |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| authorId  | no       | Filter for :term:`Delivery Service Requests` submitted by the user                       |
	|           |          | identified by this integral, unique identifier                                           |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| changeType| no       | Filter for :term:`Delivery Service Requests` of the change type specified.               |
	|           |          | Can be ``create``, ``update``, or ``delete``.                                            |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| id        | no       | Filter for the :term:`Delivery Service Request` identified by this                       |
	|           |          | integral, unique identifier.                                                             |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| status    | no       | Filter for :term:`Delivery Service Requests` whose status is the status                  |
	|           |          | specified. The status can be ``draft``, ``submitted``, ``pending``, ``rejected``, or     |
	|           |          | ``complete``.                                                                            |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| xmlId     | no       | Filter for :term:`Delivery Service Requests` that have the given                         |
	|           |          | :ref:`ds-xmlid`.                                                                         |
	+-----------+----------+------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/deliveryservice_requests?status=draft HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
The response is an array of representations of :term:`Delivery Service Requests`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Last-Modified: Mon, 01 Jan 0001 00:00:00 UTC
	Set-Cookie: mojolicious=...
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 25 Sep 2020 07:54:14 GMT
	Content-Length: 1092

	{ "response": [
		{
			"assignee": "admin",
			"author": "admin",
			"changeType": "update",
			"createdAt": "2020-09-25T06:52:23.758877Z",
			"id": 2,
			"lastEditedBy": "admin",
			"lastUpdated": "2020-09-25T07:13:28.753352Z",
			"original": {
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
				"lastUpdated": "2020-09-25 02:09:54+00",
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
				"protocol": 2,
				"qstringIgnore": 0,
				"rangeRequestHandling": 0,
				"regexRemap": null,
				"regionalGeoBlocking": false,
				"remapText": null,
				"routingName": "video",
				"signed": false,
				"sslKeyVersion": 1,
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
				"maxOriginConnections": 0,
				"ecsEnabled": false,
				"rangeSliceBlockSize": null,
				"topology": "demo1-top",
				"firstHeaderRewrite": null,
				"innerHeaderRewrite": null,
				"lastHeaderRewrite": null,
				"serviceCategory": null
			},
			"requested": {
				"active": true,
				"anonymousBlockingEnabled": false,
				"cacheurl": null,
				"ccrDnsTtl": 30,
				"cdnId": 2,
				"cdnName": null,
				"checkPath": null,
				"displayName": "Demo 1 but modified by a DSR",
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
				"initialDispersion": 3,
				"ipv6RoutingEnabled": null,
				"lastUpdated": null,
				"logsEnabled": false,
				"longDesc": "long desc",
				"longDesc1": null,
				"longDesc2": null,
				"matchList": null,
				"maxDnsAnswers": null,
				"midHeaderRewrite": null,
				"missLat": null,
				"missLong": null,
				"multiSiteOrigin": null,
				"originShield": null,
				"orgServerFqdn": null,
				"profileDescription": null,
				"profileId": null,
				"profileName": null,
				"protocol": null,
				"qstringIgnore": null,
				"rangeRequestHandling": null,
				"regexRemap": null,
				"regionalGeoBlocking": false,
				"remapText": null,
				"routingName": "cdn",
				"signed": false,
				"sslKeyVersion": null,
				"tenantId": 1,
				"type": null,
				"typeId": 8,
				"xmlId": "demo1",
				"exampleURLs": null,
				"deepCachingType": "NEVER",
				"fqPacingRate": null,
				"signingAlgorithm": null,
				"tenant": null,
				"trResponseHeaders": null,
				"trRequestHeaders": null,
				"consistentHashRegex": null,
				"consistentHashQueryParams": null,
				"maxOriginConnections": 0,
				"ecsEnabled": false,
				"rangeSliceBlockSize": null,
				"topology": null,
				"firstHeaderRewrite": null,
				"innerHeaderRewrite": null,
				"lastHeaderRewrite": null,
				"serviceCategory": null
			},
			"status": "submitted"
		}
	]}

.. _to-api-deliveryservice-requests-post:

``POST``
========
Creates a new :term:`Delivery Service Request`. "Closed" :term:`Delivery Service Requests` cannot be created, an existing :term:`Delivery Service Request` must be placed into a closed :ref:`dsr-status`. A :term:`Delivery Service Request` to create, modify or delete a :term:`Delivery Service` cannot be created if an open :term:`Delivery Service Request` exists for a :term:`Delivery Service` with the same :ref:`ds-xmlid`. Because of this, :term:`Delivery Service Requests` cannot be used to change a :term:`Delivery Service`'s :ref:`ds-xmlid`.

.. warning:: This route does NOT do the same thing as :ref:`POST deliveryservices/request <to-api-deliveryservices-request>`.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object

Request Structure
-----------------
The request must be a well-formed representation of a :term:`Delivery Service Request`, without any response-only fields, of course.

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/deliveryservice_requests HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 637
	Content-Type: application/json

	{
		"changeType": "update",
		"requested": {
			"active": true,
			"cdnId": 2,
			"ccrDnsTtl": 30,
			"deepCachingType": "NEVER",
			"displayName": "Demo 1 but modified by a DSR",
			"dscp": 0,
			"geoLimit": 0,
			"geoProvider": 0,
			"initialDispersion": 3,
			"logsEnabled": false,
			"longDesc": "long desc",
			"regionalGeoBlocking": false,
			"tenantId": 1,
			"typeId": 8,
			"xmlId": "demo1",
			"id": 1
		},
		"status": "draft"
	}


Response Structure
------------------
The response will be a representation of the created :term:`Delivery Service Request`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 201 Created
	Content-Encoding: gzip
	Content-Type: application/json
	Location: /api/3.0/deliveryservice_requests/2
	Set-Cookie: mojolicious=...
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 25 Sep 2020 02:38:04 GMT
	Content-Length: 1074

	{ "alerts": [{
		"level": "success",
		"text": "Delivery Service request created"
	}],
	"response": {
		"assignee": null,
		"author": "admin",
		"changeType": "update",
		"createdAt": "2020-09-25T02:38:04.180237Z",
		"id": 2,
		"lastEditedBy": "admin",
		"lastUpdated": "2020-09-25T02:38:04.180237Z",
		"original": {
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
			"lastUpdated": "2020-09-25 02:09:54+00",
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
			"protocol": 2,
			"qstringIgnore": 0,
			"rangeRequestHandling": 0,
			"regexRemap": null,
			"regionalGeoBlocking": false,
			"remapText": null,
			"routingName": "video",
			"signed": false,
			"sslKeyVersion": 1,
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
			"maxOriginConnections": 0,
			"ecsEnabled": false,
			"rangeSliceBlockSize": null,
			"topology": "demo1-top",
			"firstHeaderRewrite": null,
			"innerHeaderRewrite": null,
			"lastHeaderRewrite": null,
			"serviceCategory": null
		},
		"requested": {
			"active": true,
			"anonymousBlockingEnabled": false,
			"cacheurl": null,
			"ccrDnsTtl": 30,
			"cdnId": 2,
			"cdnName": null,
			"checkPath": null,
			"displayName": "Demo 1 but modified by a DSR",
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
			"initialDispersion": 3,
			"ipv6RoutingEnabled": null,
			"lastUpdated": null,
			"logsEnabled": false,
			"longDesc": "long desc",
			"longDesc1": null,
			"longDesc2": null,
			"matchList": null,
			"maxDnsAnswers": null,
			"midHeaderRewrite": null,
			"missLat": null,
			"missLong": null,
			"multiSiteOrigin": null,
			"originShield": null,
			"orgServerFqdn": null,
			"profileDescription": null,
			"profileId": null,
			"profileName": null,
			"protocol": null,
			"qstringIgnore": null,
			"rangeRequestHandling": null,
			"regexRemap": null,
			"regionalGeoBlocking": false,
			"remapText": null,
			"routingName": "cdn",
			"signed": false,
			"sslKeyVersion": null,
			"tenantId": 1,
			"type": null,
			"typeId": 8,
			"xmlId": "demo1",
			"exampleURLs": null,
			"deepCachingType": "NEVER",
			"fqPacingRate": null,
			"signingAlgorithm": null,
			"tenant": null,
			"trResponseHeaders": null,
			"trRequestHeaders": null,
			"consistentHashRegex": null,
			"consistentHashQueryParams": null,
			"maxOriginConnections": 0,
			"ecsEnabled": false,
			"rangeSliceBlockSize": null,
			"topology": null,
			"firstHeaderRewrite": null,
			"innerHeaderRewrite": null,
			"lastHeaderRewrite": null,
			"serviceCategory": null
		},
		"status": "draft"
	}}

``PUT``
=======
Updates an existing :term:`Delivery Service Request`. Note that "closed" :term:`Delivery Service Requests` are uneditable.

.. seealso:: The proper way to change a :term:`Delivery Service Request`'s :ref:`dsr-status` is by using the :ref:`to-api-deliveryservice_requests-id-status` endpoint's ``PUT`` handler.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+----------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                  |
	+===========+==========+==============================================================================================+
	| id        | yes      | The integral, unique identifier of the :term:`Delivery Service Request` that will be updated |
	+-----------+----------+----------------------------------------------------------------------------------------------+

The request body must be a representation of a :term:`Delivery Service Request` without any response-only fields.

.. code-block:: http
	:caption: Request Example

	PUT /api/3.0/deliveryservice_requests?id=2 HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 657
	Content-Type: application/json

	{
		"changeType": "update",
		"requested": {
			"active": true,
			"cdnId": 2,
			"ccrDnsTtl": 30,
			"deepCachingType": "NEVER",
			"displayName": "Demo 1 but I modified the DSR",
			"dscp": 0,
			"geoLimit": 0,
			"geoProvider": 0,
			"initialDispersion": 3,
			"logsEnabled": false,
			"longDesc": "long desc",
			"regionalGeoBlocking": false,
			"tenantId": 1,
			"typeId": 8,
			"xmlId": "demo1",
            "id": 1
		},
		"status": "draft"
	}

Response Structure
------------------
The response is a full representation of the edited :term:`Delivery Service Request`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 25 Sep 2020 06:23:30 GMT
	Content-Length: 1127

	{ "alerts": [
		{
			"text": "Delivery Service Request #2 updated",
			"level": "success"
		}
	],
	"response": {
		"assignee": null,
		"author": "",
		"changeType": "update",
		"createdAt": "2020-09-25T06:23:30.683058Z",
		"id": null,
		"lastEditedBy": "admin",
		"lastUpdated": "2020-09-25T02:38:04.180237Z",
		"original": {
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
			"lastUpdated": "2020-09-25 02:09:54+00",
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
			"protocol": 2,
			"qstringIgnore": 0,
			"rangeRequestHandling": 0,
			"regexRemap": null,
			"regionalGeoBlocking": false,
			"remapText": null,
			"routingName": "video",
			"signed": false,
			"sslKeyVersion": 1,
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
			"maxOriginConnections": 0,
			"ecsEnabled": false,
			"rangeSliceBlockSize": null,
			"topology": "demo1-top",
			"firstHeaderRewrite": null,
			"innerHeaderRewrite": null,
			"lastHeaderRewrite": null,
			"serviceCategory": null
		},
		"requested": {
			"active": true,
			"anonymousBlockingEnabled": false,
			"cacheurl": null,
			"ccrDnsTtl": 30,
			"cdnId": 2,
			"cdnName": null,
			"checkPath": null,
			"displayName": "Demo 1 but I modified the DSR",
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
			"initialDispersion": 3,
			"ipv6RoutingEnabled": null,
			"lastUpdated": null,
			"logsEnabled": false,
			"longDesc": "long desc",
			"longDesc1": null,
			"longDesc2": null,
			"matchList": null,
			"maxDnsAnswers": null,
			"midHeaderRewrite": null,
			"missLat": null,
			"missLong": null,
			"multiSiteOrigin": null,
			"originShield": null,
			"orgServerFqdn": null,
			"profileDescription": null,
			"profileId": null,
			"profileName": null,
			"protocol": null,
			"qstringIgnore": null,
			"rangeRequestHandling": null,
			"regexRemap": null,
			"regionalGeoBlocking": false,
			"remapText": null,
			"routingName": "cdn",
			"signed": false,
			"sslKeyVersion": null,
			"tenantId": 1,
			"type": null,
			"typeId": 8,
			"xmlId": "demo1",
			"exampleURLs": null,
			"deepCachingType": "NEVER",
			"fqPacingRate": null,
			"signingAlgorithm": null,
			"tenant": null,
			"trResponseHeaders": null,
			"trRequestHeaders": null,
			"consistentHashRegex": null,
			"consistentHashQueryParams": null,
			"maxOriginConnections": 0,
			"ecsEnabled": false,
			"rangeSliceBlockSize": null,
			"topology": null,
			"firstHeaderRewrite": null,
			"innerHeaderRewrite": null,
			"lastHeaderRewrite": null,
			"serviceCategory": null
		},
		"status": "draft"
	}}

``DELETE``
==========
Deletes a :term:`Delivery Service Request`.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                              |
	+===========+==========+==========================================================================================+
	| id        | yes      | The integral, unique identifier of the :ref:`Delivery Service Request <ds_requests>` that|
	|           |          | you want to delete.                                                                      |
	+-----------+----------+------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/3.0/deliveryservice_requests?id=1 HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 0

Response Structure
------------------
The response is a full representation of the deleted :term:`Delivery Service Request`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 25 Sep 2020 06:39:05 GMT
	Content-Length: 1115

	{ "alerts": [
		{
			"text": "Delivery Service Request #1 deleted",
			"level": "success"
		}
	],
	"response": {
		"assignee": null,
		"author": "admin",
		"changeType": "update",
		"createdAt": "2020-09-25T06:38:55.402641Z",
		"id": 1,
		"lastEditedBy": "admin",
		"lastUpdated": "2020-09-25T06:38:55.402641Z",
		"original": {
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
			"lastUpdated": "2020-09-25 02:09:54+00",
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
			"protocol": 2,
			"qstringIgnore": 0,
			"rangeRequestHandling": 0,
			"regexRemap": null,
			"regionalGeoBlocking": false,
			"remapText": null,
			"routingName": "video",
			"signed": false,
			"sslKeyVersion": 1,
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
			"maxOriginConnections": 0,
			"ecsEnabled": false,
			"rangeSliceBlockSize": null,
			"topology": "demo1-top",
			"firstHeaderRewrite": null,
			"innerHeaderRewrite": null,
			"lastHeaderRewrite": null,
			"serviceCategory": null
		},
		"requested": {
			"active": true,
			"anonymousBlockingEnabled": false,
			"cacheurl": null,
			"ccrDnsTtl": 30,
			"cdnId": 2,
			"cdnName": null,
			"checkPath": null,
			"displayName": "Demo 1 but modified by a DSR",
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
			"initialDispersion": 3,
			"ipv6RoutingEnabled": null,
			"lastUpdated": null,
			"logsEnabled": false,
			"longDesc": "long desc",
			"longDesc1": null,
			"longDesc2": null,
			"matchList": null,
			"maxDnsAnswers": null,
			"midHeaderRewrite": null,
			"missLat": null,
			"missLong": null,
			"multiSiteOrigin": null,
			"originShield": null,
			"orgServerFqdn": null,
			"profileDescription": null,
			"profileId": null,
			"profileName": null,
			"protocol": null,
			"qstringIgnore": null,
			"rangeRequestHandling": null,
			"regexRemap": null,
			"regionalGeoBlocking": false,
			"remapText": null,
			"routingName": "cdn",
			"signed": false,
			"sslKeyVersion": null,
			"tenantId": 1,
			"type": null,
			"typeId": 8,
			"xmlId": "demo1",
			"exampleURLs": null,
			"deepCachingType": "NEVER",
			"fqPacingRate": null,
			"signingAlgorithm": null,
			"tenant": null,
			"trResponseHeaders": null,
			"trRequestHeaders": null,
			"consistentHashRegex": null,
			"consistentHashQueryParams": null,
			"maxOriginConnections": 0,
			"ecsEnabled": false,
			"rangeSliceBlockSize": null,
			"topology": null,
			"firstHeaderRewrite": null,
			"innerHeaderRewrite": null,
			"lastHeaderRewrite": null,
			"serviceCategory": null
		},
		"status": "draft"
	}}
