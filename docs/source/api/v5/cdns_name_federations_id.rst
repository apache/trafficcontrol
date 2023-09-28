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

.. _to-api-cdns-name-federations-id:

************************************
``cdns/{{name}}/federations/{{ID}}``
************************************

``GET``
=======
Retrieves a list of federations in use by a specific CDN.

.. versionadded:: 5.0

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: CDN:READ, FEDERATION:READ, DELIVERY-SERVICE:READ
:Response Type: Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------------------+
	| Name | Description                                                                                 |
	+======+=============================================================================================+
	| name | The name of the CDN for which the :term:`Federation` identified by ``ID`` will be inspected |
	+------+---------------------------------------------------------------------------------------------+
	|  ID  | An integral, unique identifier for the :term:`Federation` to be inspected                   |
	+------+---------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/cdns/CDN-in-a-Box/federations/1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cname:           The :abbr:`CNAME (Canonical Name)` used by the :term:`Federation`
:deliveryService: An object with keys that provide identifying information for the :term:`Delivery Service` using this :term:`Federation`

	:id:    The integral, unique identifier for the :term:`Delivery Service`
	:xmlID: The :term:`Delivery Service`'s uniquely identifying :ref:`ds-xmlid`

		.. versionchanged:: 5.0
			Prior to version 5, this field was known by the name ``xmlId`` - improperly formatted camelCase.

:description: A human-readable description of the :term:`Federation`. This can be ``null`` as well as an empty string.
:lastUpdated: The date and time at which this :term:`Federation` was last modified, in :RFC:`3339` format

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was given in :ref:`non-rfc-datetime`.

:ttl: :abbr:`TTL (Time to Live)` for the ``cname``, in hours

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	content-type: application/json
	set-cookie: mojolicious=...; Path=/; HttpOnly
	whole-content-sha512: SJA7G+7G5KcOfCtnE3Dq5DCobWtGRUKSppiDkfLZoG5+paq4E1aZGqUb6vGVsd+TpPg75MLlhyqfdfCHnhLX/g==
	x-server-name: traffic_ops_golang/
	content-length: 170
	date: Wed, 05 Dec 2018 00:35:40 GMT

	{ "response": {
		"id": 1,
		"cname": "test.quest.",
		"ttl": 68,
		"description": "A test federation",
		"lastUpdated": "2018-12-05T00:05:16Z",
		"deliveryService": {
			"id": 1,
			"xmlID": "demo1"
		}
	}}

``PUT``
=======
Updates a :term:`Federation`.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: FEDERATION:UPDATE, FEDERATION:READ, CDN:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------------------------------------+
	| Name | Description                                                                               |
	+======+===========================================================================================+
	| name | The name of the CDN for which the :term:`Federation` identified by ``ID`` will be updated |
	+------+-------------------------------------------------------------------------------------------+
	|  ID  | An integral, unique identifier for the :term:`Federation` to be updated                   |
	+------+-------------------------------------------------------------------------------------------+

.. caution:: The name of the CDN doesn't actually matter. It doesn't even need to be the name of any existing CDN.

:cname: The :abbr:`CNAME (Canonical Name)` used by the :term:`Federation`

	.. note:: The CNAME must end with a "``.``"

:description: An optional description of the federation
:ttl:         Time to Live (TTL) for the name record used for ``cname`` - minimum of 60

	.. versionchanged:: 5.0
		In earlier API versions, there is no enforced minimum (although Traffic Portal would never allow a value under 60).

.. code-block:: http
	:caption: Request Example

	PUT /api/5.0/cdns/CDN-in-a-Box/federations/1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 33
	Content-Type: application/json

	{
		"cname": "foo.bar.",
		"ttl": 68
	}


Response Structure
------------------
:cname:       The :abbr:`CNAME (Canonical Name)` used by the :term:`Federation`
:description: A human-readable description of the :term:`Federation`. This can be ``null`` as well as an empty string.
:lastUpdated: The date and time at which this :term:`Federation` was last modified, in :RFC:`3339` format

	.. versionchanged:: 5.0
		In earlier versions of the API, this field was given in :ref:`non-rfc-datetime`.

:ttl: :abbr:`TTL (Time to Live)` for the ``cname``, in hours

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	content-type: application/json
	set-cookie: mojolicious=...; Path=/; HttpOnly
	whole-content-sha512: qcjfQ+gDjNxYQ1aq+dlddgrkFWnkFYxsFF+SHDqqH0uVHBVksmU0aTFgltozek/u6wbrGoR1LFf9Fr1C1SbigA==
	x-server-name: traffic_ops_golang/
	content-length: 174
	date: Wed, 05 Dec 2018 01:03:40 GMT

	{ "alerts": [
		{
			"text": "Federation was updated",
			"level": "success"
		}
	],
	"response": {
		"id": 1,
		"cname": "foo.bar.",
		"ttl": 68,
		"description": null,
		"lastUpdated": "2018-12-05T01:03:40Z"
	}}


``DELETE``
==========
Deletes a specific federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: FEDERATION:DELETE, FEDERATION:READ, CDN:READ
:Response Type:  Object

.. versionchanged:: 5.0
	In earlier API versions, no ``response`` property is present in these responses.

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-------------------------------------------------------------------------------------------+
	| Name | Description                                                                               |
	+======+===========================================================================================+
	| name | The name of the CDN for which the :term:`Federation` identified by ``ID`` will be deleted |
	+------+-------------------------------------------------------------------------------------------+
	|  ID  | An integral, unique identifier for the :term:`Federation` to be deleted                   |
	+------+-------------------------------------------------------------------------------------------+

.. caution:: The name of the CDN doesn't actually matter. It doesn't even need to be the name of any existing CDN.

.. code-block:: http
	:caption: Request Example

	DELETE /api/5.0/cdns/CDN-in-a-Box/federations/1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:cname:       The :abbr:`CNAME (Canonical Name)` used by the :term:`Federation`
:description: A human-readable description of the :term:`Federation`. This can be ``null`` as well as an empty string.
:lastUpdated: The date and time at which this :term:`Federation` was last modified, in :RFC:`3339` format
:ttl:         :abbr:`TTL (Time to Live)` for the ``cname``, in hours

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	access-control-allow-credentials: true
	access-control-allow-headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	access-control-allow-methods: POST,GET,OPTIONS,PUT,DELETE
	access-control-allow-origin: *
	content-type: application/json
	set-cookie: mojolicious=...; Path=/; HttpOnly
	whole-content-sha512: Cnkfj6dmzTD3if9oiDq33tqf7CnAflKK/SPgqJyfu6HUfOjLJOgKIZvkcs2wWY6EjLVdw5qsatsd4FPoCyjvcw==
	x-server-name: traffic_ops_golang/
	content-length: 68
	date: Wed, 05 Dec 2018 01:17:24 GMT

	{ "alerts": [
		{
			"text": "Federation was deleted",
			"level": "success"
		}
	],
	"response": {
		"id": 1,
		"cname": "foo.bar.",
		"ttl": 68,
		"description": null,
		"lastUpdated": "2018-12-05T01:03:40Z"
	}}
