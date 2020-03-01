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

.. _to-api-deliveryservice-requests:

****************************
``deliveryservice_requests``
****************************

``GET``
=======
Retrieves :ref:`ds_requests`.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                              |
	+===========+==========+==========================================================================================+
	| assignee  | no       | Filter for :ref:`ds_requests` that are assigned to the user                              |
	|           |          | identified by this username.                                                             |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| assigneeId| no       | Filter for :ref:`ds_requests` that are assigned to the user                              |
	|           |          | identified by this integral, unique identifier                                           |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| author    | no       | Filter for :ref:`ds_requests` submitted by the user                                      |
	|           |          | identified by this username                                                              |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| authorId  | no       | Filter for :ref:`ds_requests` submitted by the user                                      |
	|           |          | identified by this integral, unique identifier                                           |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| changeType| no       | Filter for :ref:`ds_requests` of the change type specified.                              |
	|           |          | Can be ``create``, ``update``, or ``delete``.                                            |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| id        | no       | Filter for the :ref:`Delivery Service Request <ds_requests>` identified by this          |
	|           |          | integral, unique identifier.                                                             |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| status    | no       | Filter for :ref:`ds_requests` whose status is the status                                 |
	|           |          | specified. The status can be ``draft``, ``submitted``, ``pending``, ``rejected``, or     |
	|           |          | ``complete``.                                                                            |
	+-----------+----------+------------------------------------------------------------------------------------------+
	| xmlId     | no       | Filter for :ref:`ds_requests` that have the given                                        |
	|           |          | :ref:`ds-xmlid`.                                                                         |
	+-----------+----------+------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/2.0/deliveryservice_requests?status=draft HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:author:            The author of the Delivery Service Request
:authorId:          The integral, unique identifier assigned to the author
:changeType:        The change type of the :term:`DSR <Delivery Service Request>`. It can be ``create``, ``update``, or ``delete``....
:createdAt:         The date and time at which the :term:`DSR <Delivery Service Request>` was created, in ISO format.
:deliveryService:   The delivery service that the :term:`DSR <Delivery Service Request>` is requesting to update.
:id:                The integral, unique identifier assigned to the :term:`DSR <Delivery Service Request>`
:lastEditedBy:      The username of user who last edited this :term:`DSR <Delivery Service Request>`
:lastEditedById:    The integral, unique identifier assigned to the user who last edited this :term:`DSR <Delivery Service Request>`
:lastUpdated:       The date and time at which the :term:`DSR <Delivery Service Request>` was last updated, in ISO format.
:status:			The status of the request. Can be "draft", "submitted", "rejected", "pending", or "complete".

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 24 Feb 2020 20:14:07 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: UBp3nklJr2x2cAW/TKbhXMVJH6+OduxUaEBGbX4P7IahDk3VkaTd9LsQj01zgFEnZLwHrikpwFfNlUO32RAZOA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 24 Feb 2020 19:14:07 GMT
	Content-Length: 872

	{
		"response": [
			{
				"authorId": 2,
				"author": "admin",
				"changeType": "update",
				"createdAt": "2020-02-24 19:11:12+00",
				"id": 1,
				"lastEditedBy": "admin",
				"lastEditedById": 2,
				"lastUpdated": "2020-02-24 19:11:12+00",
				"deliveryService": {
					"active": false,
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
					"lastUpdated": "0001-01-01 00:00:00+00",
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
					"ecsEnabled": false
				},
				"status": "draft"
			}
		]
	}

.. _to-api-deliveryservice-requests-post:

``POST``
========

.. note:: This route does NOT do the same thing as :ref:`POST deliveryservices/request <to-api-deliveryservices-request>`.

Creates a new :term:`Delivery Service Request`.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object

Request Structure
-----------------
:changeType:		The action that you want to perform on the delivery service. It can be "create", "update", or "delete".
:status:			The status of your request. Can be "draft", "submitted", "rejected", "pending", or "complete".
:deliveryService:	The :term:`Delivery Service` that you have submitted for review as part of this request.

.. code-block:: http
	:caption: Request Example

	POST /api/2.0/deliveryservice_requests HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 1979

	{
		"changeType": "update",
		"status": "draft",
		"deliveryService": {
			"active": false,
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
			"lastUpdated": "2020-02-13 16:43:54+00",
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
			"ecsEnabled": false
		}
	}


Response Structure
------------------
:author:			The author of the Delivery Service Request
:authorId:			The integral, unique identifier assigned to the author
:changeType:		The change type of the :term:`DSR <Delivery Service Request>`. It can be ``create``, ``update``, or ``delete``....
:createdAt:			The date and time at which the :term:`DSR <Delivery Service Request>` was created, in ISO format.
:deliveryService:	The delivery service that the :term:`DSR <Delivery Service Request>` is requesting to update.
:id:				The integral, unique identifier assigned to the :term:`DSR <Delivery Service Request>`
:lastEditedBy:		The username of user who last edited this :term:`DSR <Delivery Service Request>`
:lastEditedById:	The integral, unique identifier assigned to the user who last edited this :term:`DSR <Delivery Service Request>`
:lastUpdated:		The date and time at which the :term:`DSR <Delivery Service Request>` was last updated, in ISO format.
:status:			The status of the request. Can be "draft", "submitted", "rejected", "pending", or "complete".

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 24 Feb 2020 20:11:12 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: aWIrFTwUGnLq56WNZPL/FgOi/NwAVUtOy4iqjFPwx4gj7RMZ6+nd++bQKIiasBl8ytAY0WmFvNnmm30Fq9mLpA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 24 Feb 2020 19:11:12 GMT
	Content-Length: 901

	{
		"alerts": [
			{
				"text": "deliveryservice_request was created.",
				"level": "success"
			}
		],
		"response": {
			"authorId": 2,
			"author": null,
			"changeType": "update",
			"createdAt": null,
			"id": 1,
			"lastEditedBy": null,
			"lastEditedById": 2,
			"lastUpdated": "2020-02-24 19:11:12+00",
			"deliveryService": {
				"active": false,
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
				"lastUpdated": "0001-01-01 00:00:00+00",
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
				"ecsEnabled": false
			},
			"status": "draft"
		}
	}

``PUT``
=======

Updates an existing :ref:`Delivery Service Request <ds_requests>`.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object

Request Structure
-----------------
:author:			The author of the Delivery Service Request
:authorId:			The integral, unique identifier assigned to the author
:changeType:		The change type of the :term:`DSR <Delivery Service Request>`. It can be ``create``, ``update``, or ``delete``....
:createdAt:			The date and time at which the :term:`DSR <Delivery Service Request>` was created, in ISO format.
:deliveryService:	The delivery service that the :term:`DSR <Delivery Service Request>` is requesting to update.
:id:				The integral, unique identifier assigned to the :term:`DSR <Delivery Service Request>`
:lastEditedBy:		The username of user who last edited this :term:`DSR <Delivery Service Request>`
:lastEditedById:	The integral, unique identifier assigned to the user who last edited this :term:`DSR <Delivery Service Request>`
:status:			The status of the request. Can be "draft", "submitted", "rejected", "pending", or "complete".

.. table:: Request Query Parameters

	+-----------+----------+------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                              |
	+===========+==========+==========================================================================================+
	| id        | yes      | The integral, unique identifier of the :ref:`Delivery Service Request <ds_requests>` that|
	|           |          | you want to update.                                                                      |
	+-----------+----------+------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	PUT /api/2.0/deliveryservice_requests?id=1 HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 2256

	{
		"authorId": 2,
		"author": "admin",
		"changeType": "update",
		"createdAt": "2020-02-24 19:11:12+00",
		"id": 1,
		"lastEditedBy": "admin",
		"lastEditedById": 2,
		"lastUpdated": "2020-02-24 19:33:26+00",
		"deliveryService": {
			"active": false,
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
			"lastUpdated": "0001-01-01 00:00:00+00",
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
			"trResponseHeaders": "",
			"trRequestHeaders": null,
			"consistentHashRegex": null,
			"consistentHashQueryParams": [
				"abc",
				"pdq",
				"xxx",
				"zyx"
			],
			"maxOriginConnections": 0,
			"ecsEnabled": false
		},
		"status": "submitted"
	}

Response Structure
------------------
:author:			The author of the Delivery Service Request
:authorId:			The integral, unique identifier assigned to the author
:changeType:		The change type of the :term:`DSR <Delivery Service Request>`. It can be ``create``, ``update``, or ``delete``....
:createdAt:			The date and time at which the :term:`DSR <Delivery Service Request>` was created, in ISO format.
:deliveryService:	The delivery service that the :term:`DSR <Delivery Service Request>` is requesting to update.
:id:				The integral, unique identifier assigned to the :term:`DSR <Delivery Service Request>`
:lastEditedBy:		The username of user who last edited this :term:`DSR <Delivery Service Request>`
:lastEditedById:	The integral, unique identifier assigned to the user who last edited this :term:`DSR <Delivery Service Request>`
:lastUpdated:		The date and time at which the :term:`DSR <Delivery Service Request>` was last updated, in ISO format.
:status:			The status of the request. Can be "draft", "submitted", "rejected", "pending", or "complete".

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 24 Feb 2020 20:36:16 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: +W0vFm96yFkZUJqa0GAX7uzIpRKh/ohyBm0uH3egpiERTcxy5OfVVtoP3h8Ee2teLu8KFooDYXJ6rpQg6UhbNQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 24 Feb 2020 19:36:16 GMT
	Content-Length: 913

	{
		"alerts": [
			{
				"text": "deliveryservice_request was updated.",
				"level": "success"
			}
		],
		"response": {
			"authorId": 0,
			"author": "admin",
			"changeType": "update",
			"createdAt": "0001-01-01 00:00:00+00",
			"id": 1,
			"lastEditedBy": "admin",
			"lastEditedById": 2,
			"lastUpdated": "2020-02-24 19:36:16+00",
			"deliveryService": {
				"active": false,
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
				"lastUpdated": "0001-01-01 00:00:00+00",
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
				"trResponseHeaders": "",
				"trRequestHeaders": null,
				"consistentHashRegex": null,
				"consistentHashQueryParams": [
					"abc",
					"pdq",
					"xxx",
					"zyx"
				],
				"maxOriginConnections": 0,
				"ecsEnabled": false
			},
			"status": "submitted"
		}
	}


``DELETE``
==========
Deletes a :term:`Delivery Service Request`.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  ``undefined``

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

	DELETE /api/2.0/deliveryservice_requests?id=1 HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 0

Response Structure
------------------

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 24 Feb 2020 20:48:55 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: jNCbNo8Tw+JMMaWpAYQgntSXPq2Xuj+n2zSEVRaDQFWMV1SYbT9djes6SPdwiBoKq6W0lNE04hOE92jBVcjtEw==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 24 Feb 2020 19:48:55 GMT
	Content-Length: 96

	{
		"alerts": [
			{
				"text": "deliveryservice_request was deleted.",
				"level": "success"
			}
		]
	}
