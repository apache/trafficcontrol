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

.. _to-api-deliveryservices-id-regexes:

***********************************
``deliveryservices/{{ID}}/regexes``
***********************************

``GET``
=======
Retrieves routing regular expressions for a specific :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: None\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE:READ, TYPE:READ
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------+
	| Name | Description                                                                     |
	+======+=================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service` being inspected |
	+------+---------------------------------------------------------------------------------+

.. table:: Request Query Parameters

	+-------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| Name        | Required | Description                                                                                                                          |
	+=============+==========+======================================================================================================================================+
	| id          | no       | Show only the :term:`Delivery Service` regular expression that has this integral, unique identifier                                  |
	+-------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| limit       | no       | Choose the maximum number of results to return                                                                                       |
	+-------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| offset      | no       | The number of results to skip before beginning to return results. Must use in conjunction with limit                                 |
	+-------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+
	| page        | no       | Return the n\ :sup:`th` page of results, where "n" is the value of this parameter, pages are ``limit`` long and the first page is 1. |
	|             |          | If ``offset`` was defined, this query parameter has no effect. ``limit`` must be defined to make use of ``page``.                    |
	+-------------+----------+--------------------------------------------------------------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/5.0/deliveryservices/1/regexes HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

Response Structure
------------------
:id:        The integral, unique identifier of this regular expression
:pattern:   The actual regular expression - ``\``\ s are escaped
:setNumber: The order in which the regular expression is evaluated against requests
:type:      The integral, unique identifier of the :term:`Type` of this regular expression
:typeName:  The :term:`Type` of regular expression - determines that against which it will be evaluated

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Access-Control-Allow-Credentials: true
	Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept, Set-Cookie, Cookie
	Access-Control-Allow-Methods: POST,GET,OPTIONS,PUT,DELETE
	Access-Control-Allow-Origin: *
	Content-Type: application/json
	Set-Cookie: mojolicious=...; Path=/; Expires=Mon, 18 Nov 2019 17:40:54 GMT; Max-Age=3600; HttpOnly
	Whole-Content-Sha512: fW9Fde4WRpp2ShRAC41P9s/PhU71LI/SEzHgYjGqfzhk45wq0kpaWy76JvPfLpowY8eDTp8Y8TL5rNGEc+bM+A==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 27 Nov 2018 20:56:43 GMT
	Content-Length: 100

	{ "response": [
		{
			"id": 1,
			"type": 31,
			"typeName": "HOST_REGEXP",
			"setNumber": 0,
			"pattern": ".*\\.demo1\\..*"
		}
	]}


``POST``
========
Creates a routing regular expression for a :term:`Delivery Service`.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [#tenancy]_
:Permissions Required: DELIVERY-SERVICE:UPDATE, DELIVERY-SERVICE:READ, TYPE:READ
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+---------------------------------------------------------------------------------+
	| Name |                Description                                                      |
	+======+=================================================================================+
	|  ID  | The integral, unique identifier of the :term:`Delivery Service` being inspected |
	+------+---------------------------------------------------------------------------------+

:pattern: The actual regular expression

	.. warning:: Be sure that ``\``\ s are escaped, or the expression may not work as intended!

:setNumber: The order in which this regular expression should be checked
:type:      The integral, unique identifier of a routing regular expression type

.. code-block:: http
	:caption: Request Example

	POST /api/5.0/deliveryservices/1/regexes HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...
	Content-Length: 55
	Content-Type: application/json

	{
		"pattern": ".*\\.foo-bar\\..*",
		"type": 31,
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
	Date: Wed, 28 Nov 2018 17:00:42 GMT
	Content-Length: 188

	{ "alerts": [
		{
			"text": "Delivery service regex creation was successful.",
			"level": "success"
		}
	],
	"response": {
		"id": 2,
		"type": 31,
		"typeName": "HOST_REGEXP",
		"setNumber": 1,
		"pattern": ".*\\.foo-bar\\..*"
	}}


.. [#tenancy] Users will only be able to view and create regular expressions for the :term:`Delivery Services` their :term:`Tenant` is allowed to see.
