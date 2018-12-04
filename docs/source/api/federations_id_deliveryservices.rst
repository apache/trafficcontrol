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

.. _to-api-federations-id-deliveryservices:

***************************************
``federations/{{ID}}/deliveryservices``
***************************************

``GET``
=======
Retrieves Delivery Services assigned to a federation.

:Auth. Required: Yes
:Roles Required: None
:Response Type:  Array

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------+
	| Name |                 Description                                        |
	+======+====================================================================+
	|  ID  | The integral, unique identifier for the federation to be inspected |
	+------+--------------------------------------------------------------------+

Response Structure
------------------
:cdn:   The CDN to which this Delivery Service Belongs
:id:    The integral, unique identifier for the Deliver Service
:type:  The routing type used by this Delivery Service
:xmlId: The 'xml_id' which uniquely identifies this Delivery Service

.. code-block:: json
	:caption: Response Example

	{ "response": [
		{
			"id": 41
			"cdn": "cdn1",
			"type": "DNS",
			"xmlId": "booya-12"
		}
	]}

``POST``
========
Assigns one or more Delivery Services to a federation.

:Auth. Required: Yes
:Roles Required: "admin"
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+--------------------------------------------------------------------+
	| Name |                 Description                                        |
	+======+====================================================================+
	|  ID  | The integral, unique identifier for the federation to be inspected |
	+------+--------------------------------------------------------------------+

:dsIds:   An array of integral, unique identifiers for Delivery Services which will be assigned to this federation
:replace: An optional boolean (default: ``false``) which, if ``true``, will cause any conflicting assignments already in place to be overridden by this request

	.. note:: If ``replace`` is not given (and/or not ``true``), then any conflicts with existing assignments will cause the entire operation to fail.

.. code-block:: json
	:caption: Request Example

	{
		"dsIds": [ 2, 3, 4, 5, 6 ],
		"replace": true
	}

Response Structure
------------------
:dsIds:   An array of integral, unique identifiers for Delivery Services which are now assigned to this federation
:replace: An optional boolean (default: ``false``) which, if ``true``, means any conflicting assignments already in place were overridden by this request

.. code-block:: json
	:caption: Response Example

	{ "alerts": [
		{
			"level": "success",
			"text": "5 delivery service(s) were assigned to the cname. federation"
		}
	],
	"response": {
		"dsIds" : [ 2, 3, 4, 5, 6 ],
		"replace" : true
	}}
