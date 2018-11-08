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

.. _to-api-deliveryservices-id-capacity:

************************************
``deliveryservices/{{ID}}/capacity``
************************************

.. seealso:: :ref:`health-proto`

``GET``
=======
Retrieves the usage percentages of a servers associated with a Delivery Service

:Auth. Required: Yes
:Roles Required: "admin" or "operations"\ [1]_
:Response Type:  Object

Request Structure
-----------------
.. table:: Request Path Parameters

	+------+----------+----------------------------------------------------------------------+
	| Name | Required | Description                                                          |
	+======+==========+======================================================================+
	| id   | yes      | The integral, unique identifier for the Delivery Service of interest |
	+------+----------+----------------------------------------------------------------------+

Response Structure
------------------
:availablePercent:   The percent of servers assigned to this Delivery Service that is available - the allowed traffic level in terms of data per time period for all cache servers that remains unused
:unavailablePercent: The percent of servers assigned to this Delivery Service that is unavailable - the allowed traffic level in terms of data per time period for all cache servers that can't be used because the servers are deemed unhealthy
:utilizedPercent:    The percent of servers assigned to this Delivery Service that is currently in use - the allowed traffic level in terms of data per time period that is currently devoted to servicing requests
:maintenancePercent: The percent of servers assigned to this Delivery Service that is unavailable due to server maintenance - the allowed traffic level in terms of data per time period that is unavailable because servers have intentionally been marked offline by administrators

.. code-block:: json
	:caption: Response Example

	{ "response": {
		"availablePercent": 99.9990303030303,
		"unavailablePercent": 0,
		"utilizedPercent": 0.00096969696969697,
		"maintenancePercent": 0
	}}


.. [1] Users with the roles "admin" and/or "operations" will be able to see details for *all* Delivery Services, whereas any other user will only see details for the Delivery Services their Tenant is allowed to see.
