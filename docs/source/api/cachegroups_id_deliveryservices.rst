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
Assigns a :term:`Cache Group` to one or more :term:`Delivery Services`

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table::Request Path Parameters

	+------+---------------------------------------------------------------------------+
	| Name |           Description                                                     |
	+======+===========================================================================+
	|  ID  | The integral, unique identifier of the :term:`Cache Group` being assigned |
	+------+---------------------------------------------------------------------------+

:deliveryServices:  The integral, unique identifiers of the :term:`Delivery Services` to which the :term:`Cache Group` is being assigned

.. code-block:: http
	:caption: Request Example

	POST /api/1.3/cachegroups/8/deliveryservices HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 25
	Content-Type: application/x-www-form-urlencoded

	{"deliveryServices": [2]}

Response Structure
------------------
:deliveryServices: An array of *all* :term:`Delivery Services` to which the :term:`Cache Group` is assigned (**not** just the one(s) to which it was assigned via the request)
:id:               The :term:`Cache Group`\ 's ID
:serverNames:      An array of the (short) hostnames of all servers in the :term:`Cache Group`

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
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

