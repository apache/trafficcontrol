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

.. _to-api-deliveryservices-id-regexes-rid:

*******************************************
``deliveryservices/{{ID}}/regexes/{{rID}}``
*******************************************

``GET``
=======
Retrieves a specific routing regular expression for a specific Delivery Service.

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------------------+
	| Name |                Description                                                        |
	+======+===================================================================================+
	|  ID  | The integral, unique identifier of the Delivery Service being inspected           |
	+------+-----------------------------------------------------------------------------------+
	| rID  | The integral, unique identifier of the routing regular expression being inspected |
	+------+-----------------------------------------------------------------------------------+

.. code-block:: http
	:caption: Request Example

	GET /api/1.4/deliveryservices/1/regexes/1 HTTP/1.1
	Host: trafficops.infra.ciab.test
	User-Agent: curl/7.47.0
	Accept: */*
	Cookie: mojolicious=...

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
	Set-Cookie: mojolicious=...; Path=/; HttpOnly
	Whole-Content-Sha512: fW9Fde4WRpp2ShRAC41P9s/PhU71LI/SEzHgYjGqfzhk45wq0kpaWy76JvPfLpowY8eDTp8Y8TL5rNGEc+bM+A==
	X-Server-Name: traffic_ops_golang/
	Date: Tue, 27 Nov 2018 21:08:34 GMT
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

.. [1] If tenancy is used, then users (regardless of role) will only be able to see the routing regular expressions used by Delivery Services their tenant has permissions to see.


``PUT``
=======
Updates a routing regular expression.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [2]_

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------------------+
	| Name |                Description                                                        |
	+======+===================================================================================+
	|  ID  | The integral, unique identifier of the Delivery Service being inspected           |
	+------+-----------------------------------------------------------------------------------+
	| rID  | The integral, unique identifier of the routing regular expression being inspected |
	+------+-----------------------------------------------------------------------------------+

:pattern: The actual regular expression

	.. warning:: Be sure that ``\``\ s are escaped, or the expression may not work as intended!

:setNumber: The order in which this regular expression should be checked
:type:      The integral, unique identifier of a routing regular expression type

.. code-block:: json
	:caption: Request Example

	{
		"pattern": ".*\\.foo-bar\\..*",
		"type": 18,
		"setNumber": 0
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

	{ "response":{
		"id": 852,
		"type": 18,
		"typeName": "HOST_REGEXP",
		"pattern": ".*\\.foo-bar\\..*",
		"setNumber": 0
	},
	"alerts":[
		{
			"level": "success",
			"text": "Delivery service regex update was successful."
		}
	]}

.. [2] If tenancy is used, then users (regardless of role) will only be able to edit the routing regular expressions used by Delivery Services their tenant has permissions to edit. Assuming tenancy is satisfied, a routing regular expression can only be edited by a user with the "admin" or "operations" role.

``DELETE``
==========
Deletes a routing regular expression.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [3]_

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+-----------------------------------------------------------------------------------+
	| Name |                Description                                                        |
	+======+===================================================================================+
	|  ID  | The integral, unique identifier of the Delivery Service being inspected           |
	+------+-----------------------------------------------------------------------------------+
	| rID  | The integral, unique identifier of the routing regular expression being inspected |
	+------+-----------------------------------------------------------------------------------+

Response Structure
------------------
.. code-block:: http
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "Delivery service regex delete was successful."
		}
	]}

.. [3] If tenancy is used, then users (regardless of role) will only be able to delete the routing regular expressions used by Delivery Services their tenant has permissions to delete. Assuming tenancy is satisfied, a routing regular expression can only be deleted by a user with the "admin" or "operations" role.
