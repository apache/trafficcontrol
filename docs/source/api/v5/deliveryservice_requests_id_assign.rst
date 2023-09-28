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
Assign a :term:`Delivery Service Request` to a user.

``GET``
=======
:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: DS-REQUEST:READ, USER:READ
:Response Type:  Object (string)

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------------------------------------------------+
	| Name | Description                                                                                                     |
	+======+=================================================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service Request` for which assignment is being retrieved |
	+------+-----------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/deliveryservice_requests/1/assign HTTP/1.1
	User-Agent: python-requests/2.24.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
The response is the username of the user to whom the :term:`Delivery Service Request` is assigned - or ``null`` if it is unassigned.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 02 Feb 2021 22:48:48 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 02 Feb 2021 21:48:48 GMT
	Content-Length: 45

	{ "response": "admin" }


``PUT``
=======
:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: DS-REQUEST:UPDATE, DS-REQUEST:READ, USER:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------------------------------------------------------------------------------------+
	| Name | Description                                                                            |
	+======+========================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service Request` being assigned |
	+------+----------------------------------------------------------------------------------------+

:assignee:   The username of the user to whom the :term:`Delivery Service Request` is assigned
:assigneeId: The integral, unique identifier of the user to whom the :term:`Delivery Service Request` is assigned

	.. note::
		It is not required to send both of these; either property is sufficient to determine an :ref:`dsr-assignee`. In most cases, it's easier to use just `assignee`. If both *are* given, then `assigneeId` will take precedence in the event that the two properties do not refer to the same user. Sending a request that sets the assignee to ``null`` un-assigns the :term:`DSR` from any assignees it previously had\ [#implicit-null]_.

.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/deliveryservice_requests/1/assign HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: access_token=...; mojolicious=...
	Content-Length: 21

	{"assignee": "admin"}

Response Structure
------------------
The response contains a full representation of the newly assigned :term:`Delivery Service Request`.

.. versionchanged:: 5.0
	Prior to version 5.0 of the API, the ``lastUpdated`` field was in :ref:`non-rfc-datetime`.

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
	Set-Cookie: mojolicious=...; Path=/; Expires=Thu, 29 Sep 2022 23:23:28 GMT; Max-Age=3600; HttpOnly, access_token=...; Path=/; Expires=Thu, 29 Sep 2022 23:23:28 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 29 Sep 2022 22:23:28 GMT
	Content-Length: 1159

	{ "alerts": [{
		"text": "Changed assignee of 'demo1' Delivery Service Request to 'admin'",
		"level": "success"
	}],
	"response": {
		"assignee": "admin",
		"author": "admin",
		"changeType": "update",
		"createdAt": "2022-09-29T22:07:15.008503-6:00",
		"id": 1,
		"lastEditedBy": "admin",
		"lastUpdated": "2022-09-29T22:23:28.034845-6:00",
		"original": {
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
			"displayName": "Demo 1",
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
			"infoUrl": null,
			"initialDispersion": 1,
			"innerHeaderRewrite": null,
			"ipv6RoutingEnabled": true,
			"lastHeaderRewrite": null,
			"lastUpdated": "2022-09-29T20:58:53.07251-6:00",
			"logsEnabled": true,
			"longDesc": "Apachecon North America 2018",
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
		},
		"requested": {
			"active": "INACTIVE",
			"anonymousBlockingEnabled": false,
			"ccrDnsTtl": 30,
			"cdnId": 2,
			"cdnName": null,
			"checkPath": null,
			"consistentHashQueryParams": null,
			"consistentHashRegex": null,
			"deepCachingType": "NEVER",
			"displayName": "Demo 1 but I modified the DSR",
			"dnsBypassCname": null,
			"dnsBypassIp": null,
			"dnsBypassIp6": null,
			"dnsBypassTtl": null,
			"dscp": 0,
			"ecsEnabled": false,
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
			"id": 1,
			"infoUrl": null,
			"initialDispersion": 3,
			"innerHeaderRewrite": null,
			"ipv6RoutingEnabled": null,
			"lastHeaderRewrite": null,
			"lastUpdated": "2022-09-29T22:23:28-6:00",
			"logsEnabled": false,
			"longDesc": "long desc",
			"matchList": null,
			"maxDnsAnswers": null,
			"maxOriginConnections": 0,
			"maxRequestHeaderBytes": 0,
			"midHeaderRewrite": null,
			"missLat": null,
			"missLong": null,
			"multiSiteOrigin": false,
			"originShield": null,
			"orgServerFqdn": null,
			"profileDescription": null,
			"profileId": null,
			"profileName": null,
			"protocol": null,
			"qstringIgnore": null,
			"rangeRequestHandling": null,
			"rangeSliceBlockSize": null,
			"regexRemap": null,
			"regional": false,
			"regionalGeoBlocking": false,
			"remapText": null,
			"routingName": "cdn",
			"serviceCategory": null,
			"signed": false,
			"signingAlgorithm": null,
			"sslKeyVersion": null,
			"tenant": null,
			"tenantId": 1,
			"tlsVersions": null,
			"topology": null,
			"trResponseHeaders": null,
			"trRequestHeaders": null,
			"type": null,
			"typeId": 8,
			"xmlId": "demo1"
		},
		"status": "submitted"
	}}

.. [#implicit-null] Because of how the Traffic Ops API parses requests, there is no distinction between ``null`` and ``undefined``/missing properties. This means that sending the request payload ``{}`` in this case will result in the :term:`DSR` being unassigned.
