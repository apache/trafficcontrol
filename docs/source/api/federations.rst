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

.. _to-api-federations:

***************
``federations``
***************

``GET``
=======
Retrieves a list of federation mappings (aka federation resolvers) for the current user.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Array

Request Structure
-----------------
No parameters available.

Response Structure
------------------
:deliveryService: The ``xml_id`` that uniquely identifies the :term:`Delivery Service` that uses the federation mappings in ``mappings``
:mappings:        An array of objects that represent the mapping of a federation's Canonical Name (CNAME) to one or more resolvers

	:cname:    The actual CNAME used by the federation
	:resolve4: An array of IPv4 addresses capable of resolving the federation's CNAME
	:resolve6: An array of IPv6 addresses capable of resolving the federation's CNAME
	:ttl:      The Time To Live (TTL) of the CNAME in hours

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
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
Allows a user to create federation resolvers for :term:`Delivery Service`\ s, providing the :term:`Delivery Service` is within a CDN that has some associated federation.

.. warning:: Confusingly, this endpoint does **not** create a new federation; to do that, the :ref:`to-api-cdns-name-federations` endpoint must be used. Furthermore, the federation must properly be assigned to a :term:`Delivery Service` using the :ref:`to-api-federations-id-deliveryservices` and assigned to the user creating resolvers using :ref:`to-api-federations-id-users`.

.. seealso:: The :ref:`to-api-federations-id-federation_resolvers` endpoint duplicates this functionality.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object (string)

Request Structure
-----------------
:federations: The top-level key that must exist - an array of objects that each describe a set of resolvers for a :term:`Delivery Service`'s federation

	:deliveryService: The 'xml_id' of the :term:`Delivery Service` which will use the federation resolvers specified in ``mappings``
	:mappings:        An object containing two arrays of IP addresses to use as federation resolvers

		:resolve4: An array of IPv4 addresses that can resolve the :term:`Delivery Service`'s federation
		:resolve6: An array of IPv6 addresses that can resolve the :term:`Delivery Service`'s federation

.. code-block:: http
	:caption: Request Example

	POST /api/1.1/federations HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 119
	Content-Type: application/json

	{ "federations": [{
		"deliveryService": "demo1",
		"mappings": {
			"resolve4": ["0.0.0.0"],
			"resolve6": ["::"]
		}
	}]}

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Mon, 03 Dec 2018 17:00:29 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Mon, 03 Dec 2018 21:00:29 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: dXg86uD2Un1AeBCeeBLSo2rsYgl6NOHHQEc5oMlpw1THOh2HwGdjwB3rPd/qoYIhOxcnnHoEstrEiHmucFev4A==
	Content-Length: 63

	{ "response": "admin successfully created federation resolvers." }


``DELETE``
==========
Deletes **all** federation resolvers associated with the logged-in user's federations.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object (string)

Request Structure
-----------------
No parameters available

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Cache-Control: no-cache, no-store, max-age=0, must-revalidate
	Content-Type: application/json
	Date: Mon, 03 Dec 2018 17:55:10 GMT
	Server: Mojolicious (Perl)
	Set-Cookie: mojolicious=...; expires=Mon, 03 Dec 2018 21:55:10 GMT; path=/; HttpOnly
	Vary: Accept-Encoding
	Whole-Content-Sha512: b84HraJH6Kiqrz7i1L1juDBJWdkdYbbClnWM0lZDljvpSkVT9adFTTrHiv7Mjtt2RKquGdzFZ6tqt9s+ODxqsw==
	Content-Length: 93

	{ "response": "admin successfully deleted all federation resolvers: [ 0.0.0.0/32, ::/128 ]." }


``PUT``
=======
Replaces **all** federations associated with a user's :term:`Delivery Service`\ (s) with those defined inside the request payload.

:Auth. Required: Yes
:Roles Required: "admin", "Federation", "operations", "Portal", or "Steering"
:Response Type:  Object (string)

Request Structure
-----------------
:federations: The top-level key that must exist - an array of objects that each describe a set of resolvers for a :term:`Delivery Service`'s federation

	:deliveryService: The 'xml_id' of the :term:`Delivery Service` which will use the federation resolvers specified in ``mappings``
	:mappings:        An object containing two arrays of IP addresses to use as federation resolvers

		:resolve4: An array of IPv4 addresses that can resolve the :term:`Delivery Service`'s federation
		:resolve6: An array of IPv6 addresses that can resolve the :term:`Delivery Service`'s federation

.. code-block:: http
	:caption: Request Example

	PUT /api/1.4/federations HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 113
	Content-Type: application/json

	{ "federations": [{
		"deliveryService": "demo1",
		"mappings": {
			"resolve4": ["0.0.0.1"],
			"resolve6": ["::1"]
		}
	}]}

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	cache-control: no-cache, no-store, max-age=0, must-revalidate
	content-type: application/json
	date: Wed, 05 Dec 2018 00:52:31 GMT
	server: Mojolicious (Perl)
	set-cookie: mojolicious=...; expires=Wed, 05 Dec 2018 04:52:30 GMT; path=/; HttpOnly
	vary: Accept-Encoding, Accept-Encoding
	whole-content-sha512: dXg86uD2Un1AeBCeeBLSo2rsYgl6NOHHQEc5oMlpw1THOh2HwGdjwB3rPd/qoYIhOxcnnHoEstrEiHmucFev4A==
	content-length: 63

	{"response": "admin successfully created federation resolvers."}
