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

	**Request Example** ::

		{
				"userId": 50,
				"deliveryServices": [ 23, 34, 45, 56, 67 ],
				"replace": true
		}

	**Response Properties**

	+------------------------------------+--------+-------------------------------------------------------------------+
	| Parameter                          | Type   | Description                                                       |
	+====================================+========+===================================================================+
	| ``userId``                         | int    | The ID of the user.                                               |
	+------------------------------------+--------+-------------------------------------------------------------------+
	| ``deliveryServices``               | array  | An array of delivery service IDs.                                 |
	+------------------------------------+--------+-------------------------------------------------------------------+
	| ``replace``                        | array  | Existing user/ds assignments replaced? (true|false).              |
	+------------------------------------+--------+-------------------------------------------------------------------+

	**Response Example** ::

		{
				"alerts": [
									{
													"level": "success",
													"text": "Delivery service assignments complete."
									}
					],
				"response": {
						"userId" : 50,
						"deliveryServices": [ 23, 34, 45, 56, 67 ],
						"replace": true
				}
		}

|
