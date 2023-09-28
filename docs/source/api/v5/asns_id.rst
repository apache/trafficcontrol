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


.. _to-api-asns-id:

***************
``asns/{{id}}``
***************
.. seealso:: `The Autonomous System Wikipedia page <https://en.wikipedia.org/wiki/Autonomous_system_%28Internet%29>`_ for an explanation of what an :abbr:`ASN (Autonomous System Number)` actually is.

``PUT``
=======
Allows user to edit an existing :abbr:`ASN (Autonomous System Number)`-to-:term:`Cache Group` association.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: ASN:UPDATE, ASN:READ, CACHE-GROUP:UPDATE, CACHE-GROUP:READ
:Response Type: Object

Request Structure
-----------------
:asn:          The new :abbr:`ASN (Autonomous System Number)` which will be associated with the identified :term:`Cache Group` - must not conflict with existing associations
:cachegroup: An optional field which, if present, is a string that specifies the :ref:`cache-group-name` of a :term:`Cache Group` to which this :abbr:`ASN (Autonomous System Number)` will be assigned

	.. note:: While this endpoint accepts the ``cachegroup`` field, sending this in the request payload has no effect except that the response will (erroneously) name the :term:`Cache Group` to which the :abbr:`ASN (Autonomous System Number)` was assigned. Any subsequent requests will reveal that, in fact, the :term:`Cache Group` is set entirely by the ``cachegroupId`` field, and so the actual :ref:`cache-group-name` may differ from what was in the request.

:cachegroupId: An integer that is the :ref:`cache-group-id` of a :term:`Cache Group` to which this :abbr:`ASN (Autonomous System Number)` will be assigned - must not conflict with existing associations


.. table:: Request Path Parameters

	+------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| Name | Required | Description                                                                                                              |
	+======+==========+==========================================================================================================================+
	| id   | yes      | The integral, unique identifier of the desired :abbr:`ASN (Autonomous System Number)`-to-:term:`Cache Group` association |
	+------+----------+--------------------------------------------------------------------------------------------------------------------------+


.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/asns/1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 29
	Content-Type: application/json

	{"asn": 2, "cachegroupId": 1}

Response Structure
------------------
:asn:          An :abbr:`ASN (Autonomous System Number)` as specified by IANA for identifying a service provider
:cachegroup:   A string that is the :ref:`cache-group-name` of the :term:`Cache Group` that is associated with this :abbr:`ASN (Autonomous System Number)`
:cachegroupId: An integer that is the :ref:`cache-group-id` of the :term:`Cache Group` that is associated with this :abbr:`ASN (Autonomous System Number)`
:id:           An integral, unique identifier for this association between an :abbr:`ASN (Autonomous System Number)` and a :term:`Cache Group`
:lastUpdated:  The time and date this server entry was last updated in :ref:`non-rfc-datetime`

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: /83P4LJVsrQx9BKHFxo5pbhQMlB4o3a9v3PpkspyOJcpNx1S/GJhCPpiANBki547sbY+0vTq76IriHZ4GYp8bA==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 25 May 2023 21:59:33 GMT
	Content-Length: 160

	{ "alerts": [
		{
			"text": "asn was updated.",
			"level": "success"
		}
	],
	"response": {
		"asn": 2,
		"cachegroup": null,
		"cachegroupId": 1,
		"id": 1,
		"lastUpdated": "2023-05-25T15:59:33.7096-06:00"
	}}

``DELETE``
==========
Deletes an association between an :abbr:`ASN (Autonomous System Number)` and a :term:`Cache Group`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: ASN:DELETE, ASN:READ, CACHE-GROUP:READ, CACHE-GROUP:UPDATE
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------+--------------------------------------------------------------------------------------------------------------------------+
	| Name | Required | Description                                                                                                              |
	+======+==========+==========================================================================================================================+
	| id   | yes      | The integral, unique identifier of the desired :abbr:`ASN (Autonomous System Number)`-to-:term:`Cache Group` association |
	+------+----------+--------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/5.0/asns/1 HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 0

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 02 Dec 2019 23:06:24 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: 6t3WA+DOcfPJB5UnvDpzEVx5ySfmJgEV9wgkO71U5k32L1VXpxcaTdDVLNGgDDl9sdNftmYnKXf5jpfWUuFYJQ==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 02 Dec 2019 22:06:24 GMT
	Content-Length: 81

	{ "alerts": [
		{
			"text": "asn was deleted.",
			"level": "success"
		}
	]}
