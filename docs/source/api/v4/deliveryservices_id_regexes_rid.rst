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

.. _to-api-v4-deliveryservices-id-regexes-rid:

*******************************************
``deliveryservices/{{ID}}/regexes/{{rID}}``
*******************************************

``PUT``
=======
Updates a routing regular expression.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE:UPDATE, DELIVERY-SERVICE:READ, TYPE:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------------------+
	| Name |                Description                                                        |
	+======+===================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service` being inspected   |
	+------+-----------------------------------------------------------------------------------+
	| rID  | The integral, unique identifier of the routing regular expression being inspected |
	+------+-----------------------------------------------------------------------------------+

:pattern: The actual regular expression

	.. warning:: Be sure that ``\``\ s are escaped, or the expression may not work as intended!

:setNumber: The order in which this regular expression should be checked
:type:      The integral, unique identifier of a routing regular expression type

.. code-block:: http
	:caption: Request Example

	PUT /api/4.0/deliveryservices/1/regexes/2 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 55
	Content-Type: application/json

	{
		"pattern": ".*\\.foo-bar\\..*",
		"type": 33,
		"setNumber": 1
	}

Response Structure
------------------
:id:        The integral, unique identifier of this regular expression
:pattern:   The actual regular expression - ``\``\ s are escaped
:setNumber: The order in which the regular expression is evaluated against requests
:type:      The integral, unique identifier of the type of this regular expression
:typeName:  The type of regular expression - determines that against which it will be evaluated

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: kS5dRzAhFKE7vfzHK7XVIwpMOjztksk9MU+qtj5YU/1oxVHmqNbJ12FeOOIJsZJCXbYlnBS04sCI95Sz5wed1Q==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 29 Nov 2018 17:54:58 GMT
	Content-Length: 188

	{ "alerts": [
		{
			"text": "Delivery service regex creation was successful.",
			"level": "success"
		}
	],
	"response": {
		"id": 2,
		"type": 33,
		"typeName": "PATH_REGEXP",
		"setNumber": 1,
		"pattern": ".*\\.foo-bar\\..*"
	}}



``DELETE``
==========
Deletes a routing regular expression.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE:UPDATE, DELIVERY-SERVICE:READ, TYPE:READ
:Response Type:  ``undefined``

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------------------+
	| Name |                Description                                                        |
	+======+===================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service` being inspected   |
	+------+-----------------------------------------------------------------------------------+
	| rID  | The integral, unique identifier of the routing regular expression being inspected |
	+------+-----------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	DELETE /api/4.0/deliveryservices/1/regexes/2 HTTP/1.1
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
	Whole-Content-Sha512: 8oEa78x7f/o39LIS98W6G+UqE6cX/Iw4v3mMHvbAs1iWHALuDYRz3VOtA6jzfGQKpB04Om8qaVG+zWRrBVoCmQ==
	X-Server-Name: traffic_ops_golang/
	Date: Thu, 29 Nov 2018 18:44:00 GMT
	Content-Length: 76

	{ "alerts": [
		{
			"text": "deliveryservice_regex was deleted.",
			"level": "success"
		}
	]}

.. [#tenancy] Users will only be able to view, delete and update regular expressions for the :term:`Delivery Services` their :term:`Tenant` is allowed to see.
