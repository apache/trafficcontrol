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

.. _to-api-v3-staticdnsentries:

********************
``staticdnsentries``
********************

``GET``
=======
Retrieve all static DNS entries configured within Traffic Control

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| Name              | Required | Description                                                                                                                                |
	+===================+==========+============================================================================================================================================+
	| address           | no       | Return only static DNS entries that operate on this address/:abbr:`CNAME (Canonical Name)`                                                 |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| cachegroup        | no       | Return only static DNS entries assigned to the :term:`Cache Group` that has this :ref:`cache-group-name`                                   |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| cachegroupId      | no       | Return only static DNS entries assigned to the :term:`Cache Group` that has this :ref:`cache-group-id`                                     |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| deliveryservice   | no       | Return only static DNS entries that apply within the domain of the :term:`Delivery Service` with this :ref:`ds-xmlid`                      |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| deliveryserviceId | no       | Return only static DNS entries that apply within the domain of the :term:`Delivery Service` identified by this integral, unique identifier |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| host              | no       | Return only static DNS entries that resolve this :abbr:`FQDN (Fully Qualified Domain Name)`                                                |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| id                | no       | Return only the static DNS entry with this integral, unique identifier                                                                     |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| ttl               | no       | Return only static DNS entries with this :abbr:`TTL (Time To Live)`                                                                        |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| type              | no       | Return only static DNS entries of this type                                                                                                |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| typeId            | no       | Return only static DNS entries of the type identified by this integral, unique identifier                                                  |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| sortOrder         | no       | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                                                   |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| limit             | no       | Choose the maximum number of results to return                                                                                             |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| offset            | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit                                       |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+
	| page              | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and the first page is 1.       |
	|                   |          | If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to make use of ``page``.                          |
	+-------------------+----------+--------------------------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/3.0/staticdnsentries?address=foo.bar HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:address:    If ``typeId`` identifies a ``CNAME`` type record, this is the Canonical Name (CNAME) of the server with a trailing period, otherwise it is the IP address to which ``host`` shall be resolved
:cachegroup: An optional string containing the :ref:`Name of a Cache Group <cache-group-name>` which will service this static DNS entry

	.. note:: This field has no effect, and is not used by any part of Traffic Control. It exists for legacy compatibility reasons.

:cachegroupId: An optional, integer that is the :ref:`ID of a Cache Group <cache-group-id>` which will service this static DNS entry

	.. note:: This field has no effect, and is not used by any part of Traffic Control. It exists for legacy compatibility reasons.

:deliveryservice:   The name of a :term:`Delivery Service` under the domain of which this static DNS entry shall be active
:deliveryserviceId: The integral, unique identifier of a :term:`Delivery Service` under the domain of which this static DNS entry shall be active
:host:              If ``typeId`` identifies a ``CNAME`` type record, this is an alias for the CNAME of the server, otherwise it is the Fully Qualified Domain Name (FQDN) which shall resolve to ``address``
:id:                An integral, unique identifier for this static DNS entry
:lastUpdated:       The date and time at which this static DNS entry was last updated
:ttl:               The :abbr:`TTL (Time To Live)` of this static DNS entry in seconds
:type:              The name of the type of this static DNS entry
:typeId:            The integral, unique identifier of the :term:`Type` of this static DNS entry

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: Px1zTH3ihg+hfmdADGcap0Juuud39fGsx5Y3CzqaFNmRwFu1ZLMzOsy0EN2pb7vpOtpI6/zeIUYAC3dbsBwOmA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 20:04:33 GMT
	Content-Length: 226

	{ "response": [
		{
			"address": "foo.bar.",
			"cachegroup": null,
			"cachegroupId": null,
			"deliveryservice": "demo1",
			"deliveryserviceId": 1,
			"host": "test",
			"id": 2,
			"lastUpdated": "2018-12-10 19:59:56+00",
			"ttl": 300,
			"type": "CNAME_RECORD",
			"typeId": 41
		}
	]}

``POST``
========
Creates a new, static DNS entry.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
:address:      If ``typeId`` identifies a ``CNAME`` type record, this is the Canonical Name (CNAME) of the server with a trailing period, otherwise it is the IP address to which ``host`` shall be resolved
:cachegroupId: An optional, integer that is the :ref:`ID of a Cache Group <cache-group-id>` which will service this static DNS entry

	.. note:: This field has no effect, and is not used by any part of Traffic Control. It exists for legacy compatibility reasons.

:deliveryserviceId: The integral, unique identifier of a :term:`Delivery Service` under the domain of which this static DNS entry shall be active
:host:              If ``typeId`` identifies a ``CNAME`` type record, this is an alias for the CNAME of the server, otherwise it is the :abbr:`FQDN (Fully Qualified Domain Name)` which shall resolve to ``address``
:ttl:               The :abbr:`TTL (Time To Live)` of this static DNS entry in seconds
:typeId:            The integral, unique identifier of the :term:`Type` of this static DNS entry

.. code-block:: http
	:caption: Request Example

	POST /api/3.0/staticdnsentries HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 92
	Content-Type: application/json

	{
		"address": "test.quest.",
		"deliveryserviceId": 1,
		"host": "test",
		"ttl": 300,
		"typeId": 41
	}

Response Structure
------------------
:address:           If ``typeId`` identifies a ``CNAME`` type record, this is the Canonical Name (CNAME) of the server with a trailing period, otherwise it is the IP address to which ``host`` shall be resolved
:cachegroup: An optional string containing the :ref:`Name of a Cache Group <cache-group-name>` which will service this static DNS entry

	.. note:: This field has no effect, and is not used by any part of Traffic Control. It exists for legacy compatibility reasons.

:cachegroupId: An optional, integer that is the :ref:`ID of a Cache Group <cache-group-id>` which will service this static DNS entry

	.. note:: This field has no effect, and is not used by any part of Traffic Control. It exists for legacy compatibility reasons.

:deliveryservice:   The name of a :term:`Delivery Service` under the domain of which this static DNS entry shall be active
:deliveryserviceId: The integral, unique identifier of a :term:`Delivery Service` under the domain of which this static DNS entry shall be active
:host:              If ``typeId`` identifies a ``CNAME`` type record, this is an alias for the CNAME of the server, otherwise it is the Fully Qualified Domain Name (FQDN) which shall resolve to ``address``
:id:                An integral, unique identifier for this static DNS entry
:lastUpdated:       The date and time at which this static DNS entry was last updated
:ttl:               The :abbr:`TTL (Time To Live)` of this static DNS entry in seconds
:type:              The name of the :term:`Type` of this static DNS entry
:typeId:            The integral, unique identifier of the :term:`Type` of this static DNS entry

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: 8dcJyjw2NJZx0L9Oz16P7g/7j5A1jlpyiY6Y+rRVQ2wGcwYI3yiGPrz6ur0qKzgqEBBsh8aPF44WTHAR9jUJdg==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 19:54:19 GMT
	Content-Length: 282

	{ "alerts": [
		{
			"text": "staticDNSEntry was created.",
			"level": "success"
		}
	],
	"response": {
		"address": "test.quest.",
		"cachegroup": null,
		"cachegroupId": null,
		"deliveryservice": null,
		"deliveryserviceId": 1,
		"host": "test",
		"id": 2,
		"lastUpdated": "2018-12-10 19:54:19+00",
		"ttl": 300,
		"type": "CNAME_RECORD",
		"typeId": 41
	}}

``PUT``
=======
Updates a static DNS entry.

:Auth. Required:   Yes
:Role(s) Required: "admin" or "operator"
:Response Type:    Object

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+-------------------------------------------------------------------+
	| Name | Description                                                       |
	+======+===================================================================+
	|  id  | The integral, unique identifier of the static DNS entry to modify |
	+------+-------------------------------------------------------------------+

:address:           If ``typeId`` identifies a ``CNAME`` type record, this is the Canonical Name (CNAME) of the server with a trailing period, otherwise it is the IP address to which ``host`` shall be resolved
:cachegroupId: An optional, integer that is the :ref:`ID of a Cache Group <cache-group-id>` which will service this static DNS entry

	.. note:: This field has no effect, and is not used by any part of Traffic Control. It exists for legacy compatibility reasons.

:deliveryserviceId: The integral, unique identifier of a :term:`Delivery Service` under the domain of which this static DNS entry shall be active
:host:              If ``typeId`` identifies a ``CNAME`` type record, this is an alias for the CNAME of the server, otherwise it is the Fully Qualified Domain Name (FQDN) which shall resolve to ``address``
:ttl:               The :abbr:`TTL (Time To Live)` of this static DNS entry in seconds
:typeId:            The integral, unique identifier of the :term:`Type` of this static DNS entry

.. code-block:: http
	:caption: Request Example

	PUT /api/3.0/staticdnsentries?id=2 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 89
	Content-Type: application/json

	{
		"address": "foo.bar.",
		"deliveryserviceId": 1,
		"host": "test",
		"ttl": 300,
		"typeId": 41
	}

Response Structure
------------------
:address:    If ``typeId`` identifies a ``CNAME`` type record, this is the Canonical Name (CNAME) of the server with a trailing period, otherwise it is the IP address to which ``host`` shall be resolved
:cachegroup: An optional string containing the :ref:`Name of a Cache Group <cache-group-name>` which will service this static DNS entry

	.. note:: This field has no effect, and is not used by any part of Traffic Control. It exists for legacy compatibility reasons.

:cachegroupId: An optional, integer that is the :ref:`ID of a Cache Group <cache-group-id>` which will service this static DNS entry

	.. note:: This field has no effect, and is not used by any part of Traffic Control. It exists for legacy compatibility reasons.

:deliveryservice:   The name of a :term:`Delivery Service` under the domain of which this static DNS entry shall be active
:deliveryserviceId: The integral, unique identifier of a :term:`Delivery Service` under the domain of which this static DNS entry shall be active
:host:              If ``typeId`` identifies a ``CNAME`` type record, this is an alias for the CNAME of the server, otherwise it is the :abbr:`FQDN (Fully Qualified Domain Name)` which shall resolve to ``address``
:id:                An integral, unique identifier for this static DNS entry
:lastUpdated:       The date and time at which this static DNS entry was last updated
:ttl:               The :abbr:`TTL (Time To Live)` of this static DNS entry in seconds
:type:              The name of the :term:`Type` of this static DNS entry
:typeId:            The integral, unique identifier of the :term:`Type` of this static DNS entry

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: +FaYmpnlIIzVSBq0nosw29NZcV9xFhlVgWuUqXUyiDihVUSzX4jrdAloRDgzDvKsYQB8LSkPdGHwt1zjgSzUtA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 19:59:56 GMT
	Content-Length: 279

	{ "alerts": [
		{
			"text": "staticDNSEntry was updated.",
			"level": "success"
		}
	],
	"response": {
		"address": "foo.bar.",
		"cachegroup": null,
		"cachegroupId": null,
		"deliveryservice": null,
		"deliveryserviceId": 1,
		"host": "test",
		"id": 2,
		"lastUpdated": "2018-12-10 19:59:56+00",
		"ttl": 300,
		"type": "CNAME_RECORD",
		"typeId": 41
	}}


``DELETE``
==========
Delete staticdnsentries.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Query Parameters

	+------+-------------------------------------------------------------------+
	| Name | Description                                                       |
	+======+===================================================================+
	|  id  | The integral, unique identifier of the static DNS entry to delete |
	+------+-------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/3.0/staticdnsentries?id=2 HTTP/1.1
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
	Whole-Content-Sha512: g6uqHPU44LuTtqU2ahtazrVCpcpNWVc9kWJQOYRuiVLDnsm39KOB/xt3XM6j0/X3WYiIawnNspkxRC85LJHwFA==
	X-Server-Name: traffic_ops_golang/
	Date: Mon, 10 Dec 2018 20:05:52 GMT
	Content-Length: 69

	{ "alerts": [
		{
			"text": "staticDNSEntry was deleted.",
			"level": "success"
		}
	]}
