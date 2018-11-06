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

***********************************
cachegroups/{{ID}}/deliveryservices
***********************************

``POST``
========
Assigns a Cache Group to one or more Delivery Services

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
.. table::Request Path Parameters

	+------------------+----------+------------------------------------------------------------------------------+
	|      Name        | Required |           Description                                                        |
	+==================+==========+==============================================================================+
	|      id          |   yes    | The integral, unique identifier of the Cache Group being assigned            |
	+------------------+----------+------------------------------------------------------------------------------+

.. table:: Request Data Parameters

	+------------------+----------+--------------------------------------------------------------------------------------+
	|    Parameter     |   Type   |           Description                                                                |
	+==================+==========+======================================================================================+
	| deliveryServices |  array   | The integral IDs of the Delivery Services to which the Cache Group is being assigned |
	+------------------+----------+--------------------------------------------------------------------------------------+

Response Structure
------------------
:deliveryServices: An array of *all* Delivery Services to which the Cache Group is assigned (**not** just the one(s) to which it was assigned via the request)
:id:               The Cache Group's ID
:serverNames:      An array of the (short) hostnames of all servers in the Cache Group

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"text": "Delivery services successfully assigned to all the servers of cache group 7.",
			"level": "success"
		}
	],
	"response": {
		"id": 7,
		"serverNames": [ "edge" ],
		"deliveryServices": [ 1 ]
	}}
