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

.. _to-api-v4-federations-id-federation_resolvers:

*******************************************
``federations/{{ID}}/federation_resolvers``
*******************************************

``GET``
=======
Retrieves federation resolvers assigned to a federation.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: FEDERATION:READ, FEDERATION-RESOLVER:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------------------------+
	| Name |                 Description                                                              |
	+======+==========================================================================================+
	|  ID  | The integral, unique identifier for the federation for which resolvers will be retrieved |
	+------+------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/4.0/federations/1/federation_resolvers HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:id:        The integral, unique identifier of this federation resolver
:ipAddress: The IP address of the federation resolver - may be IPv4 or IPv6
:type:      The type of resolver - one of:

	RESOLVE4
		This resolver is for IPv4 addresses (and ``ipAddress`` is IPv4)
	RESOLVE6
		This resolver is for IPv6 addresses (and ``ipAddress`` is IPv6)

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	cache-control: no-cache, no-store, max-age=0, must-revalidate
	content-type: application/json
	date: Wed, 05 Dec 2018 00:49:50 GMT
	X-Server-Name: traffic_ops_golang/
	set-cookie: mojolicious=...; expires=Wed, 05 Dec 2018 04:49:50 GMT; path=/; HttpOnly
	vary: Accept-Encoding
	whole-content-sha512: csC18kE3YjiILHP1wmJg7V4h/XWY8HUMKyPuZWnde2g7HJ4gTY51HfjCSqhyKvIJQ8Rl7uEqshF3Ey6xIMOX4A==
	content-length: 63

	{ "response": [
		{
			"ipAddress": "0.0.0.0",
			"type": "RESOLVE4",
			"id": 1
		}
	]}

``POST``
========
Assigns one or more resolvers to a federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: FEDERATION:UPDATE, FEDERATION:READ, FEDERATION-RESOLVER:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------------------------+
	| Name |                 Description                                                              |
	+======+==========================================================================================+
	|  ID  | The integral, unique identifier for the federation for which resolvers will be retrieved |
	+------+------------------------------------------------------------------------------------------+

:fedResolverIds: An array of integral, unique identifiers for federation resolvers
:replace:        An optional boolean (default: ``false``) which, if ``true``, will cause any conflicting assignments already in place to be overridden by this request

	.. note:: If ``replace`` is not given (and/or not ``true``), then any conflicts with existing assignments will cause the entire operation to fail.

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/federations/1/federation_resolvers HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 41
	Content-Type: application/json

	{
		"fedResolverIds": [1],
		"replace": true
	}

Response Structure
------------------
:fedResolverIds: An array of integral, unique identifiers for federation resolvers
:replace:        An optionally-present boolean (default: ``false``) which, if ``true``, any conflicting assignments already in place were overridden by this request

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	cache-control: no-cache, no-store, max-age=0, must-revalidate
	content-type: application/json
	date: Wed, 05 Dec 2018 00:47:47 GMT
	X-Server-Name: traffic_ops_golang/
	set-cookie: mojolicious=...; expires=Wed, 05 Dec 2018 04:47:47 GMT; path=/; HttpOnly
	vary: Accept-Encoding
	whole-content-sha512: +JDcRByS3HO6pMg3Gzkvn0w7/v5oRul9e+RxyFIOKJKNHOkZILyQBS+PJpxDeCgwI19+0poW5dyHPPR9SwbNCA==
	content-length: 148

	{ "alerts": [
		{
			"level": "success",
			"text": "1 resolver(s) were assigned to the test.quest. federation"
		}
	],
	"response": {
		"replace": true,
		"fedResolverIds": [
			1
		]
	}}
