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

.. _to-api-federations-all:

*******************
``federations/all``
*******************

``GET``
=======
Retrieves a list of :term:`Federation` mappings (also called :term:`Federation` Resolvers) for the current user.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: FEDERATION-RESOLVER:READ, DELIVERY-SERVICE:READ
:Response Type:  Array

Request Structure
-----------------
No parameters available.

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/federations/all HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:deliveryService:       The :ref:`ds-xmlid` of the delivery service.
:mappings:              An array of objects that represent the mapping of a :term:`Federation`'s :abbr:`CNAME (Canonical Name)` to one or more Resolvers

	:cname:                 The actual CNAME used by the :term:`Federation`
	:ttl:                   The :abbr:`TTL (Time To Live)` of the CNAME in hours

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Sun, 23 Feb 2020 21:38:06 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: UQBlGVPJytYMkv0V42EAIoJUnXjBTCXnOGpOberxte6TtnX63LTAKFfD2LejBVYXkKtnCdkBbs+SzhA0H1zdog==
	X-Server-Name: traffic_ops_golang/
	Date: Sun, 23 Feb 2020 20:38:06 GMT
	Content-Length: 138

	{
		"response": [
			{
				"mappings": [
					{
						"ttl": 60,
						"cname": "img1.mcdn.ciab.test."
					},
					{
						"ttl": 60,
						"cname": "img2.mycdn.ciab.test."
					}
				],
				"deliveryService": "demo1"
			},
			{
				"mappings": [
					{
						"ttl": 60,
						"cname": "static.mycdn.ciab.test."
					}
				],
				"deliveryService": "demo2"
			}
		]
	}
