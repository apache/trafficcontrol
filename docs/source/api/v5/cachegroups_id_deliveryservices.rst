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

.. _to-api-cachegroups-id-deliveryservices:

***************************************
``cachegroups/{{ID}}/deliveryservices``
***************************************

``POST``
========
Assigns all of the "assignable" servers within a :term:`Cache Group` to one or more :term:`Delivery Services`.

.. note:: "Assignable" here means all of the :ref:`Cache Group's servers <cache-group-servers>` that have a :term:`Type` that matches one of the glob patterns ``EDGE*`` or ``ORG*``. If even one server of any :term:`Type` exists within the :term:`Cache Group` that is not assigned to the same CDN as the :term:`Delivery Service` to which an attempt is being made to assign them, the request will fail.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Permissions Required: CACHE-GROUP:UPDATE, DELIVERY-SERVICE:UPDATE, CACHE-GROUP:READ, DELIVERY-SERVICE:READ
:Response Type:  Object

Request Structure
-----------------
.. table::Request Path Parameters

	+------+-----------------------------------------------------------------------------------+
	| Name | Description                                                                       |
	+======+===================================================================================+
	|  ID  | The :ref:`cache-group-id` of the :term:`Cache Group` from which to assign servers |
	+------+-----------------------------------------------------------------------------------+

:deliveryServices:  The integral, unique identifiers of the :term:`Delivery Services` to which the :ref:`Cache Group's servers <cache-group-servers>` are being assigned

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/cachegroups/8/deliveryservices HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 25
	Content-Type: application/json

	{"deliveryServices": [2]}

Response Structure
------------------
:deliveryServices: An array of integral, unique identifiers for :term:`Delivery Services` to which the :ref:`Cache Group's servers <cache-group-servers>` have been assigned
:id:               An integer that is the :ref:`Cache Group's ID <cache-group-id>`
:serverNames:      An array of the (short) hostnames of all of the :term:`Cache Group`'s "assignable" :ref:`cache-group-servers`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: j/yH0gvJoaGjiLZU/0MA8o5He20O4aJ5wh1eF9ex6F6IBO1liM9Wk9RkWCw7sdiUHoy13/mf7gDntisZwzP7yw==
	X-Server-Name: traffic_ops_golang/
	Date: Wed, 14 Nov 2018 19:54:17 GMT
	Content-Length: 183

	{ "alerts": [
		{
			"text": "Delivery services successfully assigned to all the servers of cache group 8.",
			"level": "success"
		}
	],
	"response": {
		"id": 8,
		"serverNames": [
			"foo"
		],
		"deliveryServices": [
			2
		]
	}}
