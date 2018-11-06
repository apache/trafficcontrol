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

.. _to-api-deliveryservices-id-servers-eligible:

******************************
deliveryservices/{{ID}}/health
******************************

.. seealso:: :ref:`health-proto`

``GET``
=======
Retrieves the health of all Cache Groups assigned to a particular Delivery Service

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [1]_
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+-----------------+----------+--------------------------------------------------------------------------------------------------+
	| Name            | Required | Description                                                                                      |
	+=================+==========+==================================================================================================+
	| ``ID``          | yes      | The integral, unique identifier of the Delivery service for which Cache Groups will be displayed |
	+-----------------+----------+--------------------------------------------------------------------------------------------------+


Response Structure
------------------
:cachegroups: An array of objects that represent the health of each Cache Group assigned to this Delivery Service

	:name:    The name of the Cache Group represented by this object
	:offline: The number of offline cache servers within this Cache Group
	:online:  The number of online cache servers within this Cache Group

:totalOffline: Total number of offline cache servers assigned to this Delivery Service
:totalOnline:  Total number of online cache servers assigned to this Delivery Service

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"totalOffline": 0,
		"totalOnline": 1,
		"cachegroups": [
			{
				"offline": 0,
				"name": "CDN_in_a_Box_Edge",
				"online": 1
			}
		]
	}}

.. [1] Users with the roles "admin" and/or "operations" will be able to the see Cache Groups associated with *any* Delivery Services, whereas any other user will only be able to see the Cache Groups associated with Delivery Services their Tenant is allowed to see.
