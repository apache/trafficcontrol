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

.. _to-api-federation_resolvers:

************************
``federation_resolvers``
************************

``POST``
========
Creates a new federation resolver.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
:ipAddress: The IP address of the resolver - may be IPv4 or IPv6
:typeId:    The integral, unique identifier of the type of resolver being created - will *represent* one of:

	RESOLVE4
		Resolver is for IPv4 addresses and ``ipAddress`` is IPv4
	RESOLVE6
		Resolver is for IPv6 addresses and ``ipAddress`` is IPv6

.. code-block:: http
	:caption: Request Example

	POST /api/1.4/federation_resolvers HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 39
	Content-Type: application/json

	{
		"ipAddress": "0.0.0.0",
		"typeId": 35
	}

Response Structure
------------------
:id:        The integral, unique identifier of the resolver
:ipAddress: The IP address of the resolver - may be IPv4 or IPv6
:type:      The type of the resolver - one of:

	RESOLVE4
		Resolver is for IPv4 addresses and ``ipAddress`` is IPv4
	RESOLVE6
		Resolver is for IPv6 addresses and ``ipAddress`` is IPv6

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	cache-control: no-cache, no-store, max-age=0, must-revalidate
	content-type: application/json
	date: Wed, 05 Dec 2018 00:41:42 GMT
	server: Mojolicious (Perl)
	set-cookie: mojolicious=...; expires=Wed, 05 Dec 2018 04:41:42 GMT; path=/; HttpOnly
	vary: Accept-Encoding
	whole-content-sha512: JaXUP+onwOhGs+H/w7u2bNm9a7bqGLGDGJRutFsByTODBAfNr+X7NZ4aO+5w3RyDji1Ih1z5sLadQeEcdZj8vw==
	content-length: 151

	{ "alerts": [
		{
			"level": "success",
			"text": "Federation Resolver created [ IP = 0.0.0.0 ] with id: 1"
		}
	],
	"response": {
		"ipAddress": "0.0.0.0",
		"id": 1,
		"typeId": 35
	}}
