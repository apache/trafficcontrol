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

.. _to-api-v3-federation_resolvers:

************************
``federation_resolvers``
************************

``GET``
=======
Retrieves :term:`Federation` Resolvers.

:Auth. Required: Yes
:Roles Required: None
:Response Type: Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| Name       | Required | Description                                                                                         |
	+============+==========+=====================================================================================================+
	| id         | no       | Return only the Federation Resolver identified by this integral, unique identifier                  |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| ipAddress  | no       | Return only the Federation Resolver(s) that has/have this IP Address                                |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| type       | no       | Return only the Federation Resolvers of this :term:`Type`                                           |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| orderby    | no       | Choose the ordering of the results - must be the name of one of the fields of the objects in the    |
	|            |          | ``response`` array                                                                                  |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| sortOrder  | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")            |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| limit      | no       | Choose the maximum number of results to return                                                      |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| offset     | no       | The number of results to skip before beginning to return results. Must use in conjunction with      |
	|            |          | limit                                                                                               |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| page       | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are        |
	|            |          | ``limit`` long and the first page is 1. If ``offset`` was defined, this query parameter has no      |
	|            |          | effect. ``limit`` must be defined to make use of ``page``.                                          |
	+------------+----------+-----------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/federation_resolvers?type=RESOLVE6 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.63.0
	Accept: */*
	Cookie: mojolicious=...


Response Structure
------------------
:id:          The integral, unique identifier of the resolver
:ipAddress:   The IP address or :abbr:`CIDR (Classless Inter-Domain Routing)`-notation subnet of the resolver - may be IPv4 or IPv6
:lastUpdated: The date and time at which this resolver was last updated, in :ref:`non-rfc-datetime`
:type:        The :term:`Type` of the resolver

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: 4TLkULAOAuap47H+hpwyf2lHjDbHbSNQHLMj7BCTHtps2CQxCuq7mwctbwqmPdmAjLOUXAIRsHmvSuAp4K64jw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 06 Nov 2019 00:03:56 GMT
	Content-Length: 101

	{ "response": [
		{
			"id": 1,
			"ipAddress": "::1/1",
			"lastUpdated": "2019-11-06 00:00:40+00",
			"type": "RESOLVE6"
		}
	]}


``POST``
========
Creates a new federation resolver.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
:ipAddress: The IP address of the resolver - may be IPv4 or IPv6
:typeId:    The integral, unique identifier of the :term:`Type` of resolver being created

	.. caution:: This field should only ever be an identifier for one of the :term:`Types` "RESOLVE4" or "RESOLVE6", but there is **no protection for this built into Traffic Ops** and therefore **any valid** :term:`Type` **identifier will be silently accepted by Traffic Ops** and so care should be taken to ensure that these :term:`Types` are properly identified. If any :term:`Type` besides "RESOLVE4" or "RESOLVE6" is identified, the resulting resolver *will* **not** *work*.

	.. seealso:: :ref:`to-api-v3-types` is the endpoint that can be used to determine the identifier for various :term:`Types`

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/federation_resolvers HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.63.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 36
	Content-Type: application/json

	{
		"ipAddress": "::1/1",
		"typeId": 37
	}

Response Structure
------------------
:id:        The integral, unique identifier of the resolver
:ipAddress: The IP address or :abbr:`CIDR (Classless Inter-Domain Routing)`-notation subnet of the resolver - may be IPv4 or IPv6
:type:      The :term:`Type` of the resolver
:typeId:    The integral, unique identifier of the :term:`Type` of the resolver


.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: e9D8JNrQb64xpuDwoBwbISSWUkDGCL2l37NuDXsXsPYof2EqmeHondD8NzxDSwWNJ8d9B9DXpZDbRUtgdXR8BQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 06 Nov 2019 00:00:40 GMT
	Content-Length: 153

	{ "alerts": [
		{
			"text": "Federation Resolver created [ IP = ::1/1 ] with id: 1",
			"level": "success"
		}
	],
	"response": {
		"id": 1,
		"ipAddress": "::1/1",
		"type": "RESOLVE6",
		"typeId": 37
	}}

``DELETE``
==========
Deletes a federation resolver.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+----------------------------------------------------------------------------------+
	| Name | Required | Description                                                           |
	+======+==========+=======================================================================+
	|  id  | yes      | Integral, unique identifier for the federation resolver to be deleted |
	+------+----------+-----------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/3.0/federation_resolvers?id=4 HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 0

Response Structure
------------------
:id:        The integral, unique identifier of the resolver
:ipAddress: The IP address or :abbr:`CIDR (Classless Inter-Domain Routing)`-notation subnet of the resolver - may be IPv4 or IPv6
:type:      The :term:`Type` of the resolver

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: 2v4LYQdRVhaFJVd86Iv1BWVYzNPSlzpQ222bUB7Zz+Ss8A48FNyHZjPlq5a+a4g9KAQCTUIytWnIQk+L1fF6FQ==
	X-Server-Name: traffic_ops_golang/
	Date: Fri, 08 Nov 2019 23:19:01 GMT
	Content-Length: 161

	{ "alerts": [
		{
			"text": "Federation resolver deleted [ IP = 1.2.6.4/22 ] with id: 4",
			"level": "success"
		}
	],
	"response": {
		"id": 4,
		"ipAddress": "1.2.6.4/22",
		"type": "RESOLVE6"
	}}
