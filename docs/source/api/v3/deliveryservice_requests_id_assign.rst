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

.. _to-api-deliveryservice_requests-id-assign:

******************************************
``deliveryservice_requests/{{ID}}/assign``
******************************************

.. caution:: In many cases, it's much easier to simply use :ref:`to-api-deliveryservice_requests`.

``GET``
=======
Retrieves the :ref:`dsr-assignee` of a particular :term:`DSR`.

:Auth. Required: Yes
:Roles Required: None
:Response Type: Object (string)

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------------------------+
	| Name | Description                                                                             |
	+======+=========================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service Request` being inspected |
	+------+-----------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/deliveryservice_requests/6/assign HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
The response is the username of the :term:`DSR`'s :ref:`dsr-assignee`.

.. code-block:: http
	:caption: Request Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 25 Sep 2020 07:41:52 GMT
	Content-Length: 45

	{"response": "admin"}

``PUT``
=======
Assign a :term:`Delivery Service Request` to a user.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------------------------------------------------------------------------------------+
	| Name | Description                                                                            |
	+======+========================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service Request` being assigned |
	+------+----------------------------------------------------------------------------------------+

:assignee:   The username of the user to whom the :term:`Delivery Service Request` is assigned.
:assigneeId: The integral, unique identifier assigned to the :term:`DSR`.

It is not required to send both of these; either property is sufficient to determine an :ref:`dsr-assignee`. In most cases, it's easier to use just `assignee`. If both *are* given, then `assigneeId` will take precedence in the event that the two properties do not refer to the same user.

.. code-block:: http
	:caption: Request Example

	PUT /api/3.0/deliveryservice_requests/6/assign HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 21
	Content-Type: application/json

	{"assignee": "admin"}


Response Structure
------------------
The response contains a full representation of the newly assigned :term:`Delivery Service Request`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 25 Sep 2020 07:01:24 GMT
	Content-Length: 1145

	{ "alerts": [
		{
			"text": "Changed assignee of 'demo1' Delivery Service Request to 'admin'",
			"level": "success"
		}
	],
	"response": {
		"assignee": "admin",
		"author": "admin",
		"changeType": "update",
		"createdAt": "2020-09-25T06:52:23.758877Z",
		"id": 6,
		"lastEditedBy": "admin",
		"lastUpdated": "2020-09-25T07:01:24.600029Z",
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
