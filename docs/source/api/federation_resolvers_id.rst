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

.. _to-api-federation_resolvers-id:

*******************************
``federation_resolvers/{{ID}}``
*******************************

``DELETE``
==========
Deletes a federation resolver.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------+
	| Name | Description                                                           |
	+======+=======================================================================+
	|  ID  | Integral, unique identifier for the federation resolver to be deleted |
	+------+-----------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/1.4/federation_resolvers/3 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...

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
	date: Wed, 05 Dec 2018 01:06:51 GMT
	server: Mojolicious (Perl)
	set-cookie: mojolicious=...; expires=Wed, 05 Dec 2018 05:06:51 GMT; path=/; HttpOnly
	vary: Accept-Encoding
	whole-content-sha512: NqAZuZYlF1UWOaazbj/j4gWX7ye0kGGakRRFEkK6ShxqXvCxE0dCTyu75qiLPN2wSgr3FGQnp2Sq345sE7In9g==
	content-length: 98

	{ "alerts": [
		{
			"level": "success",
			"text": "Federation resolver deleted [ IP = ::1/128 ] with id: 3"
		}
	]}

