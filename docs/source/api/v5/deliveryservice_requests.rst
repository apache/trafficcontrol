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
Retrieves :term:`Delivery Service Requests`.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: DS-REQUEST:READ, DELIVERY-SERVICE:READ, USER:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                                                                            |
	+===========+==========+========================================================================================================================================================+
	| assignee  | no       | Filter for :term:`Delivery Service Requests` that are assigned to the user identified by this username.                                                |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| assigneeId| no       | Filter for :term:`Delivery Service Requests` that are assigned to the user identified by this integral, unique identifier                              |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| author    | no       | Filter for :term:`Delivery Service Requests` submitted by the user identified by this username                                                         |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| authorId  | no       | Filter for :term:`Delivery Service Requests` submitted by the user identified by this integral, unique identifier                                      |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| changeType| no       | Filter for :term:`Delivery Service Requests` of the change type specified. Can be ``create``, ``update``, or ``delete``.                               |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| createdAt | no       | Filter for :term:`Delivery Service Requests` created on a certain date/time. Value must be :rfc:`3339` compliant. Eg. ``2019-09-19T19:35:38.828535Z``  |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| id        | no       | Filter for the :term:`Delivery Service Requests` identified by this integral, unique identifier.                                                       |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| status    | no       | Filter for :term:`Delivery Service Requests` whose status is the status specified. The status can be ``draft``, ``submitted``, ``pending``,            |
	|           |          | ``rejected``, or ``complete``.                                                                                                                         |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| xmlId     | no       | Filter for :term:`Delivery Service Requests` that have the given :ref:`ds-xmlid`.                                                                      |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| orderby   | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` array                                    |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| sortOrder | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                                                               |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| limit     | no       | Choose the maximum number of results to return                                                                                                         |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| offset    | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit                                                   |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+
	| page      | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and the first page is 1. If ``offset``     |
	|           |          | was defined, this query parameter has no effect. ``limit`` must be defined to make use of ``page``.                                                    |
	+-----------+----------+--------------------------------------------------------------------------------------------------------------------------------------------------------+

.. versionadded:: ATCv6
	The ``createdAt`` query parameter was added to this in endpoint across all API versions in :abbr:`ATC (Apache Traffic Control)` version 6.0.0.

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/deliveryservice_requests?status=draft HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: access_token=...; mojolicious=...

Response Structure
------------------
The response is an array of representations of :term:`Delivery Service Requests`.

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
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 09 Jun 2023 06:32:57 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 13 Jun 2023 17:01:41 GMT
	Content-Length: 988

	{ "response": [{
		"assignee": null,
		"author": "admin",
		"changeType": "update",
		"createdAt": "2023-06-09T10:55:00.918782+05:30",
		"id": 1,
		"lastEditedBy": "admin",
		"lastUpdated": "2023-06-13T22:31:30.122247+05:30",
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
			"lastUpdated": "2023-05-19T09:52:13.3131+05:30",
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
			"lastUpdated": "2023-05-19T08:40:13Z",
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
			"multiSiteOrigin": false,
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
			"topology": null,
			"trResponseHeaders": null,
			"trRequestHeaders": null,
			"type": "HTTP",
			"typeId": 1,
			"xmlId": "demo1"
		},
		"status": "draft"
	}]}

.. _to-api-deliveryservice-requests-post:

``POST``
========
Creates a new :term:`Delivery Service Request`. "Closed" :term:`Delivery Service Requests` cannot be created, an existing :term:`Delivery Service Request` must be placed into a closed :ref:`dsr-status`. A :term:`Delivery Service Request` to create, modify or delete a :term:`Delivery Service` cannot be created if an open :term:`Delivery Service Request` exists for a :term:`Delivery Service` with the same :ref:`ds-xmlid`. Because of this, :term:`Delivery Service Requests` cannot be used to change a :term:`Delivery Service`'s :ref:`ds-xmlid`.

:Auth. Required:       Yes
:Roles Required:       "admin", "Federation", "operations", "Portal", or "Steering"
:Permissions Required: DS-REQUEST:CREATE, DELIVERY-SERVICE:READ, USER:READ
:Response Type:        Object

Request Structure
-----------------
The request must be a well-formed representation of a :term:`Delivery Service Request`, without any response-only fields, of course.

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/deliveryservice_requests HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: access_token=...; mojolicious=...
	Content-Length: 2011
	Content-Type: application/json

	{
		"changeType": "update",
		"status": "draft",
		"requested": {
			"active": "INACTIVE",
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
			"firstHeaderRewrite": null,
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
			"lastUpdated": "2023-06-09T10:51:00+05:30",
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
			"regional": false,
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
			"topology": null,
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
			"serviceCategory": null,
			"tlsVersions": null
		}
	}


Response Structure
------------------
The response will be a representation of the created :term:`Delivery Service Request`.

.. versionchanged:: 5.0
	Prior to version 5.0 of the API, the ``lastUpdated`` field was in :ref:`non-rfc-datetime`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 201 Created
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Location: /api/5.0/deliveryservice_requests/1
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 09 Jun 2023 06:25:00 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 09 Jun 2023 05:25:00 GMT
	Content-Length: 1027

	{ "alerts": [{
		"text": "Delivery Service request created",
		"level": "success"
	}],
	"response": {
		"assignee": null,
		"author": "admin",
		"changeType": "update",
		"createdAt": "2023-06-09T10:55:00.918782+05:30",
		"id": 1,
		"lastEditedBy": "admin",
		"lastUpdated": "2023-06-09T10:55:00.918782+05:30",
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
			"lastUpdated": "2023-05-19T09:52:13.3131+05:30",
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
			"lastUpdated": "2023-06-09T10:51:00+05:30",
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
			"multiSiteOrigin": false,
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
			"topology": null,
			"trResponseHeaders": null,
			"trRequestHeaders": null,
			"type": "HTTP",
			"typeId": 1,
			"xmlId": "demo1"
		},
		"status": "draft"
	}}

``PUT``
=======
Updates an existing :term:`Delivery Service Request`. Note that "closed" :term:`Delivery Service Requests` are uneditable.

.. seealso:: The proper way to change a :term:`Delivery Service Request`'s :ref:`dsr-status` is by using the :ref:`to-api-deliveryservice_requests-id-status` endpoint's ``PUT`` handler.

:Auth. Required:       Yes
:Roles Required:       "admin", "Federation", "operations", "Portal", or "Steering"
:Permissions Required: DS-REQUEST:UPDATE, DELIVERY-SERVICE:READ, USER:READ
:Response Type:        Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+--------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                      |
	+===========+==========+==================================================================================================+
	| id        | yes      | The integral, unique identifier of the :term:`Delivery Service Request` that you want to update. |
	+-----------+----------+--------------------------------------------------------------------------------------------------+

The request body must be a representation of a :term:`Delivery Service Request` without any response-only fields.

.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/deliveryservice_requests?id=1 HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: access_token=...; mojolicious=...
	Content-Length: 426

	{
		"changeType": "update",
		"requested": {
			"active": "INACTIVE",
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
			"regional": false,
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
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 09 Jun 2023 06:24:20 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 09 Jun 2023 05:24:20 GMT
	Content-Length: 1142

	{ "alerts": [{
		"text": "Delivery Service Request #1 updated",
		"level": "success"
	}],
	"response": {
		"assignee": null,
		"author": "",
		"changeType": "update",
		"createdAt": "2023-06-09T10:54:20.435475+05:30",
		"id": null,
		"lastEditedBy": "admin",
		"lastUpdated": "2023-06-09T10:51:39.552061+05:30",
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
			"lastUpdated": "2023-05-19T09:52:13.3131+05:30",
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
			"lastUpdated": "0001-01-01T00:00:00Z",
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
		"status": "draft"
	}}


``DELETE``
==========
Deletes a :term:`Delivery Service Request`.

:Auth. Required:       Yes
:Roles Required:       "admin", "Federation", "operations", "Portal", or "Steering"
:Permissions Required: DS-REQUEST:DELETE, DELIVERY-SERVICE:READ, USER:READ
:Response Type:        Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+--------------------------------------------------------------------------------------------------+
	| Name      | Required | Description                                                                                      |
	+===========+==========+==================================================================================================+
	| id        | yes      | The integral, unique identifier of the :term:`Delivery Service Request` that you want to delete. |
	+-----------+----------+--------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/5.0/deliveryservice_requests?id=1 HTTP/1.1
	User-Agent: python-requests/2.25.1
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: access_token=...; mojolicious=...
	Content-Length: 0

Response Structure
------------------
The response is a full representation of the deleted :term:`Delivery Service Request`.

.. versionchanged:: 5.0
	Prior to version 5.0 of the API, the ``lastUpdated`` field was in :ref:`non-rfc-datetime`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Permissions-Policy: interest-cohort=()
	Set-Cookie: mojolicious=...; Path=/; Expires=Fri, 09 Jun 2023 06:24:53 GMT; Max-Age=3600; HttpOnly
	Vary: Accept-Encoding
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 09 Jun 2023 05:24:53 GMT
	Content-Length: 1147

	{ "alerts": [{
		"text": "Delivery Service Request #1 deleted",
		"level": "success"
	}],
	"response": {
		"assignee": "admin",
		"author": "admin",
		"changeType": "update",
		"createdAt": "2023-06-09T10:51:39.552061+05:3",
		"id": 1,
		"lastEditedBy": "admin",
		"lastUpdated": "2023-06-09T10:54:20.435475+05:30",
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
			"lastUpdated": "2023-05-19T09:52:13.3131+05:30",
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
			"lastUpdated": "0001-01-01T00:00:00Z",
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
