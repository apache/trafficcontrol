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

.. _to-api-deliveryservices-id-state:

*****************************
deliveryservices/{{ID}}/state
*****************************

``GET``
=======
Retrieves the fail-over state for a Delivery Service.

:Auth. Required: Yes
:Roles Required: None\ [1]_
:Response Type:
	**Response Properties**

	+------------------+---------+-------------------------------------------------+
	|    Parameter     |  Type   |                   Description                   |
	+==================+=========+=================================================+
	| ``failover``     |  hash   |                                                 |
	+------------------+---------+-------------------------------------------------+
	| ``>locations``   |  array  |                                                 |
	+------------------+---------+-------------------------------------------------+
	| ``>destination`` |  hash   |                                                 |
	+------------------+---------+-------------------------------------------------+
	| ``>>location``   |  string |                                                 |
	+------------------+---------+-------------------------------------------------+
	| ``>>type``       |  string |                                                 |
	+------------------+---------+-------------------------------------------------+
	| ``>configured``  | boolean |                                                 |
	+------------------+---------+-------------------------------------------------+
	| ``>enabled``     | boolean |                                                 |
	+------------------+---------+-------------------------------------------------+
	| ``enabled``      | boolean |                                                 |
	+------------------+---------+-------------------------------------------------+

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"failover": {
			"locations": [],
			"destination": null,
			"configured": false,
			"enabled": false
		},
		"enabled": false
	}}


.. [1] If a user does not have either the "admin" nor "operations" role, then only Delivery Services assigned to the user's Tenant will be able to be queried with this endpoint
