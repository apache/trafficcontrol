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


.. _to-api-asns:

********
``asns``
********
.. seealso:: `The Autonomous System Wikipedia page <https://en.wikipedia.org/wiki/Autonomous_system_%28Internet%29>`_ for an explanation of what an :abbr:`ASN (Autonomous System Number)` actually is.

``GET``
=======
List all :abbr:`ASNs (Autonomous System Numbers)`.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: ASN:READ, CACHE-GROUP:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| Parameter  | Required | Description                                                                                         |
	+============+==========+=====================================================================================================+
	| cachegroup | no       | The :ref:`cache-group-id` of a :term:`Cache Group` - only :abbr:`ASNs (Autonomous System Numbers)`  |
	|            |          | for this :term:`Cache Group` will be returned.                                                      |
	+------------+----------+-----------------------------------------------------------------------------------------------------+
	| id         | no       | The integral, unique identifier of the desired                                                      |
	|            |          | :abbr:`ASN (Autonomous System Number)`-to-:term:`Cache Group` association                           |
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

	GET /api/5.0/asns HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...

Response Structure
------------------
:asn:          An :abbr:`ASN (Autonomous System Number)` as specified by IANA for identifying a service provider
:cachegroup:   A string that is the :ref:`cache-group-name` of the :term:`Cache Group` that is associated with this :abbr:`ASN (Autonomous System Number)`
:cachegroupId: An integer that is the :ref:`cache-group-id` of the :term:`Cache Group` that is associated with this :abbr:`ASN (Autonomous System Number)`
:id:           An integral, unique identifier for this association between an :abbr:`ASN (Autonomous System Number)` and a :term:`Cache Group`
:lastUpdated:  The time and date this server entry was last updated in :rfc:`3339` Format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 02 Dec 2019 22:51:14 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: F2NmDbTpXqrIQDX7IBKH9+1drtTL4XedSfJv6klMgLEZwbLCkddIXuSLpmgVCID6kTVqy3fTKjZS3U+HJ3YUEQ==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 25 May 2023 21:59:33 GMT
	Content-Length: 128

	{ "response": [
		{
			"asn": 1,
			"cachegroup": "TRAFFIC_ANALYTICS",
			"cachegroupId": 1,
			"id": 1,
			"lastUpdated": "2023-05-25T15:59:33.7096-06:00"
		}
	]}



``POST``
========
Creates a new :abbr:`ASN (Autonomous System Number)`.

.. note:: There cannot be two different ASN object with the same ``asn``. An ASN may only belong to one cachegroup, but a cachegroup can have zero or more ASNs.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: ASN:CREATE, ASN:READ, CACHE-GROUP:READ, CACHE-GROUP:UPDATE
:Response Type: Object

Request Structure
-----------------
:asn:        The value of the new :abbr:`ASN (Autonomous System Number)`
:cachegroup: An optional field which, if present, is a string that specifies the :ref:`cache-group-name` of a :term:`Cache Group` to which this :abbr:`ASN (Autonomous System Number)` will be assigned

	.. note:: While this endpoint accepts the ``cachegroup`` field, sending this in the request payload has no effect except that the response will (erroneously) name the :term:`Cache Group` to which the :abbr:`ASN (Autonomous System Number)` was assigned. Any subsequent requests will reveal that, in fact, the :term:`Cache Group` is set entirely by the ``cachegroupId`` field, and so the actual :ref:`cache-group-name` may differ from what was in the request.

:cachegroupId: An integer that is the :ref:`cache-group-id` of a :term:`Cache Group` to which this :abbr:`ASN (Autonomous System Number)` will be assigned

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/asns HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 29

	{"asn": 1, "cachegroupId": 1}


Response Structure
------------------
:asn:          An :abbr:`ASN (Autonomous System Number)` as specified by IANA for identifying a service provider
:cachegroup:   A string that is the :ref:`cache-group-name` of the :term:`Cache Group` that is associated with this :abbr:`ASN (Autonomous System Number)`
:cachegroupId: An integer that is the :ref:`cache-group-id` of the :term:`Cache Group` that is associated with this :abbr:`ASN (Autonomous System Number)`
:id:           An integral, unique identifier for this association between an :abbr:`ASN (Autonomous System Number)` and a :term:`Cache Group`
:lastUpdated:  The time and date this server entry was last updated in :rfc:`3339` Format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 02 Dec 2019 22:49:08 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: mx8b2GTYojz4QtMxXCMoQyZogCB504vs0yv6WGly4dwM81W3XiejWNuUwchRBYYi8QHaWsMZ3DaiGGfQi/8Giw==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 25 May 2023 21:59:33 GMT
	Content-Length: 150

	{ "alerts": [
		{
			"text": "asn was created.",
			"level": "success"
		}
	],
	"response": {
		"asn": 1,
		"cachegroup": null,
		"cachegroupId": 1,
		"id": 1,
		"lastUpdated": "2023-05-25T15:59:33.7096-06:00"
	}}

``PUT``
=======
Updates an existing :abbr:`ASN (Autonomous System Number)`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: ASN:UPDATE, ASN:READ, CACHE-GROUP:READ, CACHE-GROUP:UPDATE
:Response Type: Object

Request Structure
-----------------
:asn:           The value of the new :abbr:`ASN (Autonomous System Number)`.
:cachegroup:    A string that specifies the :ref:`cache-group-name` of a :term:`Cache Group` to which this :abbr:`ASN (Autonomous System Number)` will be assigned. If you do not pass this field, the cachegroup will be ``null``.
:cachegroupId:  The integral, unique identifier of the status of the :term:`Cache Group`.

.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/asns?id=1 HTTP/1.1
	User-Agent: python-requests/2.22.0
	Accept-Encoding: gzip, deflate
	Accept: */*
	Connection: keep-alive
	Cookie: mojolicious=...
	Content-Length: 53

	{
		"asn": 1,
		"cachegroup": "TRAFFIC_OPS",
		"cachegroupId": 2
	}

Response Structure
------------------
:asn:          An :abbr:`ASN (Autonomous System Number)` as specified by IANA for identifying a service provider
:cachegroup:   A string that is the :ref:`cache-group-name` of the :term:`Cache Group` that is associated with this :abbr:`ASN (Autonomous System Number)`
:cachegroupId: An integer that is the :ref:`cache-group-id` of the :term:`Cache Group` that is associated with this :abbr:`ASN (Autonomous System Number)`
:id:           An integral, unique identifier for this association between an :abbr:`ASN (Autonomous System Number)` and a :term:`Cache Group`
:lastUpdated:  The time and date this server entry was last updated in :rfc:`3339` Format

	.. versionchanged:: 5.0
		Prior to version 5.0 of the API, this field was in :ref:`non-rfc-datetime`.

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Encoding: gzip
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 25 Feb 2020 07:21:10 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: MjvwJg6AFbdqGPlAhK+2pfiN+VFjzgeNnhXoMVbh6+fRQYKeej6CCj3x09hwOl4uhp9d9RySrE/CQ3+L1b2VGQ==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 25 May 2023 21:59:33 GMT
	Content-Length: 164

	{
		"alerts": [
			{
				"text": "asn was updated.",
				"level": "success"
			}
		],
		"response": {
			"asn": 1,
			"cachegroup": "TRAFFIC_OPS",
			"cachegroupId": 2,
			"id": 1,
			"lastUpdated": "2023-05-25T15:59:33.7096-06:00"
		}
	}

``DELETE``
==========
Deletes an existing :abbr:`ASN (Autonomous System Number)`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: ASN:DELETE, ASN:READ, CACHE-GROUP:READ, CACHE-GROUP:UPDATE
:Response Type: ``undefined``

Request Structure
-----------------

.. code-block:: http
	:caption: Request Example

	DELETE /api/5.0/asns?id=1 HTTP/1.1
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
	Set-Cookie: mojolicious=...; Path=/; Expires=Tue, 25 Feb 2020 08:27:33 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: Woz8NSHIYVpX4V5X4xZWZIX1hvGL2uian7nUhjZ8F23Nb9RWQRMIg/cc+1vXEzkT/ehKV9t11FKRLX+avSae0g==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 25 May 2023 21:59:33 GMT
	Content-Length: 83

	{
		"alerts": [
			{
				"text": "asn was deleted.",
				"level": "success"
			}
		]
	}
