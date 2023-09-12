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

.. _to-api-v4-deliveryserviceserver:

*************************
``deliveryserviceserver``
*************************

``GET``
=======
Retrieve information about the assignment of servers to :term:`Delivery Services`

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Permissions Required: SERVER:READ, DELIVERY-SERVICE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Query Parameters

	+-----------+----------+-------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	|    Name   | Required | Default           |                                                       Description                                                                                                            |
	+===========+==========+===================+==============================================================================================================================================================================+
	| cdn       | no       | None              | Limit the results to delivery service servers for the given CDN name                                                                                                         |
	+-----------+----------+-------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| page      | no       | 0                 | The page number for use in pagination - ``0`` means "no pagination"                                                                                                          |
	+-----------+----------+-------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| limit     | no       | 20                | Limits the results to a maximum of this number - if pagination is used, this defines the number of results per page                                                          |
	+-----------+----------+-------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
	| orderby   | no       | "deliveryService" | Choose the ordering of the results - the value must either be the name of one of the fields of the objects in the ``response`` array or be empty to skip ordering altogether |
	+-----------+----------+-------------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/4.0/deliveryserviceserver?page=1&limit=2&orderby=lastUpdated HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...


Response Structure
------------------
Unlike most API endpoints, this will return a JSON response body containing both a "response" object as well as other, top-level fields (besides the optional "alerts" field). For this reason, this section contains a "response" key, which normally is implicit.

.. seealso:: :ref:`to-api-response-structure`

:limit:    The maximum size of the ``response`` array, also indicative of the number of results per page using the pagination requested by the query parameters (if any) - this should be the same as the ``limit`` query parameter (if given)
:orderby:  A string that names the field by which the elements of the ``response`` array are ordered - should be the same as the ``orderby`` request query parameter (if given)
:response: An array of objects, each of which represents a server's :term:`Delivery Service` assignment

	:deliveryService: The integral, unique identifier of the :term:`Delivery Service` to which the server identified by ``server`` is assigned
	:lastUpdated:     The date and time at which the server's assignment to a :term:`Delivery Service` was last updated, in :ref:`non-rfc-datetime`
	:server:          The integral, unique identifier of a server which is assigned to the :term:`Delivery Service` identified by ``deliveryService``

:size: The page number - if pagination was requested in the query parameters, else ``0`` to indicate no pagination - of the results represented by the ``response`` array. This is named "size" for legacy reasons


.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: J7sK8PohQWyTpTrMjjrWdlJwPj+Zyep/xutM25uVosL6cHgi30nXa6VMyOC5Y3vd9r5KLES8rTgR+qUQcZcJ/A==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 01 Nov 2018 14:27:45 GMT
	Content-Length: 129

	{ "orderby": "lastUpdated",
	"response": [
		{
			"server": 8,
			"deliveryService": 1,
			"lastUpdated": "2018-11-01 14:10:38+00"
		}
	],
	"size": 1,
	"limit": 2
	}

.. [1] While no roles are required, this endpoint *does* respect tenancy permissions (pending `GitHub Issue #2978 <https://github.com/apache/trafficcontrol/issues/2978>`_\ ).

``POST``
========
Assign a set of one or more servers to a :term:`Delivery Service`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [2]_
:Permissions Required: DELIVERY-SERVICE:READ, SERVER:READ, SERVER:UPDATE, DELIVERY-SERVICE:UPDATE
:Response Type:  Object

Request Structure
-----------------
:dsId:    The integral, unique identifier of the :term:`Delivery Service` to which the servers identified in the ``servers`` array will be assigned
:replace: If ``true``, any existing assignments for a server identified in the ``servers`` array will be overwritten by this request
:servers: An array of integral, unique identifiers for servers which are to be assigned to the :term:`Delivery Service` identified by ``deliveryService``

.. code-block:: http
	:caption: Request Example

	POST /api/4.0/deliveryserviceserver HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 46
	Content-Type: application/x-www-form-urlencoded

	dsId=1&replace=true&servers=12

Response Structure
------------------
:dsId:    The integral, unique identifier of the :term:`Delivery Service` to which the servers identified by the elements of the ``servers`` array have been assigned
:replace: If ``true``, any existing assignments for a server identified in the ``servers`` array have been overwritten by this request
:servers: An array of integral, unique identifiers for servers which have been assigned to the :term:`Delivery Service` identified by ``deliveryService``

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: D+HhGhoxzaxvka9vZIStoaOZUpX23nz7zZnMbpFHNRO3MawyEaSb3GVUHQyCv6sDgwhpZZjRggDmctGCw88flg==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 01 Nov 2018 14:12:49 GMT
	Content-Length: 123

	{ "alerts": [
		{
			"text": "server assignements complete",
			"level": "success"
		}
	],
	"response": {
		"dsId": 1,
		"replace": false,
		"servers": [ 12 ]
	}}


.. [2] Users with the "admin" or "operations" roles will be able to modify ALL server-to-Delivery-Service assignments, whereas all other users can only assign servers to the :term:`Delivery Services` their Tenant has permissions to edit.
