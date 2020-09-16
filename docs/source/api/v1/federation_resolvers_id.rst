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

.. _to-api-v1-federation_resolvers-id:

*******************************
``federation_resolvers/{{ID}}``
*******************************
.. deprecated:: 1.5


``DELETE``
==========
Deletes a federation resolver.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

	.. versionchanged:: 1.5
		Older versions of this endpoint did not return a ``response`` object (i.e. it was ``undefined``).

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

	DELETE /api/1.4/federation_resolvers/2 HTTP/1.1
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
	Set-Cookie: mojolicious=; Path=/; HttpOnly
	Whole-Content-Sha512: Gpvrv90AlWWlHrLsPyqm0Jhn+bDMZ8X/MwLXH8xIXyi9DFsLDqSd1si8LjNM7fiXP0j5lbRJXpNH4dKxsvSrTg==
	X-Server-Name: traffic_ops_golang/
	Date: Sat, 09 Nov 2019 07:30:09 GMT
	Content-Length: 215

	{ "alerts": [
		{
			"text": "Federation resolver deleted [ IP = 1.2.6.4/22 ] with id: 2",
			"level": "success"
		},
		{
			"text": "This endpoint is deprecated, please use '/federation_resolvers' instead",
			"level": "warning"
		}
	],
	"response": {
		"id": 2,
		"ipAddress": "1.2.6.4/22",
		"type": "RESOLVE4"
	}}

