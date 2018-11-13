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

.. _to-api-deliveryservice_user:

************************
``deliveryservice_user``
************************


``POST``
========
Assigns one or more Delivery Services to a user.

:Auth. Required: Yes
:Roles Required: "admin" or "operations"
:Response Type:  Object

Request Structure
-----------------
:userId:           An integral, unique identifier for the user to whom the Delivery Service(s) identified in ``deliveryServices`` will be assigned
:deliveryServices: An array of integral, unique identifiers for the Delivery Service(s) being assigned to the user identified by ``userId``
:replace:          An optional field which, when present and ``true`` will replace existing user/ds assignments? (true|false)

.. code-block:: http
	:caption: Request Example

	POST /api/1.3/deliveryservice_user HTTP/1.1
	Content-Type: application/json
	Content-Length: 81
	Accept: application/json

	{
		"userId": 50,
		"deliveryServices": [ 23, 34, 45, 56, 67 ],
		"replace": true
	}

Response Structure
------------------
:userId:           The integral, unique identifier of the user to whom the Delivery Service(s) identified in ``deliveryServices`` are assigned
:deliveryServices: An array of integral, unique identifiers of Delivery Services assigned to the user identified by ``userId``
:replace:          If ``true``, any and all existing, conflicting Delivery Service assignments were overwritten by this assignment operation

.. code-block:: http
	:caption: Response Example

	HTTP/1.1 200 OK
	Content-Type: application/json

	{ "alerts": [
		{
			"level": "success",
			"text": "Delivery service assignments complete."
		}
	],
	"response": {
			"userId" : 50,
			"deliveryServices": [ 23, 34, 45, 56, 67 ],
			"replace": true
	}}
