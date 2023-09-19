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

.. _to-api-cdns-name-federations:

*****************************
``cdns/{{name}}/federations``
*****************************

``GET``
=======
Retrieves a list of federations in use by a specific CDN.

:Auth. Required: Yes
:Roles Required: None
:Permissions Required: CDN:READ, FEDERATION:READ, DELIVERY-SERVICE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------+------------------------------------------------------------------+
	| Name      | Description                                                      |
	+===========+==================================================================+
	| name      | The name of the CDN for which :term:`Federations` will be listed |
	+-----------+------------------------------------------------------------------+

.. table:: Request Query Parameters

	+-----------+---------------------------------------------------------------------------------------------------------------+
	| Name      | Description                                                                                                   |
	+===========+===============================================================================================================+
	| id        | Return only the :term:`Federation` that has this ID                                                           |
	+-----------+---------------------------------------------------------------------------------------------------------------+
	| cname     | Return only those :term:`Federations` that have this CNAME                                                    |
	+-----------+---------------------------------------------------------------------------------------------------------------+
	| dsID      | Return only those :term:`Federations` assigned to a :term:`Delivery Service` that has this ID                 |
	+-----------+---------------------------------------------------------------------------------------------------------------+
	| xmlID     | Return only those :term:`Federations` assigned to a :term:`Delivery Service` that has this :ref:`ds-xmlid`    |
	+-----------+---------------------------------------------------------------------------------------------------------------+
	| orderby   | Choose the ordering of the results - must be the name of one of the fields of the objects in the ``response`` |
	|           | array                                                                                                         |
	+-----------+---------------------------------------------------------------------------------------------------------------+
	| sortOrder | Changes the order of sorting. Either ascending (default or "asc") or descending ("desc")                      |
	+-----------+---------------------------------------------------------------------------------------------------------------+
	| limit     | Choose the maximum number of results to return                                                                |
	+-----------+---------------------------------------------------------------------------------------------------------------+
	| offset    | The number of results to skip before beginning to return results. Must use in conjunction with limit          |
	+-----------+---------------------------------------------------------------------------------------------------------------+
	| page      | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long   |
	|           | and the first page is 1. If ``offset`` was defined, this query parameter has no effect. ``limit`` must be     |
	|           | defined to make use of ``page``.                                                                              |
	+-----------+---------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/cdns/CDN-in-a-Box/federations HTTP/1.1
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

	{ "response": [
		{
			"id": 1,
			"cname": "test.quest.",
			"ttl": 68,
			"description": "A test federation",
			"lastUpdated": "2018-12-05T00:05:16Z",
			"deliveryService": {
				"id": 1,
				"xmlID": "demo1"
			}
		}
	]}

``POST``
========
Creates a new :term:`Federation`.

.. caution:: Despite the URL of this endpoint, this does `**not**` create a :term:`Federation` within any particular CDN. A :term:`Federation` is associated with a CDN purely because any :term:`Delivery Service` to which it is assigned is scoped to a CDN. Therefore, upon creation a :term:`Federation` is not associated with any CDN in particular.

.. warning:: There is no restriction on using the special "ALL" CDN for :term:`Federations` - but this is highly discouraged, because many things treat that CDN specially and may not work properly if it is used as though it were a normal CDN.

:Auth. Required: Yes
:Roles Required: "admin"
:Permissions Required: FEDERATION:CREATE, FEDERATION:READ, CDN:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+------------------------------------------------------------------------+
	| Name | Description                                                            |
	+======+========================================================================+
	| name | The name of the CDN for which a new :term:`Federation` will be created |
	+------+------------------------------------------------------------------------+

:cname: The :abbr:`CNAME (Canonical Name)` used by the :term:`Federation`

	.. tip:: The CNAME must end with a "``.``"

:description: An optional description of the federation
:ttl:         Time to Live (TTL) for the name record used for ``cname`` - minimum of 60

	.. versionchanged:: 5.0
		In earlier API versions, there is no enforced minimum (although Traffic Portal would never allow a value under 60).

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/cdns/CDN-in-a-Box/federations HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.62.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 72
	Content-Type: application/json

	{
		"cname": "test.quest.",
		"ttl": 68,
		"description": "A test federation"
	}


Response Structure
------------------
:id:          The integral, unique identifier of the :term:`Federation`
:cname:       The :abbr:`CNAME (Canonical Name)` used by the :term:`Federation`
:description: The description of the :term:`Federation`
:lastUpdated: The date and time at which this federation was last modified, in :RFC:`3339` format

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
	whole-content-sha512: rRsWAIhXzVlj8Hy+8aFjp4Jo1QGTK49m0N1AP5QDyyAZ1TfNIdgtcgiuehu7FiN1IPWRFiv6D9CygFYKGcVDOw==
	x-server-name: traffic_ops_golang/
	content-length: 192
	date: Wed, 05 Dec 2018 00:05:16 GMT

	{ "alerts": [
		{
			"text": "Federation was created",
			"level": "success"
		}
	],
	"response": {
		"id": 1,
		"cname": "test.quest.",
		"ttl": 68,
		"description": "A test federation",
		"lastUpdated": "2018-12-05T00:05:16Z"
	}}
