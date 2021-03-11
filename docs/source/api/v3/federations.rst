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

.. _to-api-v3-federations:

***************
``federations``
***************

``GET``
=======
Retrieves a list of :term:`Federation` mappings (i.e. :term:`Federation` Resolvers) for the current user.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Array

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:deliveryService: The ``xml_id`` that uniquely identifies the :term:`Delivery Service` that uses the federation mappings in ``mappings``
:mappings:        An array of objects that represent the mapping of a :term:`Federation`'s :abbr:`CNAME (Canonical Name)` to one or more Resolvers

	:cname:    The actual CNAME used by the :term:`Federation`
	:resolve4: An array of IPv4 addresses (or subnets in :abbr:`CIDR (Classless Inter-Domain Routing)` notation) capable of resolving the :term:`Federation`'s CNAME
	:resolve6: An array of IPv6 addresses (or subnets in :abbr:`CIDR (Classless Inter-Domain Routing)` notation) capable of resolving the :term:`Federation`'s CNAME
	:ttl:      The :abbr:`TTL (Time To Live)` of the CNAME in hours

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: d6Llm5qNc2sfgVH9IimW7hA4wvtBUq6EzUmpJf805kB0k6v2WysNgFEWK4hBXNdAYkr8hYuKPrwDy3tCx0OZ8Q==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 03 Dec 2018 17:19:13 GMT
	Content-Length: 136

	{ "response": [
		{
			"mappings": [
				{
					"ttl": 300,
					"cname": "blah.blah.",
					"resolve4": [
						"0.0.0.0/32"
					],
					"resolve6": [
						"::/128"
					]
				}
			],
			"deliveryService": "demo1"
		}
	]}


``POST``
========
Allows a user to create :term:`Federation` Resolvers for :term:`Delivery Services`, providing the :term:`Delivery Service` is within a CDN that has some associated :term:`Federation`.

.. warning:: Confusingly, this method of this endpoint does **not** create a new :term:`Federation`; to do that, the :ref:`to-api-v3-cdns-name-federations` endpoint must be used. Furthermore, the :term:`Federation` must properly be assigned to a :term:`Delivery Service` using the :ref:`to-api-v3-federations-id-deliveryservices` and assigned to the user creating Resolvers using :ref:`to-api-v3-federations-id-users`.

.. seealso:: The :ref:`to-api-v3-federations-id-federation_resolvers` endpoint duplicates this functionality.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object (string)

Request Structure
-----------------

The request payload is an array of objects that describe Delivery Service :term:`Federation` Resolver mappings. Each object in the array must be in the following format.

:deliveryService: The :ref:`ds-xmlid` of the :term:`Delivery Service` which will use the :term:`Federation` Resolvers specified in ``mappings``
:mappings:        An object containing two arrays of IP addresses (or subnets in :abbr:`CIDR (Classless Inter-Domain Routing)` notation) to use as :term:`Federation` Resolvers

	:resolve4: An array of IPv4 addresses (or subnets in :abbr:`CIDR (Classless Inter-Domain Routing)` notation) that can resolve the :term:`Delivery Service`'s :term:`Federation`
	:resolve6: An array of IPv6 addresses (or subnets in :abbr:`CIDR (Classless Inter-Domain Routing)` notation) that can resolve the :term:`Delivery Service`'s :term:`Federation`

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/federations HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 118
	Content-Type: application/json


	[{
		"deliveryService":"demo1",
		"mappings":{
			"resolve4":["127.0.0.1", "0.0.0.0/32"],
			"resolve6":["::1", "5efa::ff00/128"]
		}
	}]

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
	Whole-Content-Sha512: B7TSUOYZPRPyi3mVy+CuxiXR5k/d0s07w4i6kYzpWS+YL79juEfkuSqfedaYG/kMA8O9XbjkWRjcBAdxOVrdTQ==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 23 Oct 2019 22:28:02 GMT
	Content-Length: 152

	{ "alerts": [
		{
			"text": "admin successfully created federation resolvers.",
			"level": "success"
		}
	],
	"response": "admin successfully created federation resolvers."
	}


``DELETE``
==========
Deletes **all** :term:`Federation` Resolvers associated with the logged-in user's :term:`Federations`.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object (string)

Request Structure
-----------------
No parameters available

.. code-block:: http
	:caption: Request Example

	DELETE /api/3.0/federations HTTP/1.1
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
	Whole-Content-Sha512: fd7P45mIiHuYqZZW6+8K+YjY1Pe504Aaw4J4Zp9AhrqLX72ERytTqWtAp1msutzNSRUdUSC72+odNPtpv3O8uw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 23 Oct 2019 23:34:53 GMT
	Content-Length: 184

	{ "alerts": [
		{
			"text": "admin successfully deleted all federation resolvers: [ 8.8.8.8 ]",
			"level": "success"
		}
	],
	"response": "admin successfully deleted all federation resolvers: [ 8.8.8.8 ]"
	}

``PUT``
=======
Replaces **all** :term:`Federations` associated with a user's :term:`Delivery Service`\ (s) with those defined inside the request payload.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object (string)

Request Structure
-----------------
The request payload is an array of objects that describe Delivery Service :term:`Federation` Resolver mappings. Each object in the array must be in the following format.

:deliveryService: The :ref:`ds-xmlid` of the :term:`Delivery Service` which will use the :term:`Federation` Resolvers specified in ``mappings``
:mappings:        An object containing two arrays of IP addresses (or subnets in :abbr:`CIDR (Classless Inter-Domain Routing)` notation) to use as :term:`Federation` Resolvers

	:resolve4: An array of IPv4 addresses (or subnets in :abbr:`CIDR (Classless Inter-Domain Routing)` notation) that can resolve the :term:`Delivery Service`'s :term:`Federation`
	:resolve6: An array of IPv6 addresses (or subnets in :abbr:`CIDR (Classless Inter-Domain Routing)` notation) that can resolve the :term:`Delivery Service`'s :term:`Federation`

.. code-block:: http
	:caption: Request Example

	PUT /api/3.0/federations HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 95
	Content-Type: application/json

	[{ "mappings": {
		"resolve4": ["8.8.8.8"],
		"resolve6": []
	},
	"deliveryService":"demo1"
	}]

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: dQ5AvQULhc254zQwgUpBl1/CHbLr/clKtkbs0Ju9f1BM4xIfbbO3puFNN9zaEaZ1iz0lBvHFp/PgfUqisD3QHA==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 23 Oct 2019 23:22:03 GMT
	Content-Length: 258
	Content-Type: application/json

	{ "alerts": [
		{
			"text": "admin successfully deleted all federation resolvers: [ 8.8.8.8 ]",
			"level": "success"
		},
		{
			"text": "admin successfully created federation resolvers.",
			"level": "success"
		}
	],
	"response": "admin successfully created federation resolvers."
	}
